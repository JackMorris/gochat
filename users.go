// Tracking connected users.
package main

//////////
// user //
//////////

type User struct {
	outputChan chan string
	name       string
}

/////////////
// userSet //
/////////////

type UserSet struct {
	set map[*User]bool
}

func newUserSet() UserSet {
	// Return a new userSet
	var userSet UserSet
	userSet.set = make(map[*User]bool)
	return userSet
}

func (us *UserSet) add(u *User) {
	// Add user `u` to the user set.
	// If `u` is already in the userSet, there is no effect.
	us.set[u] = true
}

func (us *UserSet) remove(u *User) {
	// Remove user `u` from the user set.
	// If `u` is not in the userSet, there is no effect.
	_, member := us.set[u]; if member {
		delete(us.set, u)
	}
}

func (us *UserSet) broadcast(msg string) {
	// Broadcast `msg` to all users in the userSet.
	for user := range us.set {
		user.outputChan <- msg
	}
}
