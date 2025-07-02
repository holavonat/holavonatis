package api

import (
	"time"
)

type Upstream struct {
	Client *Client
	Source Source
}

func (e *Upstream) Fetch() (Holavonat, error) {
	serviceDay := time.Now().Format("20060102")
	details, err := e.Client.AllDetails(serviceDay)
	if err != nil {
		return Holavonat{}, err
	}

	return Holavonat{
		VehiclePositions: details.Data.VehiclePositions,
		LastUpdated:      time.Now().Unix(),
		Timestamp:        time.Now().Format(time.RFC3339),
	}, nil
}

func (e *Upstream) FetchByServiceDay(serviceDay string) (Holavonat, error) {
	details, err := e.Client.AllDetails(serviceDay)
	if err != nil {
		return Holavonat{}, err
	}

	e.Source.Schema.Format = "json"
	return Holavonat{
		VehiclePositions: details.Data.VehiclePositions,
		LastUpdated:      time.Now().Unix(),
		Timestamp:        time.Now().Format(time.RFC3339),
		Source:           e.Source,
	}, nil
}
