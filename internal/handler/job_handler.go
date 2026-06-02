package handler

import (
	"errors"
	"net/http"
	"strconv"

	"job-tracker-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	jobUsecase usecase.JobUsecase
}

func NewJobHandler(jobUsecase usecase.JobUsecase) *JobHandler {
	return &JobHandler{jobUsecase: jobUsecase}
}

func (h *JobHandler) GetAll(c *gin.Context) {
	userID := c.GetInt64("user_id")
	jobs, err := h.jobUsecase.GetAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

func (h *JobHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req usecase.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.jobUsecase.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"job": job})
}

func (h *JobHandler) GetByID(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	job, err := h.jobUsecase.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"job": job})
}

func (h *JobHandler) UpdateStatus(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req usecase.UpdateJobStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.jobUsecase.UpdateStatus(c.Request.Context(), id, userID, req)
	if err != nil {
		if errors.Is(err, usecase.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"job": job})
}
