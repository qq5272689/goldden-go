package logger

import "testing"

func TestLoggerInit(t *testing.T) {
	Debug("default logger")
	LoggerInit("local", "/tmp/testlog", "test", "")
	Debug("test logger")
	Info("test info")
	Warn("test warn")
	Error("test error")
	Closer()
}

func TestJsonLoggerInit(t *testing.T) {
	Debug("default logger")
	JsonLoggerInit("local")
	Debug("test logger")
	Info("test info")
	Warn("test warn")
	Error("test error")
	Closer()
}
