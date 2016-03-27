# gochat
Simple text-based tcp chat server, written in Go.

# Usage
`go get github.com/jackmorris/gochat; gochat <listen addr>` (eg. localhost:8000, :8000, ...)

Clients can then connect directly to the socket (eg. by using github.com/jackmorris/tcpio)

# Commands
`/name <newname>`: change your name

`/bell`: ring the bell
