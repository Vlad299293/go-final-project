package api

import (
	"net/http"
	"time"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now := time.Now()
	if nowParam != "" {
		parsedNow, err := time.Parse(dateFormat, nowParam)
		if err != nil {
			http.Error(w, "некорректный параметр now", http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	next, err := NextDate(now, dateParam, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
