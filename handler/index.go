package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/component/page/index"
)

type IndexHandler struct {
	DB *sqlx.DB
}

func (h IndexHandler) HandleIndexShow(c echo.Context) error {
	return Render(c, index.Show(InitActivePage(c)))
}
