package main

import (
	"fmt"
	"os"
	"strconv"
	"taskboard/controller/users"
	"taskboard/model"
	"taskboard/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Init database
	fmt.Println("Initializing database...")
	databasePath := "./taskboard.sqlite3"
	err := orm.Init(databasePath)
	if err != nil {
		fmt.Printf("Failed to initialize database. error:%+v\n", err)
		return
	}

	// Create or Update tables
	err = orm.Migrate(
		&model.User{},
		&model.Task{},
		&model.Board{},
		&model.TaskEntry{},
	)
	if err != nil {
		fmt.Printf("Failed to update tables. error:%+v\n", err)
		return
	}

	// Init router of REST apis
	router := gin.Default()
	//config := cors.DefaultConfig()
	//config.AllowAllOrigins = true
	//config.AllowHeaders = []string{"Content-Type"}
	//router.Use(cors.New(config))
	router.Use(cors.Default())
	// Register api path
	routeGroup := router.Group("/taskboard")
	users.RegisterRoute(routeGroup)

	// Set listening port
	port := getListeningPort()
	// Start server
	fmt.Printf("Taskboard api server is starting... listening port:%d\n", port)
	err = router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Failed to start api server. error:%+v\n", err)
		return
	}
}

func getListeningPort() int {
	portEnv := os.Getenv("TASKBOARD_API_SERVER_PORT")
	port := 0
	if portEnv != "" {
		port, err := strconv.Atoi(portEnv)
		if err != nil || port < 0 {
			port = 0
		}
	}
	if port == 0 {
		fmt.Println("Environment variable [TASKBOARD_API_SERVER_PORT] is not set or invalid, 7000 port is used as defalt.")
		port = 7000
	}
	return port
}
