package chat

import (
	"io"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

type Client struct {
	server      *Server
	ws          *websocket.Conn
	writeSignal chan interface{}
	doneSignal  chan bool
	uid         uint32
	errCnt      int
	expiryTime  int64
	sendTime    int64
}

func NewClient(conn *websocket.Conn, server *Server) *Client {
	return &Client{
		server,
		conn,
		make(chan interface{}, 10),
		make(chan bool, 10),
		0,
		0,
		0,
		0,
	}
}

func (c *Client) Write(data interface{}) {
	c.writeSignal <- data
}

func (c *Client) listenWrite() {
	heartTicker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-c.doneSignal:
			c.doneSignal <- true
			c.server.Del(c)
			return
		case <-heartTicker.C:
			if c.errCnt > 4 {
				heartTicker.Stop()
				c.doneSignal <- true
			}
			// c.errCnt++
			c.expiryTime = time.Now().Unix()
			c.Heart()
		case data := <-c.writeSignal:
			if err := websocket.JSON.Send(c.ws, data); err != nil {
				log.Printf("Websocket send data error:%v,data:%v", err, data)
			}
			// if err := websocket.Message.Send(c.ws, data); err != nil {
			// 	log.Printf("Websocket send data error:%v,data:%v", err, data)
			// }
		}
	}
	close(c.writeSignal)
}

func (c *Client) listenRead() {
	for {
		select {
		case <-c.doneSignal:
			c.doneSignal <- true
			return
		default:
			var query string
			err := websocket.Message.Receive(c.ws, &query)
			if err == io.EOF {
				c.doneSignal <- true
				log.Printf("Client %v read data EOF:%v", &c, err)
			} else if err != nil {
				log.Printf("Client %v read data error:%v", &c, err)
			} else if err := c.Protocol(query); err != nil {
				// c.Write(NewResponse("type", "error", "msg", err.Error()))
				c.Write(&Response{"error", err.Error()})
			}
		}
	}
}

// func (c *Client) listenHeart() {

// 	heartTicker := time.NewTicker(time.Second * 5)
// 	for {
// 		select {
// 		case <-heartTicker.C:
// 			if c.hearts.errCnt > 3 {
// 				heartTicker.Stop()
// 				return
// 			}

// 			if time.Now().Unix()-c.hearts.ping > 9 {
// 				c.hearts.ping = time.Now().Unix()
// 				c.Write(fmt.Sprintf(`{"ping":%d}`, c.hearts.ping)})
// 			} else if time.Now().Unix()-c.hearts.ping > 4 {
// 				c.hearts.errCnt++
// 			}
// 		}
// 	}
// }

func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}
