package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

//Server A websocket server
type Server struct {
	ListenAddr string
	conns      map[*websocket.Conn]bool
	mtx        sync.RWMutex
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	s.mtx.Lock()
	s.conns[ws] = true
	s.mtx.Unlock()

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				//Chat ended, someone disconnected
				break
			}
			//TODO: Use a context or w/e so we can say who had an error
			log.Printf("error reading message for conn...: %s\n", err.Error())
			continue //Allowing users to continue if they send a malformed error
		}
		msg := buf[:n]
		fmt.Println(string(msg))
		// could broadcast in go routine to not block at the cost of more memory
		s.broadcast(msg)
	}
}

func (s *Server) broadcast(b []byte) {
	s.mtx.RLock()
	for ws := range s.conns {
		ws := ws
		go func() {
			if _, err := ws.Write(b); err != nil {
				log.Println("error writing message to connection", err)
			}
		}()
	}
	s.mtx.RUnlock()
}

func (s *Server) addRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})
	http.Handle("/chat", websocket.Handler(s.handleWS))
}

//Start the server
func (s *Server) Start() error {
	s.addRoutes()
	fmt.Printf("Starting server on: %s\n", s.ListenAddr)
	return http.ListenAndServe(s.ListenAddr, nil)
}

//NewServer Create a new instance of a websocket server
func NewServer(addr string) *Server {
	return &Server{
		ListenAddr: addr,
		conns:      make(map[*websocket.Conn]bool),
	}
}
