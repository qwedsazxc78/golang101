// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"encoding/json"
	"time"
	"context"

	"github.com/joho/godotenv" // 引入 godotenv 套件
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/go-redis/redis/v8"

	"golang101/controllers"
	"golang101/models"
)

var (
	redisClient *redis.Client
	db          *gorm.DB
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("無法讀取 .env 檔案")
	}
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

	// 初始化Redis連線
	initRedis()

	// 初始情況下刷新一次 Redis 快取
	refreshRedisCache()

	// 設定RESTful API路由
	setupAPIRoutes(router)

	// 啟動HTTP伺服器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server started on :" + port)
	router.Run(":" + port)
}


func setupAPIRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		todoController := controllers.NewTodoController(db, redisClient)
		api.GET("/todos", todoController.GetTodos)
		api.POST("/todos", todoController.CreateTodo)
		api.PUT("/todos/:id", todoController.UpdateTodo)
		api.DELETE("/todos/:id", todoController.DeleteTodo)
	}
}

func initDatabase() {
	dsn := os.Getenv("DB_DSN")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("無法連接資料庫")
	}

	// 自動建立資料表
	db.AutoMigrate(&models.Todo{})
}

func initRedis() {
	// 初始化Redis連線
	redisClient = InitRedis()
}


func InitRedis() *redis.Client {
	// 初始化並返回Redis連線客戶端
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return rdb
}

func refreshRedisCache() {
	var todos []models.Todo
	db.Find(&todos)

	todosJSON, _ := json.Marshal(todos)
	redisClient.Set(context.Background(), "todos", todosJSON, time.Hour)
}
