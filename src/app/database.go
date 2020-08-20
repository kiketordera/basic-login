package app

import (
	"github.com/gin-gonic/gin"
	"github.com/timshannon/bolthold"
	"gopkg.in/mgo.v2/bson"
)

// getUserByID finds and return the User in database by his ID
func (db *DB) getUserByID(id bson.ObjectId) User {
	var user User
	checkError(db.DataBase.Get(id, &user))
	return user
}

// getUserByUsername finds the User in database by his username and returns it.
func (db *DB) getUserByMail(email string) (User, bool) {
	var users []User
	err := db.DataBase.Find(&users, bolthold.Where("Email").Eq(email))
	checkError(err)
	if users == nil {
		return User{}, false
	}
	var user = users[0]
	return user, true
}

// getUserByID finds and return the User in database by his ID
func (db *DB) existUser(email, password string) bool {
	var user User
	err := db.DataBase.Find(&user, bolthold.Where("Email").Eq(email).And("Password").Eq(password).And("IsActive").Eq(true))
	checkError(err)
	if err != nil {
		return true
	}
	return false
}

// addUserToDatabase adds (or updates) the User created to the Database
// Returns true if the operatios was sucessfully and false if it fails
func (db *DB) addUserToDatabase(c *gin.Context, u User) bool {
	// We check if the username is available, if not send error mesage and break
	var users []User
	if err := db.DataBase.Find(&users, bolthold.Where("Email").Eq(u.Email)); err != nil {
		return false
	}
	if len(users) > 0 {
		if users[0].ID != u.ID {
			wrongFeedback(c, "The Email is already in use, try with another one")
			return false
		}
	}
	err := db.DataBase.Upsert(u.ID, u)
	checkError(err)
	return true
}

// This method returns all the Users made with the DBs of the villages
func (db *DB) getAllUsers(c *gin.Context) []User {
	var users []User
	db.DataBase.Find(&users, nil)
	return users
}
