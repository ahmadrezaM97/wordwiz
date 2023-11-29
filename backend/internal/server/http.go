package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"wordwiz/config"
	"wordwiz/internal/server/handler"
	"wordwiz/internal/storage"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"wordwiz/internal/server/docs" // docs is generated by Swag CLI, you have to import it.

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	srv *http.Server
}

func New(cfg config.Config, stg storage.Storage) *Server {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:           route(cfg, stg),
		ReadHeaderTimeout: cfg.HTTP.DefaultTimeout,
	}

	return &Server{srv: srv}
}

func (s *Server) Serve(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := s.srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("failed to shutdown srv")
		}
	}()

	return s.srv.ListenAndServe()
}

func route(cfg config.Config, stg storage.Storage) http.Handler {
	docs.SwaggerInfo.Title = cfg.App.Name
	docs.SwaggerInfo.Version = cfg.App.Version

	ginEngine := gin.New()

	config := cors.DefaultConfig()

	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	config.AllowHeaders = append(config.AllowHeaders, "*")
	config.AllowMethods = append(config.AllowMethods, "OPTIONS")
	ginEngine.Use(cors.New(config))

	initialGinMiddleware(ginEngine, cfg.App.Environment)

	h := handler.New(cfg, stg)

	ginEngine.GET("/health", h.Health)

	ginEngine.POST("/words/add", h.AuthMiddleware(), h.AddWord)
	ginEngine.GET("/user/words", h.AuthMiddleware(), h.GetUserWords)

	ginEngine.GET("/auth/google/login", h.GoogleLogin)
	ginEngine.GET("/auth/google/callback", h.HandleGoogleCallback)

	// use ginSwagger middleware to serve the API docs
	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return ginEngine
}

func initialGinMiddleware(ginEngine *gin.Engine, envName string) {
	switch envName {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
		ginEngine.Use(gin.Recovery())
	case "stage":
		gin.SetMode(gin.TestMode)
		ginEngine.Use(gin.Recovery())
	case "dev":
		fallthrough
	case "debug":
		gin.SetMode(gin.DebugMode)
		ginEngine.Use(gin.Logger())
	default:
		gin.SetMode(gin.DebugMode)
	}
}
