package internal

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func NewWSServer() *WSServer {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &WSServer{
		wsUpGrader: upgrader,
	}
}

type WSServer struct {
	wsMaxConnNum   int
	wsUpGrader     *websocket.Upgrader
	wsClientToConn map[uint64]map[uint8]*wsClient
}

type wsClient struct {
	*websocket.Conn
	clinetType uint8
}

func (ws *WSServer) wsHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("send"))
	if err != nil {
		log.Println("readMsg send error:", err)
	}
	wsc := &wsClient{
		conn,
		1,
	}
	go ws.readMessage(wsc)
}

// 通用消息头解析，解析出消息类型和data，传入logic进行解析data
type CommonMessage struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}

func (ws *WSServer) readMessage(wsc *wsClient) {
	for {
		msgType, message, err := wsc.Conn.ReadMessage()
		if err != nil {
			log.Println("readMsg read error:", err)
			return
		}
		log.Printf("recv: %s", message)
		log.Printf("msgType: %d", msgType)

		err = wsc.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("readMsg send error:", err)
		}

		// 对消息进行解析
		cmsg := &CommonMessage{}
		err = json.Unmarshal(message, cmsg)

		if err != nil {
			log.Println("readMsg unmarshal error:", err)
			err = wsc.Conn.WriteMessage(websocket.TextMessage, []byte("readMsg unmarshal error"))
		} else {
			wss.parseMessage(wsc, cmsg.Type, message)
		}

	}

}
