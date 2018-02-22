package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"github.com/vitor-caetano/goblueprints/_chapter1/trace"
	"github.com/stretchr/objx"
)

type room struct {
	// foward is a channel that holds incoming messages
	// that should be fowarded to the other clients
	foward chan *message

	// join is a channel for clients whising to join the room
	join chan *client

	// leave is a channel for clients whising to leave the room
	leave chan *client

	// clients hold all the clients in this room
	clients map[*client]bool

	// tracer will receive trace information of activity
	tracer trace.Tracer
}

func newRoom() *room {
	return &room{
		foward:  make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//joining
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			//leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.foward:
			r.tracer.Trace("Message received:", msg.Message)
			//foward message to all clients
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace("-- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}

	client := &client{socket: socket,
		send: make(chan *message, messageBufferSize),
		room: r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()

	go client.write()
	client.read()
}
