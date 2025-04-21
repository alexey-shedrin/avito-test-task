package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexey-shedrin/avito-test-task/internal/utils/token"

	"github.com/alexey-shedrin/avito-test-task/internal/handler"
	"github.com/alexey-shedrin/avito-test-task/internal/middleware"
	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/request"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/alexey-shedrin/avito-test-task/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupRouter(h *handler.Handler, setup func(*gin.Engine)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	setup(router)
	return router
}

func TestPostReceptions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionService := mocks.NewMockReceptionService(ctrl)

	h := handler.New(nil, nil, mockReceptionService)

	pvzID := uuid.New()
	input := request.Reception{PvzId: pvzID}
	expected := &entity.Reception{PvzId: pvzID}

	mockReceptionService.EXPECT().CreateReception(gomock.Any()).Return(expected, nil)

	body, _ := json.Marshal(input)
	router := setupRouter(h, func(r *gin.Engine) {
		r.POST("/receptions", func(c *gin.Context) {
			middleware.Auth(entity.EmployeeRole)(c)
			h.PostReceptions(c)
		})
	})

	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	jwt, _ := token.GenerateJWT(entity.EmployeeRole)
	req.Header.Set("Authorization", jwt)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
}

func TestPostProducts_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := handler.New(nil, nil, mocks.NewMockReceptionService(ctrl))

	router := setupRouter(h, func(r *gin.Engine) {
		r.POST("/products", func(c *gin.Context) {
			middleware.Auth(entity.EmployeeRole)(c)
			h.PostProducts(c)
		})
	})

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	jwt, _ := token.GenerateJWT(entity.EmployeeRole)
	req.Header.Set("Authorization", jwt)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostPvzPvzIdDeleteLastProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockReceptionService(ctrl)
	h := handler.New(nil, nil, mockService)

	pvzID := uuid.New()
	mockService.EXPECT().DeleteLastProduct(pvzID).Return(nil)

	router := setupRouter(h, func(r *gin.Engine) {
		r.POST("/pvz/:id/delete-last", func(c *gin.Context) {
			middleware.Auth(entity.EmployeeRole)(c)
			id, _ := uuid.Parse(c.Param("id"))
			h.PostPvzPvzIdDeleteLastProduct(c, id)
		})
	})

	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete-last", nil)
	jwt, _ := token.GenerateJWT(entity.EmployeeRole)
	req.Header.Set("Authorization", jwt)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestPostPvzPvzIdCloseLastReception_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockReceptionService(ctrl)
	h := handler.New(nil, nil, mockService)

	pvzID := uuid.New()
	mockService.EXPECT().CloseLastReception(pvzID).Return(nil, errors.New("some error"))

	router := setupRouter(h, func(r *gin.Engine) {
		r.POST("/pvz/:id/close", func(c *gin.Context) {
			middleware.Auth(entity.EmployeeRole)(c)
			id, _ := uuid.Parse(c.Param("id"))
			h.PostPvzPvzIdCloseLastReception(c, id)
		})
	})

	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/close", nil)
	jwt, _ := token.GenerateJWT(entity.EmployeeRole)
	req.Header.Set("Authorization", jwt)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
