package database

import (
	"context"
	"fmt"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxDB struct {
	Host          string
	Port          int
	Username      string
	Password      string
	Org           string
	Bucket        string
	Blocking      bool
	client        influxdb2.Client
	writeAPIAsync api.WriteAPI
	writeAPISync  api.WriteAPIBlocking
}

// host:     "127.0.0.1:8086",
// Username: "admin",
// Password: "",
func NewInfluxDB(host string, port int, username, password string, org, bucket string, blocking bool) *InfluxDB {
	i := &InfluxDB{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Org:      org,
		Bucket:   bucket,
		Blocking: blocking,
	}

	serverURL := fmt.Sprintf("http://%s:%d", i.Host, i.Port)
	authToken := fmt.Sprintf("%s:%s", i.Username, i.Password)
	authToken = "yA1dAZx9t-fn7J4fCryJurEdVC8xPQM0esSqftx6hpfT0JST0BfEnCnbFKO5lxrE-ilZBxpvTSKfK0eLsrdWaQ=="
	
	if blocking {
		i.client = influxdb2.NewClient(
			serverURL,
			authToken,
		)
		i.writeAPISync = i.client.WriteAPIBlocking(i.Org, i.Bucket)
	} else {
		// Create a new client using an InfluxDB server base URL and an authentication token
		// and set batch size to 20
		i.client = influxdb2.NewClientWithOptions(
			serverURL,
			authToken,
			influxdb2.DefaultOptions().SetBatchSize(20),
		)
		i.writeAPIAsync = i.client.WriteAPI(i.Org, i.Bucket)
	}

	return i
}

func (i *InfluxDB) Query() {
	// Get query client
	queryAPI := i.client.QueryAPI("my-org")
	// get QueryTableResult
	result, err := queryAPI.Query(
		context.Background(),
		`from(bucket:"my-bucket")
		|> range(start: -1h) 
		|> filter(fn: (r) => r._measurement == "stat")`,
	)

	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			fmt.Printf("value: %v\n", result.Record().Value())
		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
}

// write asynchronously 在循环外补
//
// Force all unwritten data to be sent
//
// writeAPI.Flush()
//
// Ensures background processes finishes
//
// client.Close()
func (i *InfluxDB) Write(point *write.Point) {
	if i.Blocking {
		err := i.writeAPISync.WritePoint(context.Background(), point)
		if err!=nil{
			panic(err)
		}
	} else {
		i.writeAPIAsync.WritePoint(point)
	}
}

func (i *InfluxDB) Flush() {
	i.writeAPIAsync.Flush()
}

func (i *InfluxDB) Close() {
	i.client.Close()
}
