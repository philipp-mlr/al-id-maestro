package handler

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/internal/claim"
	"github.com/philipp-mlr/al-id-maestro/internal/database"
	"github.com/philipp-mlr/al-id-maestro/internal/model"
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

		err = orderClaimsByDate(&claims)
		if err != nil {
			log.Println("Error: ", err)
		}

		return Render(c, history.TableItem(claims, uint64(page+1)))
	}

	return Render(c, history.Show(InitActivePage(c)))
}

func orderClaimsByDate(claims *[]model.ClaimedObject) error {
	for i := 0; i < len(*claims); i++ {
		for j := i + 1; j < len(*claims); j++ {
			t1, err := time.Parse(time.RFC1123, (*claims)[i].CreatedAt)
			if err != nil {
				return err
			}

			t2, err := time.Parse(time.RFC1123, (*claims)[j].CreatedAt)
			if err != nil {
				return err
			}

			if t1.Before(t2) {
				(*claims)[i], (*claims)[j] = (*claims)[j], (*claims)[i]
			}
		}
	}
	return nil
}

func (h *HistoryHandler) HandlePostQuery(c echo.Context) error {
	fmt.Println("HandlePostQuery")

	return nil //Render(c, history.TableItem(InitActivePage(c)))
}
