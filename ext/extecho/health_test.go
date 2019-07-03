package extecho

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

type testResource int

func (t *testResource) HealthCheck() error {
	*t++
	if *t%2 == 0 {
		return nil
	}
	return errors.New("BAD")
}

func TestHealth(t *testing.T) {
	var i testResource
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/_health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	m := NewHealth(&i)
	h := m(handler)
	h(c)

	assert.Equal(t, "BAD", rec.Body.String())
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	handler = func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	m = NewHealth(&i)
	h = m(handler)
	h(c)

	assert.Equal(t, "test", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/_health", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	handler = func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	m = NewHealth(&i)
	h = m(handler)
	h(c)

	assert.Equal(t, "OK", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
}
