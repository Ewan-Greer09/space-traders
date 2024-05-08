package client

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"space-traders/service/config"
)

var (
	ErrCannotGetShipLocation = fmt.Errorf("unable to get ship location")
	ErrCannotOribitShip      = fmt.Errorf("unable to orbit ship")
)

// interface for interacting with the space traders API endpoints
type SpaceTradersClient interface {
	// get the server status of the space traders API
	GetStatus() (map[string]any, error)
	// Get the callers agent agent data //? how will this work with an account that has many API keys?
	GetMyAgent() (map[string]any, error)
	// Get a list of agents contracts
	ListContracts() (map[string]any, error)
	// Get a specific contract
	GetContract(contractID string) (map[string]any, error)
	// Accept a contract
	AcceptContract(contractID string) (map[string]any, error)
	// Deliver Cargo to a Contract
	DeliverCargo(contractID string) (map[string]any, error)
}

type Client struct {
	*resty.Client
	logger *slog.Logger
	cfg    *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		Client: resty.New().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			SetBaseURL(cfg.ClientBaseURL),
		logger: slog.Default(),
		cfg:    cfg,
	}
}

func (c *Client) GetStatus() (map[string]any, error) {
	c.SetHeader("Authorization", "")

	resp, err := c.R().Get("/")
	if err != nil {
		return nil, fmt.Errorf("error getting status: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error getting status: %s", resp.Status())
	}

	var status map[string]any
	err = json.Unmarshal(resp.Body(), &status)
	if err != nil {
		return nil, fmt.Errorf("error decoding status: %w", err)
	}

	return status, nil
}
