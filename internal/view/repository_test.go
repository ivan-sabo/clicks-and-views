package view

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
	gormDB := setupDatabase(t)
	defer func() {
		teardownDatabase(t)
	}()

	sqliteRepo := SQLiteRepository{db: gormDB}

	view := View{
		URL: "test.url",
	}
	view, err := sqliteRepo.Create(context.Background(), view)
	assert.NoError(t, err)
	assert.Equal(t, view.ID, uint(1))
	assert.Equal(t, view.URL, "test.url")
}

func TestFilter(t *testing.T) {
	gormDB := setupDatabase(t)
	defer func() {
		teardownDatabase(t)
	}()

	// setup test data
	time1, _ := time.Parse(time.DateOnly, "2024-01-02")
	time2, _ := time.Parse(time.DateOnly, "2024-04-02")
	clicks := ViewCollection{
		{
			URL:       "test.url1",
			CreatedAt: time1,
		},
		{
			URL:       "test.url2",
			CreatedAt: time2,
		},
	}
	gormDB.Create(clicks)
	assert.NoError(t, gormDB.Error)

	tests := []struct {
		testName       string
		param          Filter
		expectedResult ViewCollection
		err            error
	}{
		{
			testName: "filter by URL - success",
			param:    Filter{URL: "test.url1"},
			expectedResult: ViewCollection{
				{
					ID:        1,
					URL:       "test.url1",
					CreatedAt: time1,
				},
			},
			err: nil,
		},
		{
			testName: "filter by created_at - after",
			param:    Filter{After: time2.Add(-1 * time.Hour)},
			expectedResult: ViewCollection{
				{
					ID:        2,
					URL:       "test.url2",
					CreatedAt: time2,
				},
			},
			err: nil,
		},
		{
			testName: "filter by created_at - before",
			param:    Filter{Before: time1.Add(1 * time.Hour)},
			expectedResult: ViewCollection{
				{
					ID:        1,
					URL:       "test.url1",
					CreatedAt: time1,
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			sqliteRepository := SQLiteRepository{db: gormDB}

			clicksResult, err := sqliteRepository.Filter(context.Background(), test.param)

			assert.Equal(t, test.expectedResult, clicksResult)
			if test.err == nil {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}

func setupDatabase(t *testing.T) *gorm.DB {
	t.Helper()

	gormDB, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	assert.NoError(t, err)
	gormDB.AutoMigrate(&View{})

	return gormDB
}

func teardownDatabase(t *testing.T) {
	t.Helper()

	os.Remove("gorm.db")
}
