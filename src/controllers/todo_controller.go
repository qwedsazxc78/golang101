// controllers/todo_controller.go
package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"golang101/models"
)

type TodoController struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewTodoController(db *gorm.DB, rdb *redis.Client) *TodoController {
	return &TodoController{
		db:          db,
		redisClient: rdb,
	}
}

func (c *TodoController) CreateTodo(ctx *gin.Context) {
	var todo models.Todo
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.db.Create(&todo)
	ctx.JSON(http.StatusOK, todo)
	c.refreshRedisCache(ctx)
}

func (c *TodoController) GetTodos(ctx *gin.Context) {
	cachedData, err := c.redisClient.Get(ctx, "todos").Result()
	if err == nil {
		var todos []models.Todo
		json.Unmarshal([]byte(cachedData), &todos)

		ctx.JSON(http.StatusOK, todos)
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))

	var todos []models.Todo
	var total int64

	offset := (page - 1) * size
	c.db.Offset(offset).Limit(size).Find(&todos)
	c.db.Model(&models.Todo{}).Count(&total)

	todosJSON, _ := json.Marshal(todos)
	ctx.JSON(http.StatusOK, gin.H{
		"page":  page,
		"size":  size,
		"total": total,
		"data":  todos,
	})
	c.redisClient.Set(ctx, "todos", todosJSON, time.Hour)
}


func (c *TodoController) UpdateTodo(ctx *gin.Context) {
	var todo models.Todo
	todoID := ctx.Param("id")

	if err := c.db.First(&todo, todoID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "找不到該待辦事項"})
		return
	}

	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.db.Save(&todo)
	ctx.JSON(http.StatusOK, todo)
	c.refreshRedisCache(ctx)
}

func (c *TodoController) DeleteTodo(ctx *gin.Context) {
	var todo models.Todo
	todoID := ctx.Param("id")

	if err := c.db.First(&todo, todoID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "找不到該待辦事項"})
		return
	}

	c.db.Delete(&todo)
	ctx.JSON(http.StatusOK, gin.H{"message": "已刪除待辦事項"})
	c.refreshRedisCache(ctx)
}

// redis update strategy: after updating data, refresh cache
// note: it doesn't fit big table
func (c *TodoController) refreshRedisCache(ctx *gin.Context) {
	var todos []models.Todo
	c.db.Find(&todos)

	todosJSON, _ := json.Marshal(todos)
	c.redisClient.Set(ctx, "todos", todosJSON, time.Hour)
}