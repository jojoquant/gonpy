package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	URI    string
	client *mongo.Client
}

type QueryParam struct{
	Db string
	Collection string
	Filter bson.D
}

type BarCollection struct {
	Open         float64
	High         float64
	Low          float64
	Close        float64
	OpenInterest int
	Volume       int
	Datetime     time.Time
}

type TickCollection struct{}

func NewMongoDB(host string, port int)*MongoDB{
	m := &MongoDB{
		URI: fmt.Sprintf("mongodb://%s:%d",host, port),
	}
	// Set client options
	// "mongodb://192.168.0.113:27017"
	clientOptions := options.Client().ApplyURI(m.URI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	
	m.client = client
	return m
}



func (m *MongoDB) Query(q *QueryParam) []*BarCollection {
	collection := m.client.Database(q.Db).Collection(q.Collection)
	cur, err := collection.Find(context.TODO(), q.Filter)
	if err != nil {
		log.Fatal(err)
	}

	var r []*BarCollection
	for cur.Next(context.TODO()) {
		var elem BarCollection
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, &elem)
	}

	// 完成后关闭游标
	cur.Close(context.TODO())
	return r
}
