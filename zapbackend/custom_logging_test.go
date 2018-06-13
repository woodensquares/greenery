package zapbackend_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
	"github.com/woodensquares/greenery/testhelper"
	"github.com/woodensquares/greenery/zapbackend"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var testCustomOptions = testhelper.TestRunnerOptions{
	ConfigGen:                testhelper.NewSimpleConfig,
	StructuredLogger:         zapbackend.StructuredLogger,
	PrettyLogger:             zapbackend.PrettyLogger,
	TraceLogger:              zapbackend.TraceLogger,
	OverrideStructuredLogger: true,
	OverridePrettyLogger:     true,
	OverrideTraceLogger:      true,
	OverrideBuiltinHandlers:  true,
	BuiltinHandlers:          map[string]greenery.Handler{"version": customLog},
}

func customLog(cfg greenery.Config, args []string) error {
	log := cfg.GetLogger()
	lfields := []zapcore.Field{
		zap.String("string", "hi there"),
		zap.Int("int", 3),
		zap.Duration("duration", time.Second),
	}

	log.Custom("Debug", zapbackend.Debug, lfields)
	log.Custom("Info", zapbackend.Info, lfields)
	log.Custom("Warn", zapbackend.Warn, lfields)
	log.Custom("Error", zapbackend.Error, lfields)
	return nil
}

func TestZapCustomLogging(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "debug",
			CmdLine: []string{
				"--log-level",
				"debug",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapCustomLogging.debug.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "info",
			CmdLine: []string{
				"--log-level",
				"info",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapCustomLogging.info.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "warn",
			CmdLine: []string{
				"--log-level",
				"warn",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapCustomLogging.warn.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "error",
			CmdLine: []string{
				"--log-level",
				"error",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapCustomLogging.error.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "bad level",
			CmdLine: []string{
				"version",
			},

			NoValidateConfigValues:  true,
			Parallel:                true,
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"version": func(cfg greenery.Config, args []string) error {
					log := cfg.GetLogger()
					log.Custom("Bad level", 34, []zapcore.Field{zap.Duration("time", time.Second)})
					return nil
				},
			},
			GoldLog: filepath.Join("testdata", "logging_test.TestZapCustomLogging.badlevel.log"),
		},
		testhelper.TestCase{
			Name: "bad data",
			CmdLine: []string{
				"version",
			},

			NoValidateConfigValues:  true,
			Parallel:                true,
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"version": func(cfg greenery.Config, args []string) error {
					log := cfg.GetLogger()
					log.Custom("Bad level", zapbackend.Debug, []int{3, 4})
					return nil
				},
			},
			GoldLog: filepath.Join("testdata", "logging_test.TestZapCustomLogging.baddata.log"),
		},
	}

	err := testhelper.RunTestCases(t, tcs, testCustomOptions)
	require.NoError(t, err)
}
