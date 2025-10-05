package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notes-project/internal/metrics"
	"notes-project/internal/ws"
	"os"
	"time"

	"notes-project/internal/handlers"
	_ "notes-project/internal/logger"
	"notes-project/internal/repository"
	"notes-project/internal/service"

	_ "notes-project/docs"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Trello Clone API
// @version         1.0
// @description     Это API для клона Trello, созданного на Go.
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey ApiKeyAuth
// @in              header
// @name            Authorization
func main() {

	dbHost := env("DB_HOST", "localhost")
	dbPort := env("DB_PORT", "5433")
	dbUser := env("DB_USER", "notes_user")
	dbPassword := env("DB_PASSWORD", "notes_password")
	dbName := env("DB_NAME", "notes_db")
	sslMode := env("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}
	log.Println("Connected to DB")

	redisAddr := env("REDIS_ADDR", "localhost:6380")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("cannot connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")

	hub := ws.NewHub()
	go hub.Run()
	appMetrics := metrics.NewAppMetrics()

	userRepo := repository.NewUserRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	listRepo := repository.NewListRepository(db)
	cardRepo := repository.NewCardRepository(db)

	userService := service.NewUserService(userRepo)
	boardService := service.NewBoardService(boardRepo, listRepo, cardRepo, userRepo, hub, rdb)
	cacheInvalidator := boardService.InvalidateBoardCache
	listService := service.NewListService(listRepo, boardRepo, hub, cacheInvalidator)
	cardService := service.NewCardService(cardRepo, listRepo, boardRepo, hub, cacheInvalidator)

	userHandler := handlers.NewUserHandler(userService)
	boardHandler := handlers.NewBoardHandler(boardService)
	listHandler := handlers.NewListHandler(listService)
	cardHandler := handlers.NewCardHandler(cardService)
	wsHandler := ws.NewWsHandler(hub, boardService)

	router := setupRouter(userHandler, boardHandler, listHandler, cardHandler, wsHandler, appMetrics)

	port := env("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func setupRouter(
	userHandler *handlers.UserHandler,
	boardHandler *handlers.BoardHandler,
	listHandler *handlers.ListHandler,
	cardHandler *handlers.CardHandler,
	wsHandler *ws.WsHandler,
	appMetrics *metrics.AppMetrics,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(handlers.MetricsMiddleware(appMetrics))
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {})

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		userHandler.RegisterPublicRoutes(api.Group("/users"))
		protectedRoutes := api.Group("/")
		protectedRoutes.Use(handlers.AuthMiddleware())
		{
			userHandler.RegisterProtectedRoutes(protectedRoutes)
			boardHandler.RegisterBoardRoutes(protectedRoutes)
			listHandler.RegisterListRoutes(protectedRoutes)
			cardHandler.RegisterCardRoutes(protectedRoutes)
			wsHandler.RegisterWsRoutes(protectedRoutes)
		}
	}

	return r
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
