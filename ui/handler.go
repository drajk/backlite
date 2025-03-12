package ui

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/drajk/backlite/internal/task"
	"github.com/labstack/echo/v4"
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

func NewHandler(db *sql.DB) *echo.Echo {
	h := &Handler{db: db}
	e := echo.New()

	e.GET("/", h.Running)
	e.GET("/upcoming", h.Upcoming)
	e.GET("/succeeded", h.Succeeded)
	e.GET("/failed", h.Failed)
	e.GET("/task/:task", h.Task)               // Echo supports dynamic params using `:param`
	e.GET("/completed/:task", h.TaskCompleted) // Same for completed tasks

	return e
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
	id := c.Param("task") // Retrieve path parameter
	
	tasks, err := task.GetTasks(c.Request().Context(), h.db, selectTask, id)
	if err != nil {
		return h.error(c, err)
	}
	
	if len(tasks) > 0 {
		return h.render(c, tmplTask, tasks[0])
	}
	
	// If no task found, try fetching it as a completed task
	return h.TaskCompleted(c)
}

func (h *Handler) TaskCompleted(c echo.Context) error {
	id := c.Param("task") // Retrieve path parameter
	
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
		Content: data,
	})
}
