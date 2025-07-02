package cloudflare

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	cloudflareURL = "https://1.1.1.1/cdn-cgi/trace"
)

type Trace struct {
	Fl          string `json:"fl"`
	Host        string `json:"h"`
	Ip          string `json:"ip"`
	Timestamp   string `json:"ts"`
	VisitScheme string `json:"visit_scheme"`
	UserAgent   string `json:"uag"`
	Colocation  string `json:"colo"`
	Sliver      string `json:"sliver"`
	Http        string `json:"http"`
	Location    string `json:"loc"`
	Tls         string `json:"tls"`
	Sni         string `json:"sni"`
	Warp        string `json:"warp"`
	Gateway     string `json:"gateway"`
	Rbi         string `json:"rbi"`
	Kex         string `json:"kex"`
}

func ParseCloudflareTrace(reader io.Reader) (Trace, error) {
	var trace Trace
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "fl":
			trace.Fl = value
		case "h":
			trace.Host = value
		case "ip":
			trace.Ip = value
		case "ts":
			trace.Timestamp = value
		case "visit_scheme":
			trace.VisitScheme = value
		case "uag":
			trace.UserAgent = value
		case "colo":
			trace.Colocation = value
		case "sliver":
			trace.Sliver = value
		case "http":
			trace.Http = value
		case "loc":
			trace.Location = value
		case "tls":
			trace.Tls = value
		case "sni":
			trace.Sni = value
		case "warp":
			trace.Warp = value
		case "gateway":
			trace.Gateway = value
		case "rbi":
			trace.Rbi = value
		case "kex":
			trace.Kex = value
		}
	}

	if err := scanner.Err(); err != nil {
		return Trace{}, fmt.Errorf("error scanning response: %w", err)
	}

	return trace, nil
}

func GetTrace(client *http.Client) (Trace, error) {
	resp, err := client.Get(cloudflareURL)
	if err != nil {
		return Trace{}, fmt.Errorf("error fetching Cloudflare trace: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Trace{}, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	trace, err := ParseCloudflareTrace(resp.Body)
	if err != nil {
		return Trace{}, err
	}

	return trace, nil
}

func GetTraceProxy(client *http.Client, proxyURL string) (Trace, error) {
	transport := &http.Transport{}
	if proxyURL != "" {
		parsedProxyURL, err := url.Parse(proxyURL)
		if err != nil {
			return Trace{}, fmt.Errorf("failed to parse proxy URL: %w", err)
		}
		transport.Proxy = http.ProxyURL(parsedProxyURL)
	}

	client.Transport = transport

	return GetTrace(client)
}
