package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/philipp-mlr/al-id-maestro/component/page/index"
)

type IndexHanlder struct {
}

func (h IndexHanlder) HandleIndexShow(c echo.Context) error {
	return Render(c, index.Show())
}
