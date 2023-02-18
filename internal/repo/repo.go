package repo

import (
	// "time"
	// "fmt"
	// "sort"
) 


type User struct {
	Name string `json:"login"`
	Password string `json:"password"`
}

type Repo interface {
	UserRegistrate(u User) error
	GetSeckretKey() []byte
	UserIsPresent(name string) bool
	UserAuthenticate(u User) (bool, error)
	UserGet(name string) (*User, error)
	CheckExistOrderNumber(naumber int) bool
	GetIdOrderOwner(naumber int) int
	GetUserNameById(userId int) string
	ObtainNewOrder(userId, number int) error
	// GetOk()

}


