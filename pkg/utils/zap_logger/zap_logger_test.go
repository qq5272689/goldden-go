package zap_logger

import (
	"go.uber.org/zap"
	"testing"
	//"time"
)

func TestGetLogger(t *testing.T) {
	devlogger, closer, err := GetDevLogger("/tmp", "testdev", "M")
	if err != nil {
		t.Error(err)
	}
	devlogger.Error("log", zap.Bool("debug", true), zap.Int("d", 3))
	defer closer()

	prodlogger, prodcloser, err := GetProdLogger("/tmp", "testprod", "M")
	if err != nil {
		t.Error(err)
	}
	prodlogger.Error("log", zap.Bool("debug", true), zap.Int("d", 3))
	defer prodcloser()
}
