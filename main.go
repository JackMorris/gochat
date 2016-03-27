// gochat
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// Check program args.
	if len(os.Args) != 2 {
		fmt.Println("Usage: gochat <interface>")
		os.Exit(2)
	}

	// Listen for new connections.
	listener, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening for connections on %s\n", os.Args[1])

	// Spin up the event handler
	var eventChan = make(chan interface{})
	go eventHandler(eventChan)

	// Handle each new connection.
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Error connecting. Continue to next connection.
			log.Print(err)
			continue
		}

		// Setup a new user.
		user := new(User)
		user.outputChan = make(chan string)
		user.name = conn.RemoteAddr().String()

		go handleConn(conn, user, eventChan)
	}
}

func eventHandler(eventChan <-chan interface{}) {
	// React to events on eventChan as required.
	// Keeps track of connected users for broadcasting purposes.

	// Set of connected Users.
	users := make(map[*User]bool)

	// Utility function to broadcast a message to all users.
	broadcastMsg := func(msg string) {
		for user := range users {
			user.outputChan <- msg
		}
	}

	// React to each event.
	for event := range eventChan {
		switch event := event.(type) {
		case BellEvent:
			broadcastMsg("<<< " + event.user.name + " rang the bell >>>")
		case JoinEvent:
			users[event.user] = true
			broadcastMsg(event.user.name + " has joined")
		case LeaveEvent:
			delete(users, event.user)
			broadcastMsg(event.user.name + " has left")
		case MsgEvent:
			broadcastMsg(event.user.name + ": " + event.msg)
		case NameChangeEvent:
			event.user.name = event.newName
			broadcastMsg(event.previousName + " is now known as " + event.newName)
		default:
			fmt.Printf("Error: unknown event received: %v\n", event)
		}
	}
}

func handleConn(conn net.Conn, user *User, eventChan chan<- interface{}) {
	// Handles a single connection between the server and a user
	// This involves passing output to the user (from the user's outputChan),
	// and interpreting user input into events that can be placed onto the
	// eventChan.

	// Handle output to user.
	go func() {
		// For each message to the outputChan, print to user.
		for msg := range user.outputChan {
			fmt.Fprintln(conn, msg)
		}
	}()

	// Send the join event
	eventChan <- JoinEvent{user: user}

	// Handle input from user.
	input := bufio.NewScanner(conn)
	for input.Scan() {
		if input.Err() != nil {
			// Error processing input. Skip to the next message.
			continue
		}

		if event := constructEvent(user, input.Text()); event != nil {
			eventChan <- event
		}
	}

	// Send leave event
	eventChan <- LeaveEvent{user: user}
	close(user.outputChan)
	conn.Close()
}
