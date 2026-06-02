package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HealthHandler struct {
	db *sqlx.DB
}

func NewHealthHandler(db *sqlx.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	dbStatus := "ok"
	if err := h.db.PingContext(ctx); err != nil {
		dbStatus = "unreachable"
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "ng",
			"db":        dbStatus,
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"db":        dbStatus,
		"timestamp": time.Now().UTC(),
	})
}
