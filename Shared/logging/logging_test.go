package logging

import (
	"github.com/ARMmaster17/Captain/Shared/config"
	"testing"
)

func Test_InitLogging(t *testing.T) {
	config.InitConfiguration("shared")
	config.SetAppString("LOG_PATH", "/etc/captain/shared/test.log")
	InitLogging()
}
