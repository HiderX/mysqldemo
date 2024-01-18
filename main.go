package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

var db *gorm.DB

func InitMySQL() {
	dsn := "root:root1234@(127.0.0.1:3306)/genshin?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ = gorm.Open("mysql", dsn)
}

type Character struct {
	Name          string `gorm:"type:varchar(20);not null;PRIMARY_KEY" json:"name"`
	Level         int    `gorm:"type:int;not null" json:"level"`
	Constellation int    `gorm:"type:int(1);not null" json:"constellation"`
	Skill         string `gorm:"not null" json:"skill"`
}

func main() {
	InitMySQL()
	defer db.Close()
	db.AutoMigrate(&Character{})
	//db.Create(&Character{Name: "神里绫华", Level: 90, Constellation: 6, Skill: "神里流·霜灭"})

	r := gin.Default()
	r.POST("/c", func(c *gin.Context) {
		var character Character
		c.ShouldBindJSON(&character)
		db.Create(&character)
		c.JSON(http.StatusOK, gin.H{
			"message":   "success",
			"character": character,
		})
	})
	r.GET("/r/:name", func(c *gin.Context) {
		name := c.Param("name")
		var character Character
		err := db.Where("name = ?", name).First(&character).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "not found",
			})
			return
		}
		c.JSON(http.StatusOK, character)
	})
	r.PUT("/u/:name", func(c *gin.Context) {
		name := c.Param("name")
		var character Character
		err := db.Where("name = ?", name).First(&character).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "not found",
			})
			return
		}
		c.ShouldBindJSON(&character)
		db.Save(&character)
		c.JSON(http.StatusOK, gin.H{
			"message":   "success",
			"character": character,
		})
	})
	r.DELETE("/d/:name", func(c *gin.Context) {
		name := c.Param("name")
		var character Character
		err := db.Where("name = ?", name).Delete(&character).Error
		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "not found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	r.Run("127.0.0.1:8080")
}
