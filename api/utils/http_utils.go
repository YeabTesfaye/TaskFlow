package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SendError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func SendJSON(w http.ResponseWriter, data any) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func VerifyOwnership(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID, userID string, result interface{}) error {
	return collection.FindOne(ctx, bson.M{"_id": id, "user_id": userID}).Decode(result)
}

func GetObjectIDFromRequest(r *http.Request, param string) (primitive.ObjectID, error) {
	params := mux.Vars(r)
	return primitive.ObjectIDFromHex(params[param])
}
