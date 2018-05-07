package paas_lager

import (
	"fmt"
	"log"
	"os"
	"strings"

	"paas_lager/lager"
	"paas_lager/syslog"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
)

type Config struct {
	LoggerLevel string
	LoggerFile  string

	EnableRsyslog  bool
	RsyslogNetwork string
	RsyslogAddr    string

	LogFormatText bool
}

var config *Config = DefaultConfig()

func DefaultConfig() *Config {
	return &Config{
		LoggerLevel:    INFO,
		LoggerFile:     "",
		EnableRsyslog:  false,
		RsyslogNetwork: "udp",
		RsyslogAddr:    "127.0.0.1:5140",
		LogFormatText:  false,
	}
}

func Init(c Config) {
	if c.LoggerLevel != "" {
		config.LoggerLevel = c.LoggerLevel
	}

	if c.LoggerFile != "" {
		config.LoggerFile = c.LoggerFile
	}

	if c.EnableRsyslog {
		config.EnableRsyslog = c.EnableRsyslog
	}

	if c.RsyslogNetwork != "" {
		config.RsyslogNetwork = c.RsyslogNetwork
	}

	if c.RsyslogAddr != "" {
		config.RsyslogAddr = c.RsyslogAddr
	}

	config.LogFormatText = c.LogFormatText
}

func NewLogger(component string) lager.Logger {
	return NewLoggerExt(component, component)
}

func NewLoggerExt(component string, app_guid string) lager.Logger {
	var lagerLogLevel lager.LogLevel
	switch strings.ToUpper(config.LoggerLevel) {
	case DEBUG:
		lagerLogLevel = lager.DEBUG
	case INFO:
		lagerLogLevel = lager.INFO
	case WARN:
		lagerLogLevel = lager.WARN
	case ERROR:
		lagerLogLevel = lager.ERROR
	case FATAL:
		lagerLogLevel = lager.FATAL
	default:
		panic(fmt.Errorf("unknown logger level: %s", config.LoggerLevel))
	}

	logger := lager.NewLoggerExt(component, config.LogFormatText)

	if config.LoggerFile == "" {
		sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), lagerLogLevel)
		logger.RegisterSink(sink)
	} else {
		file, err := os.OpenFile(config.LoggerFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}

		sink := lager.NewReconfigurableSink(lager.NewWriterSink(file, lager.DEBUG), lagerLogLevel)
		logger.RegisterSink(sink)
	}

	if config.EnableRsyslog {
		syslog, err := syslog.Dial(component, app_guid, config.RsyslogNetwork, config.RsyslogAddr)
		if err != nil {
			//warn, not panic
			log.Println(err.Error())
		} else {
			sink := lager.NewReconfigurableSink(lager.NewWriterSink(syslog, lager.DEBUG), lagerLogLevel)
			logger.RegisterSink(sink)
		}
	}

	return logger
}
