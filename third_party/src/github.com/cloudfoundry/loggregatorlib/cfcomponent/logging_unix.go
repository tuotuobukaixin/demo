// +build !windows,!plan9

package cfcomponent

import (
	"github.com/cloudfoundry/gosteno"
	"os"
	"os/signal"
	"syscall"
)

func GetNewSyslogSink(namespace string) *gosteno.Syslog {
	return gosteno.NewSyslogSink(namespace)
}

func RegisterGoRoutineDumpSignalChannel() chan os.Signal {
	threadDumpChan := make(chan os.Signal)
	signal.Notify(threadDumpChan, syscall.SIGUSR1)

	return threadDumpChan
}
