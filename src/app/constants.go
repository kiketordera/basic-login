package app

import (
	"os"
)

var (
	// BasePath is the path to the project
	BasePath = os.Getenv("GOPATH") + "/src/github.com/kiketordera/basic-login"
	// TokenSigningKey is the key to sign the cookies
	tokenSigningKey = []byte("SuperFancyToken:D")
)

const (
	// SessionTime is the seconds the session will be active
	sessionTime = 108000
	// Port is the port where the server will be listening
	port = "8080"
	// This is the name of the project
	projectName = "demo"
)
