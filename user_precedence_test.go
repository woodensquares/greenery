package greenery_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
	"github.com/woodensquares/greenery/testhelper"
)

func TestPrecedence(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Check for precedence 1 between cmd > env > config, cmd",
			Env: map[string]string{
				"SIMPLE_VERBOSITY": "0",
				"SIMPLE_LOGLEVEL":  "debug",
			},
			CmdLine: []string{
				"version",
				"-v=2",
				"-l=info",
			},

			CfgContents: `# Set verbosity
verbosity = 3
log-level = "debug"
`,

			ExpectedValues: map[string]testhelper.Comparer{
				"Verbosity": testhelper.Comparer{Value: 2, Accessor: "GetTyped"},
				"LogLevel":  testhelper.Comparer{Value: "info", Accessor: "GetTyped"},
			},

			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.NopNoArgs},
		},

		testhelper.TestCase{
			Name: "Check for precedence 2 between cmd > env > config, env",
			Env: map[string]string{
				"SIMPLE_VERBOSITY": "2",
				"SIMPLE_LOGLEVEL":  "info",
			},
			CmdLine: []string{
				"version",
			},

			CfgContents: `# Set verbosity
verbosity = 3
log-level = "debug"
`,

			ExpectedValues: map[string]testhelper.Comparer{
				"Verbosity": testhelper.Comparer{Value: 2, Accessor: "GetTyped"},
				"LogLevel":  testhelper.Comparer{Value: "info", Accessor: "GetTyped"},
			},

			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.NopNoArgs},
		},

		testhelper.TestCase{
			Name: "Check for precedence 3 between cmd > env > config, config",
			CmdLine: []string{
				"version",
			},

			CfgContents: `# Set verbosity
verbosity = 2
log-level = "info"
`,

			ExpectedValues: map[string]testhelper.Comparer{
				"Verbosity": testhelper.Comparer{Value: 2, Accessor: "GetTyped"},
				"LogLevel":  testhelper.Comparer{Value: "info", Accessor: "GetTyped"},
			},

			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.NopNoArgs},
		},

		testhelper.TestCase{
			Name: "Check for precedence 4 between cmd > env > config + no-env, config",
			Env: map[string]string{
				"SIMPLE_VERBOSITY": "2",
				"SIMPLE_LOGLEVEL":  "info",
			},
			CmdLine: []string{
				"version",
				"--no-env",
			},

			CfgContents: `# Set verbosity
verbosity = 3
log-level = "debug"
`,

			ExpectedValues: map[string]testhelper.Comparer{
				"NoEnv":     testhelper.Comparer{Value: true},
				"Verbosity": testhelper.Comparer{Value: 3, Accessor: "GetTyped"},
				"LogLevel":  testhelper.Comparer{Value: "debug", Accessor: "GetTyped"},
			},

			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         map[string]greenery.Handler{"version": testhelper.NopNoArgs},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}

func TestPartialPrecedence(t *testing.T) {
	fmap := map[string]greenery.Handler{
		"test": testhelper.NopNoArgs,
	}

	docs := map[string]*greenery.DocSet{
		"": &greenery.DocSet{
			Usage: map[string]*greenery.CmdHelp{
				"test": &greenery.CmdHelp{
					Short: "test",
				},
			},
			CmdLine: map[string]string{
				"TestParam": "test parameter",
			},
			ConfigFile: map[string]string{
				"test.":     "test section",
				"TestParam": "test parameter",
			},
		}}

	cmdLineOnly := func() greenery.Config {
		return &struct {
			*greenery.BaseConfig
			TestParam string `greenery:"test|testParam|,,"`
		}{
			BaseConfig: greenery.NewBaseConfig("partial", fmap),
			TestParam:  "default",
		}
	}

	cmdLineEnv := func() greenery.Config {
		return &struct {
			*greenery.BaseConfig
			TestParam string `greenery:"test|testParam|, , TESTPARAM"`
		}{
			BaseConfig: greenery.NewBaseConfig("partial", fmap),
			TestParam:  "default",
		}
	}

	cmdLineConf := func() greenery.Config {
		return &struct {
			*greenery.BaseConfig
			TestParam string `greenery:"test|testParam|, test.test-param,"`
		}{
			BaseConfig: greenery.NewBaseConfig("partial", fmap),
			TestParam:  "default",
		}
	}

	confEnvConf := func() greenery.Config {
		return &struct {
			*greenery.BaseConfig
			TestParam string `greenery:"   , test.test-param, TESTPARAM"`
		}{
			BaseConfig: greenery.NewBaseConfig("partial", fmap),
			TestParam:  "default",
		}
	}

	tcs := []testhelper.TestCase{
		// Cmdline only
		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name:      "Custom only cmdline variable, default",
			ConfigGen: cmdLineOnly,
			CmdLine: []string{
				"test",
			},

			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "default"},
			},
		},
		testhelper.TestCase{
			Name:      "Custom only cmdline variable, cmdline",
			ConfigGen: cmdLineOnly,
			CmdLine: []string{
				"test",
				"--testParam",
				"cmdline",
			},

			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "cmdline"},
			},
		},

		// Cmdline + env
		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name:      "Custom cmdline and env variable, default",
			ConfigGen: cmdLineEnv,
			CmdLine: []string{
				"test",
			},

			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "default"},
			},
		},
		testhelper.TestCase{
			Name:      "Custom cmdline and env variable, cmdline",
			ConfigGen: cmdLineEnv,
			CmdLine: []string{
				"test",
				"--testParam",
				"cmdline",
			},

			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "cmdline"},
			},
		},
		testhelper.TestCase{
			Name:      "Custom cmdline and env variable, env",
			ConfigGen: cmdLineEnv,
			CmdLine: []string{
				"test",
			},
			Env: map[string]string{
				"PARTIAL_TESTPARAM": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "env"},
			},
		},
		testhelper.TestCase{
			Name:      "Custom cmdline and env variable, cmd and env",
			ConfigGen: cmdLineEnv,
			CmdLine: []string{
				"test",
				"--testParam",
				"cmdline",
			},
			Env: map[string]string{
				"PARTIAL_TESTPARAM": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "cmdline"},
			},
		},

		// Cmdline + conf
		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name:      "Custom cmdline and conf variable, default",
			ConfigGen: cmdLineConf,
			CmdLine: []string{
				"test",
			},

			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "default"},
			},
		},
		testhelper.TestCase{
			Name:      "Custom cmdline and conf variable, cmdline",
			ConfigGen: cmdLineConf,
			CmdLine: []string{
				"test",
				"--testParam",
				"cmdline",
			},

			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "cmdline"},
			},
		},
		testhelper.TestCase{
			Name:      "Custom cmdline and conf variable, conf",
			ConfigGen: cmdLineConf,
			CmdLine: []string{
				"test",
			},
			CfgContents: `# Set the test parameter
[test]
test-param = "cfg"
`,
			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "cfg"},
			},
		},
		testhelper.TestCase{
			Name:      "Custom cmdline and conf variable, cmd and conf",
			ConfigGen: cmdLineConf,
			CmdLine: []string{
				"test",
				"--testParam",
				"cmdline",
			},
			CfgContents: `# Set the test parameter
[test]
test-param = "cfg"
`,
			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "cmdline"},
			},
		},

		// Conf + env
		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name:      "Custom conf and env variable, default",
			ConfigGen: confEnvConf,
			CmdLine: []string{
				"test",
			},

			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "default"},
			},
		},
		testhelper.TestCase{
			Name:      "Custom conf and env variable, conf",
			ConfigGen: confEnvConf,
			CmdLine: []string{
				"test",
			},
			CfgContents: `# Set the test parameter
[test]
test-param = "cfg"
`,

			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "cfg"},
			},
		},
		testhelper.TestCase{
			Name:      "Custom conf and env variable, env",
			ConfigGen: confEnvConf,
			CmdLine: []string{
				"test",
			},
			Env: map[string]string{
				"PARTIAL_TESTPARAM": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "env"},
			},
		},

		testhelper.TestCase{
			Name:      "Custom conf and env variable, conf and env",
			ConfigGen: confEnvConf,
			CmdLine: []string{
				"test",
			},
			CfgContents: `# Set the test parameter
[test]
test-param = "cfg"
`,
			Env: map[string]string{
				"PARTIAL_TESTPARAM": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"TestParam": testhelper.Comparer{Value: "env"},
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:   testhelper.NewSimpleConfig,
		UserDocList: docs,
	})
	require.NoError(t, err)
}

func TestPartialPrecedenceAdditionalTypes(t *testing.T) {
	fmap := map[string]greenery.Handler{
		"test": testhelper.NopNoArgs,
	}

	docs := map[string]*greenery.DocSet{
		"": &greenery.DocSet{
			Usage: map[string]*greenery.CmdHelp{
				"test": &greenery.CmdHelp{
					Short: "test",
				},
			},
			CmdLine: map[string]string{
				"TestParam": "test parameter",
			},
			ConfigFile: map[string]string{
				"test.":     "test section",
				"TestParam": "test parameter",
			},
		}}

	confEnvConf := func() greenery.Config {
		cfg := &struct {
			*greenery.BaseConfig
			FlagIP   *greenery.IPValue  `greenery:",    flag.ip,      FLAGIP"`
			FlagPort greenery.PortValue `greenery:",    flag.port,    FLAGPORT"`
		}{
			BaseConfig: greenery.NewBaseConfig("partial", fmap),
			FlagIP:     greenery.NewDefaultIPValue("FlagIP", "127.0.0.1"),
			FlagPort:   greenery.NewPortValue(),
		}
		_ = cfg.FlagPort.SetInt(80)
		return cfg

	}

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name:      "Custom conf and env variable, env",
			ConfigGen: confEnvConf,
			CmdLine: []string{
				"test",
			},
			Env: map[string]string{
				"PARTIAL_FLAGIP":   "10.0.0.1",
				"PARTIAL_FLAGPORT": "443",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"FlagIP":   testhelper.Comparer{Value: "10.0.0.1", Accessor: "GetTyped"},
				"FlagPort": testhelper.Comparer{Value: greenery.PortValue(443)},
			},
		},
		testhelper.TestCase{
			Name:      "Custom conf and env variable, badenv",
			ConfigGen: confEnvConf,
			CmdLine: []string{
				"test",
			},
			Env: map[string]string{
				"PARTIAL_FLAGIP": "aaaaa",
			},
			ExecError: "should be an IP address",
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:   testhelper.NewSimpleConfig,
		UserDocList: docs,
	})
	require.NoError(t, err)
}
