package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amirzayi/clean_architect/api/http/handler/v2/dto"
	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestListUserV2(t *testing.T) {
	for _, tc := range []struct {
		name         string
		headers      map[string]string
		expectedCode int
		expectedBody string
	}{
		{"no token", nil, http.StatusUnauthorized, request.ErrNoTokenInRequest.Error()},
		{"invalid token", map[string]string{"Authorization": "Invalid Token"}, http.StatusUnauthorized, request.ErrNoTokenInRequest.Error()},
		{"bad token", map[string]string{"Authorization": userToken}, http.StatusUnauthorized, request.ErrNoTokenInRequest.Error()},
		{"invalid role", map[string]string{"Authorization": "Bearer " + userToken}, http.StatusForbidden, http.StatusText(http.StatusForbidden)},
		{"valid role", map[string]string{"Authorization": "Bearer " + adminToken}, http.StatusOK, "data"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/v2/users", http.NoBody)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			mux.ServeHTTP(rec, req)
			require.Equal(t, tc.expectedCode, rec.Code)
			require.Contains(t, rec.Body.String(), tc.expectedBody)
		})
	}
}

func TestCreateUserV2(t *testing.T) {
	for _, tc := range []struct {
		name         string
		body         *dto.CreateUserRequest
		headers      map[string]string
		expectedCode int
	}{
		{"no token", nil, nil, http.StatusUnauthorized},
		{"invalid token", nil, map[string]string{"Authorization": "Invalid Token"}, http.StatusUnauthorized},
		{"bad token", nil, map[string]string{"Authorization": userToken}, http.StatusUnauthorized},
		{"invalid role", nil, map[string]string{"Authorization": "Bearer " + userToken}, http.StatusForbidden},
		{"empty input", nil, map[string]string{"Authorization": "Bearer " + adminToken}, http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			var b []byte
			var err error
			if tc.body != nil {
				b, err = json.Marshal(tc.body)
				require.NoError(t, err)
			}
			req, _ := http.NewRequest(http.MethodPost, "/v2/users", bytes.NewReader(b))
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			mux.ServeHTTP(rec, req)
			require.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	t.Run("valid", func(t *testing.T) {
		testCreateUserV2(t)
	})
}

func TestGetUserV2(t *testing.T) {
	for _, tc := range []struct {
		name         string
		id           string
		headers      map[string]string
		expectedCode int
	}{
		{"no token", "123", nil, http.StatusUnauthorized},
		{"invalid token", "123", map[string]string{"Authorization": "Invalid Token"}, http.StatusUnauthorized},
		{"bad token", "123", map[string]string{"Authorization": userToken}, http.StatusUnauthorized},
		{"invalid role", "123", map[string]string{"Authorization": "Bearer " + userToken}, http.StatusForbidden},
		{"bad id parameter", "123", map[string]string{"Authorization": "Bearer " + adminToken}, http.StatusBadRequest},
		{"not found", uuid.New().String(), map[string]string{"Authorization": "Bearer " + adminToken}, http.StatusNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/v2/users/"+tc.id, http.NoBody)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			mux.ServeHTTP(rec, req)
			require.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	t.Run("valid", func(t *testing.T) {
		user := testCreateUserV2(t)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/v2/users/"+user.ID.String(), http.NoBody)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		var newUser domain.User
		err := json.Unmarshal(rec.Body.Bytes(), &newUser)
		// some issue with sqlite save datetime
		newUser.CreatedAt = user.CreatedAt
		require.NoError(t, err)
		require.Equal(t, user, newUser)
	})
}

func TestDeleteUserV2(t *testing.T) {
	for _, tc := range []struct {
		name         string
		id           string
		headers      map[string]string
		expectedCode int
	}{
		{"no token", "123", nil, http.StatusUnauthorized},
		{"invalid token", "123", map[string]string{"Authorization": "Invalid Token"}, http.StatusUnauthorized},
		{"bad token", "123", map[string]string{"Authorization": userToken}, http.StatusUnauthorized},
		{"invalid role", "123", map[string]string{"Authorization": "Bearer " + userToken}, http.StatusForbidden},
		{"bad id parameter", "123", map[string]string{"Authorization": "Bearer " + adminToken}, http.StatusBadRequest},
		{"not found", uuid.New().String(), map[string]string{"Authorization": "Bearer " + adminToken}, http.StatusInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/v2/users/"+tc.id, http.NoBody)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			mux.ServeHTTP(rec, req)
			require.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	t.Run("valid", func(t *testing.T) {
		user := testCreateUserV2(t)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/v2/users/"+user.ID.String(), http.NoBody)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNoContent, rec.Code)
	})
}

func testCreateUserV2(t *testing.T) domain.User {
	rec := httptest.NewRecorder()

	body := dto.CreateUserRequest{"amir", "09101234567", "mirzayi994@gmail.com", "password", string(domain.UserRoleNormal)}
	b, err := json.Marshal(body)
	require.NoError(t, err)

	req, _ := http.NewRequest(http.MethodPost, "/v2/users", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+adminToken)

	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var user domain.User
	err = json.Unmarshal(rec.Body.Bytes(), &user)
	require.NoError(t, err)
	return user
}
