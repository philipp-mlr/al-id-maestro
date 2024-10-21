package handler

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/component/page/duplicates"
	"github.com/philipp-mlr/al-id-maestro/service"
)

type DuplicatesHandler struct {
	DB *sqlx.DB
}

func (h *DuplicatesHandler) HandleDuplicatesShow(c echo.Context) error {
	pageParam := c.QueryParam("page")

	if pageParam != "" {
		page, _ := strconv.Atoi(pageParam)
		if page <= 0 {
			page = 1
		}

		duplicatedObjects, err := service.SelectDuplicates(h.DB, uint64(page-1))
		if err != nil {
			return err
		}

		return Render(c, duplicates.TableItem(duplicatedObjects, uint64(page+1)))
	}

	return Render(c, duplicates.Show(InitActivePage(c)))
}
