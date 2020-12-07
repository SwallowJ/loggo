package test

import (
	"loggo"
	"testing"
)

func Test_config(t *testing.T) {
	logger := loggo.New("main")

	logger.Info("Test")

	// logger.Error("HAHA")

	loggo.Debug("Test123")
	loggo.Info("aaa")
}
