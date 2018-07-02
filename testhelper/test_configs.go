package testhelper

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
)

// NopNoArgs is a command handler that doesn't do anything
func NopNoArgs(lcfg greenery.Config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("Unsupported extra arg(s) %v", args)
	}
	return nil
}

// --------------------------------------------

// SimpleConfig is the simplest configuration, only base fields
type SimpleConfig struct {
	*greenery.BaseConfig
}

// NewSimpleConfig returns an instance of the configuration
func NewSimpleConfig() greenery.Config {
	return &SimpleConfig{
		BaseConfig: greenery.NewBaseConfig("simple", nil),
	}
}

// ----------------------------------------------

// ExtraDocs is a documentation list for the config with all possible types of fields
var ExtraDocs = greenery.DocSet{
	Usage: map[string]*greenery.CmdHelp{
		"int": &greenery.CmdHelp{
			Short: "int flags",
		},
		"bool": &greenery.CmdHelp{
			Short: "bool flags",
		},
		"float": &greenery.CmdHelp{
			Short: "float flags",
		},
		"string": &greenery.CmdHelp{
			Short: "string flags",
		},
		"flag": &greenery.CmdHelp{
			Short: "flag flags",
		},
		"custom": &greenery.CmdHelp{
			Short: "custom config values",
		},
	},
	CmdLine: map[string]string{
		"Bool": "cmdline bool",

		"FlagCString": "cmdline cstring",
		"FlagEnum":    "cmdline enum",
		"FlagIP":      "cmdline ip",
		"FlagInt":     "cmdline int",
		"FlagPort":    "cmdline port",

		"Float32": "cmdline float32",
		"Float64": "cmdline float64",

		"Int":    "cmdline int",
		"Int16":  "cmdline int16",
		"Int32":  "cmdline int32",
		"Int64":  "cmdline int64",
		"Int8":   "cmdline int8",
		"Uint":   "cmdline uint",
		"Uint16": "cmdline uint16",
		"Uint32": "cmdline uint32",
		"Uint64": "cmdline uint64",
		"Uint8":  "cmdline uint8",

		"String": "cmdline string",
	},
	ConfigFile: map[string]string{
		greenery.DocConfigHeader: "Config generated while testing",

		"bool.": "config bool section",
		"Bool":  "config bool",

		"custom.": "config custom section",
		"NameValue": `# some example custom parameter name values
[[custom.namevalue]]
name = "k1"
value = "v1"

[[custom.namevalue]]
name = "k2"
value = "v2"

[[custom.namevalue]]
name = "k3"
value = "v3"`,
		"NameEnum": `# some example custom parameter name enum, valid values a,b,c
[[custom.nameenum]]
name = "k1"
enum = "a"

[[custom.nameenum]]
name = "k2"
enum = "b"

[[custom.nameenum]]
name = "k3"
enum = "c"`,

		"flag.":       "config flag section",
		"FlagCString": "config cstring",
		"FlagEnum":    "config enum",
		"FlagIP":      "config ip",
		"FlagInt":     "config int",
		"FlagPort":    "config port",

		"Float32": "config float32",
		"Float64": "config float64",

		"int.":     "config int section",
		"Int":      "config int",
		"Int16":    "config int16",
		"Int32":    "config int32",
		"Int64":    "config int64",
		"Int8":     "config int8",
		"Uint":     "config uint",
		"Uint16":   "config uint16",
		"Uint32":   "config uint32",
		"Uint64":   "config uint64",
		"Uint8":    "config uint8",
		"Time":     "config time",
		"PTime":    "config ptime",
		"Duration": "config duration",

		"string.": "config string section",
		"String":  "config string",

		"slice.":      "config slice section",
		"SliceString": "config slice string",
		"SliceInt":    "config slice int",
	},
}

// ExtraConfigCustomNameValue is a sample struct for custom deserialization
type ExtraConfigCustomNameValue struct {
	Name  string
	Value string
}

// GetTyped is an accessor returning the name/value pair as a string list
func (e *ExtraConfigCustomNameValue) GetTyped() []string {
	return []string{e.Name, e.Value}
}

// ExtraConfigCustomNameEnum is a sample struct for custom deserialization
type ExtraConfigCustomNameEnum struct {
	Name string
	Enum *greenery.EnumValue
}

// GetTyped is an accessor returning the name/value pair as a string list
func (e *ExtraConfigCustomNameEnum) GetTyped() []string {
	return []string{e.Name, e.Enum.GetTyped()}
}

// ExtraConfigCustomVars is the names of the custom structs to deserialize
var ExtraConfigCustomVars = []string{
	"custom.namevalue",
	"custom.nameenum",
}

// ExtraConfigCustomComparers is a set of compare functions used for the
// built-in base config types
var ExtraConfigCustomComparers = map[string]CompareFunc{
	"FlagIP":      CompareGetterToGetter,
	"FlagEnum":    CompareGetterToGetter,
	"FlagInt":     CompareGetterToGetter,
	"FlagCString": CompareGetterToGetter,
	"NameValue":   ExtraConfigCustomCompareNameValue,
	"NameEnum":    ExtraConfigCustomCompareNameEnum,
}

// ExtraConfig is a config with all possible types of field
type ExtraConfig struct {
	*greenery.BaseConfig

	Bool bool `greenery:"bool|bool|,    bool.bool, BOOL"`

	NameValue []ExtraConfigCustomNameValue `greenery:"||custom,  custom.namevalue,"`
	NameEnum  []ExtraConfigCustomNameEnum  `greenery:"||custom,  custom.nameenum,"`

	FlagIP      *greenery.IPValue           `greenery:"flag|ip|,      flag.ip,      FLAGIP"`
	FlagPort    greenery.PortValue          `greenery:"flag|port|,    flag.port,    FLAGPORT"`
	FlagEnum    *greenery.EnumValue         `greenery:"flag|enum|,    flag.enum,    FLAGENUM"`
	FlagInt     *greenery.IntValue          `greenery:"flag|int|,     flag.int,     FLAGINT"`
	FlagCString *greenery.CustomStringValue `greenery:"flag|cstring|, flag.cstring, FLAGCSTRING"`

	Float32 float32 `greenery:"float|float32|,        float.float32,     FLOAT32"`
	Float64 float64 `greenery:"float|float64|,        float.float64,     FLOAT64"`

	Int      int           `greenery:"int|int|,          int.int,       INT"`
	Int8     int8          `greenery:"int|int8|,         int.int8,      INT8"`
	Int16    int16         `greenery:"int|int16|,        int.int16,     INT16"`
	Int32    int32         `greenery:"int|int32|,        int.int32,     INT32"`
	Int64    int64         `greenery:"int|int64|,        int.int64,     INT64"`
	Uint     uint          `greenery:"int|uint|,         int.uint,      UINT"`
	Uint8    uint8         `greenery:"int|uint8|,        int.uint8,     UINT8"`
	Uint16   uint16        `greenery:"int|uint16|,       int.uint16,    UINT16"`
	Uint32   uint32        `greenery:"int|uint32|,       int.uint32,    UINT32"`
	Uint64   uint64        `greenery:"int|uint64|,       int.uint64,    UINT64"`
	Time     time.Time     `greenery:"||none,            int.time,      TIME"`
	PTime    *time.Time    `greenery:"||none,            int.ptime,     PTIME"`
	Duration time.Duration `greenery:"||none,            int.duration,  DURATION"`

	String string `greenery:"string|string|,    string.string, STRING"`

	SliceString []string `greenery:"||none,    slice.string, "`
	SliceInt    []int    `greenery:"||none,    slice.int, "`

	internal string
}

// UcFirstAndUpcase is an arbitrary validation function used by the custom
// string flag in ExtraConfig, it will fail if the value does not start with
// an uppercase character, and will return an uppercased version of the value
func UcFirstAndUpcase(name, s string, data interface{}) (string, error) {
	r, _ := utf8.DecodeRuneInString(s)
	if !unicode.IsUpper(r) {
		return "", fmt.Errorf("Invalid string value %s for variable %s, must have an uppercase first letter", s, name)
	}

	return strings.ToUpper(s), nil
}

// NewExtraConfig instantiates the configuration with all possible types of field
func NewExtraConfig() greenery.Config {
	ptime := time.Date(2017, time.June, 3, 12, 8, 32, 454, time.UTC)
	cfg := &ExtraConfig{
		BaseConfig: greenery.NewBaseConfig("extra", map[string]greenery.Handler{
			"bool":   NopNoArgs,
			"custom": NopNoArgs,
			"flag":   NopNoArgs,
			"float":  NopNoArgs,
			"int":    NopNoArgs,
			"string": NopNoArgs,
		}),

		Bool: true,

		NameValue: []ExtraConfigCustomNameValue{
			ExtraConfigCustomNameValue{Name: "k1", Value: "v1"},
			ExtraConfigCustomNameValue{Name: "k2", Value: "v2"},
			ExtraConfigCustomNameValue{Name: "k3", Value: "v3"},
		},
		NameEnum: []ExtraConfigCustomNameEnum{
			ExtraConfigCustomNameEnum{Name: "k1", Enum: greenery.NewDefaultEnumValue("", "a", "a", "b", "c")},
			ExtraConfigCustomNameEnum{Name: "k2", Enum: greenery.NewDefaultEnumValue("", "b", "a", "b", "c")},
			ExtraConfigCustomNameEnum{Name: "k3", Enum: greenery.NewDefaultEnumValue("", "c", "a", "b", "c")},
		},

		FlagIP:      greenery.NewDefaultIPValue("FlagIP", "127.0.0.1"),
		FlagEnum:    greenery.NewDefaultEnumValue("FlagEnum", "b", "a", "b", "c"),
		FlagInt:     greenery.NewDefaultIntValue("FlagInt", 400, 0, 1000),
		FlagCString: greenery.NewDefaultCustomStringValue("FlagCString", "HELLO", UcFirstAndUpcase, nil),

		Float32: 12.34,
		Float64: -56.78,
		Int:     -1,
		Int8:    -2,
		Int16:   -3,
		Int32:   -4,
		Int64:   -5,
		Uint:    1,
		Uint8:   2,
		Uint16:  3,
		Uint32:  4,
		Uint64:  5,

		String: "init",

		SliceString: []string{"first", "second", "third"},
		SliceInt:    []int{1, 2, 3, 4},

		Time:     time.Date(2018, time.June, 3, 12, 8, 32, 454, time.UTC),
		PTime:    &ptime,
		Duration: time.Hour*48 + time.Minute*16 + time.Second*32 + time.Millisecond*45,

		internal: "internal",
	}

	_ = cfg.FlagPort.SetInt(80)
	return cfg
}

// ExtraConfigCustomParse is a parser for the custom structs in the config
func ExtraConfigCustomParse(lcfg greenery.Config, vals map[string]interface{}) (processed []string, err error) {
	cfg := lcfg.(*ExtraConfig)

	for k, vv := range vals {
		switch k {
		case "custom.namevalue":
			cfg.NameValue = []ExtraConfigCustomNameValue{}
			err = lcfg.Unmarshal(k, &cfg.NameValue)
			if err != nil {
				return
			}
			processed = append(processed, k)
		case "custom.nameenum":
			cfg.NameEnum = []ExtraConfigCustomNameEnum{}
			tt := reflect.TypeOf(vv)
			if tt.Kind() != reflect.Slice {
				panic("Unexpected type, not a slice")
			}

			v := reflect.ValueOf(vv)
			for i := 0; i < v.Len(); i++ {
				ev := v.Index(i).Elem()
				if ev.Kind() != reflect.Map {
					panic("Unexpected type, not a map")
				}

				got := 0
				converted := ExtraConfigCustomNameEnum{
					Name: "",
					Enum: greenery.NewEnumValue("Enum", "a", "b", "c"),
				}

				keys := ev.MapKeys()
				for _, key := range keys {
					evelemkey := key.String()
					evelemvalue := ev.MapIndex(key)
					if _, ok := evelemvalue.Elem().Interface().(string); !ok {
						panic(fmt.Sprintf("Unexpected value type, not a string %v", evelemvalue.Interface()))
					}

					switch evelemkey {
					case "name":
						converted.Name = evelemvalue.Elem().String()
						got++
					case "enum":
						if err := converted.Enum.Set(evelemvalue.Elem().String()); err != nil {
							panic(err.Error())
						}

						got++
					default:
						panic(fmt.Sprintf("Unknown key %s", evelemkey))
					}

				}
				if got != 2 {
					panic("Didn't get both name and elem")
				}

				cfg.NameEnum = append(cfg.NameEnum, converted)
			}
			processed = append(processed, k)
		}
	}
	return
}

// ExtraConfigCustomCompareNameEnum is a comparer function for the
// ExtraConfigCustomNameEnum type
func ExtraConfigCustomCompareNameEnum(t *testing.T, name string, vl, vr interface{}) error {
	right, ok := vr.([]ExtraConfigCustomNameEnum)
	if !ok {
		return fmt.Errorf("Wrong enum type for variable %s, %v", name, spew.Sdump(vr))
	}

	left, ok := vl.([][]string)
	if ok {
		require.Equal(t, len(left), len(right), "Lengths differ for variable %s", name)
		for i, lv := range left {
			require.Equalf(t, lv, right[i].GetTyped(), "Variable %s element differs", name)
		}
	} else {
		left, ok := vl.([]ExtraConfigCustomNameEnum)
		if !ok {
			return fmt.Errorf("Wrong enum type for variable %s, %v", name, spew.Sdump(vl))
		}
		require.Equal(t, len(left), len(right), "Lengths differ for variable %s", name)
		for i, lv := range left {
			require.Equalf(t, lv.GetTyped(), right[i].GetTyped(), "Variable %s element differs", name)
		}
	}
	return nil
}

// ExtraConfigCustomCompareNameValue is a comparer function for the
// ExtraConfigCustomNameValue type
func ExtraConfigCustomCompareNameValue(t *testing.T, name string, vl, vr interface{}) error {
	right, ok := vr.([]ExtraConfigCustomNameValue)
	if !ok {
		return fmt.Errorf("Wrong value type for variable %s, %v", name, spew.Sdump(vr))
	}

	left, ok := vl.([][]string)
	if ok {
		require.Equal(t, len(left), len(right), "Lengths differ for variable %s", name)
		for i, lv := range left {
			require.Equalf(t, lv, right[i].GetTyped(), "Variable %s element differs", name)
		}
	} else {
		left, ok := vl.([]ExtraConfigCustomNameValue)
		if !ok {
			return fmt.Errorf("Wrong value type for variable %s, %v", name, spew.Sdump(vl))
		}
		require.Equal(t, len(left), len(right), "Lengths differ for variable %s", name)
		for i, lv := range left {
			require.Equalf(t, lv.GetTyped(), right[i].GetTyped(), "Variable %s element differs", name)
		}
	}
	return nil
}
