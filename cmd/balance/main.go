package main

import (
	"context"
	//"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

)

type Config struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBSource   string `mapstructure:"DB_SOURCE"`
	ServerPort string `mapstructure:"SERVER_PORT"`
}

func main() {
	//db, err := sql.Open("postgres", )
	//ur := users.NewRepo(db)
	sm := mux.NewRouter()

	server := &http.Server{
		Addr:         ":3333",
		Handler:      sm,
		IdleTimeout:  20*time.Second,
		ReadTimeout:  1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	go func() {
		fmt.Printf("Server is listening on port%s\n", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("Error starting server: %s", err)
			os.Exit(1)
		}
	}()

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt)
	signal.Notify(sigChannel, os.Kill)

	sig := <-sigChannel
	fmt.Println("Command to terminate received, shutdown", sig)

	timeoutContext, finish := context.WithTimeout(context.Background(), 30*time.Second)
	defer finish()
	server.Shutdown(timeoutContext)
}