package negroni_logger_middleware

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"text/template"
	"time"
)

type LoggerEntry struct {
	StartTime string
	Status    int
	Duration  time.Duration
	Hostname  string
	Method    string
	Path      string
	Request   *http.Request
}

type Logger struct {
	ALogger  *logrus.Logger
	template *template.Template
}

var LoggerDefaultFormat = "{{.StartTime}} | {{.Status}} | \t {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}}"

// NewLogger returns a new Logger instance
func NewLogger(l *logrus.Logger) *Logger {
	logger := &Logger{ALogger: l}
	logger.SetFormat(LoggerDefaultFormat)
	return logger
}

func (l *Logger) SetFormat(format string) {
	l.template = template.Must(template.New("negroni_parser").Parse(format))
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	log := LoggerEntry{
		StartTime: start.Format(time.RFC3339),
		Status:    res.Status(),
		Duration:  time.Since(start),
		Hostname:  r.Host,
		Method:    r.Method,
		Path:      r.URL.Path,
		Request:   r,
	}

	buff := &bytes.Buffer{}
	l.template.Execute(buff, log)
	l.ALogger.Infoln(buff.String())
}
