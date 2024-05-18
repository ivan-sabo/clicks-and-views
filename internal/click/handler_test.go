package click

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

type ClickRepositoryMock struct {
	mock.Mock
}

func (m *ClickRepositoryMock) Create(ctx context.Context, click Click) (Click, error) {
	args := m.Called(ctx, click)
	return args.Get(0).(Click), args.Error(1)
}

func (m *ClickRepositoryMock) Filter(ctx context.Context, filter Filter) (ClickCollection, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(ClickCollection), args.Error(1)
}

func TestHandlerCreate(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/clicks", strings.NewReader(`{"url":"test.url1"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	timeNow := time.Now()

	clickRepository := &ClickRepositoryMock{}
	clickRepository.
		On("Create", c.Request().Context(), Click{URL: "test.url1"}).
		Return(
			Click{
				ID:        1,
				URL:       "test.url1",
				CreatedAt: timeNow,
			},
			nil,
		).Once()

	h := &Handler{clickRepository: clickRepository}

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
	req := httptest.NewRequest(http.MethodGet, "/clicks?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	timeNow := time.Now()

	clickRepository := &ClickRepositoryMock{}
	clickRepository.
		On("Filter", c.Request().Context(), Filter{URL: "test.url1"}).
		Return(
			ClickCollection{
				{
					ID:        1,
					URL:       "test.url1",
					CreatedAt: timeNow,
				},
			},
			nil,
		).Once()

	h := &Handler{clickRepository: clickRepository}

	if assert.NoError(t, h.Filter(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedJSON := fmt.Sprintf(`[{"id":1,"url":"test.url1","createdAt":"%s"}]`+"\n", timeNow.Format(time.DateTime))
		assert.Equal(t, expectedJSON, rec.Body.String())
	}
}
