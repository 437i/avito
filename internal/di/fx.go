package di

import (
	"avito/internal/handler/common"
	handlerCourier "avito/internal/handler/courier"
	repoCourier "avito/internal/repository/courier"
	serviceCourier "avito/internal/service/courier"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Module("courier",
	fx.Provide(
		initDBPool, 
		repoCourier.NewRepository, 
		serviceCourier.NewService,
		handlerCourier.NewHandler,
	),
	fx.Invoke(startServer),
)

func startServer(lc fx.Lifecycle, dbPool *pgxpool.Pool, courierHandler *handlerCourier.Handler, port string) {
	var srv *http.Server
	serverErr := make(chan error, 1)

	lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				srv = &http.Server{
					Addr:    ":" + port,
					Handler: getRouter(courierHandler),
				}

				log.Printf("Server started on %s\n", srv.Addr)
				go func() {
					if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						serverErr <- err
						close(serverErr)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				gracefulShutdown(ctx, srv, dbPool, serverErr)
				log.Println("Shutting down service-courier")
				return nil
			},
	})
}

func gracefulShutdown(ctx context.Context, srv *http.Server, dbPool *pgxpool.Pool, serverErr <-chan error) {
	var reason string
	select {
	case err := <-serverErr:
		reason = "server error: " + err.Error()
	default:
	}
	log.Printf("Shutdown initiated (%s)", reason)

	log.Println("Shutting down HTTP server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("HTTP server stopped")
	}

	log.Println("Closing DB pool...")
	dbPool.Close()
	log.Println("DB pool closed")
}

func waitGracefulShutdown(srv *http.Server, dbPool *pgxpool.Pool, serverErr <-chan error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var reason string
	select {
	case <-ctx.Done():
		reason = "signal"
	case err := <-serverErr:
		reason = "server error: " + err.Error()
	}
	log.Printf("Shutdown initiated (%s)", reason)

	log.Println("Shutting down HTTP server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("HTTP server stopped")
	}

	log.Println("Closing DB pool...")
	dbPool.Close()
	log.Println("DB pool closed")
}

func getRouter(handler *handlerCourier.Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/ping", common.Ping).Methods("GET")
	r.HandleFunc("/healthcheck", common.Healthcheck).Methods("HEAD")

	r.HandleFunc("/couriers", handler.GetCouriers).Methods("GET")

	courier := r.PathPrefix("/courier").Subrouter()
	courier.HandleFunc("/{id}", handler.GetCourier).Methods("GET")
	courier.HandleFunc("", handler.CreateCourier).Methods("POST")
	courier.HandleFunc("", handler.UpdateCourier).Methods("PUT")
	courier.HandleFunc("/{id}", handler.Delete).Methods("DELETE")
	return r
}

func initDBPool() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(connString())
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
	}
	config.MaxConns = 20                      // Максимальное количество соединений
	config.MinConns = 5                       // Минимальное количество соединений
	config.MaxConnLifetime = time.Hour        // Время жизни соединения
	config.MaxConnIdleTime = time.Minute * 30 // Время бездействия
	config.HealthCheckPeriod = time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	pingAttemptsLimit := 3
	var pingErr error

	for i := 1; i <= pingAttemptsLimit; i++ {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pingErr = dbPool.Ping(pingCtx)
		pingCancel()
		if pingErr == nil {
			break
		}
		log.Printf("db ping attempt %d failed: %v", i, pingErr)
		if i < pingAttemptsLimit {
			time.Sleep(500 * time.Millisecond)
		}
	}

	if pingErr != nil {
		log.Fatalf("Unable to ping database")
	}

	log.Println("Database connection pool established")
	return dbPool
}

func connString() string {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pswd := os.Getenv("POSTGRES_PASSWORD")
	db :=   os.Getenv("POSTGRES_DB")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pswd, host, port, db)
}