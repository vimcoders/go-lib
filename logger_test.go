package lib

import "testing"

func TestSysloggerDebug(t *testing.T) {
	logger, _ := NewSyslogger()

	logger.Debug("debug")
}

func TestSysloggerInfo(t *testing.T) {
	logger, _ := NewSyslogger()

	logger.Info("info")
}

func TestSysloggerWaring(t *testing.T) {
	logger, _ := NewSyslogger()

	logger.Warning("Warning")
}

func TestSysloggerError(t *testing.T) {
	logger, _ := NewSyslogger()

	logger.Error("err")
}
