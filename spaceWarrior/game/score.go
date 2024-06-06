package game

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PlayerScore struct {
	Name  string
	Score int
}

func insertPlayerScore(playerName string, score int) error {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	collection := client.Database("game").Collection("players")
	_, err = collection.InsertOne(ctx, bson.M{"name": playerName, "score": score})
	if err != nil {
		return err
	}
	return nil
}

func getPlayerScores() ([]PlayerScore, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	var playerScores []PlayerScore
	collection := client.Database("game").Collection("players")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var player PlayerScore
		err := cursor.Decode(&player)
		if err != nil {
			return nil, err
		}
		playerScores = append(playerScores, player)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return playerScores, nil
}
