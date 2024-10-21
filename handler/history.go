package handler

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/component/page/history"
	"github.com/philipp-mlr/al-id-maestro/service"
)

const PageSize = 10

type HistoryHandler struct {
	DB *sqlx.DB
}

func (h *HistoryHandler) HandleHistoryShow(c echo.Context) error {
	pageParam := c.QueryParam("page")

	if pageParam != "" {
		page, _ := strconv.Atoi(pageParam)
		if page <= 0 {
			page = 1
		}

		if err := service.UpdateClaimed(h.DB); err != nil {
			return err
		}

		claims, err := service.SelectClaimedObjects(h.DB, uint64(page-1))
		if err != nil {
			return err
		}

		return Render(c, history.TableItem(claims, uint64(page+1)))
	}

	return Render(c, history.Show(InitActivePage(c)))
}

func (h *HistoryHandler) HandlePostQuery(c echo.Context) error {
	fmt.Println("HandlePostQuery")

	return nil
}
