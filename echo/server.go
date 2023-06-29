package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//WsServer A websocket server
type WsServer struct {
	upgrader   *websocket.Upgrader
	ListenAddr string
}

func (ws *WsServer) makeUpgrader() {
	ws.upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}

func (ws *WsServer) addRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := ws.upgrader.Upgrade(w, r, nil)
		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Websocket err:", err)
				return
			}
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
			if err = conn.WriteMessage(mt, msg); err != nil {
				log.Println("Websocket err:", err)
				return
			}
		}
	})

}

//Start the server
func (ws *WsServer) Start() error {
	ws.makeUpgrader()
	ws.addRoutes()
	fmt.Printf("Starting server on: %s\n", ws.ListenAddr)
	return http.ListenAndServe(ws.ListenAddr, nil)
}

//NewServer Create a new instance of a websocket server
func NewServer(addr string) *WsServer {
	//TODO: Check Addr is a proper bind address
	return &WsServer{
		ListenAddr: addr,
	}
}
