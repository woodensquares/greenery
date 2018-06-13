package greenery_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/woodensquares/greenery"
)

type exampleCustomNameValue struct {
	Name  string
	Value string
}

type exampleCSVStruct struct {
	First  string
	Second string
	Third  string
}

type exampleCustomConfig struct {
	*greenery.BaseConfig
	NameValue []exampleCustomNameValue `greenery:"||custom,           custom-app.namevalue,"`
	CSV       exampleCSVStruct         `greenery:"||custom,           custom-app.csv,"`
	List      []string                 `greenery:"||custom,           custom-app.csv,"`
}

func exampleNewCustomConfig() *exampleCustomConfig {
	cfg := &exampleCustomConfig{
		BaseConfig: greenery.NewBaseConfig("custom", map[string]greenery.Handler{
			"display": exampleCustomDisplay,
		}),
		NameValue: []exampleCustomNameValue{},
		CSV:       exampleCSVStruct{},
		List:      []string{},
	}

	cfg.RegisterExtraParse(exampleCustomParse, []string{
		"custom-app.namevalue",
		"custom-app.csv",
		"custom-app.list",
	})

	return cfg
}

var exampleCustomDocs = map[string]*greenery.DocSet{
	"en": &greenery.DocSet{
		Short: "URL fetcher",
		Usage: map[string]*greenery.CmdHelp{
			"display": &greenery.CmdHelp{
				Short: "Show some custom variables",
			},
		},
		ConfigFile: map[string]string{
			"NameValue": "a set of name/value pairs",
			"CSV":       "a comma separated value variable",
			"List":      "a list of strings",
		},
	},
}

func exampleCustomParse(lcfg greenery.Config, vals map[string]interface{}) ([]string, error) {
	cfg := lcfg.(*exampleCustomConfig)
	var processed []string

	for k, v := range vals {
		switch k {
		case "custom-app.namevalue":
			if err := lcfg.Unmarshal(k, &cfg.NameValue); err != nil {
				return nil, err
			}
			processed = append(processed, k)
		case "custom-app.list":
			if err := lcfg.Unmarshal(k, &cfg.List); err != nil {
				return nil, err
			}
			processed = append(processed, k)
		case "custom-app.csv":
			vv, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("Unexpected type %T for key %s, not a string: %v", v, k, v)
			}

			parts := strings.Split(vv, ",")
			if len(parts) != 3 {
				return nil, fmt.Errorf("Invalid value for key %s: %v", k, v)
			}

			cfg.CSV.First = parts[0]
			cfg.CSV.Second = parts[1]
			cfg.CSV.Third = parts[2]

			processed = append(processed, k)
		}
	}

	return processed, nil
}

func exampleCustomDisplay(lcfg greenery.Config, args []string) error {
	cfg := lcfg.(*exampleCustomConfig)
	fmt.Printf("NameValue is %v\n", cfg.NameValue)
	fmt.Printf("CSV is %v\n", cfg.CSV)
	fmt.Printf("List is %v\n", cfg.List)
	return nil
}

func exampleCustomMain() {
	cfg := exampleNewCustomConfig()
	defer cfg.Cleanup()

	if err := cfg.Execute(cfg, exampleCustomDocs); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing: %s\n", err)
		os.Exit(1)
	}
}

// Example_custom is an example application showing how to enable greenery to
// parse custom configuration variables.
func Example_custom() {
	fmt.Println("----------------------------------------------------------")
	cfgFile, err := ioutil.TempFile("", "custom")
	if err != nil {
		fmt.Println("Cannot create a temporary file")
		os.Exit(1)
	}
	defer func() {
		_ = os.Remove(cfgFile.Name())
	}()

	if _, err = cfgFile.Write([]byte(`# Sample config file
[custom-app]
csv = "first,second,third"
list = [ "10.0.0.1", "10.0.0.2", "10.0.0.3" ]

[[custom-app.namevalue]]
name = "k1"
value = "v1"
[[custom-app.namevalue]]
name = "k2"
value = "v2"
`)); err != nil {
		fmt.Println("Cannot write the config file")
		os.Exit(1)
	}
	if err := cfgFile.Close(); err != nil {
		fmt.Println("Cannot close the config file")
		os.Exit(1)
	}

	os.Args = []string{"custom", "-c", cfgFile.Name(), "display"}
	fmt.Printf("Displaying some custom values set in the config file\n\n")
	exampleCustomMain()

	// Output: ----------------------------------------------------------
	// Displaying some custom values set in the config file
	//
	// NameValue is [{k1 v1} {k2 v2}]
	// CSV is {first second third}
	// List is [10.0.0.1 10.0.0.2 10.0.0.3]
}
