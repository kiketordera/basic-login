package app

import (
	"github.com/timshannon/bolthold"
	"gopkg.in/mgo.v2/bson"
)

// DB is the DataBase which contains the information.
type DB struct {
	DataBase *bolthold.Store
}

// User is the user of the system
type User struct {
	ID       bson.ObjectId `json:"id"`
	Name     string        `json:"name"`
	Surname  string        `json:"surname"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
	Role     Role          `json:"role"`
}

// Role is the close list of different roles that the users can have
type Role string

// The options for RoleType
const (
	Admin    Role = "admin"
	Manager  Role = "manager"
	Customer Role = "customer"
)
