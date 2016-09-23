package main

// Active use for define Pack.Active
type Active string

const (
	// Init send on connect
	Init Active = "init"
	// SystemMessage mean this pack come from system
	SystemMessage Active = "system-msg"
	// Message mean this pack is message
	Message Active = "msg"
	// Join mean user come in room
	Join Active = "join"
	// Leave mean user leaved
	Leave Active = "leave"
	// Clear mean admin want clear all message
	Clear Active = "clear"
)

// Profile contain user infomation
type Profile struct {
	Icon string `json:"icon"`
	Name string `json:"name"`
}

// Pack contain event, data, profile
type Pack struct {
	Profile Profile `json:"profile"`
	Active  Active  `json:"active"`
	Msg     string  `json:"msg,omitempty"`
	Time    int64   `json:"time"`
}
