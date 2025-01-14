package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"strings"
)

func (ws *WSServer) parseMessage(client *wsClient, mType int, data []byte) {
	switch mType {
	case 1:
		ws.handleConnectMsg(client, data)
	case 3:
		ws.handleSendGroupMsg(client, data)
	case 5:
		ws.handleGetClientPubkeyMsg(client, data)
	case 7:
		ws.handleSendSingleMsg(client, data)
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
		Pubkey string `json:"pubkey"`
	} `json:"data"`
}

type ConnectMsgResp struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
	Data    struct {
		UUID string `json:"uuid"`
	} `json:"data"`
}

// 客户端初次连接WSS消息,发送uuid和pubKey
func (wss *WSServer) handleConnectMsg(wsc *wsClient, data []byte) {
	connMsgReq := &ConnectMsgReq{}

	err := json.Unmarshal(data, connMsgReq)
	if err != nil {
		log.Println(err)
		return
	}
	// 生成uuid
	u := uuid.New()
	my_uuid := strings.ReplaceAll(u.String(), "-", "")
	my_uuid = my_uuid[0:10]
	// 将公钥存储再WSS
	// log.Println("公钥地址：", connMsgReq.Data.Pubkey)\
	wsc.UUID = my_uuid
	wsc.Pubkey = connMsgReq.Data.Pubkey
	wss.wsClientToConn[my_uuid] = wsc
	// 返回消息
	resp_struct := &ConnectMsgResp{
		Type:    2,
		Message: "WebSocket连接成功",
		Data: struct {
			UUID string `json:"uuid"`
		}{UUID: my_uuid},
	}

	resp_json, _ := json.Marshal(resp_struct)

	wsc.WriteMessage(1, resp_json)
}

type SendGroupMsgReq struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
	Data    struct {
		SendUUID string `json:"send_uuid"`
		Content  string `json:"content"`
	} `json:"data"`
}

func (wss *WSServer) handleSendGroupMsg(wsc *wsClient, data []byte) {
	groupMsgReq := &SendGroupMsgReq{}
	err := json.Unmarshal(data, groupMsgReq)
	if err != nil {
		log.Println(err)
		return
	}
	// log.Println(groupMsgReq.Type, groupMsgReq.Message, groupMsgReq.Data.SendUUID)
	broadcastMsg := BroadcastData{
		SendUUID: groupMsgReq.Data.SendUUID,
		Content:  groupMsgReq.Data.Content,
	}
	wss.broadcast <- broadcastMsg
}

// ==================获取公钥===================
type GetClientPubkeyReq struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
	Data    struct {
		AddClient string `json:"add_client"`
	} `json:"data"`
}

type GetClientPubkeyResp struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
	Data    struct {
		AddClient string `json:"add_client"`
		Pubkey    string `json:"pubkey"`
	} `json:"data"`
}

func (wss *WSServer) handleGetClientPubkeyMsg(wsc *wsClient, data []byte) {
	clinetMsgReq := &GetClientPubkeyReq{}
	err := json.Unmarshal(data, clinetMsgReq)
	if err != nil {
		log.Println(err)
		return
	}
	//--如果查询不到，这里会出错
	pubkey := wss.wsClientToConn[clinetMsgReq.Data.AddClient].Pubkey
	clinetMsgResp := &GetClientPubkeyResp{
		Type:    6,
		Message: "响应公钥",
		Data: struct {
			AddClient string `json:"add_client"`
			Pubkey    string `json:"pubkey"`
		}{AddClient: clinetMsgReq.Data.AddClient, Pubkey: pubkey},
	}
	resp_json, _ := json.Marshal(clinetMsgResp)
	wsc.WriteMessage(1, resp_json)
	// 发送信号，响应一条消息给被连接的用户
	pubkey2 := wsc.Pubkey
	clinetMsgResp2 := &GetClientPubkeyResp{
		Type:    6,
		Message: "响应公钥",
		Data: struct {
			AddClient string `json:"add_client"`
			Pubkey    string `json:"pubkey"`
		}{AddClient: wsc.UUID, Pubkey: pubkey2},
	}
	resp_json2, _ := json.Marshal(clinetMsgResp2)

	wss.wsClientToConn[clinetMsgReq.Data.AddClient].WriteMessage(1, resp_json2)
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

// 点对点加密的消息
func (wss *WSServer) handleSendSingleMsg(wsc *wsClient, data []byte) {
	sendSingleMsgReq := &SendSingleMsgReq{}
	err := json.Unmarshal(data, sendSingleMsgReq)
	if err != nil {
		log.Println(err)
		return
	}
	wss.wsClientToConn[sendSingleMsgReq.Data.RecvUUID].WriteMessage(1, data)

}
