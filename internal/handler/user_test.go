package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexey-shedrin/avito-test-task/internal/handler"
	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/request"
	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/response"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/alexey-shedrin/avito-test-task/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupUserRouter(h *handler.Handler, setup func(*gin.Engine)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	setup(r)
	return r
}

func TestPostDummyLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUserService(ctrl)
	h := handler.New(mockUser, nil, nil)

	input := request.DummyLogin{Role: "moderator"}
	expected := &response.DummyLogin{Token: "token"}

	mockUser.EXPECT().DummyLogin(&input).Return(expected, nil)

	body, _ := json.Marshal(input)

	r := setupUserRouter(h, func(r *gin.Engine) {
		r.POST("/dummy-login", h.PostDummyLogin)
	})

	req := httptest.NewRequest(http.MethodPost, "/dummy-login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestPostDummyLogin_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := handler.New(mocks.NewMockUserService(ctrl), nil, nil)

	r := setupUserRouter(h, func(r *gin.Engine) {
		r.POST("/dummy-login", h.PostDummyLogin)
	})

	req := httptest.NewRequest(http.MethodPost, "/dummy-login", bytes.NewReader([]byte("{bad json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUserService(ctrl)
	h := handler.New(mockUser, nil, nil)

	input := request.Register{Email: "test@example.com", Password: "pass", Role: "employee"}
	expected := &entity.User{Email: input.Email, Role: input.Role}

	mockUser.EXPECT().Register(&input).Return(expected, nil)

	body, _ := json.Marshal(input)

	r := setupUserRouter(h, func(r *gin.Engine) {
		r.POST("/register", h.PostRegister)
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
}

func TestPostRegister_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := handler.New(mocks.NewMockUserService(ctrl), nil, nil)

	r := setupUserRouter(h, func(r *gin.Engine) {
		r.POST("/register", h.PostRegister)
	})

	badJSON := []byte(`{"email": "user", "password": "123"}`) // role пустой

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(badJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUserService(ctrl)
	h := handler.New(mockUser, nil, nil)

	input := request.Login{Email: "user@mail.com", Password: "secret"}
	expected := &response.Login{Token: "jwt"}

	mockUser.EXPECT().Login(&input).Return(expected, nil)

	body, _ := json.Marshal(input)

	r := setupUserRouter(h, func(r *gin.Engine) {
		r.POST("/login", h.PostLogin)
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestPostLogin_InvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := handler.New(mocks.NewMockUserService(ctrl), nil, nil)

	r := setupUserRouter(h, func(r *gin.Engine) {
		r.POST("/login", h.PostLogin)
	})

	badJSON := []byte(`{"password": "123"}`) // без email

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(badJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostLogin_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUserService(ctrl)
	h := handler.New(mockUser, nil, nil)

	input := request.Login{Email: "wrong@mail.com", Password: "wrong"}
	errMsg := "unauthorized"

	mockUser.EXPECT().Login(&input).Return(nil, errors.New(errMsg))

	body, _ := json.Marshal(input)

	r := setupUserRouter(h, func(r *gin.Engine) {
		r.POST("/login", h.PostLogin)
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}
