package main

import (
	"go-api-2/models"

	"net/http"

	"github.com/gin-gonic/gin"

	_ "go-api-2/docs" // import 檔案夾的名稱要以 go.mod 的名稱來命名，再加檔案夾名稱
)

func main() {
	router := gin.Default()
	// v1 := r.Group("/api/v1")

	// router := v1.Group("/mysql")

	router.GET("/logs", getLogs)
	router.GET("/logs/:ip", getLog)
	router.POST("/log", models.PostData)

	// Use the Run function to attach the router to an http.Server and start the server.
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8088") // router.Run("localhost:8088") 會因資安吽報錯
}

func getLogs(c *gin.Context) {
	logs := models.GetLogs()

	if logs == nil || len(logs) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, logs)
	}
}

func getLog(c *gin.Context) {
	ip := c.Param("ip")

	log := models.GetLog(ip)

	if log == nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, log)
	}
}

func addLog(c *gin.Context) {
	var logIp models.Log

	if err := c.BindJSON(&logIp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		models.AddLog(logIp)
		c.IndentedJSON(http.StatusCreated, logIp)
	}
}
