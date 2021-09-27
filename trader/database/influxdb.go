package database

import (
	"context"
	"fmt"

	"github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDB struct {
	Host string
	Port int
	Username string
	Password string
	client influxdb2.Client
}

// host:     "127.0.0.1:8086",
// Username: "admin",
// Password: "",
func NewInfluxDB(host string, port int, username, password string) *InfluxDB {
	i := &InfluxDB{
		Host: host,
		Port: port,
		Username: username,
		Password: password,
	}

	i.client = influxdb2.NewClient(
		fmt.Sprintf("http://%s:%d", i.Host, i.Port), 
		fmt.Sprintf("%s:%s",i.Username,i.Password),
	)

	return i
}

func(i *InfluxDB)Query(){
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

func(i *InfluxDB)Close(){
	i.client.Close()
}

