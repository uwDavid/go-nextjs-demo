package api

import (
	"fmt"
	"log"
	"nextjs/backend/utils"

	"github.com/gin-gonic/gin"
)

func NewServer(port int) {
	// Load env configs
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Could not load env config", err)
	}

	// Initialize Gin
	g := gin.Default()
	g.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Welcome to Fintech"})
	})

	//g.Run(":3000")
	g.Run(fmt.Sprintf(":%v", port))
}
