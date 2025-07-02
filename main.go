package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/D3vl0per/crypt/compression"
	"github.com/holavonat/holavonatis/internal/api"
	"github.com/holavonat/holavonatis/internal/cloudflare"
	"github.com/holavonat/holavonatis/internal/cloudflare/r2"
	"github.com/holavonat/holavonatis/internal/config"
	log "github.com/holavonat/holavonatis/internal/logger"
)

func main() {

	l := log.New("main")
	cfg, err := config.GetConfig()
	if err != nil {
		l.DPanicw("Failed to load config", "error", err)
		return
	}

	trace, err := cloudflare.GetTrace(&http.Client{})
	if err != nil {
		l.DPanicw("Failed to get trace", "error", err)
		return
	}

	l.Infow("Public IP address", "ip", trace.Ip, "location", trace.Location, "colo", trace.Colocation)

	app := config.App{
		Cfg: cfg,
	}

	if cfg.File.Path != "" {
		_, err := os.Stat(cfg.File.Path)
		if os.IsNotExist(err) {
			err = os.MkdirAll(cfg.File.Path, 0755)
			if err != nil {
				l.DPanicw("Failed to create folder", "error", err)
				return
			}
			l.Infow("Created folder for file output", "path", cfg.File.Path)
		} else if err != nil {
			l.DPanicw("Failed to check folder", "error", err)
			return
		}
	}

	if cfg.ObjectStorage.AccessKeyID != "" || cfg.ObjectStorage.SecretAccessKey != "" || cfg.ObjectStorage.EndpointURL != "" {
		s3, err := r2.NewClient(r2.Cloudflare{
			AccessKeyID:       cfg.ObjectStorage.AccessKeyID,
			SecretAccessKey:   cfg.ObjectStorage.SecretAccessKey,
			BucketName:        cfg.ObjectStorage.BucketName,
			ObjectPath:        cfg.ObjectStorage.ObjectPath,
			EndpointURL:       cfg.ObjectStorage.EndpointURL,
			PublicEndpointURL: cfg.ObjectStorage.PublicEndpointURL,
		})
		if err != nil {
			l.DPanicw("Failed to create R2 client", "error", err)
			return
		}
		app.ObjectStorage = s3
		l.Infow("Using R2 client for object storage", "bucket", cfg.ObjectStorage.BucketName)
	}

	var headers map[string]string
	if cfg.Headers != nil {
		headers = cfg.Headers
	}

	var client *api.Client

	timeout := 60 * time.Second
	if cfg.Network.Timeout > 0 {
		timeout = time.Duration(cfg.Network.Timeout) * time.Second
	}

	if cfg.Network.Proxy != "" {
		client, err = api.NewClientProxy(cfg.GraphqlEndpoint, headers, &http.Client{
			Timeout: timeout,
		}, cfg.Network.Proxy)
		if err != nil {
			l.DPanicw("Failed to create client with proxy", "error", err)
			return
		}

		proxyTrace, err := cloudflare.GetTraceProxy(&http.Client{}, cfg.Network.Proxy)
		if err != nil {
			l.DPanicw("Failed to get trace with proxy", "error", err)
			return
		}
		if proxyTrace.Ip == trace.Ip {
			l.Warnw("Proxy IP matches public IP, this may not be a valid proxy", "proxy", cfg.Network.Proxy, "public_ip", trace.Ip)
		}

		l.Infow("Using proxy for API requests", "proxy", cfg.Network.Proxy, "public_ip", trace.Ip)

	} else {
		client, err = api.NewClientCustomHTTP(cfg.GraphqlEndpoint, headers, &http.Client{
			Timeout: timeout,
		})
		if err != nil {
			l.DPanicw("Failed to create custom HTTP client", "error", err)
			return
		}
	}

	upstream := &api.Upstream{
		Client: client,
	}

	switch cfg.Cron.Mode {
	case "fix":
		l.Infow("Starting fix cron job")
		FixCron(&app, upstream)
	case "window":
		l.Infow("Starting window cron job")
		WindowCron(&app, upstream)
	}
}

func FixCron(app *config.App, upstream *api.Upstream) {
	var multiplier time.Duration
	switch app.Cfg.Cron.Duration {
	case config.Second:
		multiplier = time.Second
	case config.Minute:
		multiplier = time.Minute
	case config.Hour:
		multiplier = time.Hour
	default:
		multiplier = time.Second
	}

	interval := time.Duration(app.Cfg.Cron.Fix.Interval) * multiplier

	for {
		l := log.New("FixCron")
		l.Infow("Starting scheduled fetch/upload cycle")
		err := Task(app, upstream)
		if err != nil {
			l.Errorw("Failed to complete scheduled task", "error", err)
		} else {
			l.Infow("Scheduled task completed successfully")
		}

		l.Infow("Sleeping until next scheduled run", "interval", interval.Seconds(), "date", time.Now().Add(interval).Format(time.RFC3339))
		time.Sleep(interval)
	}
}

func WindowCron(app *config.App, upstream *api.Upstream) {
	var multiplier time.Duration
	switch app.Cfg.Cron.Duration {
	case config.Second:
		multiplier = time.Second
	case config.Minute:
		multiplier = time.Minute
	case config.Hour:
		multiplier = time.Hour
	default:
		multiplier = time.Second
	}

	min := app.Cfg.Cron.Window.Min
	max := app.Cfg.Cron.Window.Max

	for {
		l := log.New("WindowCron")
		l.Infow("Starting scheduled fetch/upload cycle")
		err := Task(app, upstream)
		if err != nil {
			l.Errorw("Failed to complete scheduled task", "error", err)
		} else {
			l.Infow("Scheduled task completed successfully")
		}

		interval := time.Duration(rand.Intn(max-min+1)+min) * multiplier

		l.Infow("Sleeping until next scheduled run", "interval", interval.Seconds(), "date", time.Now().Add(interval).Format(time.RFC3339))
		time.Sleep(interval)
	}
}

func Task(app *config.App, upstream *api.Upstream) error {
	data, err := upstream.Fetch()
	if err != nil {
		return err
	}

	archiveName := app.Cfg.Output.NamePrefix + "_" + time.Now().Format(time.RFC3339) + ".json"

	data.Source = app.Cfg.Source
	data.Source.DirectLink = data.Source.Latest + archiveName
	data.Source.Latest += app.Cfg.Output.NamePrefix + ".json"

	raw, err := data.Json()
	if err != nil {
		return err
	}

	if app.ObjectStorage.BucketName != "" {
		var payload []byte
		var compressionMime string = ""
		switch app.Cfg.ObjectStorage.Compression {
		case config.Brotli:
			brotli := compression.Brotli{
				Level: compression.BrotliBestCompression,
			}
			compressedData, err := brotli.Compress(raw)
			if err != nil {
				return err
			}
			payload = compressedData
			compressionMime = "br"

		case config.Gzip:
			gzip := compression.Gzip{
				Level: compression.BestCompression,
			}
			compressedData, err := gzip.Compress(raw)
			if err != nil {
				return err
			}
			payload = compressedData
			compressionMime = "gzip"

		case config.Zstd:
			zstd := compression.Zstd{
				Level: compression.ZstdSpeedBestCompression,
			}
			compressedData, err := zstd.Compress(raw)
			if err != nil {
				return err
			}
			payload = compressedData
			compressionMime = "zstd"
		default:
			payload = raw
		}

		_, err = app.ObjectStorage.UploadFile(app.Cfg.Output.NamePrefix+".json", payload, "application/json", compressionMime)
		if err != nil {
			return err
		}

		if app.Cfg.Output.Archive {
			archiveName := app.Cfg.Output.NamePrefix + "_" + time.Now().Format(time.RFC3339) + ".json"
			_, err = app.ObjectStorage.UploadFile(archiveName, payload, "application/json", compressionMime)
			if err != nil {
				return err
			}
		}
	}

	if app.Cfg.File.Path != "" {
		filePath := app.Cfg.File.Path + "/" + app.Cfg.Output.NamePrefix + ".json"
		err = os.WriteFile(filePath, raw, 0600)
		if err != nil {
			return err
		}

		if app.Cfg.Output.Archive {
			archivePath := app.Cfg.File.Path + "/" + app.Cfg.Output.NamePrefix + "_" + time.Now().Format(time.RFC3339) + ".json"
			err = os.WriteFile(archivePath, raw, 0600)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
