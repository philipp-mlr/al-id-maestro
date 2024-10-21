package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/component/page/about"
)

type AboutHandler struct {
	DB *sqlx.DB
}

func (h AboutHandler) HandleAboutShow(c echo.Context) error {
	return Render(c, about.Show(InitActivePage(c)))
}
