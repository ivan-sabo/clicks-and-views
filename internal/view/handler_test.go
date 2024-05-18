package view

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ViewRepositoryMock struct {
	mock.Mock
}

func (m *ViewRepositoryMock) Create(ctx context.Context, view View) (View, error) {
	args := m.Called(ctx, view)
	return args.Get(0).(View), args.Error(1)
}

func (m *ViewRepositoryMock) Filter(ctx context.Context, filter Filter) (ViewCollection, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(ViewCollection), args.Error(1)
}

func TestHandlerCreate(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/views", strings.NewReader(`{"url":"test.url1"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	timeNow := time.Now()

	viewRepository := &ViewRepositoryMock{}
	viewRepository.
		On("Create", c.Request().Context(), View{URL: "test.url1"}).
		Return(
			View{
				ID:        1,
				URL:       "test.url1",
				CreatedAt: timeNow,
			},
			nil,
		).Once()

	h := &Handler{viewRepository: viewRepository}

	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		expectedJSON := fmt.Sprintf(`{"id":1,"url":"test.url1","createdAt":"%s"}`+"\n", timeNow.Format(time.DateTime))
		assert.Equal(t, expectedJSON, rec.Body.String())
	}
}

func TestHandlerFilter(t *testing.T) {
	e := echo.New()
	q := make(url.Values)
	q.Set("url", "test.url1")
	req := httptest.NewRequest(http.MethodGet, "/views?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	timeNow := time.Now()

	viewRepository := &ViewRepositoryMock{}
	viewRepository.
		On("Filter", c.Request().Context(), Filter{URL: "test.url1"}).
		Return(
			ViewCollection{
				{
					ID:        1,
					URL:       "test.url1",
					CreatedAt: timeNow,
				},
			},
			nil,
		).Once()

	h := &Handler{viewRepository: viewRepository}

	if assert.NoError(t, h.Filter(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedJSON := fmt.Sprintf(`[{"id":1,"url":"test.url1","createdAt":"%s"}]`+"\n", timeNow.Format(time.DateTime))
		assert.Equal(t, expectedJSON, rec.Body.String())
	}
}
