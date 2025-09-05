package cmd

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


var MainCollection *mongo.Collection

func InitMongoCollection(db *mongo.Database){
	MainCollection = db.Collection("Chethu")
}
func ConnectToMongoDB(mongoURI string)(*mongo.Client,error){
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	client,err := mongo.Connect(ctx,options.Client().ApplyURI(mongoURI))
	if err != nil{
		return nil,fmt.Errorf("failed to connect mongo dB :%w",err)
	}
	if err := client.Ping(ctx,readpref.Primary());err !=nil{
		return nil,fmt.Errorf("failed to pin mongo db: %w",err)
	}
	return client,nil
}