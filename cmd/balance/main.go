package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/ghost171/avito_currency_deposition/tree/main/pkg/config"
	"github.com/ghost171/avito_currency_deposition/tree/main/pkg/handler"
	"github.com/ghost171/avito_currency_deposition/tree/main/pkg/users"

)

func main() {
	cfg, err := config.Load("configs")
	if err != nil {
		log.Fatal("Error while reading config file:", err)
	}
	db, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal("Error while sql.Open: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to db, err: %v\n", err)
	}
	
	ur := users.NewRepo(db)

	sm := mux.NewRouter()

	uh := handler.UserHandler(ur)

	sm.HandleFunc("/deposit", uh.Deposit).Methods("POST")
	sm.HandleFunc("/cashout", uh.Cashout).Methods("POST")
	sm.HandleFunc("/transfer", uh.Transfer).Methods("POST")
	sm.HandleFunc("/value", uh.GetValue).Methods("GET")
	sm.HandleFunc("/operations", uh.ListOperations).Methods("GET")

	server := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
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

	timeoutContext, finish := context.WithTimeout(context.Background(), 30 * time.Second)
	defer finish()
	server.Shutdown(timeoutContext)
}