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
		wsMaxConnNum:   2048,
		wsUpGrader:     upgrader,
		wsClientToConn: make(map[string]*wsClient),
		broadcast:      make(chan BroadcastData),
	}
}

type WSServer struct {
	wsMaxConnNum   int
	wsUpGrader     *websocket.Upgrader
	wsClientToConn map[string]*wsClient
	broadcast      chan BroadcastData
	connect        chan string
}

type wsClient struct {
	*websocket.Conn
	UUID   string
	Pubkey string
}

type SendBroadcastMsgResp struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
	Data    struct {
		SendUUID string `json:"send_uuid"`
		Content  string `json:"content"`
	} `json:"data"`
}

type BroadcastData struct {
	SendUUID string
	Content  string
}

// 多线程使用一个函数没问题，例如http的句柄函数
// 多线程改变一个实例的数据可能存在问题，比如WSServer，多个连接一起来会改变wsClientToConn字段的值。后期再看，可能需要加锁。
// chan map类型：https://blog.csdn.net/yamanda/article/details/100057006
func (wss *WSServer) ListenBroadcast() {
	go func() {
		for {
			select {
			case msg_data := <-wss.broadcast:
				broadcastMsg := SendBroadcastMsgResp{
					Type:    4,
					Message: "广播回写",
					Data: struct {
						SendUUID string `json:"send_uuid"`
						Content  string `json:"content"`
					}{SendUUID: msg_data.SendUUID, Content: msg_data.Content},
				}

				msg_send, _ := json.Marshal(broadcastMsg)
				for _, conn := range wss.wsClientToConn {
					err := conn.WriteMessage(websocket.TextMessage, msg_send)
					// log.Println("广播消息发送给:", ud)
					if err != nil {
						log.Println("broadcast error")
					}
				}

			}
		}

	}()
}

func (ws *WSServer) wsHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if err != nil {
		log.Println("readMsg send error:", err)
	}
	wsc := &wsClient{
		Conn: conn,
	}
	// 如果一个客户端读取到多条消息，这情况没有用多线程去处理
	go ws.readMessage(wsc)
}

// 通用消息头解析，解析出消息类型和data，传入logic进行解析data
type CommonMessage struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}

func (wss *WSServer) readMessage(wsc *wsClient) {
	for {
		_, message, err := wsc.Conn.ReadMessage()
		// 需要知道web端，刷新网页后，客户端是主动关闭、异常关闭、还是没关闭
		if err != nil {
			// log.Println("readMsg read error:", err)
			// 刷新网页造成的关闭
			delete(wss.wsClientToConn, wsc.UUID)

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("异常关闭：", err.Error())
			}
			return
		}

		// log.Printf("recv: %s", message)

		//err = wsc.Conn.WriteMessage(websocket.TextMessage, message)
		//if err != nil {
		//	log.Println("readMsg send error:", err)
		//	continue
		//}
		// 如果客户端下线了，需要再WSS上删除客户端信息
		//if msgType != websocket.CloseMessage {
		//
		//}

		// 对消息进行解析
		cmsg := &CommonMessage{}
		err = json.Unmarshal(message, cmsg)

		if err != nil {
			log.Println("readMsg unmarshal error:", err)
			// 消息json规范化后，不能这这样写入错误信息
			err = wsc.Conn.WriteMessage(websocket.TextMessage, []byte("readMsg unmarshal error"))
		} else {
			wss.parseMessage(wsc, cmsg.Type, message)
		}

	}
}
