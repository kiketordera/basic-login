package app

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Gets the String information of the HTML in a POST request
// returns true if it finds the string and false in case of problems
func getStringFromHTML(c *gin.Context, name, element string, printError bool) (string, bool) {
	data := c.PostForm(name)
	// We remove the spaces before the first letter and after the last one
	dataClean := strings.TrimSpace(data)
	// Render error if the parameter was empty, this way we have a server-side validation
	if data == "" {
		if printError {
			wrongFeedback(c, "Sorry, there was a problem trying to process the "+name+" of the "+element)
		}
		return "", false
	}
	return dataClean, true
}
