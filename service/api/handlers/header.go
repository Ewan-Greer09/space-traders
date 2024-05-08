package handlers

import (
	"github.com/labstack/echo/v4"

	"space-traders/service/api/models"
	"space-traders/service/views/index"
)

func (h *ViewHandler) HandlerHeader(c echo.Context) error {
	status, err := h.Client.GetStatus()
	if err != nil {
		return err
	}

	// This is a demo of how data will be passed to the view handler. In the real application,
	// data will be generated from translations, database calls, and API requests. It will be highly dynamic.
	// I think storing it as a map is the best way to handle it. As changes in API responses, translations, and database
	// calls will be reflected in the data map. Without manual changes to keys, such as in a struct.
	var data = make(map[string]interface{})
	headerArgs := models.HeaderArgs{
		Username: "Demo User",
		Credits:  "1000",
	}

	if status["status"] != "" {
		headerArgs.Status = true
	} else {
		headerArgs.Status = false
	}

	data["HeaderArgs"] = headerArgs

	return index.Page(data).Render(c.Request().Context(), c.Response())
}
