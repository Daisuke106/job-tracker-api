package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type InterviewHandler struct{}

func NewInterviewHandler() *InterviewHandler {
	return &InterviewHandler{}
}

func (h *InterviewHandler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "not implemented yet"})
}

func (h *InterviewHandler) Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "not implemented yet"})
}

func (h *InterviewHandler) GetByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "not implemented yet"})
}

func (h *InterviewHandler) CreateNote(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "not implemented yet"})
}
