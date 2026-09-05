package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go_final_project/pkg/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	if _, err := time.Parse(dateFormat, task.Date); err != nil {
		return errors.New("дата представлена в неверном формате")
	}

	today := now.Format(dateFormat)

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		if task.Date < today {
			task.Date = next
		}
	} else if task.Date < today {
		task.Date = today
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}
