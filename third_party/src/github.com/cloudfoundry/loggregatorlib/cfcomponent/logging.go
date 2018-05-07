package cfcomponent

import (
	"github.com/cloudfoundry/gosteno"
	"os"
	"runtime/pprof"
	"strings"
)

var Logger *gosteno.Logger

func NewLogger(verbose bool, logFilePath, name string, config Config) *gosteno.Logger {
	level := gosteno.LOG_INFO

	if verbose {
		level = gosteno.LOG_DEBUG
	}

	loggingConfig := &gosteno.Config{
		Sinks:     make([]gosteno.Sink, 1),
		Level:     level,
		Codec:     gosteno.NewJsonCodec(),
		EnableLOC: true}

	if strings.TrimSpace(logFilePath) == "" {
		loggingConfig.Sinks[0] = gosteno.NewIOSink(os.Stdout)
	} else {
		loggingConfig.Sinks[0] = gosteno.NewFileSink(logFilePath)
	}

	if config.Syslog != "" {
		loggingConfig.Sinks = append(loggingConfig.Sinks, GetNewSyslogSink(config.Syslog))
	}

	gosteno.Init(loggingConfig)
	logger := gosteno.NewLogger(name)
	logger.Debugf("Component %s in debug mode!", name)
	setGlobalLogger(logger)
	return logger
}

func setGlobalLogger(logger *gosteno.Logger) {
	Logger = logger
}

func DumpGoRoutine() {
	goRoutineProfiles := pprof.Lookup("goroutine")
	goRoutineProfiles.WriteTo(os.Stdout, 2)
}
