package main

import (
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
)

// test route ping
func TestRouter(t *testing.T) {
	r := gofight.New()
	r.GET("/ping").
		Run(Router(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			t.Log(r.Body.String())
			assert.Equal(t, "pong", r.Body.String(), "should be equal")
			assert.Equal(t, 200, r.Code, "should be equal")
		})
}

// test route create user
func TestCreateUser(t *testing.T) {
	r := gofight.New()
	r.POST("/users").
		SetJSON(gofight.D{
			"name": "test",
			"username": "test",
			"email": "test@gmail.com",
			"password": "123456",
			"age": 20,
		}).
		Run(Router(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code, "response code is not 200")
		})
	}

// test route get all user
func TestUsers(t *testing.T) {
	r := gofight.New()
	r.GET("/users").
		Run(Router(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			// data := []byte(r.Body.String())

			// assert.Equal(t, "name", name)
			assert.Equal(t, 200, r.Code, "response code is not 200")
		})
}

// test hash password
func TestHashPassword(t *testing.T) {
	var data User
	data.Password = "123456"
	
	assert.Equal(t, nil, data.HashPassword(data.Password))
	
}

// test route update user
func TestUpdateUser(t *testing.T) {
	r := gofight.New()
	r.PUT("/users/5").
		SetJSON(gofight.D{
			"name": "test update",
			"username": "test",
			"email": "testupdate@gmail.com",
			"password": "123456",
			"age": 20,
		}).
		Run(Router(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code, "response code is not 200")
		})
}

// test route delete user
func TestDeleteUser(t *testing.T) {
	r := gofight.New()
	r.DELETE("/users/4").
		Run(Router(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code, "response code is not 200")
		})
}

// test route get user by username
func TestGetUserByUsername(t *testing.T) {
	r := gofight.New()
	r.GET("/users/login").
		SetJSON(gofight.D{
			"username": "test",
			"password": "123456",
		}).
		Run(Router(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code, "response code is not 200")
		})
}

// test cek password
func TestCheckPassword(t *testing.T) {
	password := "123456"
	hash := "$2a$14$BBc5BMgs7Q5N5P4t2Jj3Bulr6vg2Ea4Iaw.YNkpDSKl3aj15gs97u"
	match := CheckPassword(password, hash)
	assert.Equal(t, true, match, "password is not match")
}