package handler

import (
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/component/page/claim"
	"github.com/philipp-mlr/al-id-maestro/model"
	"github.com/philipp-mlr/al-id-maestro/service"
)

type ClaimHandler struct {
	DB          *sqlx.DB
	AllowedList *model.AllowedList
}

type ClaimRequest struct {
	ObjectType model.ObjectType `json:"query" validate:"required"`
}

func (h ClaimHandler) HandlePageShow(c echo.Context) error {
	return Render(c, claim.Show(InitActivePage(c)))
}

func (h *ClaimHandler) HandleRequestID(c echo.Context) error {
	claimRequest := ClaimRequest{}

	if err := c.Bind(&claimRequest); err != nil {
		return err
	}

	if err := c.Validate(claimRequest); err != nil {
		return err
	}

	if model.MapObjectType(string(claimRequest.ObjectType)) == "" {
		return echo.ErrBadRequest
	}

	if err := service.UpdateClaimed(h.DB); err != nil {
		return err
	}

	claimed, err := service.ClaimObjectID(h.DB, h.AllowedList, claimRequest.ObjectType)
	if err != nil {
		return Render(c, claim.ClaimedID(err.Error()))
	}

	idString := strconv.Itoa(int(claimed.ID))

	return Render(c, claim.ClaimedID(idString))
}

type ObjectTypeQuery struct {
	ObjectType model.ObjectType `json:"query"`
}

func (h ClaimHandler) HandleObjectTypeQuery(c echo.Context) error {
	objectTypeQuery := ObjectTypeQuery{}
	if err := c.Bind(&objectTypeQuery); err != nil {
		return err
	}

	objectTypes := model.GetObjectTypes()

	if objectTypeQuery.ObjectType == "" {
		return Render(c, claim.QueryResults(objectTypes))
	}

	results := []model.ObjectType{}
	for _, objectType := range objectTypes {
		if strings.Contains(strings.ToLower(string(objectType)), strings.ToLower(string(objectTypeQuery.ObjectType))) {
			results = append(results, objectType)
		}
	}

	return Render(c, claim.QueryResults(results))
}
