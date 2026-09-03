package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

const tasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	var (
		tasks []*db.Task
		err   error
	)

	if search != "" {
		tasks, err = db.TasksSearch(search, tasksLimit)
	} else {
		tasks, err = db.Tasks(tasksLimit)
	}

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, TasksResp{Tasks: tasks})
}
