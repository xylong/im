package ctrl

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net/http"
	"strconv"
	"sync"
)

const (
	CMD_SINGLE_MSG = 10
	CMD_ROOM_MSG   = 11
	CMD_HEART      = 0
)

var (
	clientMap map[int64]*Node = make(map[int64]*Node, 0)
	rwlocker  sync.RWMutex
)

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"`           //消息ID
	Userid  int64  `json:"userid,omitempty" form:"userid"`   //谁发的
	Cmd     int    `json:"cmd,omitempty" form:"cmd"`         //群聊还是私聊
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`     //对端用户ID/群ID
	Media   int    `json:"media,omitempty" form:"media"`     //消息按照什么样式展示
	Content string `json:"content,omitempty" form:"content"` //消息的内容
	Pic     string `json:"pic,omitempty" form:"pic"`         //预览图片
	Url     string `json:"url,omitempty" form:"url"`         //服务的URL
	Memo    string `json:"memo,omitempty" form:"memo"`       //简单描述
	Amount  int    `json:"amount,omitempty" form:"amount"`   //其他和数字相关的
}

// Chat 聊天
// ws://127.0.0.1/chat?id=1&token=123
func Chat(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		conn *websocket.Conn
	)
	// 检验接入是否合法
	query := r.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64)
	ok := checkToken(userId, token)

	if conn, err = (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return ok
		},
	}).Upgrade(w, r, nil); err != nil {
		log.Println(err.Error())
		return
	}

	// 获取conn
	node := &Node{
		Conn: conn,
		// 并行转串行
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	comIds := contactService.SearchComunityIds(userId)
	for _, v := range comIds {
		node.GroupSets.Add(v)
	}
	// user和node绑定
	rwlocker.Lock()
	clientMap[userId] = node
	rwlocker.Unlock()
	// 发送
	go send(node)
	// 接收
	go receive(node)
	sendMsg(userId, []byte("hello world!"))
}

// send 发送消息
func send(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			if err := node.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// receive 接收消息
func receive(node *Node) {
	for {
		if _, data, err := node.Conn.ReadMessage(); err == nil {
			dispatch(data)
		} else {
			log.Println(err.Error())
			return
		}
	}
}

// checkToken 检测token是否有效
func checkToken(userId int64, token string) bool {
	user := userService.Find(userId)
	return user.Token == token
}

// sendMsg 发送消息
func sendMsg(userId int64, msg []byte) {
	rwlocker.RLock()
	node, ok := clientMap[userId]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

// dispatch 派遣
func dispatch(data []byte) {
	msg := Message{}
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Println(err.Error())
		return
	}
	switch msg.Cmd {
	case CMD_SINGLE_MSG:
		sendMsg(msg.Dstid, data)
	case CMD_ROOM_MSG:
		for _, v := range clientMap {
			if v.GroupSets.Has(msg.Dstid) {
				v.DataQueue <- data
			}
		}
	case CMD_HEART:

	}
}
