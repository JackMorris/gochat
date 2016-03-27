// Different events that can be sent by a connected user.
package main

type JoinEvent struct {
	user *User
}

type LeaveEvent struct {
	user *User
}

type MsgEvent struct {
	user *User
	msg  string
}

func constructEvent(user *User, eventString string) interface{} {
	// From the specified `eventString` construct the correct event
	// for the given `user`, and return it.
	if len(eventString) == 0 {
		// Empty string -> no event.
		return nil
	}
	if eventString[0] == '/' {
		// Attempt to parse command.
		// Currently, no commands.
	}

	// Otherwise, fail through to standard message
	return MsgEvent{user: user, msg: eventString}
}
