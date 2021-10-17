package main

import (
	"context"
	"fmt"
	"gonpy/trader/database"
	"gonpy/trader/util"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SHFE_d_AUL8 struct {
	Open         float64
	High         float64
	Low          float64
	Close        float64
	OpenInterest int
	Volume       int
	Datetime     time.Time
}

func Temp() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://192.168.0.113:27017")

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

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("vnpy").Collection("SHFE_d_AUL8")

	var result SHFE_d_AUL8

	filter := bson.D{{}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	var result2 []*SHFE_d_AUL8
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem SHFE_d_AUL8
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		result2 = append(result2, &elem)
		fmt.Printf("Found multiple documents (array of pointers): %#v\n", elem)
	}

	// 完成后关闭游标
	cur.Close(context.TODO())
	// fmt.Printf("Found multiple documents (array of pointers): %#v\n", result2)

}

func main() {
	// m := database.NewMongoDB("192.168.0.113", 27017)
	// fmt.Println(m)
	// r := m.Query(&database.QueryParam{Db: "vnpy", Collection: "SHFE_d_AUL8", Filter: bson.M{}})
	// fmt.Println(r)
	// fmt.Println("r length: ", len(r), r[0], r[0].Close, r[0].Datetime)

	// 创建数据源和展示数据库
	host := "localhost"
	dk := util.GetDockerComposeYml("../../docker-compose.yml")

	// authToken := "yA1dAZx9t-fn7J4fCryJurEdVC8xPQM0esSqftx6hpfT0JST0BfEnCnbFKO5lxrE-ilZBxpvTSKfK0eLsrdWaQ=="
	// authToken := "lAAs7a0buXNb88a4rb_ZKB1M9SxKqHKvGl_fSkwKAAIH-80NAw4SbzEpOPAWUgqr_KlxgurP4cNqolHggmN0pg=="
	displayDB := database.NewInfluxDB(
		host, dk.Services.Influxdb.PortSrc,
		dk.Services.Influxdb.Env.Username, dk.Services.Influxdb.Env.Password,
		dk.Services.Influxdb.Env.AdminToken,
		dk.Services.Influxdb.Env.Org, dk.Services.Influxdb.Env.Bucket,
		false,
	)

	measurement := "bar_data_DualMA"
	flux := fmt.Sprintf(`from(bucket:"%s")
	|> range(start: -7d) 
	|> filter(fn: (r) => r._measurement == "%s")
	|> filter(fn: (r) => r._field == "close_price")`, displayDB.Bucket, measurement)
	displayDB.Query(flux)

	// displayDB.Delete(measurement)

	displayDB.Close()
}
