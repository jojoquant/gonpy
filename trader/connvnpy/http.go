package connvnpy

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func getToken(addr string) *ResTokenData {
	//这里添加post的body内容
	data := make(url.Values)
	data["username"] = []string{"vnpy"}
	data["password"] = []string{"vnpy"}

	//把post表单发送给目标服务器
	res, err := http.PostForm(fmt.Sprintf("http://%s/token", addr), data)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(string(body))
	tokenData := &ResTokenData{}
	if err := json.Unmarshal(body, tokenData); err != nil {
		log.Fatal(err)
	}

	return tokenData
}

func httpDo(method string, url string, body io.Reader, tokenData *ResTokenData, headerContentType string) []byte {
	fmt.Println("----", url, "----")
	client := &http.Client{}
	req, err := http.NewRequest(method,
		fmt.Sprintf("http://%s", url),
		body)
	if err != nil {
		log.Fatal(err)
	}

	if headerContentType != "" {
		// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Content-Type", headerContentType)
	}

	if tokenData != nil {
		x := fmt.Sprintln(tokenData.TokenType, tokenData.AccessToken)
		x = strings.Trim(x, "\n")
		req.Header.Set("Authorization", x)
	}

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
