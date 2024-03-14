package actions

import (
	"context"
	"os"

	openAPI "github.com/UnseenBook/spacetraders-go-sdk"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"space-traders/repository/postgres"
	"space-traders/service/views/components/index"
)

type ActionHandler struct {
	Client *openAPI.APIClient
	userDB *postgres.Queries
}

func NewActionHandler() *ActionHandler {
	cfg := openAPI.NewConfiguration()
	cfg.AddDefaultHeader("Content-Type", "application/json")
	cfg.AddDefaultHeader("Accept", "application/json")

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	return &ActionHandler{
		Client: openAPI.NewAPIClient(cfg),
		userDB: postgres.New(pool),
	}
}

func (ah *ActionHandler) HandleGetSystemDataFromAllShips(c echo.Context) error {
	resp, _, err := ah.Client.FleetAPI.GetMyShips(c.Request().Context()).Execute()
	if err != nil {
		return err
	}

	for _, ship := range resp.Data {
		resp, _, err := ah.Client.SystemsAPI.GetSystem(context.Background(), ship.Nav.SystemSymbol).Execute()
		if err != nil {
			return err
		}

		_, err = ah.userDB.GetSystemBySymbol(context.Background(), pgtype.Text{String: resp.Data.Symbol, Valid: true})
		if err == nil {
			continue
		}

		_, err = ah.userDB.CreateSystem(context.Background(), postgres.CreateSystemParams{
			Symbol:       pgtype.Text{String: resp.Data.Symbol, Valid: true},
			SectorSymbol: pgtype.Text{String: resp.Data.SectorSymbol, Valid: true},
			Type:         pgtype.Text{String: string(resp.Data.Type), Valid: true},
			X:            pgtype.Int4{Int32: int32(resp.Data.X), Valid: true},
			Y:            pgtype.Int4{Int32: int32(resp.Data.Y), Valid: true},
		})
		if err != nil {
			return err
		}

	}

	return index.GetShipLocationsSuccess().Render(c.Request().Context(), c.Response())
}
