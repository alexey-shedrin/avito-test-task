package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	suite.Suite
	ctx            context.Context
	cancel         context.CancelFunc
	client         *resty.Client
	moderatorToken string
	employeeToken  string
}

func (s *IntegrationSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 30*time.Second)

	s.client = resty.New().SetBaseURL("http://localhost:8080")
}

func (s *IntegrationSuite) TearDownTest() {
	s.cancel()
}

func (s *IntegrationSuite) dummyLoginHelper(role string) string {
	body := map[string]interface{}{
		"role": role,
	}
	var tokenResp struct {
		Token string `json:"token"`
	}
	r, err := s.client.R().
		SetBody(body).
		Post("/dummyLogin")
	s.T().Logf("DummyLogin (%s) response: %s", role, r.Body())
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, r.StatusCode())

	err = json.Unmarshal(r.Body(), &tokenResp)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), tokenResp.Token)

	return tokenResp.Token
}

func (s *IntegrationSuite) TestPVZPipeline() {
	s.moderatorToken = s.dummyLoginHelper("moderator")

	pvzID := uuid.New()
	registrationDate := time.Now().Format(time.RFC3339)
	pvzBody := map[string]interface{}{
		"id":               pvzID,
		"registrationDate": registrationDate,
		"city":             "Казань",
	}

	var pvzResp struct {
		ID string `json:"id"`
	}
	r, err := s.client.R().
		SetHeader("Authorization", s.moderatorToken).
		SetBody(pvzBody).
		Post("/pvz")
	s.T().Logf("Create PVZ response: %v", string(r.Body()))
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, r.StatusCode())

	err = json.Unmarshal(r.Body(), &pvzResp)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), pvzResp.ID)

	s.employeeToken = s.dummyLoginHelper("employee")

	receptionBody := map[string]interface{}{
		"pvzId": pvzResp.ID,
	}

	var receptionResp struct {
		ID string `json:"id"`
	}
	r, err = s.client.R().
		SetHeader("Authorization", s.employeeToken).
		SetBody(receptionBody).
		Post("/receptions")
	s.T().Logf("Create Reception response: %v", string(r.Body()))
	require.NoError(s.T(), err)
	require.Equal(s.T(), r.StatusCode(), http.StatusCreated)

	err = json.Unmarshal(r.Body(), &receptionResp)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), receptionResp.ID)

	for i := 1; i <= 50; i++ {
		productBody := map[string]interface{}{
			"type":  "одежда",
			"pvzId": pvzResp.ID,
		}

		r, err = s.client.R().
			SetHeader("Authorization", s.employeeToken).
			SetBody(productBody).
			Post("/products")
		s.T().Logf("Add Product %d response: %s", i, r.Body())
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.StatusCreated, r.StatusCode())
	}

	closeReceptionURL := fmt.Sprintf("/pvz/%s/close_last_reception", pvzResp.ID)
	r, err = s.client.R().
		SetHeader("Authorization", s.employeeToken).
		Post(closeReceptionURL)
	s.T().Logf("Close Reception response: %v", string(r.Body()))
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, r.StatusCode())
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
