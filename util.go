package gome

import (
	"log"
	"runtime/debug"

	"github.com/veandco/go-sdl2/sdl"
)

// A Vector holds an X and a Y value and can be
// used for positions or directions.
type Vector struct {
	X uint
	Y uint
}

// A FloatVector holds an X and a Y float value and can be
// used for positions or directions.
type FloatVector struct {
	X float32
	Y float32
}

// Throw outputs an error to console and stops execution
// it should only be used in top-level functions, as returning the
// error is required for unit testing
func Throw(err error, msg string) {
	debug.PrintStack()
	log.Fatalf("=> %v: %s", err, msg)
}

/*
	MailBox
*/

// A Message is a piece of information sendable through the MailBox.
type Message interface {
	Name() string
}

type mailBox struct {
	listeners map[string][]func(Message)
}

// The MailBox is used to communicate between systems. Through the MailBox, one
// can send Messages and listen for them.
var MailBox *mailBox = &mailBox{make(map[string][]func(Message))}

// Send sends a Message through the MailBox to functions listening for
// that type of Message.
func (mb *mailBox) Send(msg Message) {
	for _, fun := range mb.listeners[msg.Name()] {
		fun(msg)
	}
}

// Listen adds the function to the group listening for a Message of a specific type.
func (mb *mailBox) Listen(msgName string, fun func(Message)) {
	mb.listeners[msgName] = append(mb.listeners[msgName], fun)
}

/*
	Default Messages
*/

type KeyboardMessage struct {
	Key       sdl.Keysym
	Timestamp uint32
}

func (KeyboardMessage) Name() string { return "Keyboard" }
