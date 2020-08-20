package app

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Redirects when a user introduces a route that does not exists
func redirect(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

// Give a good feedback to the User
func goodFeedback(c *gin.Context) {
	render(c, gin.H{}, "good-feedback.html")
}

// Give a bad feedback to the User
func wrongFeedback(c *gin.Context, t string) {
	render(c, gin.H{
		"text": t,
	}, "wrong-feedback.html")
}

// Shows the home page
func showLogin(c *gin.Context) {
	render(c, gin.H{}, "login.html")
}

// Shows the home page
func (db *DB) login(c *gin.Context) {
	if aemail, b := getStringFromHTML(c, "aemail", "form", false); b {
		// Is a Log in
		user, exist := db.getUserByMail(aemail)
		if !exist {
			wrongFeedback(c, "The username introduced is invalid")
			return
		}
		password, b := getStringFromHTML(c, "apassword", "form", true)
		if !b {
			wrongFeedback(c, "The was a problem with the password")
			return
		}
		// Check if the username/password combination is valid
		fmt.Print("these are the passwords:", user.Password, password)
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err == nil {
			setSession(c, user.Email, user.ID.Hex())
			c.Redirect(http.StatusFound, "/")
			return
		}
		fmt.Print("these is the error:", err)
		wrongFeedback(c, "The password introduced is invalid")
		return
	}
	// Is a Sign up
	name, a := getStringFromHTML(c, "name", "form", true)
	surname, b := getStringFromHTML(c, "surname", "form", true)
	email, ce := getStringFromHTML(c, "bemail", "form", true)
	password, d := getStringFromHTML(c, "bpassword", "form", true)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	checkError(err)
	if !a || !b || !ce || !d {
		return
	}
	user := User{
		ID:       bson.NewObjectId(),
		Name:     name,
		Surname:  surname,
		Password: string(passwordHash),
		Email:    email,
	}
	if !db.addUserToDatabase(c, user) {
		return
	}
	setSession(c, user.Email, string(user.Role))
	c.Redirect(http.StatusFound, "/users")
	goodFeedback(c)
}

// Shows the home page
func (db *DB) showUsers(c *gin.Context) {
	users := db.getAllUsers(c)
	render(c, gin.H{
		"users": users,
	}, "users.html")
}

// Give a bad feedback to the User
func wrongFeedbackTest(c *gin.Context) {
	render(c, gin.H{
		"text": "The User mail already exits in the database! Try with another BLABLBAL",
	}, "good-feedback.html")
}
