package jobs

import (
	"api/configs"
	"api/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var taskCollection = configs.GetCollection(configs.DB, "tasks")

func StartReminderJob() {
	go func() {
		for {
			checkReminders()
			time.Sleep(5 * time.Minute) // Check every 5 minutes
		}
	}()
}

func checkReminders() {
	ctx := context.Background()

	// Find tasks with unsent reminders that are due
	filter := bson.M{
		"reminders": bson.M{
			"$elemMatch": bson.M{
				"reminder_time": bson.M{"$lte": time.Now()},
				"sent":          false,
			},
		},
		"status": bson.M{"$ne": "Completed"},
	}

	cursor, err := taskCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error finding tasks for reminders: %v", err)
		return
	}

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		log.Printf("Error decoding tasks: %v", err)
		return
	}

	for _, task := range tasks {
		// Here you would implement your notification system
		// For example, sending emails, push notifications, etc.
		log.Printf("Sending reminder for task: %s", task.Title)

		// Mark reminder as sent
		_, err = taskCollection.UpdateOne(
			ctx,
			bson.M{"_id": task.ID},
			bson.M{
				"$set": bson.M{
					"reminders.$[elem].sent": true,
				},
			},
			options.Update().SetArrayFilters(
				options.ArrayFilters{
					Filters: []interface{}{
						bson.M{"elem.reminder_time": bson.M{"$lte": time.Now()}},
					},
				},
			),
		)

		if err != nil {
			log.Printf("Error updating reminder status: %v", err)
		}
	}
}
