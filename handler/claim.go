package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/component/page/claim"
	"github.com/philipp-mlr/al-id-maestro/model"
)

type ClaimHandler struct {
	DB *model.DB
}

func (h ClaimHandler) HandleClaimShow(c echo.Context) error {
	return Render(c, claim.Show())
}

func (h ClaimHandler) HandleIDClaim(c echo.Context) error {
	// generate random int

	return Render(c, claim.ClaimedID("12345"))
}

func (h ClaimHandler) HandleClaimTypeQuery(c echo.Context) error {
	objectTypeQuery := new(model.ObjectTypeQuery)
	if err := c.Bind(objectTypeQuery); err != nil {
		return err
	}

	// if err := c.Validate(objectTypeQuery); err != nil {
	// 	return err
	// }

	objectTypes, err := objectTypeQuery.GetResults(h.DB)
	if err != nil {
		return err
	}

	return Render(c, claim.Result(objectTypes))
}
