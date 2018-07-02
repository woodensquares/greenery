package greenery

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/shibukawa/configdir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func getSetterStringer(field reflect.Value) (reflect.Value, reflect.Value) {
	if field.Addr().Elem().Kind() == reflect.Ptr {
		if field.Type().Implements(flagInterface) {
			return field.Addr().Elem().MethodByName("Set"), field.Addr().Elem().MethodByName("String")
		}
	} else {
		if reflect.PtrTo(field.Type()).Implements(flagInterface) {
			return field.Addr().MethodByName("Set"), field.Addr().MethodByName("String")
		}
	}

	return reflect.Value{}, reflect.Value{}
}

// setString will set a string value in a config variable using setters as
// needed.
func (bcfg *BaseConfig) setString(cfg Config, k, v string) error {
	p := []reflect.Value{reflect.ValueOf(v)}
	baseValue := reflect.ValueOf(bcfg).Elem()

	userType := reflect.TypeOf(cfg).Elem()
	userValue := reflect.ValueOf(cfg).Elem()

	var field reflect.Value
	var ok bool

	if _, ok = userType.FieldByName(k); !ok {
		// Current base config cannot exercise this as we don't have any
		// non-cmd env fields
		if _, ok = baseType.FieldByName(k); !ok {
			return fmt.Errorf(
				"Internal error, cannot find field %s (setting value %s)", k, v)
		}
		field = baseValue.FieldByName(k)
	} else {
		field = userValue.FieldByName(k)
	}

	setter, _ := getSetterStringer(field)
	if setter.Kind() != reflect.Invalid {
		cfg.Tracef("Calling the setter")
		rv := setter.Call(p)
		rve := rv[len(rv)-1]
		if !rve.IsNil() {
			if rerr, ok := rve.Interface().(error); ok {
				return rerr
			}
			return fmt.Errorf("Internal error in setString, wrong type returned by the setter")

			// Another possibility maybe
			// eValues := rve.MethodByName("Error").Call([]reflect.Value{})
			// err = fmt.Errorf(eValues[0].String())
		}
		return nil
	}

	cfg.Tracef("Calling setField with %v", v)
	return setField(cfg, field, nil, v, k)
}

// loadHelper is the main load worker, which is executed both on the embedded
// and user structs
func loadHelper(cfg Config, vp *viper.Viper, viperKeys map[string]bool,
	x reflect.StructField, v reflect.Value) error {

	cobra, vipername, _, _ := parseTags(x)
	if vipername != "" && !strings.HasSuffix(cobra, sepCmdParts+"custom") {
		cfg.Tracef("Get value for %s", vipername)
		if vp.Get(vipername) == nil {
			// If there is no cobra binding, and no configuration file, the
			// variable will not be present in viper, so don't clobber the
			// value already in there.
			cfg.Tracef("Viper value not found, skipping")
			return nil
		}
		viperKeys[vipername] = true
		field := v.FieldByName(x.Name)
		// Assume that anything that is a flag wants its .Set method to be
		// called rather than assigning the variable directly (for
		// validation purposes).
		setter, _ := getSetterStringer(field)

		// If the setter exists, call it to set the value instead of using
		// reflect.Set. Note Viper, unlike cobra, does not validate, so we
		// can't ignore the return value from .Call(). Assume error is the
		// last returned value from our called method and fail if it's not
		// nil.
		if setter.Kind() != reflect.Invalid {
			cfg.Trace("Have a flag setter")
			// For our custom flags, in cfg we have one more level of
			// indirection as we don't have the straight
			// value, but a pointer to it, so dereference
			vs := vp.GetString(vipername)
			cfg.Tracef("Will set string: %v", vs)
			p := []reflect.Value{reflect.ValueOf(vs)}
			rv := setter.Call(p)
			rve := rv[len(rv)-1]
			if !rve.IsNil() {
				if rerr, ok := rve.Interface().(error); ok {
					return rerr
				}
				return fmt.Errorf("Internal error in loadHelper, wrong type returned by the setter")

				// Another possibility maybe
				// eValues := rve.MethodByName("Error").Call([]reflect.Value{})
				// err = fmt.Errorf(eValues[0].String())
			}
		} else {
			cfg.Tracef("Set as-is")
			r, _ := utf8.DecodeRuneInString(x.Name)
			if !unicode.IsUpper(r) {
				return fmt.Errorf("Trying to set value %v to a non-exported field %s", vp.Get(vipername), x.Name)
			}

			return setField(cfg, field, vp, nil, vipername)
		}
	}
	return nil
}

// load will load configuration values from file and environment and set in
// the configuration structure implementing the Config interface.
func (bcfg *BaseConfig) load(cfg Config, defaultCfgFileName string, ccmd *cobra.Command,
	vp *viper.Viper) (err error) {
	var baseConf *BaseConfig
	var cfgFile = bcfg.ConfFile
	var noCfg = bcfg.NoCfg

	// No cfg has precedence over -c no matter what, if it is set via env or
	// cmdline don't even try to find a cfg file.
	if !noCfg {
		// Either no passed config file, or non-absolute name, so allow viper to
		// search in the places it should
		if cfgFile == "" || !strings.Contains(cfgFile, string(os.PathSeparator)) {
			for _, x := range configdir.New("", bcfg.s_appName).QueryFolders(configdir.All) {
				bcfg.Tracef("Adding path %s", x.Path)
				vp.AddConfigPath(x.Path)
			}

			if cwd, lerr := os.Getwd(); lerr == nil {
				bcfg.Tracef("Adding path %s", cwd)
				vp.AddConfigPath(cwd)
			}
		}

		if cfgFile == "" {
			bcfg.Tracef("Cfg file not set, trying with the default %s", defaultCfgFileName)
			vp.SetConfigFile(defaultCfgFileName)
		} else {
			vp.SetConfigFile(cfgFile)
		}

		// Only TOML is supported for now
		vp.SetConfigType("toml")

		// In case viper/cobra debugging is needed
		// jwalterweatherman.SetLogThreshold(jwalterweatherman.LevelTrace)
		// jwalterweatherman.SetStdoutThreshold(jwalterweatherman.LevelTrace)

		if err = vp.ReadInConfig(); err == nil {
			bcfg.Tracef("Loaded a valid configuration file from %s", vp.ConfigFileUsed())
			bcfg.s_usedConf = vp.ConfigFileUsed()
			bcfg.s_loaded = true
			bcfg.s_cfgDir = path.Dir(vp.ConfigFileUsed())
		} else {
			bcfg.Tracef("Could not load config file / config file unset: \"%s\"", cfgFile)
			// Most commands' options are available on the command line, so it is
			// possible to run them without a config file, for the ones that do
			// need some variables set they will validate as needed.
			_, viperNotFound := err.(viper.ConfigFileNotFoundError)
			_, golangNotFound := err.(*os.PathError)

			if viperNotFound || golangNotFound {
				// Not an error unless the user did actually want a config file
				if cfgFile != "" {
					err = errors.WithMessage(err, fmt.Sprintf("Could not load config file."))
					return
				}
				err = nil
				noCfg = true
			}

			if !noCfg {
				bcfg.Trace("Parse issue")
				err = errors.WithMessage(err, fmt.Sprintf("Could not parse config file %s", cfgFile))
				return
			}
		}
	}

	// Viper has all the up-to-date values, need to put them in cfg, viper
	// already does the correct overrides if the user sets them on cmdline vs
	// env vs config file. Note the assumption is we have at most one level of
	// pointers/references in the cfg, so when type is reflect.Ptr no need to
	// walk back further than one level.
	bcfg.Trace("Get values from Viper")
	t := reflect.TypeOf(cfg).Elem()
	v := reflect.ValueOf(cfg).Elem()
	viperKeys := make(map[string]bool)

	// Need to get our base configuration first, so we have access to the s_
	// internals for noclobber.
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type == basePType {
			var ok bool
			baseConf, ok = v.Field(i).Interface().(*BaseConfig)
			if !ok {
				// Should really not happen
				err = fmt.Errorf("Internal error in load, cannot cast")
				return
			}
			break
		}
	}

	// If the cmdline for a variable is set, no matter if it's shared multiple
	// times in different ways, it always should take precedence over
	// configuration and environment, so let's save the names of the variables
	// that were set on the cmdline first and their cmdline values using the
	// annotations we set up when binding.
	noclobber := map[string]string{}
	ccmd.Flags().Visit(func(fl *pflag.Flag) {
		if fl.Changed {
			if ann, ok := fl.Annotations["greeneryVar"]; ok {
				noclobber[ann[0]] = fl.Value.String()
			}
		}
	})

	for i := 0; i < t.NumField(); i++ {
		x := t.Field(i)

		if x.Type == basePType {
			bcfg.Trace("Get values for our base struct")
			v2 := v.Field(i).Elem()
			for i2 := 0; i2 < baseType.NumField(); i2++ {
				x2 := baseType.Field(i2)

				if clb, ok := noclobber[x.Name]; ok {
					// Not going to happen for now as we don't have shared
					// base values at the moment, everything is a straight
					// flag.
					cfg.Tracef("%s is on cmdline, not touching it as it's already set to %s", x2.Name, clb)
					if _, vipername, _, _ := parseTags(x2); vipername != "" {
						viperKeys[vipername] = true
					}
					continue
				}

				if err = loadHelper(cfg, vp, viperKeys, x2, v2); err != nil {
					return
				}
			}
		} else {
			if clb, ok := noclobber[x.Name]; ok {
				cfg.Tracef("%s is on cmdline, not touching it as it's already set to %s", x.Name, clb)
				if _, vipername, _, _ := parseTags(x); vipername != "" {
					viperKeys[vipername] = true
				}
				continue
			}

			if err = loadHelper(cfg, vp, viperKeys, x, v); err != nil {
				return
			}
		}
	}

	extraKeys := map[string]bool{}
	for _, vv := range bcfg.s_extraWanted {
		extraKeys[vv] = true
	}

	for _, vv := range baseConf.s_additionalEnv {
		cfg.Tracef("Processing env overrides for %v", vv)
		var evalue string
		if vv.Env != "" {
			ename := baseConf.s_ucAppName + "_" + vv.Env
			if evalue = os.Getenv(ename); evalue != "" {
				baseConf.s_env[ename] = evalue
			}
		}

		if vv.Cmdline == "" {
			cfg.Tracef("%s not on cmdline, only conf %s & env %s", vv.Name, vv.Viper, vv.Env)
			// If it has a viper (config) and no cobra, it's either conf or
			// env depending on which is set, env takes precedence.
			if evalue != "" {
				if _, ok := noclobber[vv.Name]; !ok {
					cfg.Tracef("Assign env %s to %s", evalue, vv.Name)
					if err = bcfg.setString(cfg, vv.Name, evalue); err != nil {
						return
					}
				}
			} else {
				cfg.Tracef("No env, leave it be conf")
			}
		} else if vv.Viper == "" {
			// If it has a cobra name (cmdline) and no viper (conf), if there was no
			// cmdline the env, if present, takes precedence
			if _, ok := noclobber[vv.Name]; evalue != "" && !ok {
				cfg.Tracef("Assign env %s to %s", evalue, vv.Name)
				if err = bcfg.setString(cfg, vv.Name, evalue); err != nil {
					return
				}
			}
		} else {
			err = fmt.Errorf("Internal error, cmd & env both not empty for %s: %s, %s", vv.Name, vv.Cmdline, vv.Viper)
			return
		}
	}

	// No config file, nothing else to do
	if noCfg {
		cfg.Trace("No config file, nothing else to do")
		return
	}

	// Check if we read any configuration keys that do not belong there (or
	// are mis-spelled) and fail.
	readKeys := vp.AllKeys()
	otherKeys := map[string]interface{}{}

	for _, v := range readKeys {
		cfg.Tracef("Extra looking at %s", v)
		if _, ok := viperKeys[v]; !ok && v != "" {
			cfg.Tracef("%s is not in viperkeys", v)
			if _, ok := extraKeys[v]; ok {
				bcfg.Tracef("----- Extra key found %s", v)
				otherKeys[v] = vp.Get(v)
			} else {
				cfg.Trace("Possible syntax error")
				err = fmt.Errorf("Invalid key(s) in the configuration file: %v", v)
				return
			}
		}
	}

	if len(otherKeys) != 0 {
		var processed []string

		if bcfg.s_extraParser != nil {
			processed, err = bcfg.s_extraParser(cfg, otherKeys)

			if err != nil || len(processed) == len(otherKeys) {
				return
			}

			// If we are here, no error but we haven't processed everything,
			// so continue below so we fail.
		}

		var upk []string
	outer:
		for pk := range otherKeys {
			// Slow but shouldn't happen during normal usage
			for _, cpk := range processed {
				if pk == cpk {
					continue outer
				}
			}
			upk = append(upk, pk)
		}
		err = fmt.Errorf("Unprocessed key(s) in the configuration file, needed custom processing for %s. The parser has processed: %s", strings.Join(upk, ","), strings.Join(processed, ","))
	}

	return
}
