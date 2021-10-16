package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	URI    string
	client *mongo.Client
	ctx    context.Context
}

type QueryParam struct {
	Db         string
	Collection string
	Filter     bson.M
}

type InsertParam struct {
	Db         string
	Collection string
	Ordered    bool
	Docs       []interface{} // for InsertMany
	Doc        bson.M        // for InsertOne
}

func NewMongoDB(host string, port int, username, password string) *MongoDB {
	m := &MongoDB{
		URI: fmt.Sprintf("mongodb://%s:%d", host, port),
		ctx: context.TODO(),
	}

	credential := options.Credential{
		Username: username,
		Password: password,
	}
	// Set client options
	// "mongodb://192.168.0.113:27017"
	clientOptions := options.Client().ApplyURI(m.URI).SetAuth(credential)

	// Connect to MongoDB
	client, err := mongo.Connect(m.ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(m.ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to MongoDB!")

	m.client = client
	return m
}

func (m *MongoDB) Query(q *QueryParam) []*BarData {
	collection := m.client.Database(q.Db).Collection(q.Collection)
	cur, err := collection.Find(m.ctx, q.Filter)
	if err != nil {
		log.Fatal(err)
	}

	var r []*BarData

	// 经过简单对比测试, 单独cur 和 cur.All 差距不大, 后续需要继续对比测试
	// for cur.Next(m.ctx) {
	// 	var elem BarData
	// 	err := cur.Decode(&elem)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	r = append(r, &elem)
	// }

	if err = cur.All(m.ctx, &r); err != nil {
		log.Fatal(err)
	}

	// 完成后关闭游标
	cur.Close(m.ctx)
	return r
}

func (m *MongoDB) InsertMany(i *InsertParam) {
	collection := m.client.Database(i.Db).Collection(i.Collection)
	opts := options.InsertMany().SetOrdered(i.Ordered)

	_, err := collection.InsertMany(m.ctx, i.Docs, opts)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *MongoDB) InsertOne(i *InsertParam) {
	collection := m.client.Database(i.Db).Collection(i.Collection)

	_, err := collection.InsertOne(m.ctx, i.Doc)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *MongoDB) Close() {
	if err := m.client.Disconnect(m.ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to MongoDB closed.")
}
