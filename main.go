package main

import (
	"log"
	"os"

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

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
