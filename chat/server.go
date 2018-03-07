package chat

import (
	"errors"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var (
	ErrPackageHeaderLength = errors.New("package header length error")
	ErrContentLength       = errors.New("content too long or empty")
	server                 = NewServer()
)

type Clients map[*Client]struct{}

func (clients *Clients) Uids() []uint32 {
	var uids []uint32
	for client := range *clients {
		uids = append(uids, client.uid)
	}
	return uids
}

type Server struct {
	clients         Clients
	addSignal       chan *Client
	delSignal       chan *Client
	broadcastSignal chan interface{}
	assignSignal    chan *Assign
	doneSignal      chan bool
}

func NewServer() *Server {
	return &Server{
		make(Clients),
		make(chan *Client, 10),
		make(chan *Client, 10),
		make(chan interface{}, 10),
		make(chan *Assign, 10),
		make(chan bool),
	}
}

func (s *Server) Add(c *Client) {
	s.addSignal <- c
}
func (s *Server) add(c *Client) {
	s.clients[c] = struct{}{}
	log.Printf("Number of client connections:%d", len(s.clients))
}

func (s *Server) Del(c *Client) {
	s.delSignal <- c
}
func (s *Server) del(c *Client) {
	delete(s.clients, c)

	uid := c.uid
	exit := true
	for c := range s.clients {
		if uid == c.uid || c.uid == 0 {
			exit = false
		}
	}
	if exit {
		c.CloseExit()
	}

	log.Println("del", c, exit)
}

func (s *Server) Broadcast(data interface{}) {
	s.broadcastSignal <- data
}
func (s *Server) broadcast(data interface{}) {
	for c := range s.clients {
		go c.Write(data)
	}
}

func (s *Server) Assign(data *Assign) {
	s.assignSignal <- data
}
func (s *Server) assign(data *Assign) {
	for c := range data.Clients {
		go c.Write(data.Data)
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen(mux *http.ServeMux) {
	// log.Println("start init websocker...")
	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				log.Println("websocker close:", err)
				s.doneSignal <- true
			}

			if err := recover(); err != nil {
				// if err is io.EOF, maye be beacause of client closing
				if err != io.EOF {
					log.Println("websocker recover:", err)
				}
			}
		}()
		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	mux.Handle("/chat", websocket.Handler(onConnected))
	// log.Println("finish init websocker.")

	for {
		select {
		case client := <-s.addSignal:
			s.add(client)
		case client := <-s.delSignal:
			s.del(client)
		case data := <-s.broadcastSignal:
			s.broadcast(data)
		case data := <-s.assignSignal:
			s.assign(data)
		case <-s.doneSignal:
			return
		}
	}
}
