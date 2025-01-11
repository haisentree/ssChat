package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

func main() {

	u := url.URL{Scheme: "ws", Host: "localhost:8081", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	// 将json数据反序列化
	jsonExample := `{
	"type": 1,
	"message": "WebSocket连接成功",
	"data": {
		"uuid": "123",
		"pubkey": "123"
	}
}`
	type JSONData struct {
		Type    int    `json:"type"`
		Message string `json:"message"`
		Data    struct {
			UUID   string `json:"uuid"`
			Pubkey string `json:"pubkey"`
		} `json:"data"`
	}
	jsonStruct := &JSONData{}
	_ = json.Unmarshal([]byte(jsonExample), jsonStruct)
	fmt.Println(jsonStruct.Data.UUID)
	err = c.WriteMessage(websocket.TextMessage, []byte(jsonExample))
	if err != nil {
		log.Println("write:", err)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}

}
