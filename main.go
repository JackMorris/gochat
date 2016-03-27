// gochat
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var eventChan = make(chan interface{})

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
	go eventHandler()

	// Handle each new connection.
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Error connecting. Continue to next connection.
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func eventHandler() {
	// Listen on eventChan, and react to each event.
	// Keeps track of users so it can broadcast messages, if needed.

	// Set of connected users.
	users := newUserSet()

	for event := range eventChan {
		switch event := event.(type) {
		case JoinEvent:
			users.add(event.user)
			users.broadcast(event.user.name + " has joined")
		case LeaveEvent:
			users.remove(event.user)
			users.broadcast(event.user.name + " has left")
		default:
			fmt.Printf("Error: unknown event received: %v\n", event)
		}
	}
}

func handleConn(conn net.Conn) {
	// Called for each new user. Creates a struct to represent the user, and notifies
	// that the user has joined through the eventChan.
	// Messages from this user are also sent thrgough eventChan, including the final
	// leaveMessage when the user disconnects.
	// Also handles sending data from the users outputChan to that user through conn.

	// Setup a new user.
	user := new(User)
	user.outputChan = make(chan string)
	user.name = conn.RemoteAddr().String()
	go func() {
		// For each message to the outputChan, print to user.
		for msg := range user.outputChan {
			fmt.Fprintln(conn, msg)
		}
	}()

	// Send the join event
	eventChan <- JoinEvent{user}

	// Process input
	input := bufio.NewScanner(conn)
	for input.Scan() {
		if input.Err() != nil {
			// Error processing input. Skip to the next message.
			continue
		}
		// Don't do anything with message for now
	}

	// Send leave event
	eventChan <- LeaveEvent{user}
	close(user.outputChan)
	conn.Close()
}