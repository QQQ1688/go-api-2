package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql" // _ = 只使用 mysql 裡面的 init()
)

const (
	USERNAME = "root"
	PASSWORD = "dc0906708652"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "iplog"
)

// const dbuser = "root"
// const dbpass = "dc0906708652"
// const dbname = "ip_log"

func GetLogs() []Log {

	// db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp(127.0.0.1:3306)/"+dbname)
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)

	// if there is an error opening the connection, handle it
	if err != nil {
		// simply print the error to the console
		fmt.Println("Err connecting mysql db", err.Error())
		// returns nil on error
		return nil
	}

	defer db.Close() // defer 越早寫得越晚執行
	results, err := db.Query("SELECT * FROM ip_log")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	logs := []Log{}
	for results.Next() {
		var logIp Log
		// for each row, scan into the Product struct&logIp.IP, &logIp.Time, &logIp.URL, &logIp.Status
		err = results.Scan()
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// append the product into products array
		logs = append(logs, logIp)
	}

	return logs

}

func GetLog(IP string) *Log {

	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)

	logIp := &Log{}
	if err != nil {
		// simply print the error to the console
		fmt.Println("Err", err.Error())
		// returns nil on error
		return nil
	}

	defer db.Close()

	results, err := db.Query("SELECT * FROM ip_log where IP=?", IP)

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	if results.Next() {
		err = results.Scan(&logIp.IP, &logIp.Time, &logIp.URL, &logIp.Status)
		if err != nil {
			return nil
		}
	} else {

		return nil
	}

	return logIp
}

func AddLog(log Log) {

	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)

	if err != nil {
		panic(err.Error())
	}

	// defer the close till after this function has finished
	// executing
	defer db.Close()

	insert, err := db.Query(
		"INSERT INTO ip_log (IP,Time,URL,Status) VALUES (?,?,?,?)",
		log.IP, log.Time, log.URL, log.Status)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

}

// PostData godoc
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
func PostData(c *gin.Context) {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		c.JSON(http.StatusBadGateway, "Connect DB failed")
		return
	}
	if err := db.Ping(); err != nil {
		c.JSON(http.StatusBadGateway, "Connect DB failed")
		return
	}
	fmt.Println("DB iplog connected")

	decoder := json.NewDecoder(c.Request.Body)
	var newData Log // 要收input變數的 newData 以 data 為 struct 用 BindJSON 對應鍵值
	// Call BindJSON to bind the received JSON to
	// newData.
	err2 := decoder.Decode(&newData)
	if err2 != nil {
		panic(err.Error())
	}
	fmt.Println(newData.IP)

	c.JSON(http.StatusOK, gin.H{
		"IP":     newData.IP,
		"Time":   newData.Time,
		"URL":    newData.URL,
		"Status": newData.Status,
	})

	// Add the newData to the mysql db.
	if _, err := db.Exec(
		"INSERT INTO ip_log (IP, Time, Url, Status) VALUES (?, ?, ?, ?)", &newData.IP, &newData.Time, &newData.URL, &newData.Status); err != nil {
		fmt.Printf("add newData: %v", err)
		return
	}
	fmt.Println("iplogs posted successed")
}
