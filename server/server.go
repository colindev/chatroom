package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/websocket"
)

var env = struct {
	wsAddr     string
	clientAddr string
	clientDir  string
}{}

func main() {

	flag.StringVar(&env.wsAddr, "ws", ":9090", "WebSocket serve port")
	flag.StringVar(&env.clientAddr, "c", ":9000", "Chatroom UI serve port")
	flag.StringVar(&env.clientDir, "ui", "client", "Chatroom UI dir")
	flag.Parse()

	go func() {

		cm := NewChatroom()

		err := http.ListenAndServe(env.wsAddr, websocket.Handler(func(c *websocket.Conn) {
			roomName := c.Request().URL.Path
			connName := c.Request().URL.Query().Get("name")

			fmt.Printf("[%s] %s joined\n", roomName, connName)

			// 取得聊天室
			r := cm.Checkin(roomName, connName, c)

			// 讀取資料
			for {
				var s string
				err := websocket.Message.Receive(c, &s)
				if err == io.EOF {
					break
				} else if err != nil {
					log.Println("[ws receive]", err)
				} else {
					fmt.Printf("[%s] %s: %s\n", roomName, connName, s)
					r.Broadcast(struct {
						Name string `json:"name"`
						Msg  string `json:"msg"`
					}{connName, s})
				}
			}
			fmt.Printf("[%s] %s leaved\n", roomName, connName)

		}))
		log.Println("[ws] shutdown:", err)
	}()

	go func() {
		err := http.ListenAndServe(env.clientAddr, http.FileServer(http.Dir(env.clientDir)))
		log.Println("[ui] shutdown:", err)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	sign := <-shutdown
	log.Printf("[shutdown] %s\n", sign)
}
