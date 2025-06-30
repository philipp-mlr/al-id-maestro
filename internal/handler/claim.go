package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/internal/claim"
	"github.com/philipp-mlr/al-id-maestro/internal/model"
	"github.com/philipp-mlr/al-id-maestro/internal/objectType"
	claimPage "github.com/philipp-mlr/al-id-maestro/website/component/page/claim"
)

type ClaimHandler struct {
	DB          *sqlx.DB
	AllowedList *model.LicensedObjectList
}

type ClaimRequest struct {
	ObjectType objectType.Type `json:"query" validate:"required"`
}

func (h ClaimHandler) HandlePageShow(c echo.Context) error {
	return Render(c, claimPage.Show(InitActivePage(c)))
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

	t := objectType.MapObjectType(string(claimRequest.ObjectType))
	if t == objectType.Unknown {
		return echo.ErrBadRequest
	}

	if err := claim.UpdateClaimed(h.DB); err != nil {
		log.Println(err)
		return err
	}

	claimed, err := claim.ClaimObjectID(h.DB, h.AllowedList, t, model.GUI)
	if err != nil {
		log.Println(err)
		return Render(c, claimPage.ClaimedID(err.Error()))
	}

	return Render(c, claimPage.ClaimedID(strconv.Itoa(int(claimed.ID))))
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

	t := objectType.MapObjectType(string(claimRequest.ObjectType))
	if t == objectType.Unknown {
		return c.JSON(http.StatusInternalServerError, ObjectClaimAPIError{
			Message: "Invalid object type",
		})
	}

	if err := claim.UpdateClaimed(h.DB); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, ObjectClaimAPIError{
			Message: "Error updating claimed objects",
		})
	}

	claimed, err := claim.ClaimObjectID(h.DB, h.AllowedList, t, model.API)
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
	ObjectType objectType.Type `json:"query"`
}

func (h ClaimHandler) HandleObjectTypeQuery(c echo.Context) error {
	objectTypeQuery := ObjectTypeQuery{}
	if err := c.Bind(&objectTypeQuery); err != nil {
		return err
	}

	objectTypes := objectType.GetObjectTypes()

	if objectTypeQuery.ObjectType == "" {
		return Render(c, claimPage.QueryResults(objectTypes))
	}

	results := []objectType.Type{}
	for _, objectType := range objectTypes {
		if strings.Contains(strings.ToLower(string(objectType)), strings.ToLower(string(objectTypeQuery.ObjectType))) {
			results = append(results, objectType)
		}
	}

	return Render(c, claimPage.QueryResults(results))
}
