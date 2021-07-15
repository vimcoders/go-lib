package lib

import (
	"fmt"
	"io"
	"log"
	"log/syslog"

	driver "github.com/vimcoders/go-driver"
)

type Syslogger struct {
	io.Closer
	logger *log.Logger
}

func (log *Syslogger) Debug(format string, v ...interface{}) {
	log.logger.Output(2, fmt.Sprintf("[DEBUG] %v", fmt.Sprintf(format, v...)))
	fmt.Println("[DEBUG] ", fmt.Sprintf(format, v...))
}

func (log *Syslogger) Info(format string, v ...interface{}) {
	log.logger.Output(2, fmt.Sprintf("[INFO] %v", fmt.Sprintf(format, v...)))
	fmt.Println("[INFO]", "\033[32m", fmt.Sprintf(format, v...), "\033[0m")
}

func (log *Syslogger) Warning(format string, v ...interface{}) {
	log.logger.Output(2, fmt.Sprintf("[WARNING] %v", fmt.Sprintf(format, v...)))
	fmt.Println("[WARNING] ", "\033[33m", fmt.Sprintf(format, v...), "\033[0m")
}

func (log *Syslogger) Error(format string, v ...interface{}) {
	log.logger.Output(2, fmt.Sprintf("[ERROR] %v", fmt.Sprintf(format, v...)))
	fmt.Println("[Error] ", "\033[31m", fmt.Sprintf(format, v...), "\033[0m")
}

func NewSyslogger() (driver.Logger, error) {
	sysLog, err := syslog.New(syslog.LOG_NOTICE, "syslog")

	if err != nil {
		return nil, err
	}

	return &Syslogger{
		logger: log.New(sysLog, "", log.Ldate|log.Ltime|log.Lshortfile),
		Closer: sysLog,
	}, nil
}
