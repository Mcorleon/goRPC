package api

import (
	"fmt"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserHandler struct{}

func (u *UserHandler) GetUserById(id int) (User, error) {
	if u, ok := userList[id]; ok {
		fmt.Println("valid user")
		return u, nil
	}
	fmt.Println("empty user")
	return User{}, nil
}
var userList = map[int]User{
	1: User{1, "hero", 11},
	2: User{2, "kavin", 12},
}