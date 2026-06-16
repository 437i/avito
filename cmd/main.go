package main

import (
	"avito/internal/di"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

// проверить валидацию имени
// как сделать jitter для ретрай пинга базы
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error loading env: %v", err)
	}
	port := getPort()

	app := fx.New(
		fx.Provide(func() string {return port}),
		di.Module,
	)
	app.Run()
}

func getPort() string {
	port := os.Getenv("PORT")
	portFlag := getFlag()
	if portFlag != "" {
		port = portFlag
	}
	if port == "" {
		port = "8080"
	}
	return port
}

func getFlag() string {
	var port string
	flag.StringVar(&port, "port", "", "sets server port; default empty")
	flag.Parse()
	return port
}

// func startServer(dbPool *pgxpool.Pool, courierHandler *handler.Handler, port string) {
// 	srv := &http.Server{
// 		Addr:    ":" + port,
// 		Handler: getRouter(courierHandler),
// 	}
// 	log.Printf("Server started on %s\n", srv.Addr)
// 	serverErr := make(chan error, 1)
// 	go func() {
// 		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
// 			serverErr <- err
// 		}
// 	}()

// 	waitGracefulShutdown(srv, dbPool, serverErr)

// 	log.Println("Shutting down service-courier")
// }