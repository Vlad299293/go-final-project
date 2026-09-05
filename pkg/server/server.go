package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"go_final_project/pkg/api"
)

const defaultPort = 7540
const webDir = "./web"

func Run() error {
	port := defaultPort

	envPort := os.Getenv("TODO_PORT")
	if envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := ":" + strconv.Itoa(port)
	log.Printf("Сервер запущен на порту %d", port)

	return http.ListenAndServe(addr, nil)
}
