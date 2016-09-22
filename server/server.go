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
	"time"

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
			now := time.Now()
			roomName := strings.TrimLeft(c.Request().URL.Path, "/")
			connName := c.Request().URL.Query().Get("name")

			// user profile
			userProfile := Profile{
				Name: connName,
				Icon: sumName(connName),
			}

			fmt.Printf("%s [%s] %s joined\n", time.Now(), roomName, connName)

			// 取得聊天室
			r := cm.Checkin(roomName, connName, c)

			// 告知剛進房使用者當前人數
			r.Each(func(conn *websocket.Conn, name string) {
				if c == conn {
					return
				}

				r.Send(c, Pack{
					Active: Join,
					Profile: Profile{
						Name: name,
						Icon: sumName(name),
					},
					Time: now.Unix(),
				})
			})

			// 廣播使用者進房
			r.Broadcast(Pack{
				Active:  Join,
				Profile: userProfile,
				Time:    now.Unix(),
			})

			// 讀取資料
			for {
				var s string
				err := websocket.Message.Receive(c, &s)
				if err == io.EOF {
					break
				} else if err != nil {
					log.Println("[ws receive]", err)
				} else {
					fmt.Printf("%s [%s] %s: %s\n", time.Now(), roomName, connName, strings.TrimSpace(s))
					r.Broadcast(Pack{
						Active:  Message,
						Profile: userProfile,
						Msg:     s,
						Time:    time.Now().Unix(),
					})
				}
			}
			fmt.Printf("%s [%s] %s leaved\n", time.Now(), roomName, connName)
			cm.GC(c)

			// 廣播使用者退出
			r.Broadcast(Pack{
				Active:  Leave,
				Profile: userProfile,
				Time:    time.Now().Unix(),
			})
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

func sumName(name string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(name)))
}
