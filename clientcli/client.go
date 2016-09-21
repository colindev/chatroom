package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

var env = struct {
	addr, room, name string
}{}

type message struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

func main() {

	flag.StringVar(&env.addr, "addr", "ws://localhost:9090", "Chatroom WS addr")
	flag.StringVar(&env.room, "room", "", "room name (default is random)")
	flag.StringVar(&env.name, "name", "ghost", "user name")
	flag.Parse()

	if env.room == "" {
		env.room = strconv.FormatInt(time.Now().Unix(), 10)
	}

	wsURL := env.addr + "/" + env.room + "?name=" + env.name

	fmt.Println("[dial]", wsURL)

	ws, err := websocket.Dial(wsURL, "", "http://127.0.0.1")
	if err != nil {
		log.Fatal("[ws dial]", err)
	}

	go func() {
		var msg = []byte{}
		for {
			var buf = make([]byte, 16)
			if n, err := ws.Read(buf); err != nil {
				log.Fatal("[read]", err)
			} else {
				msg = append(msg, buf[:n]...)
				if n < len(buf) {
					// json decode
					var m message
					if e := json.Unmarshal(msg, &m); err != nil {
						log.Fatal("[json unmarshal]", e)
					}
					fmt.Println(time.Now(), m.Name, ">", strings.TrimSpace(m.Msg))
					msg = []byte{}
				}
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if text != "" {
			if _, err := ws.Write([]byte(text)); err != nil {
				log.Fatal("[write]", err)
			}
		}
	}
}
