package client

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"

	"space-traders/models"
)

type Client struct {
	resty *resty.Client
}

func NewClient() *Client {
	r := resty.New()
	r.SetHeader("Content-Type", "application/json")
	r.SetHeader("Accept", "application/json")
	r.SetBaseURL("https://api.spacetraders.io/v2")
	r.SetDoNotParseResponse(true)

	return &Client{
		resty: r,
	}
}

func (c *Client) GetMyShips(symbol string) (*models.GetMyShip200Response, error) {
	resp, err := c.resty.R().Get("/my/ships/" + symbol)
	if err != nil {
		return nil, err
	}

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
	resp, err := c.resty.R().Post("/my/ships/" + symbol + "/orbit")
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to orbit ship")
	}

	return nil
}

func (c *Client) NavigateShip(symbol, waypointSymbol string) (*models.Waypoint, error) {
	type navigateReq struct {
		WaypointSymbol string `json:"waypointSymbol"`
	}

	navResp, err := c.resty.R().SetDoNotParseResponse(true).SetBody(navigateReq{WaypointSymbol: waypointSymbol}).
		Post(fmt.Sprintf("/my/ships/%s/navigate", symbol))

	if err != nil {
		return nil, err
	}

	if navResp.StatusCode() != 200 {
		var data map[string]interface{}
		err = json.NewDecoder(navResp.RawBody()).Decode(&data)
		if err != nil {
			return nil, err
		}

		log.Printf("navigate ship error: %v", data)

		return nil, fmt.Errorf("unable to navigate ship:" + navResp.Status())
	}

	var waypointData models.Waypoint
	err = json.NewDecoder(navResp.RawBody()).Decode(&waypointData)
	if err != nil {
		return nil, err
	}

	return &waypointData, nil
}

func (c *Client) SetHeader(key, value string) {
	c.resty.SetHeader(key, value)
}
