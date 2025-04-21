package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexey-shedrin/avito-test-task/internal/handler"
	"github.com/alexey-shedrin/avito-test-task/internal/middleware"
	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/request"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/alexey-shedrin/avito-test-task/internal/service/mocks"
	"github.com/alexey-shedrin/avito-test-task/internal/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupPvzRouter(h *handler.Handler, setup func(*gin.Engine)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	setup(r)
	return r
}

func TestPostPvz_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPvzService(ctrl)
	h := handler.New(nil, mockService, nil)

	input := request.Pvz{City: "Москва"}
	expected := &entity.Pvz{City: "Москва"}

	mockService.EXPECT().CreatePvz(gomock.Any()).Return(expected, nil)

	body, _ := json.Marshal(input)

	r := setupPvzRouter(h, func(r *gin.Engine) {
		r.POST("/pvz", func(c *gin.Context) {
			middleware.Auth(entity.ModeratorRole)(c)
			h.PostPvz(c)
		})
	})

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	jwt, _ := token.GenerateJWT(entity.ModeratorRole)
	req.Header.Set("Authorization", jwt)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
}

func TestPostPvz_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := handler.New(nil, mocks.NewMockPvzService(ctrl), nil)

	r := setupPvzRouter(h, func(r *gin.Engine) {
		r.POST("/pvz", func(c *gin.Context) {
			middleware.Auth(entity.ModeratorRole)(c)
			h.PostPvz(c)
		})
	})

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader([]byte("{bad json")))
	req.Header.Set("Content-Type", "application/json")
	jwt, _ := token.GenerateJWT(entity.ModeratorRole)
	req.Header.Set("Authorization", jwt)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
