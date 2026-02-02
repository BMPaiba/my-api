package router

import (
	"github/mbpaiba/my-api/internal/db/sqlc"
	"github/mbpaiba/my-api/internal/user"

	"github.com/gin-gonic/gin"
)

func Setup(queries *sqlc.Queries) *gin.Engine {
	r := gin.Default()

	userH := user.NewHandler(user.NewService(user.NewRepository(queries)))

	r.GET("/users", userH.GetUsers)

	return r
}
