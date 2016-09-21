package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	cm := NewChatroom()

	go func() {

		err := http.ListenAndServe(env.wsAddr, websocket.Handler(func(c *websocket.Conn) {
			roomName := strings.TrimLeft(c.Request().URL.Path, "/")
			connName := c.Request().URL.Query().Get("name")

			// 計算頭像
			iconHash := fmt.Sprintf("%x", md5.Sum([]byte(connName)))

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
					fmt.Printf("[%s] %s: %s\n", roomName, connName, strings.TrimSpace(s))
					r.Broadcast(struct {
						Active string `json:"active"`
						Name   string `json:"name"`
						Icon   string `json:"icon"`
						Msg    string `json:"msg"`
					}{
						Active: "msg",
						Name:   connName,
						Icon:   iconHash,
						Msg:    s,
					})
				}
			}
			fmt.Printf("[%s] %s leaved\n", roomName, connName)
			cm.GC(c)
		}))
		log.Println("[ws] shutdown:", err)
	}()

	go func() {
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(env.clientDir))))
		http.HandleFunc("/api/clear", func(w http.ResponseWriter, r *http.Request) {
			room := r.URL.Query().Get("room")

			if r := cm.Find(room); r != nil {
				r.Broadcast(struct {
					Active string `json:"active"`
				}{"clear"})
			}
		})
		err := http.ListenAndServe(env.clientAddr, nil)
		log.Println("[ui] shutdown:", err)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	sign := <-shutdown
	log.Printf("[shutdown] %s\n", sign)
}
