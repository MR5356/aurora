package ginutil

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRoute() *gin.Engine {
	r := gin.Default()
	r.GET("/page", func(c *gin.Context) {
		page, size := GetPageParams(c)
		c.String(200, fmt.Sprintf("page: %d, size: %d", page, size))
	})

	r.GET("/token", func(c *gin.Context) {
		token := GetToken(c)
		c.String(200, token)
	})
	return r
}

func TestGetPageParams(t *testing.T) {
	r := setupRoute()

	tests := []struct {
		name string
		page int
		size int
	}{
		{
			name: "test",
			page: 10,
			size: 1,
		},
		{
			name: "test-nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req, err := http.NewRequest("GET", fmt.Sprintf("/page?page=%d&size=%d", tt.page, tt.size), nil)
			if err != nil {
				t.Fatal(err)
			}
			r.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Errorf("GetId() = %v, want %v", rr.Code, http.StatusOK)
			}
			if rr.Body.String() != fmt.Sprintf("page: %d, size: %d", tt.page, tt.size) {
				t.Errorf("GetId() = %v, want %v", rr.Body.String(), fmt.Sprintf("page: %d, size: %d", tt.page, tt.size))
			}
		})
	}
}

func TestGetToken(t *testing.T) {
	r := setupRoute()

	tests := []struct {
		name   string
		token  string
		header string
		cookie string
	}{
		{
			name:   "test-cookie",
			token:  "test",
			cookie: "test",
		},
		{
			name:   "test-header",
			token:  "test",
			header: "test",
		},
		{
			name:   "test-nil",
			token:  "",
			cookie: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req, err := http.NewRequest("GET", "/token", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{Name: "token", Value: tt.cookie})
			}
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			r.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Errorf("GetToken() = %v, want %v", rr.Code, http.StatusOK)
			}
			if rr.Body.String() != tt.token {
				t.Errorf("GetToken() = %v, want %v", rr.Body.String(), tt.token)
			}
		})
	}
}
