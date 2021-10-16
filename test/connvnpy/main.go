package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// ClientManager is a websocket manager
type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// Client is a websocket client
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

// Message is an object for websocket message which is mapped to json type
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

// Manager define a ws server manager
var Manager = ClientManager{
	Broadcast:  make(chan []byte, 1),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[*Client]bool),
}

// Start is to start a ws server
func (manager *ClientManager) Start() {

	for {
		select {
		// websocket
		case conn := <-manager.Register:
			manager.Clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{Content: "A new socket has connected."})
			// log.Println("In ClientManager Start. Once Registed, send json msg:", string(jsonMessage))
			log.Println("(Register)Clients of Manager : ", manager.Clients)
			manager.Send(jsonMessage, conn)
		case conn := <-manager.Unregister:
			if _, ok := manager.Clients[conn]; ok {
				close(conn.Send)
				delete(manager.Clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "A socket has disconnected."})
				log.Println("In ClientManager Start. Unregisted, send json msg:", string(jsonMessage))
				log.Println("(Unregister)Clients of Manager : ", manager.Clients)
				manager.Send(jsonMessage, conn)
			}
		case message := <-manager.Broadcast:
			log.Println("message from broadcast channel is:", string(message))
			for conn := range manager.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}

// Send is to send ws message to ws client
func (manager *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if conn != ignore {
			log.Println("In ClientManager Send msg to conn send, send json msg:", string(message))
			conn.Send <- message
		}
	}
}

// Read 读取注册的连接发送过来的信息, 参数中增加消息处理中间件函数
func (c *Client) Read(MsgMidFunc func(*string, *map[string]string), argsMap *map[string]string) {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		messageStr := string(message)
		MsgMidFunc(&messageStr, argsMap)

		jsonMessage, _ := json.Marshal(&Message{Sender: c.ID, Content: messageStr})
		Manager.Broadcast <- jsonMessage
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////////

type ResTokenData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

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

func main() {
	addr := "192.168.0.104:8000"
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	tokenData := getToken(addr)
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws/", RawQuery: fmt.Sprintf("token=%s", tokenData.AccessToken)}
	fmt.Printf("connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}

			if t.Second()%3 == 0 {

				// reqBody := httpDo("GET", fmt.Sprintf("%s/position",addr), strings.NewReader(fmt.Sprintf("token=%s", tokenData.AccessToken)))
				reqBody := httpDo("GET", fmt.Sprintf("%s/account", addr), nil, tokenData, "application/json;charset=utf-8")

				fmt.Println(string(reqBody))
				// c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			}

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
