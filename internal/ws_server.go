package internal

import (
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
}

func (ws *WSServer) readMsg(conn *wsClient) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("send"))
	if err != nil {
		log.Println("readMsg send error:", err)
	}
}
