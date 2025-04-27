package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amirzayi/clean_architect/api/http/handler/v2/dto"
	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestRegisterV2(t *testing.T) {
	rec := httptest.NewRecorder()

	b, err := json.Marshal(dto.RegisterRequest{Name: "amir", Email: "mirzayi994@gmail.com", PhoneNumber: "+989101234567", Password: "password"})
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/v2/auth/register", bytes.NewReader(b))

	mux.ServeHTTP(rec, req)
	require.Equal(t, "", rec.Body.String())
	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestLoginV2(t *testing.T) {
	rec := httptest.NewRecorder()

	b, err := json.Marshal(domain.Auth{Email: "mirzayi994@gmail.com", PhoneNumber: "+989101234567", Password: "password"})
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/v2/auth/login", bytes.NewReader(b))

	mux.ServeHTTP(rec, req)
	require.Contains(t, rec.Body.String(), "token")
	require.Equal(t, http.StatusOK, rec.Code)
}
