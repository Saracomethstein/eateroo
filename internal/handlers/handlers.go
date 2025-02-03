package handlers

import (
	"go_day_03/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "ok\n")
}

func GetPlace(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	searchQuery := c.QueryParam("search")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 || limit > 100 {
		limit = 30
	}

	restaurants, total, err := service.GetPlace(page, limit, searchQuery)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]error{"error": err})

	}

	response := map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"total":       total,
		"restaurants": restaurants,
	}

	return c.JSON(http.StatusOK, response)
}
