package internal

import (
	"encoding/json"
	"log"
)

func (ws *WSServer) parseMessage(client *wsClient, mType int, data []byte) {
	switch mType {
	case 1:
		ws.handleConnectMsg(client, data)
	default:
		return
	}
}

// 由于先要解析获取消息类型，才能进一步解析data数据，但是data以string的形式进一步解析，无法解析。
// 解决方法是，先丢弃data字段解析一遍，根据消息type再全部解析一遍
type ConnectMsgReq struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
	Data    struct {
		UUID   string `json:"uuid"`
		Pubkey string `json:"pubkey"`
	} `json:"data"`
}

// 客户端初次连接WSS消息,发送uuid和pubKey
func (ws *WSServer) handleConnectMsg(wsc *wsClient, data []byte) {
	connMsgReq := &ConnectMsgReq{}

	err := json.Unmarshal(data, connMsgReq)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(connMsgReq.Data.UUID, connMsgReq.Data.Pubkey)

}

type SendSingleMsgReq struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
	Data    struct {
		SendUUID string `json:"send_uuid"`
		RecvUUID string `json:"recv_uuid"`
		Content  string `json:"content"`
	} `json:"data"`
}

// 对方不在线发送的消息
func (ws *WSServer) handleSendSingleMsg(wsc *wsClient, data []byte) {
	sendSingleMsgReq := &SendSingleMsgReq{}
	err := json.Unmarshal(data, sendSingleMsgReq)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(sendSingleMsgReq.Data.SendUUID, sendSingleMsgReq.Data.RecvUUID)
}
