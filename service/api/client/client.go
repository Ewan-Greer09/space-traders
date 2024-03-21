package client

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"space-traders/models"
)

type Client struct {
	*resty.Client
}

func NewClient() *Client {
	r := resty.New()
	r.SetHeader("Content-Type", "application/json")
	r.SetHeader("Accept", "application/json")
	r.SetBaseURL("https://api.spacetraders.io/v2")
	r.SetDoNotParseResponse(true)

	return &Client{
		Client: r,
	}
}

func (c *Client) GetMyShips(symbol string) (*models.GetMyShip200Response, error) {
	resp, err := c.R().Get("/my/ships/" + symbol)
	if err != nil {
		return nil, err
	}
	defer resp.RawBody().Close()

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to get ship location: " + resp.Status())
	}

	var data models.GetMyShip200Response
	err = json.NewDecoder(resp.RawBody()).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) SendToOrbit(symbol string) error {
	resp, err := c.R().Post("/my/ships/" + symbol + "/orbit")
	if err != nil {
		return err
	}
	defer resp.RawBody().Close()

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to orbit ship")
	}

	return nil
}

func (c *Client) NavigateShip(systemSymbol, waypointSymbol string) (*models.Waypoint, error) {
	resp, err := c.R().SetBody(struct {
		WaypointSymbol string `json:"waypointSymbol"`
	}{
		WaypointSymbol: waypointSymbol,
	}).
		Post(fmt.Sprintf("/my/ships/%s/navigate", systemSymbol))
	if err != nil {
		return nil, err
	}
	defer resp.RawBody().Close()

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unable to navigate ship:" + resp.Status())
	}

	var waypointData models.Waypoint
	err = json.NewDecoder(resp.RawBody()).Decode(&waypointData)
	if err != nil {
		return nil, err
	}

	return &waypointData, nil
}
