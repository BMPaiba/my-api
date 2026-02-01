package main

import (
	"fmt"
	"github/mbpaiba/my-api/internal/env"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	// agregar if condicional dependiendo si esta en produccion
	gin.SetMode(gin.ReleaseMode)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	port := fmt.Sprintf(":%v", env.GetString("PORT", "3030"))

	if err := r.Run(port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
