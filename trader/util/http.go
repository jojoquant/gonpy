package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func HttpDo(method string, url string, body io.Reader, headerMap map[string]string) []byte {

	client := &http.Client{}
	req, err := http.NewRequest(method,
		fmt.Sprintf("http://%s", url),
		body)
	if err != nil {
		log.Fatal(err)
	}

	if len(headerMap) != 0{
		for k,v := range headerMap{
			req.Header.Set(k, v)
		}
	}
	
	// if headerContentType != "" {
	// 	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// 	req.Header.Set("Content-Type", headerContentType)
	// }

	// if token != "" {
	// 	req.Header.Set("Authorization", token)
	// }

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	result_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result_body))
	return result_body
}