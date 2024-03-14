package models

import "time"

type Waypoint struct {
	// The symbol of the waypoint.
	Symbol string `json:"symbol"`
	Type   string `json:"type"`
	// The symbol of the system.
	SystemSymbol string `json:"systemSymbol"`
	// Relative position of the waypoint on the system's x axis. This is not an absolute position in the universe.
	X int32 `json:"x"`
	// Relative position of the waypoint on the system's y axis. This is not an absolute position in the universe.
	Y int32 `json:"y"`
	// Waypoints that orbit this waypoint.
	Orbitals []WaypointOrbital `json:"orbitals"`
	// The symbol of the parent waypoint, if this waypoint is in orbit around another waypoint. Otherwise this value is undefined.
	Orbits  *string          `json:"orbits,omitempty"`
	Faction *WaypointFaction `json:"faction,omitempty"`
	// The traits of the waypoint.
	Traits []WaypointTrait `json:"traits"`
	// The modifiers of the waypoint.
	Modifiers []WaypointModifier `json:"modifiers,omitempty"`
	Chart     *Chart             `json:"chart,omitempty"`
	// True if the waypoint is under construction.
	IsUnderConstruction  bool `json:"isUnderConstruction"`
	AdditionalProperties map[string]interface{}
}

type WaypointOrbital struct {
	// The symbol of the orbiting waypoint.
	Symbol               string `json:"symbol"`
	AdditionalProperties map[string]interface{}
}

type WaypointFaction struct {
	Symbol               string `json:"symbol"`
	AdditionalProperties map[string]interface{}
}

type WaypointTrait struct {
	Symbol string `json:"symbol"`
	// The name of the trait.
	Name string `json:"name"`
	// A description of the trait.
	Description          string `json:"description"`
	AdditionalProperties map[string]interface{}
}

type WaypointModifier struct {
	Symbol string `json:"symbol"`
	// The name of the trait.
	Name string `json:"name"`
	// A description of the trait.
	Description          string `json:"description"`
	AdditionalProperties map[string]interface{}
}

type Chart struct {
	// The symbol of the waypoint.
	WaypointSymbol *string `json:"waypointSymbol,omitempty"`
	// The agent that submitted the chart for this waypoint.
	SubmittedBy *string `json:"submittedBy,omitempty"`
	// The time the chart for this waypoint was submitted.
	SubmittedOn          *time.Time `json:"submittedOn,omitempty"`
	AdditionalProperties map[string]interface{}
}

type GetStatus200Response struct {
	Status      string `json:"status,omitempty"`
	Version     string `json:"version,omitempty"`
	ResetDate   string `json:"resetDate,omitempty"`
	Description string `json:"description,omitempty"`
	Stats       struct {
		Agents    int `json:"agents,omitempty"`
		Ships     int `json:"ships,omitempty"`
		Systems   int `json:"systems,omitempty"`
		Waypoints int `json:"waypoints,omitempty"`
	} `json:"stats,omitempty"`
	Leaderboards struct {
		MostCredits []struct {
			AgentSymbol string `json:"agentSymbol,omitempty"`
			Credits     int64  `json:"credits,omitempty"`
		} `json:"mostCredits,omitempty"`
		MostSubmittedCharts []struct {
			AgentSymbol string `json:"agentSymbol,omitempty"`
			ChartCount  int    `json:"chartCount,omitempty"`
		} `json:"mostSubmittedCharts,omitempty"`
	} `json:"leaderboards,omitempty"`
	ServerResets struct {
		Next      string `json:"next,omitempty"`
		Frequency string `json:"frequency,omitempty"`
	} `json:"serverResets,omitempty"`
	Announcements []struct {
		Title string `json:"title,omitempty"`
		Body  string `json:"body,omitempty"`
	} `json:"announcements,omitempty"`
	Links []struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"links,omitempty"`
}

type GetMyShip200Response struct {
	Data struct {
		Symbol       string `json:"symbol,omitempty"`
		Registration struct {
			Name          string `json:"name,omitempty"`
			FactionSymbol string `json:"factionSymbol,omitempty"`
			Role          string `json:"role,omitempty"`
		} `json:"registration,omitempty"`
		Nav struct {
			SystemSymbol   string `json:"systemSymbol,omitempty"`
			WaypointSymbol string `json:"waypointSymbol,omitempty"`
			Route          struct {
				Destination struct {
					Symbol       string `json:"symbol,omitempty"`
					Type         string `json:"type,omitempty"`
					SystemSymbol string `json:"systemSymbol,omitempty"`
					X            int    `json:"x,omitempty"`
					Y            int    `json:"y,omitempty"`
				} `json:"destination,omitempty"`
				Origin struct {
					Symbol       string `json:"symbol,omitempty"`
					Type         string `json:"type,omitempty"`
					SystemSymbol string `json:"systemSymbol,omitempty"`
					X            int    `json:"x,omitempty"`
					Y            int    `json:"y,omitempty"`
				} `json:"origin,omitempty"`
				DepartureTime time.Time `json:"departureTime,omitempty"`
				Arrival       time.Time `json:"arrival,omitempty"`
			} `json:"route,omitempty"`
			Status     string `json:"status,omitempty"`
			FlightMode string `json:"flightMode,omitempty"`
		} `json:"nav,omitempty"`
		Crew struct {
			Current  int    `json:"current,omitempty"`
			Required int    `json:"required,omitempty"`
			Capacity int    `json:"capacity,omitempty"`
			Rotation string `json:"rotation,omitempty"`
			Morale   int    `json:"morale,omitempty"`
			Wages    int    `json:"wages,omitempty"`
		} `json:"crew,omitempty"`
		Frame struct {
			Symbol         string `json:"symbol,omitempty"`
			Name           string `json:"name,omitempty"`
			Description    string `json:"description,omitempty"`
			Condition      int    `json:"condition,omitempty"`
			Integrity      int    `json:"integrity,omitempty"`
			ModuleSlots    int    `json:"moduleSlots,omitempty"`
			MountingPoints int    `json:"mountingPoints,omitempty"`
			FuelCapacity   int    `json:"fuelCapacity,omitempty"`
			Requirements   struct {
				Power int `json:"power,omitempty"`
				Crew  int `json:"crew,omitempty"`
				Slots int `json:"slots,omitempty"`
			} `json:"requirements,omitempty"`
		} `json:"frame,omitempty"`
		Reactor struct {
			Symbol       string `json:"symbol,omitempty"`
			Name         string `json:"name,omitempty"`
			Description  string `json:"description,omitempty"`
			Condition    int    `json:"condition,omitempty"`
			Integrity    int    `json:"integrity,omitempty"`
			PowerOutput  int    `json:"powerOutput,omitempty"`
			Requirements struct {
				Power int `json:"power,omitempty"`
				Crew  int `json:"crew,omitempty"`
				Slots int `json:"slots,omitempty"`
			} `json:"requirements,omitempty"`
		} `json:"reactor,omitempty"`
		Engine struct {
			Symbol       string `json:"symbol,omitempty"`
			Name         string `json:"name,omitempty"`
			Description  string `json:"description,omitempty"`
			Condition    int    `json:"condition,omitempty"`
			Integrity    int    `json:"integrity,omitempty"`
			Speed        int    `json:"speed,omitempty"`
			Requirements struct {
				Power int `json:"power,omitempty"`
				Crew  int `json:"crew,omitempty"`
				Slots int `json:"slots,omitempty"`
			} `json:"requirements,omitempty"`
		} `json:"engine,omitempty"`
		Cooldown struct {
			ShipSymbol       string    `json:"shipSymbol,omitempty"`
			TotalSeconds     int       `json:"totalSeconds,omitempty"`
			RemainingSeconds int       `json:"remainingSeconds,omitempty"`
			Expiration       time.Time `json:"expiration,omitempty"`
		} `json:"cooldown,omitempty"`
		Modules []struct {
			Symbol       string `json:"symbol,omitempty"`
			Capacity     int    `json:"capacity,omitempty"`
			Range        int    `json:"range,omitempty"`
			Name         string `json:"name,omitempty"`
			Description  string `json:"description,omitempty"`
			Requirements struct {
				Power int `json:"power,omitempty"`
				Crew  int `json:"crew,omitempty"`
				Slots int `json:"slots,omitempty"`
			} `json:"requirements,omitempty"`
		} `json:"modules,omitempty"`
		Mounts []struct {
			Symbol       string   `json:"symbol,omitempty"`
			Name         string   `json:"name,omitempty"`
			Description  string   `json:"description,omitempty"`
			Strength     int      `json:"strength,omitempty"`
			Deposits     []string `json:"deposits,omitempty"`
			Requirements struct {
				Power int `json:"power,omitempty"`
				Crew  int `json:"crew,omitempty"`
				Slots int `json:"slots,omitempty"`
			} `json:"requirements,omitempty"`
		} `json:"mounts,omitempty"`
		Cargo struct {
			Capacity  int `json:"capacity,omitempty"`
			Units     int `json:"units,omitempty"`
			Inventory []struct {
				Symbol      string `json:"symbol,omitempty"`
				Name        string `json:"name,omitempty"`
				Description string `json:"description,omitempty"`
				Units       int    `json:"units,omitempty"`
			} `json:"inventory,omitempty"`
		} `json:"cargo,omitempty"`
		Fuel struct {
			Current  int `json:"current,omitempty"`
			Capacity int `json:"capacity,omitempty"`
			Consumed struct {
				Amount    int       `json:"amount,omitempty"`
				Timestamp time.Time `json:"timestamp,omitempty"`
			} `json:"consumed,omitempty"`
		} `json:"fuel,omitempty"`
	} `json:"data,omitempty"`
}
