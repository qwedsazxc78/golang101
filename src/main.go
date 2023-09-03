// main.go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"golang101/controllers"
	"golang101/models"
)

var db *gorm.DB

func main() {
	// 初始化Gin引擎
	router := gin.Default()

	// 設定健康檢查路由
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// 初始化資料庫連線
	initDatabase()

	// 設定RESTful API路由
	api := router.Group("/api")
	{
		todoController := controllers.NewTodoController(db)
		api.GET("/todos", todoController.GetTodos)
		api.POST("/todos", todoController.CreateTodo)
		api.PUT("/todos/:id", todoController.UpdateTodo)
		api.DELETE("/todos/:id", todoController.DeleteTodo)
	}

	// 啟動HTTP伺服器
	log.Println("Server started on :8080")
	router.Run(":8080")
}

func initDatabase() {
	dsn := "host=127.0.0.1 user=golang101 password=golang101pass dbname=golang101 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("無法連接資料庫")
	}

	// 自動建立資料表
	db.AutoMigrate(&models.Todo{})
}

// 其他初始化函式（例如初始化Redis）可以放在這裡
