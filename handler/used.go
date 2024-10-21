package handler

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/component/page/used"
	"github.com/philipp-mlr/al-id-maestro/service"
)

type UsedHandler struct {
	DB *sqlx.DB
}

func (h UsedHandler) HandleUsedShow(c echo.Context) error {
	pageParam := c.QueryParam("page")

	if pageParam != "" {
		page, _ := strconv.Atoi(pageParam)
		if page <= 0 {
			page = 1
		}

		duplicatedObjects, err := service.SelectFound(h.DB, uint64(page-1))
		if err != nil {
			return err
		}

		return Render(c, used.TableItem(duplicatedObjects, uint64(page+1)))
	}

	return Render(c, used.Show(InitActivePage(c)))
}
