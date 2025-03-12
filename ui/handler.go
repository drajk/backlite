package ui

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/drajk/backlite/internal/task"
	"github.com/labstack/echo/v4"
)

type (
	Handler struct {
		db     *sql.DB
		prefix string
	}

	TemplateData struct {
		Path    string
		Prefix  string
		Content any
	}
)

// NewHandler accepts a prefix and an echo.Group
func NewHandler(g *echo.Group, prefix string, db *sql.DB) {
	h := &Handler{db: db, prefix: prefix}

	if prefix != "" && !hasLeadingSlash(prefix) {
		prefix = "/" + prefix
	}

	g.GET(prefix+"/running", h.Running)
	g.GET(prefix+"/upcoming", h.Upcoming)
	g.GET(prefix+"/succeeded", h.Succeeded)
	g.GET(prefix+"/failed", h.Failed)
	g.GET(prefix+"/task/:task", h.Task)
	g.GET(prefix+"/completed/:task", h.TaskCompleted)
}

func (h *Handler) Running(c echo.Context) error {
	tasks, err := task.GetTasks(c.Request().Context(), h.db, selectRunningTasks, itemLimit)
	if err != nil {
		return h.error(c, err)
	}
	return h.render(c, tmplTasksRunning, tasks)
}

func (h *Handler) Upcoming(c echo.Context) error {
	tasks, err := task.GetScheduledTasks(c.Request().Context(), h.db, time.Now().Add(time.Hour), itemLimit)
	if err != nil {
		return h.error(c, err)
	}
	return h.render(c, tmplTasksUpcoming, tasks)
}

func (h *Handler) Succeeded(c echo.Context) error {
	tasks, err := task.GetCompletedTasks(c.Request().Context(), h.db, selectCompletedTasks, 1, itemLimit)
	if err != nil {
		return h.error(c, err)
	}
	return h.render(c, tmplTasksCompleted, tasks)
}

func (h *Handler) Failed(c echo.Context) error {
	tasks, err := task.GetCompletedTasks(c.Request().Context(), h.db, selectCompletedTasks, 0, itemLimit)
	if err != nil {
		return h.error(c, err)
	}
	return h.render(c, tmplTasksCompleted, tasks)
}

func (h *Handler) Task(c echo.Context) error {
	id := c.Param("task")

	tasks, err := task.GetTasks(c.Request().Context(), h.db, selectTask, id)
	if err != nil {
		return h.error(c, err)
	}

	if len(tasks) > 0 {
		return h.render(c, tmplTask, tasks[0])
	}

	return h.TaskCompleted(c)
}

func (h *Handler) TaskCompleted(c echo.Context) error {
	id := c.Param("task")

	tasks, err := task.GetCompletedTasks(c.Request().Context(), h.db, selectCompletedTask, id)
	if err != nil {
		return h.error(c, err)
	}

	if len(tasks) > 0 {
		return h.render(c, tmplTaskCompleted, tasks[0])
	}

	return c.String(http.StatusNotFound, "Task not found")
}

func (h *Handler) error(c echo.Context, err error) error {
	log.Println(err)
	return c.String(http.StatusInternalServerError, err.Error())
}

func (h *Handler) render(c echo.Context, tmpl *template.Template, data any) error {
	return tmpl.ExecuteTemplate(c.Response().Writer, "layout.gohtml", TemplateData{
		Path:    c.Request().URL.Path,
		Prefix:  h.prefix,
		Content: data,
	})
}

// Helper function to ensure prefix starts with "/"
func hasLeadingSlash(s string) bool {
	return len(s) > 0 && s[0] == '/'
}
