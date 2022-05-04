package main

import (
	"go-api-2/models"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "go-api-2/docs" // import 檔案夾的名稱要以 go.mod 的名稱來命名，再加檔案夾名稱

	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// sets up an association in which GetDatass handles requests to the /mysql endpoint path.
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server mysql server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8088
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		router := v1.Group("/logs") // Group 的 /logs 要做為其他的 endpoint
		{
			router.GET("", getLogs) // swagger 寫法 @router  /logs [get]
			router.GET(":ip", getLog)
			router.POST("", addLog)
		}
	}

	// Use the Run function to attach the router to an http.Server and start the server.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8088") // router.Run("localhost:8088") 會因資安吽報錯
}

// getAllLogs godoc
// @Summary      Get all ipLogs
// @Description  get all ipLogs
// @Tags         ipLogs
// @Accept       json
// @Produce      json
// @Success      200  string   success
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /logs [get]
func getLogs(c *gin.Context) {

	logs := models.GetLogs()

	if logs == nil || len(logs) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, logs)
	}
}

// getLog godoc
// @Summary      Show a IP's logs
// @Description  get string by IP
// @Tags         Logs
// @Accept       json
// @Produce      json
// @Param        ip    query     string  true  "iplogs search by string"  Format(IP)
// @Router       /logs/:ip [get]
func getLog(c *gin.Context) {
	ip := c.Param("ip")

	log := models.GetLog(ip)

	if log == nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, log)
	}
}

// addLog godoc
// @Summary      Add an ipLog
// @Description  add by json body
// @Tags         iplogs
// @Accept       json
// @Produce      json
// @Param        ipLog   body      models.Log  true  "Add iplogs"
// @Success      200      string	success
// @Failure      400      {object}  httputil.HTTPError
// @Failure      404      {object}  httputil.HTTPError
// @Failure      500      {object}  httputil.HTTPError
// @Router       /logs [post]
func addLog(c *gin.Context) {
	var logIp models.Log

	if err := c.BindJSON(&logIp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		models.AddLog(logIp)
		c.IndentedJSON(http.StatusCreated, logIp)
	}
}
