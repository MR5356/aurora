package response

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRoute() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/new", func(context *gin.Context) {
		New(context, http.StatusOK, CodeSuccess, message(CodeSuccess), "success")
	})
	r.GET("/success", func(context *gin.Context) {
		Success(context, "success")
	})
	r.GET("/error", func(context *gin.Context) {
		Error(context, CodeNotFound)
	})
	r.GET("/error-with-msg", func(context *gin.Context) {
		ErrorWithMsg(context, CodeNotFound, MessageUnknown)
	})
	return r
}

func TestNew(t *testing.T) {
	r := setupRoute()

	req, err := http.NewRequest("GET", "/new", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := &Response{
		Code:    CodeSuccess,
		Message: message(CodeSuccess),
		Data:    "success",
	}

	expectedStr, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	if rr.Body.String() != string(expectedStr) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(expectedStr))
	}
}

func TestSuccess(t *testing.T) {
	r := setupRoute()

	req, err := http.NewRequest("GET", "/success", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := &Response{
		Code:    CodeSuccess,
		Message: message(CodeSuccess),
		Data:    "success",
	}

	expectedStr, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	if rr.Body.String() != string(expectedStr) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(expectedStr))
	}
}

func TestError(t *testing.T) {
	r := setupRoute()

	req, err := http.NewRequest("GET", "/error", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := &Response{
		Code:    CodeNotFound,
		Message: message(CodeNotFound),
		Data:    nil,
	}

	expectedStr, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	if rr.Body.String() != string(expectedStr) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(expectedStr))
	}
}

func TestErrorWithMsg(t *testing.T) {
	r := setupRoute()

	req, err := http.NewRequest("GET", "/error-with-msg", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := &Response{
		Code:    CodeNotFound,
		Message: MessageUnknown,
		Data:    nil,
	}

	expectedStr, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	if rr.Body.String() != string(expectedStr) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(expectedStr))
	}
}
