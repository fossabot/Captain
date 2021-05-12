package CaptainLib

import (
	"encoding/json"
	"fmt"
)

type Flight struct {
	ID         int
	AirspaceID int
	Name       string
}

func (c *CaptainClient) GetAllFlights() ([]Flight, error) {
	results, err := c.restGET("flights")
	if err != nil {
		return nil, fmt.Errorf("unable to get a list of flights:\n%w", err)
	}
	var flights []Flight
	err = json.Unmarshal(results, &flights)
	if err != nil {
		return nil, fmt.Errorf("unable to format response as array of Flights:\n%w", err)
	}
	return flights, nil
}

func (c *CaptainClient) GetFlightByID(id int) (Flight, error) {
	results, err := c.restGET(fmt.Sprintf("flight/%d", id))
	if err != nil {
		return Flight{}, fmt.Errorf("unable to get flight by id %d:\n%w", id, err)
	}
	var flight Flight
	err = json.Unmarshal(results, &flight)
	if err != nil {
		return Flight{}, fmt.Errorf("unable to format response as a Flight:\n%w", err)
	}
	return flight, nil
}

func (c *CaptainClient) CreateFlight(name string, airspaceID int) (Flight, error) {
	result, err := c.restPOST("flight", map[string]interface{}{
		"AirspaceID": airspaceID,
		"Name": name,
	})
	if err != nil {
		return Flight{}, fmt.Errorf("unable to create Flight:\n%w", err)
	}
	var flight Flight
	err = json.Unmarshal(result, &flight)
	if err != nil {
		return Flight{}, fmt.Errorf("unable to parse response as a Flight:\n%w", err)
	}
	return flight, nil
}

func (c *CaptainClient) UpdateFlight(flight Flight) error {
	_, err := c.restPUT(fmt.Sprintf("flight/%d", flight.ID), map[string]interface{}{
		"AirspaceID": flight.AirspaceID,
		"Name": flight.Name,
	})
	if err != nil {
		return fmt.Errorf("unable to update flight with ID %d:\n%w", flight.ID, err)
	}
	return nil
}

func (c *CaptainClient) DeleteFlight(id int) error {
	_, err := c.restDELETE(fmt.Sprintf("flight/%d", id))
	if err != nil {
		return fmt.Errorf("unable to delete flight with ID %d:\n%w", id, err)
	}
	return nil
}
