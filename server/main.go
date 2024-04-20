package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"server/api"
	"server/config"
	"server/utils"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

func main() {
	CONFIG_PATH := os.Getenv("CONFIG_PATH")
	if CONFIG_PATH == "" {
		CONFIG_PATH = "data/config.yaml"
	}
	if err := config.SetupGlobalConfig(CONFIG_PATH); err != nil {
		log.Fatal(err)
	}
	configs := config.Get()

	rootPath, _ := getPackageRootPath()
	rootLevel, err := log.ParseLevel(configs.Log.Root)
	if err != nil {
		rootLevel = log.InfoLevel
		log.Warn("Failed to parse root log level, using default level: info")
	}
	log.SetLevel(rootLevel)
	log.SetReportCaller(true)

	if configs.Log.WithColor {
		log.SetOutput(colorable.NewColorableStdout())
		log.SetFormatter(&log.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				return "\n    ", fmt.Sprintf(" %s:%d", strings.TrimPrefix(f.File, rootPath), f.Line)
			},
		})
	}

	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()

	router := gin.New()
	router.Use(cors.New(corsConfig(configs.API.AllowOrigins)))
	api.RegisterRoutes(router)

	srv := &http.Server{Addr: configs.API.Address, Handler: router}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		port := utils.GetPortFromURL(configs.API.Address)
		if port == "" {
			port = "80"
		}

		log.Infoln("Server starting...")
		log.Infof("Server listening on http://localhost:%s/\n", port)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Infoln("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Infoln("Server exiting")
}

func getPackageRootPath() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("failed to get caller information")
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	return filepath.Dir(absPath) + "/", nil
}

func corsConfig(allowOrigins []string) cors.Config {
	config := cors.DefaultConfig()

	// , "DELETE", "PUT"
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{
		"Authorization", "Content-Type", "Origin", "Content-Length",
		"Connection", "Accept-Encoding", "Accept-Language", "Host",
	}
	if gin.Mode() == gin.DebugMode {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = allowOrigins
	}

	config.MaxAge = 1 * time.Hour
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}

	return config
}
