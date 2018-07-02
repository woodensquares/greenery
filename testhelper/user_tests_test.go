package testhelper_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
	"github.com/woodensquares/greenery/testhelper"
	"github.com/woodensquares/greenery/zapbackend"
)

func TestHandlerMap(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "global override",
			CmdLine: []string{
				"int",
			},
			ConfigGen:           testhelper.NewExtraConfig,
			OverrideUserDocList: true,
			UserDocList: map[string]*greenery.DocSet{
				"": &testhelper.ExtraDocs},

			NoValidateConfigValues: true,
			OverrideHandlerMap:     true,
			HandlerMap: map[string]greenery.Handler{
				"float": func(cfg greenery.Config, args []string) error {
					fmt.Println("local overridden float called")
					return nil
				},
			},
			OutStdOut: "global overridden int called\n",
		},
		testhelper.TestCase{
			Name: "local override",
			CmdLine: []string{
				"float",
			},
			ConfigGen:           testhelper.NewExtraConfig,
			OverrideUserDocList: true,
			UserDocList: map[string]*greenery.DocSet{
				"": &testhelper.ExtraDocs},

			NoValidateConfigValues: true,
			OverrideHandlerMap:     true,
			HandlerMap: map[string]greenery.Handler{
				"float": func(cfg greenery.Config, args []string) error {
					fmt.Println("local overridden float called")
					fmt.Fprintf(os.Stderr, "hi there")
					return nil
				},
			},
			OutStdOut:      "local overridden float called\n",
			OutStdErrRegex: "there",
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		OverrideHandlerMap: true,
		HandlerMap: map[string]greenery.Handler{
			"int": func(cfg greenery.Config, args []string) error {
				fmt.Println("global overridden int called")
				return nil
			},
			"float": func(cfg greenery.Config, args []string) error {
				fmt.Println("global overridden float called")
				return nil
			},
		},
	})
	require.NoError(t, err)
}

func TestMiscFile(t *testing.T) {
	tf, err := ioutil.TempFile("", "something")
	require.NoError(t, err)
	defer func() {
		_ = tf.Close()
	}()

	fname := tf.Name()
	_, err = fmt.Fprintf(tf, "hi %s", fname)
	require.NoError(t, err)
	require.NoError(t, tf.Close())

	// Check the file is there
	_, err = ioutil.ReadFile(fname)
	require.NoError(t, err)

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "check for the file being removed",
			CmdLine: []string{
				"version",
			},
			OutStdOut:              "0.0\n",
			NoValidateConfigValues: true,
			RealFilesystem:         true,
			RemoveFiles:            []string{fname},
		},
		testhelper.TestCase{
			Name: "check for the file being created",
			CmdLine: []string{
				"version",
			},
			NoValidateConfigValues: true,
			PrecreateFiles: []testhelper.TestFile{
				testhelper.TestFile{Contents: []byte("123"), Location: "/tmp/whatever"},
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{"version": func(cfg greenery.Config, args []string) error {
				af := cfg.GetFs()
				c, lerr := afero.ReadFile(af, "/tmp/whatever")
				if lerr != nil {
					return lerr
				}

				if len(c) != 3 || c[0] != '1' || c[1] != '2' || c[2] != '3' {
					return fmt.Errorf("no match")
				}

				fmt.Println("ok")
				return nil
			}},
			OutStdOut: "ok\n",
		},
	}

	err = testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)

	// Check the file is gone
	_, err = os.Stat(fname)
	require.Error(t, err)
}

// These tests are a representative subset of the general greenery tests to
// exercise testhelper functionality, with some additional modifications to
// increase coverage

type cmdsBaseConfig struct {
	*greenery.BaseConfig
}

var cmdTestName = "cmds_test"

func newCmdsBaseConfig() greenery.Config {
	cfg := &cmdsBaseConfig{
		BaseConfig: greenery.NewBaseConfig(cmdTestName, nil),
	}

	return cfg
}

func TestStdOutRegex(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "stdout regex",
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
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}

func TestLogVarious(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "out log regex",
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
		testhelper.TestCase{
			Name: "out log lines 1",
			CmdLine: []string{
				"--log-level",
				"debug",
				"-v",
				"3",
				"version",
			},

			NoValidateConfigValues: true,
			OutLogAllLines:         false,
			OutLogLines: []testhelper.LogLine{
				testhelper.LogLine{Level: "info", Msg: "random test message"},
			},
			OverrideStructuredLogger: true,
			OverridePrettyLogger:     true,
			OverrideTraceLogger:      true,
			StructuredLogger:         zapbackend.StructuredLogger,
			PrettyLogger:             zapbackend.PrettyLogger,
			TraceLogger:              zapbackend.TraceLogger,
			OverrideBuiltinHandlers:  true,
			BuiltinHandlers: map[string]greenery.Handler{"version": func(cfg greenery.Config, args []string) error {
				cfg.Infov("random test message")
				return nil
			},
			},
		},
		testhelper.TestCase{
			Name: "out log lines 2",
			CmdLine: []string{
				"--log-level",
				"debug",
				"-v",
				"3",
				"version",
			},

			NoValidateConfigValues: true,
			OutLogAllLines:         true,
			OutLogLines: []testhelper.LogLine{
				testhelper.LogLine{Level: "info", MsgRegex: "random.*message"},
				testhelper.LogLine{Level: "info", Msg: "structured test message", CustomRegex: map[string]string{"custom": "va.*e"}},
				testhelper.LogLine{Level: "info", Msg: "structured test message 2", Custom: map[string]interface{}{"bool": true}},
			},
			OverrideStructuredLogger: true,
			OverridePrettyLogger:     true,
			OverrideTraceLogger:      true,
			StructuredLogger:         zapbackend.StructuredLogger,
			PrettyLogger:             zapbackend.PrettyLogger,
			TraceLogger:              zapbackend.TraceLogger,
			OverrideBuiltinHandlers:  true,
			BuiltinHandlers: map[string]greenery.Handler{"version": func(cfg greenery.Config, args []string) error {
				cfg.Infov("random test message")
				cfg.Infovs("structured test message", cfg.LogString("custom", "value"))
				cfg.Infovs("structured test message 2", cfg.LogGeneric("bool", true))
				return nil
			},
			},
		},
		testhelper.TestCase{
			Name: "out log lines 3",
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
			BuiltinHandlers: map[string]greenery.Handler{"version": func(cfg greenery.Config, args []string) error {
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

func TestGoldLog(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "gold log",
			CmdLine: []string{
				"--log-level",
				"debug",
				"-v",
				"3",
				"version",
			},

			NoValidateConfigValues:   true,
			OverrideStructuredLogger: true,
			OverridePrettyLogger:     true,
			OverrideTraceLogger:      true,
			StructuredLogger:         greenery.BaseStructuredLogger,
			PrettyLogger:             greenery.BasePrettyLogger,
			TraceLogger:              greenery.BaseTraceLogger,
			GoldLog:                  filepath.Join("testdata", "logging_test.TestBaseLogging.debug3.log"),
			Parallel:                 true,
		},
		testhelper.TestCase{
			Name: "gold log structured skip ts",
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
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:                testhelper.NewSimpleConfig,
		OverrideStructuredLogger: true,
		OverridePrettyLogger:     true,
		OverrideTraceLogger:      true,
		StructuredLogger:         zapbackend.StructuredLogger,
		PrettyLogger:             zapbackend.PrettyLogger,
		TraceLogger:              zapbackend.TraceLogger,
		OverrideBuiltinHandlers:  true,
		BuiltinHandlers:          map[string]greenery.Handler{"version": testhelper.LoggingTester},
	})
	require.NoError(t, err)
}

func TestRealFilesystem(t *testing.T) {
	goldDefault := filepath.Join("testdata", cmdTestName+".TestConfig.defaultcfg")

	cwd, err := os.Getwd()
	require.NoError(t, err)
	cwdConf := filepath.Join(cwd, cmdTestName+".toml")

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "real filesystem 1 create and cleanup",
			CmdLine: []string{
				"config",
				"init",
			},
			GoldFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: cmdTestName + ".toml", Source: goldDefault, Perms: 0644,
					Custom: testhelper.CompareIgnoreTmp},
			},
			OutStdOutRegex: "^Configuration file generated at " + cwdConf + "\n$",
			RealFilesystem: true,
		},

		testhelper.TestCase{
			Name: "real filesystem 2 fstat",
			CmdLine: []string{
				"config",
				"init",
			},
			PrecreateFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: cmdTestName + ".toml", Source: goldDefault, Perms: 0644},
			},
			ExecError:      "Cannot create config file",
			RealFilesystem: true,
		},
		testhelper.TestCase{
			Name: "real filesystem 3 overwrite",
			CmdLine: []string{
				"config",
				"init",
				"--force",
			},
			PrecreateFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: cmdTestName + ".toml", Source: goldDefault, Perms: 0644},
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"CfgForce":    testhelper.Comparer{Value: true},
				"CfgLocation": testhelper.Comparer{DefaultValue: true},
			},
			OutStdOutRegex: "^Configuration file generated at " + cwdConf + "\n$",
			RealFilesystem: true,
		},
	}

	err = testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: newCmdsBaseConfig,
		UserDocList: map[string]*greenery.DocSet{
			"": &greenery.DocSet{
				ConfigFile: map[string]string{
					greenery.DocConfigHeader: "Config generated while testing",
				},
			},
		}})

	require.NoError(t, err)
}

func TestAllCompare(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Compare output and custom vars/parser/comparer",
			CmdLine: []string{
				"config",
				"display",
			},
			CustomVars:   testhelper.ExtraConfigCustomVars,
			CustomParser: testhelper.ExtraConfigCustomParse,
			CfgFile:      filepath.Join("testdata", cmdTestName+".TestConfig.extracustomcfg"),
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfig.displayextracustomcfg"),
				Custom: testhelper.CompareIgnoreTmp},
			ExpectedValues: map[string]testhelper.Comparer{
				"LogLevel": testhelper.Comparer{Value: "info", Accessor: "GetTyped"},
				"NoEnv":    testhelper.Comparer{Value: true},
				"Pretty": testhelper.Comparer{Value: true, Compare: func(t *testing.T, name string, expected, actual interface{}) error {
					exbool := expected.(bool)
					acbool := actual.(bool)
					if exbool != acbool {
						return fmt.Errorf("%v is different from expected %v for variable %s", acbool, exbool, name)
					}
					return nil
				}},
				"Verbosity": testhelper.Comparer{Value: 2, Accessor: "GetTyped"},
				"Bool":      testhelper.Comparer{Value: false},
				"NameValue": testhelper.Comparer{Value: [][]string{
					[]string{"j1", "v1"},
					[]string{"j2", "v2"},
					[]string{"j3", "v3"},
				}},
				"NameEnum": testhelper.Comparer{Value: [][]string{
					[]string{"j1", "a"},
					[]string{"j2", "b"},
					[]string{"j3", "c"},
				}},
				"FlagIP":      testhelper.Comparer{Value: "192.168.1.1", Accessor: "GetTyped"},
				"FlagPort":    testhelper.Comparer{Value: greenery.PortValue(443)},
				"FlagCString": testhelper.Comparer{Value: "HIHIHI", Compare: testhelper.CompareValueToGetter},
				"FlagEnum":    testhelper.Comparer{Value: "c", Accessor: "GetTyped"},
				"FlagInt":     testhelper.Comparer{Value: 500, Accessor: "GetTyped"},
				"Float32":     testhelper.Comparer{Value: float32(12.34)},
				"Float64":     testhelper.Comparer{Value: float64(-56.78)},
				"Int":         testhelper.Comparer{Value: -11},
				"Int8":        testhelper.Comparer{Value: int8(-12)},
				"Int16":       testhelper.Comparer{Value: int16(-13)},
				"Int32":       testhelper.Comparer{Value: int32(-14)},
				"Int64":       testhelper.Comparer{Value: int64(-15)},
				"Uint":        testhelper.Comparer{Value: uint(11)},
				"Uint8":       testhelper.Comparer{Value: uint8(12)},
				"Uint16":      testhelper.Comparer{Value: uint16(13)},
				"Uint32":      testhelper.Comparer{Value: uint32(14)},
				"Uint64":      testhelper.Comparer{Value: uint64(15)},
				"String":      testhelper.Comparer{Value: "other"},
				"Time":        testhelper.Comparer{Value: time.Date(2018, time.June, 3, 12, 8, 32, 454, time.UTC)},
				"Duration":    testhelper.Comparer{Value: time.Hour*48 + time.Minute*16 + time.Second*32 + time.Millisecond*45},

				"SliceString": testhelper.Comparer{Value: []string{"first", "second", "third"}},
				"SliceInt":    testhelper.Comparer{Value: []int{1, 2, 3, 4}},
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:  testhelper.NewExtraConfig,
		CompareMap: testhelper.ExtraConfigCustomComparers,
		UserDocList: map[string]*greenery.DocSet{
			"": &testhelper.ExtraDocs}})
	require.NoError(t, err)
}

func TestValuesValidator(t *testing.T) {
	rootHelp := filepath.Join("testdata", "language_test.common.rootHelp.stdout")
	errHelp := filepath.Join("testdata", "empty.file")

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "values validator",
			ValuesValidator: func(t *testing.T, icfg greenery.Config) {
				lang, _ := icfg.GetDocs()
				require.EqualValues(t, ",en,", lang)
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelp},
			GoldStdErr: &testhelper.TestFile{Source: errHelp},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}

func TestConfigDefaults(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "config defaults",
			CmdLine: []string{
				"version",
			},
			OutStdOut: "1.2.3a-testing\n",
			ConfigDefaults: &greenery.BaseConfigOptions{
				VersionMajor:      "1",
				VersionMinor:      "2",
				VersionPatchlevel: "3a-testing",
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: newCmdsBaseConfig,
	})

	require.NoError(t, err)
}

func TestEnv(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "env setting",
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
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}

func TestCmdLineCfg(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "cmdlinecfg",
			CmdLine: []string{
				"bool",
			},
			CfgContents:    "invalid yadda yadda 1234\n\n\nsgfgfs fg = 453",
			CmdlineCfgName: "/tmp/nonexistent",
			ExecError:      "Could not load",
		},
		testhelper.TestCase{
			Name: "cmdlinecfg",
			CmdLine: []string{
				"bool",
			},
			CfgContents:    "invalid yadda yadda 1234\n\n\nsgfgfs fg = 453",
			CmdlineCfgName: "/tmp/nonexistent",
			ExecErrorRegex: "Could.*load",
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:  testhelper.NewExtraConfig,
		CompareMap: testhelper.ExtraConfigCustomComparers,
		UserDocList: map[string]*greenery.DocSet{
			"": &testhelper.ExtraDocs}})
	require.NoError(t, err)
}

func TestErrorOutput(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "fail pre-exec handler",
			CmdLine: []string{
				"version",
			},
			NoValidateConfigValues:  true,
			ExecError:               "pre-exec fail",
			ExecErrorOutput:         true,
			OutStdOut:               "",
			OutStdErrRegex:          "Usage:",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"-pre-exec-handler": func(cfg greenery.Config, args []string) error {
					return fmt.Errorf("pre-exec fail")
				},
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}
