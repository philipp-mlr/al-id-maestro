package handler

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/internal/claim"
	"github.com/philipp-mlr/al-id-maestro/internal/database"
	"github.com/philipp-mlr/al-id-maestro/website/component/page/history"
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

		if err := claim.UpdateClaimed(h.DB); err != nil {
			return err
		}

		claims, err := database.SelectClaimedObjects(h.DB, uint64(page-1))
		if err != nil {
			return err
		}

		return Render(c, history.TableItem(claims, uint64(page+1)))
	}

	return Render(c, history.Show(InitActivePage(c)))
}

func (h *HistoryHandler) HandlePostQuery(c echo.Context) error {
	fmt.Println("HandlePostQuery")

	return nil //Render(c, history.TableItem(InitActivePage(c)))
}
