package http

import (
	"go-hex/internal/adapter/config"
	"go-hex/internal/core/port"
	"log/slog"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type Router struct {
	*gin.Engine
}

func NewRouter(config *config.HTTP, token port.TokenService, userHandler UserHandler, authHandler AuthHandler) (*Router, error) {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ginConfig := cors.DefaultConfig()
	allowedOrigins := config.AllowedOrigins
	originsList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originsList

	router := gin.New()
	router.Use(sloggin.New(slog.Default()), gin.Recovery(), cors.New(ginConfig))

	// v, ok := binding.Validator.Engine().(*validator.Validate)
	// if ok {
	// 	if err := v.RegisterValidation("user_role", userRoleValidator); err != nil {
	// 		return nil, err
	// 	}
	// }

	v1 := router.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/", userHandler.Register)
			user.POST("/login", authHandler.Login)
		}
	}

	return nil, nil
}

func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
