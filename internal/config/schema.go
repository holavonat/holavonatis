package config

import (
	"github.com/holavonat/holavonatis/internal/api"
	"github.com/holavonat/holavonatis/internal/cloudflare/r2"
	log "github.com/holavonat/holavonatis/internal/logger"
)

type Compression string

type TimeFrame string

type CronMode string

const (
	Second TimeFrame   = "second"
	Minute TimeFrame   = "minute"
	Hour   TimeFrame   = "hour"
	Brotli Compression = "br"
	Gzip   Compression = "gzip"
	Zstd   Compression = "zstd"
	None   Compression = "none"
	Fix    CronMode    = "fix"
	Window CronMode    = "window"
)

type Config struct {
	Headers         map[string]string `yaml:"headers"`
	ObjectStorage   ObjectStorage     `yaml:"objectstorage"`
	Source          api.Source        `yaml:"Source"`
	Network         Network           `yaml:"Network"`
	GraphqlEndpoint string            `yaml:"graphqlendpoint"`
	File            File              `yaml:"file"`
	Output          Output            `yaml:"output"`
	Cron            Cron              `yaml:"cron"`
	Log             log.Config        `yaml:"log"`
	EulaAccepted    bool              `yaml:"eula_accepted"`
}

type Output struct {
	NamePrefix string `yaml:"nameprefix"`
	Format     Format `yaml:"format"`
	Archive    bool   `yaml:"archive"`
}

type Format struct {
	JSON bool `yaml:"json"`
}

type ObjectStorage struct {
	Compression       Compression `yaml:"compression"`
	AccessKeyID       string      `yaml:"accesskeyid"`
	SecretAccessKey   string      `yaml:"secretaccesskey"`
	BucketName        string      `yaml:"bucketname"`
	ObjectPath        string      `yaml:"objectpath"`
	EndpointURL       string      `yaml:"endpointurl"`
	PublicEndpointURL string      `yaml:"publicendpointurl"`
}

type File struct {
	Path string `yaml:"path"`
}

type Cron struct {
	Mode     CronMode   `yaml:"mode"`
	Duration TimeFrame  `yaml:"duration"`
	Fix      FixMode    `yaml:"fix"`
	Window   WindowMode `yaml:"window"`
}

type FixMode struct {
	Interval int `yaml:"interval"`
}

type WindowMode struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type Network struct {
	Proxy   string `yaml:"proxy"`
	Timeout int    `yaml:"timeout"`
}

type App struct {
	ObjectStorage r2.Cloudflare
	Cfg           Config
}
