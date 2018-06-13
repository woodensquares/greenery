package greenery_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
	"github.com/woodensquares/greenery/testhelper"
)

// Easiest to just override version to have the logging helper called.

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
			OutStdOutRegex:          "TRACE: Tracing test message",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{"version": func(cfg greenery.Config, args []string) error {
				cfg.Trace("Tracing test message")
				return nil
			},
			},
		},
		testhelper.TestCase{
			Name: "trace on off",
			CmdLine: []string{
				"version",
			},

			NoValidateConfigValues: true,
			OutStdOut: `TRACE: Tracing on
TRACE: Tracing on
TRACE: Tracing on 2
TRACE: Tracing on 2
`,
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{"version": func(cfg greenery.Config, args []string) error {
				cfg.Trace("Tracing off")
				cfg.StartTracing()
				cfg.Trace("Tracing on")
				cfg.Tracef("Tracing %s", "on")
				cfg.StopTracing()
				cfg.Trace("Tracing off")
				cfg.StartTracing()
				cfg.Trace("Tracing on 2")
				cfg.Tracef("Tracing %s 2", "on")
				cfg.StopTracing()
				cfg.StopTracing()
				cfg.Trace("Tracing off")
				return nil
			},
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
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
			OutLogRegex:             "Pretty test message",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{"version": func(cfg greenery.Config, args []string) error {
				cfg.Infov("Pretty test message")
				return nil
			},
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}

func TestBaseLogging(t *testing.T) {
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.debug3.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.debug2.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.debug1.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.debug0.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.info3.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.info2.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.info1.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.info0.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.warn3.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.warn2.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.warn1.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.warn0.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.error3.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.error2.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.error1.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			GoldLog:                 filepath.Join("testdata", "logging_test.TestBaseLogging.error0.log"),
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}

func TestNoLogging(t *testing.T) {
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
			OutLogRegex:             "^$",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{"version": func(cfg greenery.Config, args []string) error {
				cfg.Infov("Pretty test message")
				return nil
			},
			},
		},
		testhelper.TestCase{
			Name: "debug 3",
			CmdLine: []string{
				"--log-level",
				"debug",
				"-v",
				"3",
				"version",
			},

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
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

			NoValidateConfigValues:  true,
			OutLogAllLines:          true,
			OutLogLines:             []testhelper.LogLine{},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.LoggingTester},
			Parallel:                true,
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:                testhelper.NewSimpleConfig,
		StructuredLogger:         nil,
		PrettyLogger:             nil,
		TraceLogger:              nil,
		OverrideStructuredLogger: true,
		OverridePrettyLogger:     true,
		OverrideTraceLogger:      true,
	})
	require.NoError(t, err)
}
