package greenery_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
	"github.com/woodensquares/greenery/testhelper"
)

func TestBadConfig(t *testing.T) {
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
				"TestParam":  "test parameter",
				"TestParam1": "test parameter",
				"TestParam2": "test parameter",
			},
			ConfigFile: map[string]string{
				"test.": "test section",
			},
		}}

	// No appname
	require.PanicsWithValue(t, "A non-empty name for the application must be set when creating a configuration", func() {
		fmt.Printf("%v", &struct {
			*greenery.BaseConfig
		}{
			BaseConfig: greenery.NewBaseConfig("", nil),
		})
	})

	// Checking it this way, because testhelper does require a properly
	// initialized configuration to work.
	require.PanicsWithValue(t, "The configuration struct passed to execute was not initialized properly", func() {
		cfg := &struct {
			*greenery.BaseConfig
		}{
			BaseConfig: &greenery.BaseConfig{},
		}
		panic(cfg.Execute(cfg, nil).Error())
	})

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Conf variable, no command",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam string `greenery:"nonexistent|testParam|,,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "cannot find cmd for nonexistent",
		},
		testhelper.TestCase{
			Name: "Conf variable, empty command without none",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam string `greenery:"||,, TEST"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "has empty commandline",
		},
		testhelper.TestCase{
			Name: "Conf variable, bad tag, wrong number of parts",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam string `greenery:"test|testParam|,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "found 2 parts instead of 3",
		},
		testhelper.TestCase{
			Name: "Conf variable, bad cmdline tag, wrong number of parts",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam string `greenery:"test|testParam||,,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "malformed cmdline tag",
		},
		testhelper.TestCase{
			Name: "Conf variable, bad cfgfile tag 1",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam string `greenery:"test|testParam|, missingleading,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "Invalid config file tag",
		},
		testhelper.TestCase{
			Name: "Conf variable, bad cfgfile tag 2",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam string `greenery:"test|testParam|, .more.than.two,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "more than two . present in '.more.than.two'",
		},
		testhelper.TestCase{
			Name: "Conf variable, bad cfg tag, no .",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam string `greenery:"||none, wrong,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "Invalid config file tag",
		},
		testhelper.TestCase{
			Name: "Conf variable, unsupported type",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam complex64 `greenery:"test|testParam|,,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "Cannot create a flag for TestParam/complex64",
		},
		testhelper.TestCase{
			Name: "Conf variable, unsupported type 2, cmdline",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam []complex64 `greenery:"test|testParam|,,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "Cannot create a flag for TestParam/slice",
		},
		testhelper.TestCase{
			Name: "Conf variable, unsupported type 3, env",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam []complex64 `greenery:"||none ,,BLAH"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "Cannot create flag TestParam",
		},
		testhelper.TestCase{
			Name: "Conf variable, unsupported type 4, cfg",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam []complex128 `greenery:"||none , .blah,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "Cannot create configuration file variable TestParam",
		},
		testhelper.TestCase{
			Name: "Conf variable, duplicate binding",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam1 string `greenery:"test|testParam|,,"`
					TestParam2 string `greenery:"test|testParam|,,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "Error while creating flag for variable",
		},
		testhelper.TestCase{
			Name: "Conf variable, shadowing",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					LogLevel string `greenery:"||none,,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "Field collision on field LogLevel",
		},
		testhelper.TestCase{
			Name: "Conf variable, unsupported type, but not declared as custom so generate fail",
			CmdLine: []string{
				"config",
				"init",
			},
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam complex128 `greenery:", .configvar,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "Internal error, unsupported type complex128 for TestParam",
		},
		testhelper.TestCase{
			Name: "Conf variable, no docs",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					TestParam3 string `greenery:"test|testParam|,,"`
				}{
					BaseConfig: greenery.NewBaseConfig("partial", fmap),
				}
			},
			ExecError: "No documentation for command line parameter testParam (variable TestParam3)",
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:   testhelper.NewSimpleConfig,
		UserDocList: docs,
	})
	require.NoError(t, err)
}

func TestBadInvocation(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "root with extra args and no overrides",
			CmdLine: []string{
				"extra",
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelp},
		},
		testhelper.TestCase{
			Name: "config with extra args and no overrides",
			CmdLine: []string{
				"config",
				"extra",
			},
			GoldStdOut: &testhelper.TestFile{Source: configHelp},
		},
		testhelper.TestCase{
			Name: "config init with extra args and no overrides",
			CmdLine: []string{
				"config",
				"init",
				"extra",
			},
			ExecError: "does not support additional arguments",
		},
		testhelper.TestCase{
			Name: "config display with extra args and no overrides",
			CmdLine: []string{
				"config",
				"display",
				"extra",
			},
			ExecError: "does not support additional arguments",
		},
		testhelper.TestCase{
			Name: "config env with extra args and no overrides",
			CmdLine: []string{
				"config",
				"env",
				"extra",
			},
			ExecError: "does not support additional arguments",
		},
		testhelper.TestCase{
			Name: "help with extra args and no overrides",
			CmdLine: []string{
				"help",
				"extra",
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelp},
		},
		testhelper.TestCase{
			Name: "version with extra args and no overrides",
			CmdLine: []string{
				"version",
				"extra",
			},
			ExecError: "does not support additional arguments",
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}

func TestBadFlagValuesCmd(t *testing.T) {
	tcs := []testhelper.TestCase{
		// cannot set a bad bool via cmdline, it's either there or not

		testhelper.TestCase{
			Name: "Bad value int",
			CmdLine: []string{
				"int",
				"--int",
				"nonint",
			},
			ExecError: "invalid argument",
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:  testhelper.NewExtraConfig,
		CompareMap: testhelper.ExtraConfigCustomComparers,
		UserDocList: map[string]*greenery.DocSet{
			"": &testhelper.ExtraDocs}})
	require.NoError(t, err)
}

func TestBadFlagValuesEnv(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Bad value bool",
			Env: map[string]string{
				"EXTRA_BOOL": "nonbool",
			},
			CmdLine: []string{
				"bool",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name: "Bad value int",
			Env: map[string]string{
				"EXTRA_INT": "nonint",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name: "Bad value time",
			Env: map[string]string{
				"EXTRA_TIME": "not a time",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value via unmarshaling: parsing time",
		},
		testhelper.TestCase{
			Name: "Bad value int8",
			Env: map[string]string{
				"EXTRA_INT8": "128",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int8: 128 would overflow an int8",
		},
		testhelper.TestCase{
			Name: "Bad value int8 2",
			Env: map[string]string{
				"EXTRA_INT8": "-129",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int8: -129 would overflow an int8",
		},
		testhelper.TestCase{
			Name: "Bad value int8 3",
			Env: map[string]string{
				"EXTRA_INT8": "aaa",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int8: unable to cast \"aaa\" of type string to int8",
		},
		testhelper.TestCase{
			Name: "Max value int8",
			Env: map[string]string{
				"EXTRA_INT8": "127",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int8": testhelper.Comparer{Value: int8(127)},
			},
		},
		testhelper.TestCase{
			Name: "Max value int8 2",
			Env: map[string]string{
				"EXTRA_INT8": "-128",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int8": testhelper.Comparer{Value: int8(-128)},
			},
		},
		testhelper.TestCase{
			Name: "Bad value int16",
			Env: map[string]string{
				"EXTRA_INT16": "32768",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int16: 32768 would overflow an int16",
		},
		testhelper.TestCase{
			Name: "Bad value int16 2",
			Env: map[string]string{
				"EXTRA_INT16": "-32769",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int16: -32769 would overflow an int16",
		},
		testhelper.TestCase{
			Name: "Bad value int16 3",
			Env: map[string]string{
				"EXTRA_INT16": "aaa",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int16: unable to cast \"aaa\" of type string to int16",
		},
		testhelper.TestCase{
			Name: "Max value int16",
			Env: map[string]string{
				"EXTRA_INT16": "32767",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int16": testhelper.Comparer{Value: int16(32767)},
			},
		},
		testhelper.TestCase{
			Name: "Max value int16 2",
			Env: map[string]string{
				"EXTRA_INT16": "-32768",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int16": testhelper.Comparer{Value: int16(-32768)},
			},
		},
		testhelper.TestCase{
			Name: "Bad value int32",
			Env: map[string]string{
				"EXTRA_INT32": "2147483648",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int32: 2147483648 would overflow an int32",
		},
		testhelper.TestCase{
			Name: "Bad value int32 2",
			Env: map[string]string{
				"EXTRA_INT32": "-2147483649",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int32: -2147483649 would overflow an int32",
		},
		testhelper.TestCase{
			Name: "Bad value int32 3",
			Env: map[string]string{
				"EXTRA_INT32": "aaa",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int32: unable to cast \"aaa\" of type string to int32",
		},
		testhelper.TestCase{
			Name: "Max value int32",
			Env: map[string]string{
				"EXTRA_INT32": "2147483647",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int32": testhelper.Comparer{Value: int32(2147483647)},
			},
		},
		testhelper.TestCase{
			Name: "Max value int32 2",
			Env: map[string]string{
				"EXTRA_INT32": "-2147483648",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int32": testhelper.Comparer{Value: int32(-2147483648)},
			},
		},
		testhelper.TestCase{
			Name: "Bad value int64",
			Env: map[string]string{
				"EXTRA_INT64": "9223372036854775808",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int64: unable to cast \"9223372036854775808\" of type string to int64",
		},
		testhelper.TestCase{
			Name: "Bad value int64 2",
			Env: map[string]string{
				"EXTRA_INT64": "-9223372036854775809",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int64: unable to cast \"-9223372036854775809\" of type string to int64",
		},
		testhelper.TestCase{
			Name: "Max value int64",
			Env: map[string]string{
				"EXTRA_INT64": "9223372036854775807",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int64": testhelper.Comparer{Value: int64(9223372036854775807)},
			},
		},
		testhelper.TestCase{
			Name: "Max value int64 2",
			Env: map[string]string{
				"EXTRA_INT64": "-9223372036854775808",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int64": testhelper.Comparer{Value: int64(-9223372036854775808)},
			},
		},
		testhelper.TestCase{
			Name: "Bad value uint",
			Env: map[string]string{
				"EXTRA_UINT": "-2",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name: "Bad value uint8",
			Env: map[string]string{
				"EXTRA_UINT8": "-2",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.uint8: unable to cast \"-2\" to uint8",
		},
		testhelper.TestCase{
			Name: "Bad value uint8 2",
			Env: map[string]string{
				"EXTRA_UINT8": "256",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.uint8: unable to cast \"256\" to uint8",
		},
		testhelper.TestCase{
			Name: "Max value uint8",
			Env: map[string]string{
				"EXTRA_UINT8": "255",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Uint8": testhelper.Comparer{Value: uint8(255)},
			},
		},
		testhelper.TestCase{
			Name: "Bad value uint16",
			Env: map[string]string{
				"EXTRA_UINT16": "-2",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name: "Bad value uint16 2",
			Env: map[string]string{
				"EXTRA_UINT16": "65536",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.uint16: unable to cast \"65536\" to uint16",
		},
		testhelper.TestCase{
			Name: "Max value uint16",
			Env: map[string]string{
				"EXTRA_UINT16": "65535",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Uint16": testhelper.Comparer{Value: uint16(65535)},
			},
		},
		testhelper.TestCase{
			Name: "Bad value uint32",
			Env: map[string]string{
				"EXTRA_UINT32": "-2",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name: "Bad value uint32 2",
			Env: map[string]string{
				"EXTRA_UINT32": "4294967296",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.uint32: unable to cast \"4294967296\" to uint32",
		},
		testhelper.TestCase{
			Name: "Max value uint32",
			Env: map[string]string{
				"EXTRA_UINT32": "4294967295",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Uint32": testhelper.Comparer{Value: uint32(4294967295)},
			},
		},
		testhelper.TestCase{
			Name: "Bad value uint64",
			Env: map[string]string{
				"EXTRA_UINT64": "-2",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name: "Bad value uint64 2",
			Env: map[string]string{
				"EXTRA_UINT64": "18446744073709551616",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.uint64: unable to cast \"18446744073709551616\" to uint64",
		},
		testhelper.TestCase{
			Name: "Max value uint64",
			Env: map[string]string{
				"EXTRA_UINT64": "18446744073709551615",
			},
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Uint64": testhelper.Comparer{Value: uint64(18446744073709551615)},
			},
		},

		testhelper.TestCase{
			Name: "Bad value float32",
			Env: map[string]string{
				"EXTRA_FLOAT32": "3.4028236e38",
			},
			CmdLine: []string{
				"float",
			},
			ExecError: "Cannot convert flag value float.float32: 3.4028236e38 would overflow a float32 (max 3.4028234663852886e+38)",
		},
		testhelper.TestCase{
			Name: "Max value float32",
			Env: map[string]string{
				"EXTRA_FLOAT32": "3.4028234663852886e+38",
			},
			CmdLine: []string{
				"float",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Float32": testhelper.Comparer{Value: float32(3.4028234663852886e+38)},
			},
		},

		testhelper.TestCase{
			Name: "Bad value float64",
			Env: map[string]string{
				"EXTRA_FLOAT64": "1.7976931348623159e308",
			},
			CmdLine: []string{
				"float",
			},
			ExecError: "Cannot convert flag value float.float64: unable to cast \"1.7976931348623159e308\" of type string to float64",
		},
		testhelper.TestCase{
			Name: "Max value float641",
			Env: map[string]string{
				"EXTRA_FLOAT64": "1.7976931348623157e308",
			},
			CmdLine: []string{
				"float",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Float64": testhelper.Comparer{Value: float64(1.7976931348623157e308)},
			},
		},

		testhelper.TestCase{
			Name: "Bad custom flag",
			Env: map[string]string{
				"EXTRA_FLAGIP": "1.2.3.4.5",
			},
			CmdLine: []string{
				"flag",
			},
			ExecError: "Invalid value 1.2.3.4.5, should be an IP address",
		},
		testhelper.TestCase{
			Name: "Bad port",
			Env: map[string]string{
				"EXTRA_FLAGPORT": "11111111111",
			},
			CmdLine: []string{
				"flag",
			},
			ExecError: "Invalid port 11111111111, must be an integer between 0 and 65535",
		},
		testhelper.TestCase{
			Name: "Int outside range",
			Env: map[string]string{
				"EXTRA_FLAGINT": "20000",
			},
			CmdLine: []string{
				"flag",
			},
			ExecError: "Invalid value 20000 for variable FlagInt, should be between 0 and 1000",
		},
		testhelper.TestCase{
			Name: "Int not an int",
			Env: map[string]string{
				"EXTRA_FLAGINT": "hi",
			},
			CmdLine: []string{
				"flag",
			},
			ExecError: "Variable FlagInt, hi, cannot be converted to a number",
		},
		testhelper.TestCase{
			Name: "Wrong enum",
			Env: map[string]string{
				"EXTRA_FLAGENUM": "hello",
			},
			CmdLine: []string{
				"flag",
			},
			ExecError: "Invalid value hello for variable FlagEnum, should be one of a, b, c.",
		},
		testhelper.TestCase{
			Name: "Wrong cstring",
			Env: map[string]string{
				"EXTRA_FLAGCSTRING": "hello",
			},
			CmdLine: []string{
				"flag",
			},
			ExecError: "Invalid string value hello for variable FlagCString, must have an uppercase first letter",
		},

		testhelper.TestCase{
			Name: "Bad time",
			Env: map[string]string{
				"EXTRA_TIME": "some day",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value via unmarshaling: parsing time",
		},
		testhelper.TestCase{
			Name: "Bad ptime",
			Env: map[string]string{
				"EXTRA_PTIME": "some day",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value via unmarshaling: parsing time",
		},
		testhelper.TestCase{
			Name: "Bad duration",
			Env: map[string]string{
				"EXTRA_DURATION": "forever",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value Duration: unable to cast",
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:  testhelper.NewExtraConfig,
		CompareMap: testhelper.ExtraConfigCustomComparers,
		UserDocList: map[string]*greenery.DocSet{
			"": &testhelper.ExtraDocs}})
	require.NoError(t, err)
}

func TestBadFlagValuesCfg(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name:        "Bad value bool",
			CfgContents: "[bool]\nbool = \"nonbool\"",
			CmdLine: []string{
				"bool",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name:        "Bad value int",
			CfgContents: "[int]\nint = \"nonint\"",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name:        "Bad value int8",
			CfgContents: "[int]\nint8 = 128",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int8: 128 would overflow an int8",
		},
		testhelper.TestCase{
			Name:        "Bad value int8 2",
			CfgContents: "[int]\nint8 = -129",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int8: -129 would overflow an int8",
		},
		testhelper.TestCase{
			Name:        "Max value int8",
			CfgContents: "[int]\nint8 = 127",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int8": testhelper.Comparer{Value: int8(127)},
			},
		},
		testhelper.TestCase{
			Name:        "Max value int8 2",
			CfgContents: "[int]\nint8 = -128",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int8": testhelper.Comparer{Value: int8(-128)},
			},
		},
		testhelper.TestCase{
			Name:        "Bad value int16",
			CfgContents: "[int]\nint16 = 32768",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int16: 32768 would overflow an int16",
		},
		testhelper.TestCase{
			Name:        "Bad value int16 2",
			CfgContents: "[int]\nint16 = -32769",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int16: -32769 would overflow an int16",
		},
		testhelper.TestCase{
			Name:        "Max value int16",
			CfgContents: "[int]\nint16 = 32767",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int16": testhelper.Comparer{Value: int16(32767)},
			},
		},
		testhelper.TestCase{
			Name:        "Max value int16 2",
			CfgContents: "[int]\nint16 = -32768",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int16": testhelper.Comparer{Value: int16(-32768)},
			},
		},
		testhelper.TestCase{
			Name:        "Bad value int32",
			CfgContents: "[int]\nint32 = 2147483648",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int32: 2147483648 would overflow an int32",
		},
		testhelper.TestCase{
			Name:        "Bad value int32 2",
			CfgContents: "[int]\nint32 = -2147483649",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.int32: -2147483649 would overflow an int32",
		},
		testhelper.TestCase{
			Name:        "Max value int32",
			CfgContents: "[int]\nint32 = 2147483647",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int32": testhelper.Comparer{Value: int32(2147483647)},
			},
		},
		testhelper.TestCase{
			Name:        "Max value int32 2",
			CfgContents: "[int]\nint32 = -2147483648",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int32": testhelper.Comparer{Value: int32(-2147483648)},
			},
		},
		testhelper.TestCase{
			Name:        "Bad value int64",
			CfgContents: "[int]\nint64 = 9223372036854775808",
			CmdLine: []string{
				"int",
			},
			// In this case toml fails
			ExecError: "Could not parse config file",
		},
		testhelper.TestCase{
			Name:        "Bad value int64 2",
			CfgContents: "[int]\nint64 = -9223372036854775809",
			CmdLine: []string{
				"int",
			},
			// In this case toml fails
			ExecError: "Could not parse config file",
		},
		testhelper.TestCase{
			Name:        "Max value int64",
			CfgContents: "[int]\nint64 = 9223372036854775807",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int64": testhelper.Comparer{Value: int64(9223372036854775807)},
			},
		},
		testhelper.TestCase{
			Name:        "Max value int64 2",
			CfgContents: "[int]\nint64 = -9223372036854775808",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Int64": testhelper.Comparer{Value: int64(-9223372036854775808)},
			},
		},
		testhelper.TestCase{
			Name:        "Bad value uint",
			CfgContents: "[int]\nuint = -2",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name:        "Bad value uint8",
			CfgContents: "[int]\nuint8 = -2",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name:        "Bad value uint8 2",
			CfgContents: "[int]\nuint8 = 256",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.uint8: 256 would overflow an uint8",
		},
		testhelper.TestCase{
			Name:        "Max value uint8",
			CfgContents: "[int]\nuint8 = 255",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Uint8": testhelper.Comparer{Value: uint8(255)},
			},
		},
		testhelper.TestCase{
			Name:        "Bad value uint16",
			CfgContents: "[int]\nuint16 = -2",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name:        "Bad value uint16 2",
			CfgContents: "[int]\nuint16 = 65536",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.uint16: 65536 would overflow an uint16",
		},
		testhelper.TestCase{
			Name:        "Max value uint16",
			CfgContents: "[int]\nuint16 = 65535",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Uint16": testhelper.Comparer{Value: uint16(65535)},
			},
		},
		testhelper.TestCase{
			Name:        "Bad value uint32",
			CfgContents: "[int]\nuint32 = -2",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name:        "Bad value uint32 2",
			CfgContents: "[int]\nuint32 = 4294967296",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.uint32: 4294967296 would overflow an uint32",
		},
		testhelper.TestCase{
			Name:        "Max value uint32",
			CfgContents: "[int]\nuint32 = 4294967295",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Uint32": testhelper.Comparer{Value: uint32(4294967295)},
			},
		},
		testhelper.TestCase{
			Name:        "Bad value uint64",
			CfgContents: "[int]\nuint64 = -2",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value",
		},
		testhelper.TestCase{
			Name:        "Bad value uint64 2",
			CfgContents: "[int]\nuint64 = 9223372036854775808",
			CmdLine: []string{
				"int",
			},
			// In this case toml fails, unfortunately TOML assumes int64 so we
			// can't get the full range of uint64 values
			ExecError: "Could not parse config file",
		},
		testhelper.TestCase{
			Name:        "Max value uint64",
			CfgContents: "[int]\nuint64 = 9223372036854775807",
			CmdLine: []string{
				"int",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Uint64": testhelper.Comparer{Value: uint64(9223372036854775807)},
			},
		},

		testhelper.TestCase{
			Name:        "Bad value float32",
			CfgContents: "[float]\nfloat32 = 3.4028236e38",
			CmdLine: []string{
				"float",
			},
			ExecError: "Cannot convert flag value float.float32: 3.4028236e+38 would overflow a float32",
		},
		testhelper.TestCase{
			Name:        "Max value float32",
			CfgContents: "[float]\nfloat32 = 3.4028234663852886e+38",
			CmdLine: []string{
				"float",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Float32": testhelper.Comparer{Value: float32(3.4028234663852886e+38)},
			},
		},

		testhelper.TestCase{
			Name:        "Bad value float64",
			CfgContents: "[float]\nfloat64 = 1.7976931348623159e308",
			CmdLine: []string{
				"float",
			},
			// In this case toml fails
			ExecError: "Could not parse config file",
		},
		testhelper.TestCase{
			Name:        "Max value float64",
			CfgContents: "[float]\nfloat64 = 1.7976931348623157e308",
			CmdLine: []string{
				"float",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Float64": testhelper.Comparer{Value: float64(1.7976931348623157e308)},
			},
		},

		testhelper.TestCase{
			Name:        "Bad int slice",
			CfgContents: "[slice]\nint = [\"a\", \"b\"]",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value slice.int",
		},

		testhelper.TestCase{
			Name:        "Bad custom flag",
			CfgContents: "[flag]\nip = \"1.2.3.4.5\"",
			CmdLine: []string{
				"flag",
			},
			ExecError: "Invalid value 1.2.3.4.5, should be an IP address",
		},
		testhelper.TestCase{
			Name:        "Bad port",
			CfgContents: "[flag]\nport = 11111111111",
			CmdLine: []string{
				"flag",
			},
			ExecError: "Invalid port 11111111111, must be an integer between 0 and 65535",
		},
		testhelper.TestCase{
			Name:        "Int outside range",
			CfgContents: "[flag]\nint = 20000",
			CmdLine: []string{
				"flag",
			},
			ExecError: "Invalid value 20000 for variable FlagInt, should be between 0 and 1000",
		},
		testhelper.TestCase{
			Name:        "Int not an int",
			CfgContents: "[flag]\nint = \"hi\"",
			CmdLine: []string{
				"flag",
			},
			ExecError: "Variable FlagInt, hi, cannot be converted to a number",
		},
		testhelper.TestCase{
			Name:        "Wrong enum",
			CfgContents: "[flag]\nenum = \"boo\"",
			CmdLine: []string{
				"flag",
			},
			ExecError: "Invalid value boo for variable FlagEnum, should be one of a, b, c.",
		},
		testhelper.TestCase{
			Name:        "Wrong cstring",
			CfgContents: "[flag]\ncstring = \"hello\"",
			CmdLine: []string{
				"flag",
			},
			ExecError: "Invalid string value hello for variable FlagCString, must have an uppercase first letter",
		},

		testhelper.TestCase{
			Name: "Bad time",
			Env: map[string]string{
				"EXTRA_TIME": "some day",
			},
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value via unmarshaling: parsing time",
		},

		testhelper.TestCase{
			Name:        "Bad time",
			CfgContents: "[int]\ntime = \"some day\"",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value via unmarshaling: parsing time",
		},
		testhelper.TestCase{
			Name:        "Bad ptime",
			CfgContents: "[int]\nptime = \"some day\"",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value via unmarshaling: parsing time",
		},
		testhelper.TestCase{
			Name:        "Bad duration",
			CfgContents: "[int]\nduration = \"forever\"",
			CmdLine: []string{
				"int",
			},
			ExecError: "Cannot convert flag value int.duration: unable to cast",
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:  testhelper.NewExtraConfig,
		CompareMap: testhelper.ExtraConfigCustomComparers,
		UserDocList: map[string]*greenery.DocSet{
			"": &testhelper.ExtraDocs}})
	require.NoError(t, err)
}

func TestBadConfFile(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Bad config file",
			CmdLine: []string{
				"bool",
			},
			CfgContents: "invalid yadda yadda 1234\n\n\nsgfgfs fg = 453",
			ExecError:   "Could not parse",
		},
		testhelper.TestCase{
			Name: "Nonexistent config file",
			CmdLine: []string{
				"bool",
			},
			CfgContents:    "invalid yadda yadda 1234\n\n\nsgfgfs fg = 453",
			CmdlineCfgName: "/tmp/nonexistent",
			ExecError:      "Could not load",
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:  testhelper.NewExtraConfig,
		CompareMap: testhelper.ExtraConfigCustomComparers,
		UserDocList: map[string]*greenery.DocSet{
			"": &testhelper.ExtraDocs}})
	require.NoError(t, err)
}

func TestNoCustomDeser(t *testing.T) {
	goldDefault := filepath.Join("testdata", cmdTestName+".TestConfig.extracfg")

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "No custom vars listed",
			CmdLine: []string{
				"config",
				"display",
			},
			CfgFile: goldDefault,
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfig.displayextracfg"),
				Custom: testhelper.CompareIgnoreTmp},
			ExecError: "Invalid key(s)",
		},
		testhelper.TestCase{
			Name: "Invalid keys",
			CmdLine: []string{
				"config",
				"display",
			},
			CfgContents: `# config custom section
[custom]
# some example custom parameter name enum, valid values a,b,c
[[custom.nameenum]]
name = "k1"
enum = "a"

[[custom.nameenum]]
name = "k2"
enum = "b"

[[custom.nameenum]]
name = "k3"
enum = "c"
# some example custom parameter name values
[[custom.namevalue]]
name = "k1"
value = "v1"

[[custom.namevalue]]
name = "k2"
value = "v2"

[[custom.namevalue]]
name = "k3"
value = "v3"
`,
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfig.displayextracfg"),
				Custom: testhelper.CompareIgnoreTmp},
			ExecError: "Invalid key(s)",
		},
		testhelper.TestCase{
			Name: "Custom vars not all processed",
			CmdLine: []string{
				"config",
				"display",
			},
			CustomVars: testhelper.ExtraConfigCustomVars,
			CustomParser: func(lcfg greenery.Config, vals map[string]interface{}) (processed []string, err error) {
				// Fake parser, just to get an error
				for k := range vals {
					switch k {
					case "custom.namevalue":
						processed = append(processed, k)
					case "custom.nameenum":
						// Not adding to processed
					}
				}
				return
			},
			CfgFile:   goldDefault,
			ExecError: "Unprocessed key(s) in the configuration file, needed custom processing for custom.nameenum. The parser has processed: custom.namevalue",
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfig.displayextracfg"),
				Custom: testhelper.CompareIgnoreTmp},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:  testhelper.NewExtraConfig,
		CompareMap: testhelper.ExtraConfigCustomComparers,
		UserDocList: map[string]*greenery.DocSet{
			"": &testhelper.ExtraDocs}})
	require.NoError(t, err)
}
