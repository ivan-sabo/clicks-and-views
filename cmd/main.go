package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.com/ivan-sabo/clicks-and-views/internal/click"
	"google.com/ivan-sabo/clicks-and-views/internal/view"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	gormDB, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		e.Logger.Fatal(err)
	}
	gormDB.AutoMigrate(&view.View{}, &click.Click{})

	clickRepository := click.NewSQLiteRepository(gormDB)
	viewRepository := view.NewSQLiteRepository(gormDB)

	clickHandler := click.NewHandler(clickRepository)
	viewHandler := view.NewHandler(viewRepository)

	e.GET("/clicks", clickHandler.Filter)
	e.POST("/clicks", clickHandler.Create)
	e.GET("/views", viewHandler.Filter)
	e.POST("/views", viewHandler.Create)

	e.Logger.Fatal(e.Start(":8080"))
}
