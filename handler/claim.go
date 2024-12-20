package handler

import (
	"log"
	"net/http"
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

func (h *ClaimHandler) HandleNewObjectClaim(c echo.Context) error {
	claimRequest := ClaimRequest{}

	if err := c.Bind(&claimRequest); err != nil {
		log.Println(err)
		return err
	}

	if err := c.Validate(claimRequest); err != nil {
		log.Println(err)
		return err
	}

	objectType := model.MapObjectType(string(claimRequest.ObjectType))
	if objectType == model.Unknown {
		return echo.ErrBadRequest
	}

	if err := service.UpdateClaimed(h.DB); err != nil {
		log.Println(err)
		return err
	}

	claimed, err := service.ClaimObjectID(h.DB, h.AllowedList, objectType)
	if err != nil {
		log.Println(err)
		return Render(c, claim.ClaimedID(err.Error()))
	}

	return Render(c, claim.ClaimedID(strconv.Itoa(int(claimed.ID))))
}

type ObjectClaimAPIResponse struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type ObjectClaimAPIError struct {
	Message string `json:"message"`
}

func (h *ClaimHandler) HandleNewObjectClaimAPI(c echo.Context) error {
	claimRequest := ClaimRequest{}

	if err := c.Bind(&claimRequest); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, ObjectClaimAPIError{
			Message: "Invalid request body",
		})
	}

	if err := c.Validate(claimRequest); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, ObjectClaimAPIError{
			Message: "Invalid request body",
		})
	}

	objectType := model.MapObjectType(string(claimRequest.ObjectType))
	if objectType == model.Unknown {
		return c.JSON(http.StatusInternalServerError, ObjectClaimAPIError{
			Message: "Invalid object type",
		})
	}

	if err := service.UpdateClaimed(h.DB); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, ObjectClaimAPIError{
			Message: "Error updating claimed objects",
		})
	}

	claimed, err := service.ClaimObjectID(h.DB, h.AllowedList, objectType)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, ObjectClaimAPIError{
			Message: "Could not claim object",
		})
	}

	return c.JSON(http.StatusOK, ObjectClaimAPIResponse{
		ID:   strconv.Itoa(int(claimed.ID)),
		Type: string(claimed.ObjectType),
	})
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
