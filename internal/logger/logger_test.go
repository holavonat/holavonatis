package logger_test

import (
	"testing"

	cfg "github.com/holavonat/holavonatis/internal/config"
	log "github.com/holavonat/holavonatis/internal/logger"
	r "github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	config, err := cfg.GetConfig()

	r.NoError(t, err)

	log.Init(
		log.Config{
			DevelopmerMode: config.Log.DevelopmerMode,
			Level:          config.Log.Level,
			FileLog: log.FileLog{
				Filename:   config.Log.FileLog.Filename,
				MaxSize:    config.Log.FileLog.MaxSize,
				MaxBackups: config.Log.FileLog.MaxBackups,
				MaxAge:     config.Log.FileLog.MaxAge,
				Localtime:  config.Log.FileLog.Localtime,
				Compress:   config.Log.FileLog.Compress,
			},
		},
	)
	log.New("test").Info("hello")
}
