package greenery

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// getHelpCmd creates a help command with a custom runner
func getHelpCmd() *cobra.Command {
	return &cobra.Command{
		RunE: func(c *cobra.Command, args []string) error {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				// Should never happen as it seems for now cobra will return
				// the root cmd if it doesn't find anything else
				c.Printf("Unknown help topic %#q\n", args)
				if err := c.Root().Usage(); err != nil {
					return err
				}
			} else {
				cmd.InitDefaultHelpFlag() // make possible 'help' flag to be shown
				if err := cmd.Help(); err != nil {
					// Should not happen
					return err
				}
			}
			return nil
		},
	}
}

// These handlers are called with the user config type, not the base, so we
// need to use reflect to get to the base configuration struct if needed (for
// logging, say).

// getCfg gets the BaseCfg struct in the user composed struct, can also be
// called with the base cfg directly.
func getCfg(icfg interface{}) (*BaseConfig, error) {
	ti := reflect.TypeOf(icfg)

	if ti == basePType {
		return icfg.(*BaseConfig), nil
	}

	typ := reflect.TypeOf(icfg)
	val := reflect.ValueOf(icfg)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Not a struct")
	}

	for i := 0; i < val.NumField(); i++ {
		x := typ.Field(i)
		if x.Type == basePType {
			if rc, ok := val.Field(i).Interface().(*BaseConfig); ok {
				return rc, nil
			}

			// Should never happen
			return nil, fmt.Errorf("Internal error, cannot cast to *BaseConfig")
		}
	}

	// Should never happen
	return nil, fmt.Errorf("Not a struct embedding BaseConfig")
}

// configInitCmdRunner is the runner for the config init command
func configInitCmdRunner(icfg Config, args []string) error {
	cfg, err := getCfg(icfg)
	if err != nil {
		// Should never happen given the interface
		return err
	}

	if len(args) != 0 {
		return fmt.Errorf("The command does not support additional arguments")
	}

	used, err := initCfgFile(icfg, cfg)
	cfg.s_usedConf = used

	if err == nil {
		// Allow users to quiet this for use in scripts
		if cfg.Verbosity.Value != 0 {
			fmt.Printf("Configuration file generated at %s\n", used)
		}
	}
	return err
}

// configEnvCmdRunner is the runner for the config env command
func configEnvCmdRunner(icfg Config, args []string) error {
	cfg, err := getCfg(icfg)
	if err != nil {
		// Should never happen given the interface
		return err
	}

	if len(args) != 0 {
		return fmt.Errorf("The command does not support additional arguments")
	}

	_, docs := cfg.GetDocs()

	fmt.Printf("%s\n-------------------------------------------------------------------", docs.ConfigEnvMsg1)

	out := []string{}
	any := false
	for k, v := range cfg.s_env {
		// Note that viper uses getenv, not lookupenv, so empty env variables
		// count the same as unset env variables.
		if v != "" {
			out = append(out, fmt.Sprintf("\n  %s -> %s", k, v))
			any = true
		}
	}

	if !any {
		fmt.Printf("\n  %s", docs.ConfigEnvMsg2)
	} else {
		sort.Strings(out)
		fmt.Print(strings.Join(out, ""))
	}

	fmt.Printf(`
-------------------------------------------------------------------


%s
-------------------------------------------------------------------
`, docs.ConfigEnvMsg3)

	out = []string{}
	t := reflect.TypeOf(icfg).Elem()

	apHelper := func(t []string, x reflect.StructField) []string {
		_, _, viperenv, _ := parseTags(x)
		if viperenv != "" {
			ds, ok := cfg.s_docs.ConfigFile[x.Name]
			if ds == "" || !ok {
				ds, ok = cfg.s_docs.CmdLine[x.Name]
				if !ok {
					// Should not happen due to previous checks
					cfg.Errorf("Could not find any documentation, cmdline or configfile, for %s", x.Name)
				}
			}
			t = append(t, fmt.Sprintf("%s_%s: %s\n", cfg.s_ucAppName, viperenv, ds))
		}

		return t
	}

	for i := 0; i < t.NumField(); i++ {
		x := t.Field(i)
		if x.Type == basePType {
			cfg.Trace("Get values for our base struct")
			for i2 := 0; i2 < baseType.NumField(); i2++ {
				x2 := baseType.Field(i2)
				out = apHelper(out, x2)
			}
		} else {
			out = apHelper(out, x)
		}
	}

	sort.Strings(out)
	fmt.Printf("%s-------------------------------------------------------------------\n", strings.Join(out, ""))

	return nil
}

// configDisplayCmdRunner is the runner for the config display command
func configDisplayCmdRunner(icfg Config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("The command does not support additional arguments")
	}

	out, err := icfg.Dump(icfg)
	if err != nil {
		// Should never happen given the interface (Dump can only fail due to that)
		return err
	}

	fmt.Println(out)
	return nil
}

// versionCmdRunner is the runner for the version command
func versionCmdRunner(icfg Config, args []string) error {
	cfg, err := getCfg(icfg)
	if err != nil {
		// Should never happen given the interface
		return err
	}

	if len(args) != 0 {
		return fmt.Errorf("The command does not support additional arguments")
	}

	if cfg.VersionFull != "" {
		fmt.Println(cfg.VersionFull)
	} else {
		if cfg.VersionPatchlevel == "" {
			if cfg.VersionMinor != "" {
				fmt.Printf("%s.%s\n", cfg.VersionMajor, cfg.VersionMinor)
			} else {
				fmt.Printf("%s\n", cfg.VersionMajor)
			}
		} else {
			fmt.Printf("%s.%s.%s\n", cfg.VersionMajor, cfg.VersionMinor, cfg.VersionPatchlevel)
		}
	}
	return nil
}
