package app

import (
	"fmt"
	"log"

	"github.com/alexey-shedrin/avito-test-task/internal/config"
	"github.com/alexey-shedrin/avito-test-task/internal/database"
	openapi "github.com/alexey-shedrin/avito-test-task/internal/gen"
	pvzv1 "github.com/alexey-shedrin/avito-test-task/internal/grpc/pvz/v1"
	"github.com/alexey-shedrin/avito-test-task/internal/handler"
	"github.com/alexey-shedrin/avito-test-task/internal/metrics"
	"github.com/alexey-shedrin/avito-test-task/internal/repository"
	"github.com/alexey-shedrin/avito-test-task/internal/service"
	"github.com/gin-gonic/gin"
)

func Run() {
	cfg := config.New()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	pvzRepo := repository.NewPVZRepository(db)
	receptionRepo := repository.NewReceptionRepository(db)

	userService := service.NewUserService(userRepo)
	pvzService := service.NewPVZService(pvzRepo)
	receptionService := service.NewReceptionService(receptionRepo, db)

	hndlr := handler.New(userService, pvzService, receptionService)
	r := gin.Default()

	openapi.RegisterHandlers(r, hndlr)

	r.Use(metrics.GetMetricsMiddleware())

	serverAddr := fmt.Sprintf("%s:%s", cfg.HttpServer.Host, cfg.HttpServer.Port)
	go func() {
		pvzv1.Start(cfg.GrpcServer.Port)
	}()

	go func() {
		metrics.StartMetricsServer(cfg.PrometheusServer.Port)
	}()

	r.Run(serverAddr)
}
