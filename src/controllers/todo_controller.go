package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "strconv"
    "golang101/models"
)

type TodoController struct {
    db *gorm.DB
}

// 初始化TodoController
func NewTodoController(db *gorm.DB) *TodoController {
    return &TodoController{
        db: db,
    }
}

// 新增待辦事項
func (c *TodoController) CreateTodo(ctx *gin.Context) {
    var todo models.Todo
    if err := ctx.ShouldBindJSON(&todo); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.db.Create(&todo)
    ctx.JSON(http.StatusOK, todo)
}

// 取得所有待辦事項
func (c *TodoController) GetTodos(ctx *gin.Context) {
    // 獲取分頁參數
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))

    var todos []models.Todo
    var total int64

    // 分頁查詢
    offset := (page - 1) * size
    c.db.Offset(offset).Limit(size).Find(&todos)

    // 獲取總數
    c.db.Model(&models.Todo{}).Count(&total)

    ctx.JSON(http.StatusOK, gin.H{
        "page":  page,
        "size":  size,
        "total": total,
        "data":  todos,
    })
}

// 更新待辦事項
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
}

// 刪除待辦事項
func (c *TodoController) DeleteTodo(ctx *gin.Context) {
    var todo models.Todo
    todoID := ctx.Param("id")

    if err := c.db.First(&todo, todoID).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "找不到該待辦事項"})
        return
    }

    c.db.Delete(&todo)
    ctx.JSON(http.StatusOK, gin.H{"message": "已刪除待辦事項"})
}

// ... （其他操作，根據需求新增更多方法）
