// +build windows

package cfcomponent

import (
	"github.com/cloudfoundry/gosteno"
	"os"
)

func GetNewSyslogSink(namespace string) gosteno.Sink {
	panic("Syslog is not supported on windows")
	return nil
}

func RegisterGoRoutineDumpSignalChannel() chan os.Signal {
	threadDumpChan := make(chan os.Signal)

	return threadDumpChan
}
