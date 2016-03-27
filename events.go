// Different events that can be sent by a connected user.
package main

type JoinEvent struct {
	user *User
}

type LeaveEvent struct {
	user *User
}
