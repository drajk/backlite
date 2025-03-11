package ui

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/drajk/backlite/internal/task"
)

type (
	Handler struct {
		db *sql.DB
	}

	TemplateData struct {
		Path    string
		Content any
	}
)

func NewHandler(db *sql.DB) *http.ServeMux {
	h := &Handler{db: db}
	mux := http.NewServeMux()

	routes := []struct {
		Path    string
		Handler func(http.ResponseWriter, *http.Request)
	}{
		{Path: "/", Handler: h.Running},
		{Path: "/upcoming", Handler: h.Upcoming},
		{Path: "/succeeded", Handler: h.Succeeded},
		{Path: "/failed", Handler: h.Failed},
		{Path: "/task/{task}", Handler: h.Task},
		{Path: "/completed/{task}", Handler: h.TaskCompleted},
	}

	// Register routes with the provided prefix
	for _, route := range routes {
		fullPath := prefix + route.Path
		mux.HandleFunc("GET "+fullPath, route.Handler)
	}

	return mux
}

func (h *Handler) Running(w http.ResponseWriter, req *http.Request) {
	err := func() error {
		tasks, err := task.GetTasks(req.Context(), h.db, selectRunningTasks, itemLimit)
		if err != nil {
			return err
		}

		return h.render(req, w, tmplTasksRunning, tasks)
	}()

	if err != nil {
		h.error(w, err)
	}
}

func (h *Handler) Upcoming(w http.ResponseWriter, req *http.Request) {
	err := func() error {
		// TODO use actual time from the client
		tasks, err := task.GetScheduledTasks(req.Context(), h.db, time.Now().Add(time.Hour), itemLimit)
		if err != nil {
			return err
		}

		return h.render(req, w, tmplTasksUpcoming, tasks)
	}()

	if err != nil {
		h.error(w, err)
	}
}

func (h *Handler) Succeeded(w http.ResponseWriter, req *http.Request) {
	err := func() error {
		tasks, err := task.GetCompletedTasks(req.Context(), h.db, selectCompletedTasks, 1, itemLimit)
		if err != nil {
			return err
		}

		return h.render(req, w, tmplTasksCompleted, tasks)
	}()

	if err != nil {
		h.error(w, err)
	}
}

func (h *Handler) Failed(w http.ResponseWriter, req *http.Request) {
	err := func() error {
		tasks, err := task.GetCompletedTasks(req.Context(), h.db, selectCompletedTasks, 0, itemLimit)
		if err != nil {
			return err
		}

		return h.render(req, w, tmplTasksCompleted, tasks)
	}()

	if err != nil {
		h.error(w, err)
	}
}

func (h *Handler) Task(w http.ResponseWriter, req *http.Request) {
	var t *task.Task

	err := func() error {
		id := req.PathValue("task")
		tasks, err := task.GetTasks(req.Context(), h.db, selectTask, id)
		if err != nil {
			return err
		}

		if len(tasks) > 0 {
			t = tasks[0]
			return h.render(req, w, tmplTask, t)
		}

		return nil
	}()

	if err != nil {
		h.error(w, err)
	} else if t == nil {
		// If no task found, try the same ID as a completed task.
		h.TaskCompleted(w, req)
	}
}

func (h *Handler) TaskCompleted(w http.ResponseWriter, req *http.Request) {
	err := func() error {
		var t *task.Completed
		id := req.PathValue("task")
		tasks, err := task.GetCompletedTasks(req.Context(), h.db, selectCompletedTask, id)
		if err != nil {
			return err
		}

		if len(tasks) > 0 {
			t = tasks[0]
		}

		return h.render(req, w, tmplTaskCompleted, t)
	}()

	if err != nil {
		h.error(w, err)
	}
}

func (h *Handler) error(w http.ResponseWriter, err error) {
	fmt.Fprint(w, err)
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (h *Handler) render(req *http.Request, w io.Writer, tmpl *template.Template, data any) error {
	return tmpl.ExecuteTemplate(w, "layout.gohtml", TemplateData{
		Path:    req.URL.Path,
		Content: data,
	})
}
