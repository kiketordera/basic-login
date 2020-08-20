package app

// checkError checks if the operation was wrong, and if wrong, prints the error in the screen, the UI
func checkError(err error) {
	if err != nil {
		print("An error happend while dealing with the database")
	}
}
