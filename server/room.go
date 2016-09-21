package main

import (
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

// Chatroom handle all rooms
type Chatroom interface {
	Checkin(string, string, *websocket.Conn) Room
	Find(string) Room
	GC(*websocket.Conn)
}

type chatroom struct {
	*sync.Mutex
	rooms map[string]Room
}

// NewChatroom return a chatroom instance
func NewChatroom() Chatroom {
	return &chatroom{
		Mutex: &sync.Mutex{},
		rooms: make(map[string]Room),
	}
}

func (cm *chatroom) Checkin(chatName, userName string, c *websocket.Conn) Room {
	cm.Lock()
	defer cm.Unlock()

	r, ok := cm.rooms[chatName]
	if !ok {
		r = NewRoom()
		cm.rooms[chatName] = r
	}

	r.Join(userName, c)

	return r
}

func (cm *chatroom) Find(chatName string) Room {
	cm.Lock()
	defer cm.Unlock()

	return cm.rooms[chatName]
}

func (cm *chatroom) GC(c *websocket.Conn) {

	cm.Lock()
	defer cm.Unlock()
	for name, r := range cm.rooms {
		if r.Has(c) {
			if 0 >= r.Leave(c) {
				log.Println("[GC room]", name)
				delete(cm.rooms, name)
				return
			}
		}
	}
}

// Room is used for put conn of ws
type Room interface {
	Members() int
	Has(*websocket.Conn) bool
	Join(string, *websocket.Conn) int
	Leave(*websocket.Conn) int
	Broadcast(interface{})
}

type room struct {
	*sync.RWMutex
	seats map[*websocket.Conn]string
}

// NewRoom return a room instance
func NewRoom() Room {
	return &room{
		RWMutex: &sync.RWMutex{},
		seats:   make(map[*websocket.Conn]string),
	}
}

func (r *room) Members() int {
	r.RLock()
	defer r.RUnlock()

	return len(r.seats)
}

func (r *room) Has(c *websocket.Conn) bool {
	r.RLock()
	defer r.RUnlock()

	_, ok := r.seats[c]
	return ok
}

func (r *room) Join(name string, c *websocket.Conn) int {
	r.Lock()
	defer r.Unlock()
	r.seats[c] = name

	return len(r.seats)
}

func (r *room) Leave(c *websocket.Conn) int {
	r.Lock()
	defer r.Unlock()
	delete(r.seats, c)

	return len(r.seats)
}

func (r *room) Broadcast(msg interface{}) {
	r.RLock()
	defer r.RUnlock()

	for c, name := range r.seats {
		if err := websocket.JSON.Send(c, msg); err != nil {
			log.Printf("[broadcast] fail on %s\n", name)
		}
	}
}
