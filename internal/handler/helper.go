package handler

import "github.com/labstack/echo/v4"

func InitActivePage(ctx echo.Context) map[string]bool {
	activePage := make(map[string]bool)
	activePage[ctx.Path()] = true
	return activePage
}
