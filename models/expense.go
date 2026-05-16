package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	Description string             `bson:"description" json:"description"`
	Amount      float64            `bson:"amount" json:"amount"`
	Category    string             `bson:"category" json:"category"`
	Date        time.Time          `bson:"date" json:"date"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}
