package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// connect mysql database gorm
func ConnectDb() *gorm.DB{
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// migrate
	var user User
	db.AutoMigrate(&user)

	return db
}

// router
func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	users := r.Group("/users")
	{
		users.GET("", Users)
		users.POST("", CreateUser)
		users.PUT("/:id", UpdateUser)
		users.DELETE("/:id", DeleteUser)
		users.GET("/login", GetUser)
	}
	return r
}


// return all user
func Users(c *gin.Context) {
	db := ConnectDb()
	var users []User
	if db.Find(&users).Error == nil && len(users) > 0 {
		c.JSON(200, users)	
	} else {
		c.JSON(400, gin.H{"error": "users not found"})
	}
}


// create user 
func CreateUser(c *gin.Context) {
	// connect database
	db := ConnectDb()

	// get json
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// get paddword hash
	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	record := db.Create(&user)

	if record.Error != nil {
		c.JSON(400, gin.H{"error" : record.Error.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "user created succesfully"})
	
}

// hash password
func (user *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	user.Password = string(hash)
	return nil
}

// update user
func UpdateUser(c *gin.Context) {
	db := ConnectDb()

	// get json
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error" : err.Error()})
		c.Abort()
		return
	}

	// id
	id := c.Param("id")
	if  id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	// // get paddword hash
	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(400, gin.H{"error" : err.Error()})
		c.Abort()
		return
	}

	record := db.Model(&user).Where("id = ?", id).Updates(&user)

	if record.Error != nil {
		c.JSON(400, gin.H{"error" : record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"message" : "user successfully updated"})
}

// delete user
func DeleteUser(c *gin.Context) {
	db := ConnectDb()

	var user User

	// id
	id := c.Param("id")
	if  id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	// delete user
	record := db.First(&user, id).Delete(&user)

	// if error
	if record.Error != nil {
		c.JSON(400, gin.H{"error" : record.Error.Error()})
		c.Abort()
		return
	}
	// if success
	c.JSON(200, gin.H{"message" : "successfully deleted"})

}

// get user by username
func GetUser(c *gin.Context) {
	db := ConnectDb()
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error" : err.Error()})
		c.Abort()
		return
	}
	if  user.Username == "" {
		c.JSON(400, gin.H{"error": "username is required"})
		return
	}

	var u User
	record := db.First(&u).Where("username = ?", user.Username)

	if record.Error != nil {
		c.JSON(400, gin.H{"error" : record.Error.Error()})
		c.Abort()
		return
	}

	if CheckPassword(user.Password, u.Password) == false {
		c.JSON(400, gin.H{"error" : "password not match"})
		return
	}
	c.JSON(200, gin.H{"message" : "found"})
}

// check password
func CheckPassword(passwordHash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordHash))
	return err == nil
}