package chat

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"
)

/****************************************************************/

type MapPool struct {
	value map[string]interface{}
	lock  *sync.RWMutex
}

func NewMapPool() *MapPool {
	return &MapPool{make(map[string]interface{}), new(sync.RWMutex)}
}
func (c *MapPool) Set(key string, value interface{}) {
	defer c.lock.Unlock()
	c.lock.Lock()
	c.value[key] = value
}
func (c *MapPool) Del(key string) {
	defer c.lock.Unlock()
	c.lock.Lock()
	delete(c.value, key)
}
func (c *MapPool) Get(key string) interface{} {
	defer c.lock.RUnlock()
	c.lock.RLock()
	return c.value[key]
}

var (
	Sub        = NewMapPool()
	HistoryMsg = NewMapPool()
	RoomList   = []string{"roomzh_cn", "roomzh_tw", "roomen_us", "roomservice"}
)

func init() {
	for _, room := range RoomList {
		Sub.Set(room, make(Clients))
		msgs := make(map[uint32][]ChatInfo)
		HistoryMsg.Set(room, msgs)
	}
}

type HeartData struct {
	Pong string `json:"pong,omitempty"`
	Ping string `json:"ping,omitempty"`
}

type Request struct {
	url.Values
}

func (c *Request) Clients() (Clients, error) {
	roomkey := c.Get("sub")
	if roomkey == "" {
		return nil, errors.New("房号不能为空")
	}
	value := Sub.Get(roomkey)
	if value == nil {
		return nil, errors.New("聊天室不存在")
	}
	return value.(Clients), nil
}

type Response struct {
	// Sub  string      `json:"sub",omitempty`
	Task string      `json:"task,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type Result struct {
	Id     string      `json:"id,omitempty",`
	Result string      `json:"result,omitempty"`
	Chats  interface{} `json:"chats,omitempty"`
	Subs   interface{} `json:"subs,omitempty"`
}

type ChatInfo struct {
	Uid    uint32   `json:"uid,omitempty"`
	Sub    string   `json:"sub,omitempty"`
	Time   int64    `json:"time,omitempty"`
	Msg    string   `json:"msg,omitempty"`
	Action string   `json:"action,omitempty"`
	Uids   []uint32 `json:"uids,omitempty"`
}

type Assign struct {
	Clients Clients
	Data    interface{}
}

func NewResponse(list ...string) string {
	values := make(url.Values)
	var key string
	for index, info := range list {
		if index%2 == 0 {
			key = info
		}
		values.Set(key, info)
	}
	return values.Encode()
}

func (c *Client) Protocol(query string) error {
	var err error
	values, err := url.ParseQuery(query)
	if err != nil {
		return err
	}
	req := new(Request)
	req.Values = values
	log.Println("values", req.Values)
	c.errCnt = 0
	if pong := req.Get("pong"); pong != "" {
		c.errCnt = 0
		return nil
	}
	if ping := req.Get("ping"); ping != "" {
		c.Pong(ping)
		return nil
	}

	switch req.Get("task") {
	case "join":
		err = c.Join(req)
	case "msg":
		err = c.Msg(req)
	case "exit":
		err = c.Exit(req)
	case "sub":
		err = c.Sub(req)
	// case "image":
	// 	err = c.Image(req, time.Now())
	default:
		return errors.New("unknow type")
	}

	// handle above error
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Ping() error {
	c.Write(HeartData{Ping: fmt.Sprintf("%d", c.expiryTime)})
	return nil
}

func (c *Client) Pong(ping string) error {
	c.Write(HeartData{Pong: ping})
	return nil
}

func (c *Client) Sub(req *Request) error {
	result := Result{Id: req.Get("id"), Result: "success", Subs: RoomList}
	c.Write(&Response{"sub", result})
	return nil
}

func (c *Client) Join(req *Request) error {

	uidkey := req.Get("uid")
	if uidkey == "" {
		return errors.New("UID不能为空")
	}
	var uid uint32
	fmt.Sscanf(uidkey, "%d", &uid)
	c.uid = uid

	clients, err := req.Clients()
	if err != nil {
		return err
	}
	clients[c] = struct{}{}
	roomkey := req.Get("sub")
	Sub.Set(roomkey, clients)

	result := Result{Id: req.Get("id"), Result: "success"}
	index := HistoryMsg.Get(roomkey).(map[uint32][]ChatInfo)

	if roomkey == "roomservice" {
		result.Chats = index[c.uid]
		c.Write(&Response{"join", result})
	} else {
		result.Chats = index[0]
		c.Write(&Response{"join", result})
		chat := ChatInfo{Uid: c.uid, Sub: req.Get("sub"), Action: "join", Uids: clients.Uids()}
		c.server.Assign(&Assign{clients, Response{"userupdate", chat}})
	}
	// c.server.Broadcast(NewResponse("type", "join", "uid", fmt.Sprintf("%d", c.uid), "uids", room.All(), "msg", fmt.Sprintf("%d 加入聊天室", c.uid)))

	return nil
}

func (c *Client) Exit(req *Request) error {
	clients, err := req.Clients()
	if err != nil {
		return err
	}
	delete(clients, c)
	Sub.Set(req.Get("sub"), clients)
	// c.server.Broadcast(NewResponse("type", "exit", "uid", fmt.Sprintf("%d", c.uid), "uids", room.All(), "msg", fmt.Sprintf("%d 加入聊天室", c.uid)))
	result := Result{Id: req.Get("id"), Result: "success"}
	c.Write(&Response{"exit", result})

	rep := Response{"userupdate", &ChatInfo{Uid: c.uid, Sub: req.Get("sub"), Action: "exit", Uids: clients.Uids()}}
	c.server.Assign(&Assign{clients, rep})
	return nil
}

func (c *Client) CloseExit() {
	for _, room := range RoomList {
		clients := Sub.Get(room).(Clients)
		delete(clients, c)
		Sub.Set(room, clients)
		rep := Response{"userupdate", &ChatInfo{Uid: c.uid, Sub: room, Action: "exit", Uids: clients.Uids()}}
		c.server.Assign(&Assign{clients, rep})
	}
}

func (c *Client) Msg(req *Request) error {

	msg := req.Get("msg")
	if msg == "" {
		errors.New("信息不能为空")
		return nil
	}

	clients, err := req.Clients()
	if err != nil {
		return err
	}

	if time.Now().Unix()-c.sendTime < 1 {
		c.sendTime = time.Now().Unix()
		errors.New("您说话太频繁了")
		// c.Write(NewResponse("type", "error", "msg", "您说话太频繁了"))
		return nil
	}
	c.sendTime = time.Now().Unix()
	roomkey := req.Get("sub")
	chat := ChatInfo{Uid: c.uid, Sub: roomkey, Time: c.sendTime, Msg: msg}
	rep := Response{"msg", chat}
	index := HistoryMsg.Get(roomkey).(map[uint32][]ChatInfo)
	if roomkey == "roomservice" {
		index[c.uid] = append(index[c.uid], chat)
		HistoryMsg.Set(roomkey, index)
		// c.Write(NewResponse("type", "msg", "uid", fmt.Sprintf("%d", c.uid), "room", roomkey, "time", fmt.Sprintf("%d", c.sendTime), "msg", msg))
		c.Write(&rep)
	} else {
		index[0] = append(index[0], chat)
		HistoryMsg.Set(roomkey, index)
		// c.server.Broadcast(NewResponse("type", "msg", "uid", fmt.Sprintf("%d", c.uid), "room", roomkey, "time", fmt.Sprintf("%d", c.sendTime), "msg", msg))
		c.server.Assign(&Assign{clients, rep})
	}
	return nil
}

// // Image
// func (c *Client) Image(req *Request, now time.Time) error {
// 	dateDir := now.Format("060102")
// 	dirPath := filepath.Join("upload", dateDir)

// 	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
// 		err = os.Mkdir(dirPath, 0777)
// 		if err != nil {
// 			return err
// 		}

// 		// create index.html to forbidden index directory
// 		fp, err := os.Create(filepath.Join(dirPath, "index.html"))
// 		if err != nil {
// 			return err
// 		}
// 		fp.Close()
// 	}

// 	resourceId := uuid.New()

// 	fw, err := os.Create(filepath.Join(dirPath, resourceId))
// 	if err != nil {
// 		return err
// 	}
// 	defer fw.Close()

// 	outBuf := make([]byte, len(req.Body))
// 	n, err := base64.StdEncoding.Decode(outBuf, req.Body)
// 	if err != nil {
// 		return err
// 	}

// 	// write to file
// 	buffer := bytes.NewBuffer(outBuf[:n])
// 	_, err = buffer.WriteTo(fw)
// 	if err != nil {
// 		return err
// 	}

// 	// every resource has a uuid pathid
// 	pathId := fmt.Sprintf("%s/%s", dateDir, resourceId)

// 	c.Send(NewResponse("image", nil, "pathid", pathId, "index", req.Values.Get("index")).EncodeBytes())

// 	return nil
// }
