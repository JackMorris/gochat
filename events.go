// Different events that can be sent by a connected user.
package main

import "strings"

type BellEvent struct {
	user *User
}

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

type NameChangeEvent struct {
	user         *User
	previousName string
	newName      string
}

func constructEvent(user *User, eventString string) interface{} {
	// From the specified `eventString` construct the correct event
	// for the given `user`, and return it.
	if len(eventString) == 0 {
		// Empty string -> no event.
		return nil
	}
	if strings.HasPrefix(eventString, `/`) {
		// Attempt to parse command.

		// Change name
		// /name <newName>
		if strings.HasPrefix(eventString, `/name `) {
			previousName := user.name
			newName := strings.SplitN(eventString, " ", 2)[1] // Strip off the command

			if previousName == newName {
				return nil // Nothing to be done
			}

			return NameChangeEvent{
				user:         user,
				previousName: previousName,
				newName:      newName}
		}

		// Ring the bell
		// /bell
		if eventString == "/bell" {
			return BellEvent{user: user}
		}
	}

	// Otherwise, fail through to standard message
	return MsgEvent{user: user, msg: eventString}
}
