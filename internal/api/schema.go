package api

import "encoding/json"

type OTPResponse struct {
	Data Data `json:"data,omitempty"`
}
type Stop struct {
	Name         string  `json:"name"`
	PlatformCode string  `json:"platformCode"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
}
type Route struct {
	Mode      string `json:"mode,omitempty"`
	ShortName string `json:"shortName,omitempty"`
	LongName  string `json:"longName,omitempty"`
	TextColor string `json:"textColor,omitempty"`
	Color     string `json:"color,omitempty"`
}
type StopRelationship struct {
	Status string `json:"status,omitempty"`
	Stop   Stop   `json:"stop,omitempty"`
}
type NextStop struct {
	ArrivalDelay int64 `json:"arrivalDelay"`
}
type ArrivalStoptime struct {
	ArrivalDelay int64 `json:"arrivalDelay"`
}

type Stoptimes struct {
	Stop               Stop  `json:"stop,omitempty"`
	RealtimeArrival    int64 `json:"realtimeArrival"`
	RealtimeDeparture  int64 `json:"realtimeDeparture"`
	ArrivalDelay       int64 `json:"arrivalDelay"`
	DepartureDelay     int64 `json:"departureDelay"`
	ScheduledArrival   int64 `json:"scheduledArrival"`
	ScheduledDeparture int64 `json:"scheduledDeparture"`
}
type TripGeometry struct {
	Points string `json:"points,omitempty"`
	Length int    `json:"length,omitempty"`
}
type Alerts struct {
	AlertURL             any    `json:"alertUrl"`
	ID                   string `json:"id"`
	Feed                 string `json:"feed"`
	AlertHeaderText      string `json:"alertHeaderText"`
	AlertDescriptionText string `json:"alertDescriptionText"`
	AlertCause           string `json:"alertCause"`
	AlertSeverityLevel   string `json:"alertSeverityLevel"`
	AlertEffect          string `json:"alertEffect"`
	AlertHash            int    `json:"alertHash"`
	EffectiveEndDate     int    `json:"effectiveEndDate"`
	EffectiveStartDate   int    `json:"effectiveStartDate"`
}
type Pattern struct {
	ID string `json:"id,omitempty"`
}

type InfoService struct {
	Name          string `json:"name,omitempty"`
	FontCharSet   string `json:"fontCharSet,omitempty"`
	FromStopIndex int    `json:"fromStopIndex,omitempty"`
	TillStopIndex int    `json:"tillStopIndex,omitempty"`
	FontCode      int    `json:"fontCode,omitempty"`
	Displayable   bool   `json:"displayable,omitempty"`
}

type Trip struct {
	Route                  Route           `json:"route,omitempty"`
	TripGeometry           TripGeometry    `json:"tripGeometry,omitempty"`
	WheelchairAccessible   string          `json:"wheelchairAccessible,omitempty"`
	TripHeadsign           string          `json:"tripHeadsign,omitempty"`
	TripShortName          string          `json:"tripShortName,omitempty"`
	DomesticResTrainNumber string          `json:"domesticResTrainNumber,omitempty"`
	RouteShortName         string          `json:"routeShortName,omitempty"`
	BikesAllowed           string          `json:"bikesAllowed,omitempty"`
	Pattern                Pattern         `json:"pattern,omitempty"`
	TripNumber             string          `json:"tripNumber,omitempty"`
	GtfsID                 string          `json:"gtfsId,omitempty"`
	TrainCategoryID        string          `json:"trainCategoryId,omitempty"`
	ID                     string          `json:"id,omitempty"`
	InfoServices           []InfoService   `json:"infoServices,omitempty"`
	Stoptimes              []Stoptimes     `json:"stoptimes,omitempty"`
	Alerts                 []Alerts        `json:"alerts"`
	ArrivalStoptime        ArrivalStoptime `json:"arrivalStoptime"`
	TrainCategoryBaseID    int64           `json:"trainCategoryBaseId,omitempty"`
}
type VehiclePositions struct {
	StopRelationship StopRelationship `json:"stopRelationship,omitempty"`
	VehicleID        string           `json:"vehicleId,omitempty"`
	Label            string           `json:"label,omitempty"`
	Trip             Trip             `json:"trip"`
	Lat              float64          `json:"lat,omitempty"`
	Lon              float64          `json:"lon,omitempty"`
	Heading          float64          `json:"heading,omitempty"`
	LastUpdated      int              `json:"lastUpdated,omitempty"`
	Speed            float64          `json:"speed,omitempty"`
	NextStop         NextStop         `json:"nextStop"`
}
type Data struct {
	VehiclePositions []VehiclePositions `json:"vehiclePositions"`
}

type Holavonat struct {
	Source           Source             `json:"source"`
	Timestamp        string             `json:"timestamp"`
	VehiclePositions []VehiclePositions `json:"vehiclePositions"`
	LastUpdated      int64              `json:"lastUpdated"`
}

func (v *Holavonat) Json() ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Source struct {
	Origin     string `json:"origin,omitempty"`
	Latest     string `json:"latest,omitempty"`
	DirectLink string `json:"directLink,omitempty"`
	Schema     Schema `json:"schema,omitempty"`
}

type Schema struct {
	Version string `json:"version,omitempty"`
	Link    string `json:"link,omitempty"`
	Format  string `json:"format,omitempty"`
}
