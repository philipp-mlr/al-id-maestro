package handler

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/internal/chart"
	"github.com/philipp-mlr/al-id-maestro/website/component/page/index"
)

type IndexHandler struct {
	DB *sqlx.DB
}

func (h IndexHandler) HandleIndexShow(c echo.Context) error {
	return Render(c, index.Show(InitActivePage(c)))
}

func (h *IndexHandler) HandleChartShow(c echo.Context) error {
	data, err := chart.GetChartData(h.DB, 14)
	if err != nil {
		log.Println("Error getting chart data: ", err)
	}

	log.Println("Chart data: ", data)

	return Render(c, index.Chart(data))
}
