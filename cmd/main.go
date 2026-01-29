package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	providerfx "oolio/internal/app/fx"
	"oolio/internal/app/services"
	"oolio/internal/app/worker"
	"oolio/internal/config"
	"oolio/internal/database"
)

func main() {
	app := fx.New(
		providerfx.AppModule,
		fx.Options(
			fx.Provide(
				NewLogger,
				NewHTTPServer,
			),
		),
		fx.Invoke(StartServer),
	)

	app.Run()
}

func NewLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func NewHTTPServer(
	cfg *config.Config,
	ginRouter *gin.Engine,
	lc fx.Lifecycle,
	logger *zap.Logger,
) *http.Server {

	server := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      ginRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting HTTP server",
				zap.String("address", server.Addr))

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("Failed to start server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server")

			shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			return server.Shutdown(shutdownCtx)
		},
	})

	return server
}

func StartServer(
	lc fx.Lifecycle,
	server *http.Server,
	db *database.Database,
	couponService services.CouponService,
	orderWorker *worker.OrderWorker,
	logger *zap.Logger,
) {
	go func() {
		ctx := context.Background()
		if err := couponService.DownloadAndParseCouponFiles(ctx); err != nil {
			logger.Error("Failed to initialize coupon service", zap.Error(err))
		} else {
			logger.Info("Coupon service initialized successfully")
		}

		go couponService.StartPeriodicRefresh(ctx, 24*time.Hour)
	}()

	go func() {
		ctx := context.Background()
		orderWorker.Start(ctx)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Shutdown signal received")

		if err := db.Close(); err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
		}
	}()

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Application stopped gracefully")
			return nil
		},
	})

	logger.Info("Food Ordering API Backend started successfully")
}

func handleError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
