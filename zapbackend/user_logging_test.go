package zapbackend_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
	"github.com/woodensquares/greenery/testhelper"
	"github.com/woodensquares/greenery/zapbackend"
)

var testOptions = testhelper.TestRunnerOptions{
	ConfigGen:                testhelper.NewSimpleConfig,
	StructuredLogger:         zapbackend.StructuredLogger,
	PrettyLogger:             zapbackend.PrettyLogger,
	TraceLogger:              zapbackend.TraceLogger,
	OverrideStructuredLogger: true,
	OverridePrettyLogger:     true,
	OverrideTraceLogger:      true,
	OverrideBuiltinHandlers:  true,
	BuiltinHandlers:          map[string]greenery.Handler{"version": testhelper.LoggingTester},
}

func TestTracing(t *testing.T) {
	// With tracing there will be a lot of extra internal output, just check
	// we can see the custom message

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "tracing",
			CmdLine: []string{
				"--trace",
				"version",
			},

			NoValidateConfigValues:  true,
			OutStdOutRegex:          "TRACE.*Tracing test message",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{"version": func(cfg greenery.Config, args []string) error {
				cfg.Trace("Tracing test message")
				return nil
			},
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testOptions)
	require.NoError(t, err)
}

func TestPretty(t *testing.T) {
	// Pretty is the same for the base logger, just verify it works

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "pretty",
			CmdLine: []string{
				"--log-level",
				"debug",
				"-v",
				"3",
				"--pretty",
				"version",
			},

			NoValidateConfigValues:  true,
			OutLogRegex:             "INFO.*Pretty test message",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{"version": func(cfg greenery.Config, args []string) error {
				cfg.Infov("Pretty test message")
				return nil
			},
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testOptions)
	require.NoError(t, err)
}

func TestZapLogging(t *testing.T) {
	// Easiest to just override version to have the logging helper called.

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "debug 3",
			CmdLine: []string{
				"--log-level",
				"debug",
				"-v",
				"3",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.debug3.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "debug 2",
			CmdLine: []string{
				"--log-level",
				"debug",
				"-v",
				"2",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.debug2.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "debug 1",
			CmdLine: []string{
				"--log-level",
				"debug",
				"-v",
				"1",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.debug1.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "debug 0",
			CmdLine: []string{
				"--log-level",
				"debug",
				"-v",
				"0",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.debug0.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "info 3",
			CmdLine: []string{
				"--log-level",
				"info",
				"-v",
				"3",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.info3.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "info 2",
			CmdLine: []string{
				"--log-level",
				"info",
				"-v",
				"2",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.info2.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "info 1",
			CmdLine: []string{
				"--log-level",
				"info",
				"-v",
				"1",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.info1.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "info 0",
			CmdLine: []string{
				"--log-level",
				"info",
				"-v",
				"0",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.info0.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "warn 3",
			CmdLine: []string{
				"--log-level",
				"warn",
				"-v",
				"3",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.warn3.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "warn 2",
			CmdLine: []string{
				"--log-level",
				"warn",
				"-v",
				"2",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.warn2.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "warn 1",
			CmdLine: []string{
				"--log-level",
				"warn",
				"-v",
				"1",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.warn1.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "warn 0",
			CmdLine: []string{
				"--log-level",
				"warn",
				"-v",
				"0",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.warn0.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "error 3",
			CmdLine: []string{
				"--log-level",
				"error",
				"-v",
				"3",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.error3.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "error 2",
			CmdLine: []string{
				"--log-level",
				"error",
				"-v",
				"2",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.error2.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "error 1",
			CmdLine: []string{
				"--log-level",
				"error",
				"-v",
				"1",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.error1.log"),
			Parallel:               true,
		},
		testhelper.TestCase{
			Name: "error 0",
			CmdLine: []string{
				"--log-level",
				"error",
				"-v",
				"0",
				"version",
			},

			NoValidateConfigValues: true,
			GoldLog:                filepath.Join("testdata", "logging_test.TestZapLogging.error0.log"),
			Parallel:               true,
		},
	}

	err := testhelper.RunTestCases(t, tcs, testOptions)
	require.NoError(t, err)
}
