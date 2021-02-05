package main

import (
	"github.com/gin-gonic/gin"
	"github.com/test/services"
)

func main() {
	r := gin.Default()

	//Add Services
	services.InitServices(r)

	//Start Server
	r.Run(":9091")
}
