package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gonpy/trader/util"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxDB struct {
	Host          string
	Port          int
	Username      string
	Password      string
	AuthToken     string
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
func NewInfluxDB(host string, port int, username, password, authToken string, org, bucket string, blocking bool) *InfluxDB {
	i := &InfluxDB{
		Host:      host,
		Port:      port,
		Username:  username,
		Password:  password,
		AuthToken: authToken,
		Org:       org,
		Bucket:    bucket,
		Blocking:  blocking,
	}

	serverURL := fmt.Sprintf("http://%s:%d", i.Host, i.Port)

	// 1.8 以前 authToken 用以下方式拼接
	// authToken := fmt.Sprintf("%s:%s", i.Username, i.Password)
	// authToken = "yA1dAZx9t-fn7J4fCryJurEdVC8xPQM0esSqftx6hpfT0JST0BfEnCnbFKO5lxrE-ilZBxpvTSKfK0eLsrdWaQ=="

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

// flux := fmt.Sprintf(`from(bucket:"%s")
// 	|> range(start: -7d) 
// 	|> filter(fn: (r) => r._measurement == "%s")`, i.Bucket, measurement)
func (i *InfluxDB) Query(flux string) {
	// Get query client
	queryAPI := i.client.QueryAPI(i.Org)
	// get QueryTableResult

	result, err := queryAPI.Query(
		context.Background(),
		flux,
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
			fmt.Printf("value: %v\n", result.Record())
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
		if err != nil {
			panic(err)
		}
	} else {
		i.writeAPIAsync.WritePoint(point)
	}
}

// curl --request POST http://localhost:8086/api/v2/delete/?org=example-org&bucket=example-bucket \
//   --header 'Authorization: Token <YOURAUTHTOKEN>' \
//   --header 'Content-Type: application/json' \
//   --data '{
//     "start": "2020-03-01T00:00:00Z",
//     "stop": "2020-11-14T00:00:00Z"
//   }'
func (i *InfluxDB) Delete(measurement string) {
	headerMap := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Token %s", i.AuthToken),
	}

	bodyJson, _ := json.Marshal(
		map[string]string{
			"start": "2020-03-01T00:00:00Z",
    		"stop": "2021-11-14T00:00:00Z",
			"predicate": fmt.Sprintf("_measurement=\"%s\"", measurement),
		})
	body := bytes.NewReader(bodyJson)

	res := util.HttpDo(
		"POST", 
		fmt.Sprintf("%s:%d/api/v2/delete/?org=%s&bucket=%s",i.Host, i.Port, i.Org, i.Bucket), 
		body, headerMap,
	)
	fmt.Println(res)
}

func (i *InfluxDB) Flush() {
	i.writeAPIAsync.Flush()
}

func (i *InfluxDB) Close() {
	i.client.Close()
}
