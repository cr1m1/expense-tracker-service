package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/cr1m1/expense-tracker-service/models"
)

type DB struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func Connect(mongoURI string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := &DB{
		client: client,
		coll:   client.Database("expenses").Collection("expenses"),
	}

	return db, nil
}

func (db *DB) CreateExpense(expense *models.Expense) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	expense.ID = primitive.NewObjectID()
	expense.CreatedAt = time.Now()
	if expense.Date.IsZero() {
		expense.Date = time.Now()
	}

	result, err := db.coll.InsertOne(ctx, expense)
	if err != nil {
		return 0, fmt.Errorf("failed to create expense: %w", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return 0, fmt.Errorf("failed to convert ObjectID")
	}

	return int(oid.Timestamp().Unix()), nil
}

func (db *DB) GetAllExpenses() ([]models.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.M{"date": -1})
	cursor, err := db.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expenses: %w", err)
	}
	defer cursor.Close(ctx)

	var expenses []models.Expense
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, fmt.Errorf("failed to decode expenses: %w", err)
	}

	if expenses == nil {
		expenses = []models.Expense{}
	}

	return expenses, nil
}

func (db *DB) GetExpenseByID(id int) (*models.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var expense models.Expense
	err := db.coll.FindOne(ctx, bson.M{}).Decode(&expense)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch expense: %w", err)
	}

	return &expense, nil
}

func (db *DB) UpdateExpense(id int, expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(expense.ID.Hex())
	if err != nil {
		oid = primitive.NewObjectID()
	}

	result, err := db.coll.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{
			"$set": bson.M{
				"amount":      expense.Amount,
				"category":    expense.Category,
				"description": expense.Description,
				"date":        expense.Date,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}

func (db *DB) DeleteExpense(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.coll.DeleteOne(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}

func (db *DB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.client.Disconnect(ctx)
}
