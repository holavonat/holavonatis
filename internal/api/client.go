package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	_ "golang.org/x/crypto/x509roots/fallback"
)

var (
	ErrMissingCustomHTTPClient = fmt.Errorf("http.Client cannot be nil")
)

type Client struct {
	Client   *http.Client
	Headers  map[string]string
	Endpoint string
}

func NewClient(endpoint string, headers map[string]string) *Client {
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	return &Client{
		Endpoint: endpoint,
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
		Headers: headers,
	}
}

func NewClientProxy(endpoint string, headers map[string]string, client *http.Client, proxyURL string) (*Client, error) {
	if client == nil {
		return nil, ErrMissingCustomHTTPClient
	}
	if proxyURL == "" {
		return nil, fmt.Errorf("proxy URL cannot be empty")
	}

	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"

	transport := &http.Transport{}
	if proxyURL != "" {
		parsedProxyURL, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
		}
		transport.Proxy = http.ProxyURL(parsedProxyURL)
	}

	client.Transport = transport

	return &Client{
		Endpoint: endpoint,
		Client:   client,
		Headers:  headers,
	}, nil
}

func NewClientCustomHTTP(endpoint string, headers map[string]string, client *http.Client) (*Client, error) {
	if client == nil {
		return nil, ErrMissingCustomHTTPClient
	}
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	return &Client{
		Endpoint: endpoint,
		Client:   client,
		Headers:  headers,
	}, nil
}

func (c *Client) Do(query string) ([]byte, error) {
	reqBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) AllDetails(serviceDay string) (OTPResponse, error) {
	if serviceDay == "" {
		return OTPResponse{}, fmt.Errorf("serviceDay cannot be empty")
	}

	if _, err := time.Parse("20060102", serviceDay); err != nil {
		return OTPResponse{}, fmt.Errorf("invalid serviceDay format: %w", err)
	}

	query := fmt.Sprintf(`
		 {
        vehiclePositions(
            swLat: 45.5,
			swLon: 16.1,
			neLat: 48.7,
			neLon: 22.8,
            modes: [TRAM,RAIL,RAIL_REPLACEMENT_BUS,SUBURBAN_RAILWAY,TRAMTRAIN],
        ) {
            vehicleId
            lat
            lon
			label
            heading
            lastUpdated
            speed
            stopRelationship {
                status
                stop {
                    gtfsId
                    name
                }
            }
            trip {
                id
                stoptimes: stoptimesForDate(
                  serviceDate: "%s"
                ){
                  stop {
    								name
    								lat
    								lon
    								platformCode
    							}
    							realtimeArrival
    							realtimeDeparture
    							arrivalDelay
    							departureDelay
    							scheduledArrival
    							scheduledDeparture        
                }
                gtfsId
                routeShortName
                tripHeadsign
                tripShortName
                trainName
        				domesticResTrainNumber
        				wheelchairAccessible
        				bikesAllowed
                trainCategoryBaseId
                tripNumber
                tripGeometry {
                  length
                  points
                }
                alerts {
                  id
                  alertHash
                  feed
                  alertHeaderText
                  alertDescriptionText
                  alertCause
                  alertSeverityLevel
                  alertUrl
                  alertEffect
                  effectiveEndDate
                  effectiveStartDate
                                
                }
                trainCategoryId
        				infoServices {
        					name
        					fromStopIndex
        					tillStopIndex
        					fontCharSet
        					fontCode
        					displayable
        				}
                route {
                    mode
                    shortName
                    longName(language: "hu")
                    textColor
                    color
                }
                pattern {
                    id
                }
            }
        }
    }`, serviceDay)
	body, err := c.Do(query)
	if err != nil {
		return OTPResponse{}, err
	}

	var result OTPResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return OTPResponse{}, err
	}

	return result, nil
}
