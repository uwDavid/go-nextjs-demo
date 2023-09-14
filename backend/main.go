package main

import (
	"nextjs/backend/api"
	db "nextjs/backend/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	queries *db.Queries
	router  *gin.Engine
}

func main() {
	api.NewServer(3000)
}
