package main

import (
	"sync"

	a "github.com/kiketordera/basic-login/src/app"
)

func main() {

	// Initialize the app
	db := a.InitDatabase()

	// New waitgroup for sync
	var wg sync.WaitGroup

	// We inicialize 1 app, so we wait 1 process
	wg.Add(1)
	// Init the app (iun the background due to GO)
	go db.Start(&wg)

	// Wait until all apps stop
	wg.Wait()
	db.DataBase.Close()
}
