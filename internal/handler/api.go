package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/website/component/page/index"
)

type ApiHandler struct {
	DB *sqlx.DB
}

func (h *ApiHandler) HandleStatsReposCount(c echo.Context) error {
	count, err := h.DB.SelectReposCount(h.DB)
	if err != nil {
		return err
	}

	return Render(c, index.StatCount())
}
