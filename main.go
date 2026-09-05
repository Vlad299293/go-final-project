package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

const defaultDBFile = "scheduler.db"

func main() {
	dbFile := defaultDBFile
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		dbFile = envDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
	db.Close()
}
