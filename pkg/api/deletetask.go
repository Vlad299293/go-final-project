package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}
