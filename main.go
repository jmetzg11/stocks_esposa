package main

import (
	"fmt"
	"stocks/backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println("server started")

	r := gin.Default()
	routes.SetupAPIRoutes(r)

	r.Run(":8000")
}
