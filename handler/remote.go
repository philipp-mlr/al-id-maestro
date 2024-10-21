package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/component/page/remote"
)

type RemoteHandler struct {
	DB *sqlx.DB
}

func (h RemoteHandler) HandleRemoteShow(c echo.Context) error {
	return Render(c, remote.Show(InitActivePage(c)))
}
