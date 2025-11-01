package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"support-bot/config"
	"support-bot/pkg/postgres"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Run(configPath string) {
	// Configuration
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("app - Run - config error: %w", err)
	}

	// Logger
	logrusLevel, err := log.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(logrusLevel)
	}
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
	log.Info("Init application...")

	// DB
	log.Info("Init postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatalf("app - Run - postgres error: %w", err)
	}
	defer pg.Close()

	// Repositories
	log.Info("Init repositories...")

	// Services
	log.Info("Init services...")

	// TG bot
	log.Info("Init Telegarm bot...")

	// Server healthy
	log.Info("Init HTTP server...")
	srv := runServer(cfg.HTTP.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("STOP signal received")

	// TG bot stopping
	log.Info("Bot stopping...")
	log.Info("Bot stopped")

	// Wait 5 sec
	log.Info("Waitting...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Server stopping
	log.Info("Server stopping...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Info("Server stopped")
}

func runServer(port string) *http.Server {
	type MetricsResponse struct {
		Goroutines    int    `json:"goroutines"`
		MemoryUsageMB uint64 `json:"memory_usage_mb"`
		CPUCores      int    `json:"cpu_cores"`
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.GET("/metrics", func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		c.JSON(http.StatusOK, MetricsResponse{
			Goroutines:    runtime.NumGoroutine(),
			MemoryUsageMB: m.Alloc / 1024 / 1024,
			CPUCores:      runtime.NumCPU(),
		})
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("app - Run - listen error: %s\n", err)
		}
	}()

	return srv
}
