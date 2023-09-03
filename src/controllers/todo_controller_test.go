package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"golang101/models"
)

// 定義 ControllerTest 介面
type ControllerTest interface {
	Setup()
	TearDown()
	SetupRouter() *gin.Engine
}

// TodoControllerTest 實作 ControllerTest 介面
type TodoControllerTest struct {
	controller *TodoController
	db         *gorm.DB
}

func (t *TodoControllerTest) Setup() {
	// 設定測試資料庫連線
	t.db, _ = gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	// 初始化TodoController
	t.controller = NewTodoController(t.db)
}

func (t *TodoControllerTest) TearDown() {
	// 關閉資料庫連線
	sqlDB, _ := t.db.DB()
	sqlDB.Close()
}

func (t *TodoControllerTest) SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/api/todos", t.controller.CreateTodo)
	router.GET("/api/todos", t.controller.GetTodos)
	return router
}

func TestMain(m *testing.M) {
	test := &TodoControllerTest{}
	test.Setup()
	defer test.TearDown()

	// 執行測試
	m.Run()
}

func TestTodoController_CreateTodo(t *testing.T) {
	test := &TodoControllerTest{}
	test.Setup()
	defer test.TearDown()

	// 創建測試用的Gin引擎
	router := test.SetupRouter()

	// 建立HTTP測試請求
	req, _ := http.NewRequest("POST", "/api/todos", bytes.NewBuffer([]byte(`{"Title":"測試","Description":"這是一個測試","Completed":false}`)))
	w := httptest.NewRecorder()

	// 發送HTTP請求
	router.ServeHTTP(w, req)

	// 斷言回應狀態碼是否為200
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTodoController_GetTodos(t *testing.T) {
	test := &TodoControllerTest{}
	test.Setup()
	defer test.TearDown()

	// 創建測試用的Gin引擎
	router := test.SetupRouter()

	// 建立HTTP測試請求
	req, _ := http.NewRequest("GET", "/api/todos", nil)
	w := httptest.NewRecorder()

	// 發送HTTP請求
	router.ServeHTTP(w, req)

	// 斷言回應狀態碼是否為200
	assert.Equal(t, http.StatusOK, w.Code)
}
