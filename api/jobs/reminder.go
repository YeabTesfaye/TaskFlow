package jobs

import (
	"api/configs"
	"api/models"
	"api/utils"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
    currentTime := time.Now()
    oneHourFromNow := currentTime.Add(1 * time.Hour)

    // Filter tasks due within the next hour and not yet reminded
    filter := bson.M{
        "due_date": bson.M{
            "$gt":  currentTime,
            "$lte": oneHourFromNow,
        },
        "status":             bson.M{"$ne": "Completed"},
        "hour_reminder_sent": bson.M{"$ne": true},
    }

    cursor, err := taskCollection.Find(ctx, filter)
    if err != nil {
        log.Printf("Error finding tasks for reminders: %v", err)
        return
    }

    var tasks []models.Task
    if err := cursor.All(ctx, &tasks); err != nil {
        log.Printf("Error decoding tasks: %v", err)
        return
    }

    for _, task := range tasks {
        // Retrieve user email based on UserID
        userCollection := configs.GetCollection(configs.DB, "users")
        var user models.User
        err := userCollection.FindOne(ctx, bson.M{"_id": task.UserID}).Decode(&user)
        userEmail := user.Email
        if err != nil {
            log.Printf("Error retrieving user email: %v", err)
            continue
        }

        subject := fmt.Sprintf("⏳ Reminder: Task '%s' Due Soon", task.Title)
        body := fmt.Sprintf(`
            <p>Hi,</p>
            <p>Your task <strong>%s</strong> is due at <strong>%s</strong>.</p>
            <p><b>Description:</b> %s</p>
            <p>Please ensure to complete it on time.</p>
            <p>— Task Manager</p>
        `, task.Title, task.DueDate.Format("Jan 2, 2006 15:04"), task.Description)

        if err = utils.SendEmail(userEmail, subject, body); err != nil {
            log.Printf("Error sending email: %v", err)
            continue
        }

        // Mark the reminder as sent
        _, err = taskCollection.UpdateOne(
            ctx,
            bson.M{"_id": task.ID},
            bson.M{"$set": bson.M{"hour_reminder_sent": true}},
        )
        if err != nil {
            log.Printf("Error updating task reminder status: %v", err)
        }
    }
}


