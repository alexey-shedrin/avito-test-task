package handler

import (
	"log"
	"strings"

	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/request"
	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/response"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/gin-gonic/gin"
)

const (
	InvalidRole  = "invalid role"
	InvalidEmail = "invalid email"
)

type UserService interface {
	DummyLogin(request *request.DummyLogin) (*response.DummyLogin, error)
	Register(request *request.Register) (*entity.User, error)
	Login(request *request.Login) (*response.Login, error)
}

func (h *Handler) PostDummyLogin(c *gin.Context) {
	log.SetPrefix("handler.PostDummyLogin")
	var req request.DummyLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "Field validation") {
			c.JSON(400, gin.H{"error": InvalidRole})

			return
		}

		log.Printf("error: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})

		return
	}

	resp, err := h.userService.DummyLogin(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}

func (h *Handler) PostRegister(c *gin.Context) {
	log.SetPrefix("handler.PostRegister")
	var req request.Register
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "Field validation") {
			c.JSON(400, gin.H{"error": InvalidRole})

			return
		}

		log.Printf("error: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})

		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user.ToResponse())
}

func (h *Handler) PostLogin(c *gin.Context) {
	log.SetPrefix("handler.PostLogin")
	var req request.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "Field validation") {
			c.JSON(400, gin.H{"error": InvalidEmail})

			return
		}

		log.Printf("error: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})

		return
	}

	resp, err := h.userService.Login(&req)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}
