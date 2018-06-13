package greenery

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Used by various functions
var basePType = reflect.TypeOf(&BaseConfig{})
var baseType = reflect.TypeOf(BaseConfig{})
var flagInterface = reflect.TypeOf((*flag.Value)(nil)).Elem()
var baseInterface = reflect.TypeOf((*Config)(nil)).Elem()

// additionalStruct covers "special" variables in terms of bindings
type additionalStruct struct {
	Name    string
	Command string
	Cmdline string
	Viper   string
	Env     string
}

// createBindings executes the binding and sets up the environmental
// variables, as well as viper defaults as specified in the configuration
// definition tags.
func createBindings(tracer func(int, string, ...interface{}), cfg interface{},
	vp *viper.Viper, p map[string]*cobra.Command, env, docs map[string]string,
	appname string, seenFields map[string]bool) ([]additionalStruct, error) {

	additional := []additionalStruct{}
	tp := reflect.TypeOf(cfg)
	if !tp.Implements(baseInterface) {
		// Should never happen for users due to the sealed interface
		return nil, fmt.Errorf("Internal error: eindings should be called only for config objects")
	}

	t := tp.Elem()
	v := reflect.ValueOf(cfg).Elem()

	// We don't want users to be able to have fields in their structs with our
	// same names, exported or unexported.
	for i := 0; i < v.NumField(); i++ {
		x := t.Field(i)

		if seenFields[x.Name] {
			return nil, fmt.Errorf("Field collision on field %s, user configurations cannot shadow base configuration fields", x.Name)
		}
		seenFields[x.Name] = true

		// Ignore ourselves if we are in an embedded struct
		if x.Type == basePType {
			continue
		}

		// Ignore unexported fields
		r, _ := utf8.DecodeRuneInString(x.Name)
		if !unicode.IsUpper(r) {
			tracer(1, "Unexported field %s, skipping", x.Name)
			continue
		}

		field := v.FieldByName(x.Name)

		cobra, vipername, viperenv, err := parseTags(x)
		if err != nil {
			return nil, err
		}

		if viperenv != "" {
			var ok bool
			switch field.Kind() {
			case reflect.String,
				reflect.Bool,
				reflect.Float32, reflect.Float64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ok = true
			case reflect.Ptr:
				if field.Type().Implements(flagInterface) {
					ok = true
				}
				if field.Type().Implements(unmarshalInterface) {
					ok = true
				}
			default:
				if reflect.PtrTo(field.Type()).Implements(flagInterface) {
					// Typically flags seem to always be *Flag
					ok = true
				}
				if reflect.PtrTo(field.Type()).Implements(unmarshalInterface) {
					ok = true
				}
			}
			if !ok {
				return nil, fmt.Errorf("Cannot create flag %s, has environment %s but would not be able to be deserialized since it has no Set or UnmarshalText", x.Name, viperenv)
			}
		}

		tracer(1, "Looking at %s -> %s / %s / %s",
			x.Name, vipername, viperenv, cobra)

		if cobra == "" && (vipername != "" || viperenv != "") {
			tracer(1, "Not exposed on the command line, but env %s or cfg %s", vipername, viperenv)
			additional = append(additional, additionalStruct{
				Name:  x.Name,
				Viper: vipername,
				Env:   viperenv,
			})
			continue
		}

		var checked bool
		for ic, cc := range strings.Split(cobra, sepMultipleCmds) {
			if ic > 0 && vipername != "" {
				vipername = ""
			}

			// Assume it's always cobracmdname|cobraoptname|cobraoptshort
			ccobra := strings.Split(cc, sepCmdParts)
			if len(ccobra) != 3 {
				return nil, fmt.Errorf("Internal error, malformed cmdline tag %s", cc)
			}

			// For flags not exposed in cmd + cfg we neeed some special
			// processing on load, save them so we can do so. Also make sure
			// that the type can in fact be deserialized and serialized.
			if ccobra[2] == "none" || ccobra[2] == "custom" {
				if vipername != "" || viperenv != "" {
					if !checked && ccobra[2] != "custom" {
						checked = true
						var ok bool
						switch field.Kind() {
						case reflect.String,
							reflect.Bool,
							reflect.Float32, reflect.Float64,
							reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
							reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							ok = true
						case reflect.Slice:
							switch field.Type().Elem().Kind() {
							case reflect.Int:
								ok = true
							case reflect.String:
								ok = true
							}
						case reflect.Struct:
							if reflect.PtrTo(field.Type()).Implements(unmarshalInterface) && reflect.PtrTo(field.Type()).Implements(marshalInterface) {
								ok = true
							}
							if reflect.PtrTo(field.Type()).Implements(flagInterface) {
								// Same as above, not typically found
								ok = true
							}
						case reflect.Ptr:
							if field.Type().Implements(flagInterface) {
								// Same as above, not typically found
								ok = true
							}
							if field.Type().Implements(unmarshalInterface) && field.Type().Implements(marshalInterface) {
								ok = true
							}
						}

						if !ok {
							return nil, fmt.Errorf("Cannot create configuration file variable %s, unsupported type %s", x.Name, field.Type())
						}
					}

					tracer(1, "Not exposed on the command line, env: %s cfg: %s", viperenv, vipername)
					additional = append(additional, additionalStruct{
						Name:  x.Name,
						Viper: vipername,
						Env:   viperenv,
					})
				} else {
					tracer(1, "No env or cfg, skipping")
				}
				continue
			} else {
				if ccobra[1] == "" {
					return nil, fmt.Errorf("Flag %s has empty commandline name but is not set as none/custom", x.Name)
				}
			}

			if vipername == "" && viperenv != "" {
				tracer(1, "Not exposed in the configuration file but present with env: %s", viperenv)
				additional = append(additional, additionalStruct{
					Name:    x.Name,
					Cmdline: ccobra[1],
					Command: ccobra[0],
					Env:     viperenv,
				})
			}

			cmd, ok := p[ccobra[0]]
			if !ok {
				return nil, fmt.Errorf("Internal error, cannot find cmd for %s (%s, %s)",
					ccobra[0], ccobra[1], x.Name)
			}

			if err := doBind(tracer, vp, cmd, field, vipername, viperenv, x.Name,
				ccobra[1], ccobra[2], env, docs, appname); err != nil {
				return nil, err
			}

			if ccobra[1] != "" {
				if pf := cmd.Flag(ccobra[1]); pf != nil {
					pf.Annotations = map[string][]string{"greeneryVar": {x.Name}}
				} else {
					// Should not happen as we just bound it
					tracer(1, "Cannot find flag for %s on %s", ccobra[1], cmd.Name())
				}
			}

		}
	}

	return additional, nil
}

// doBind binds the specified variable, separating to make createBindings not
// as super long
func doBind(tracer func(int, string, ...interface{}), v *viper.Viper, cmd *cobra.Command,
	field reflect.Value, vipername, viperenv, varname, name, short string,
	env, docs map[string]string, appname string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Cobra/viper panic, let's catch it and override any existing error
			err = fmt.Errorf("Error while creating flag for variable %s: %v", varname, r)
		}
	}()

	var hidden bool
	if short == "hidden" {
		short = ""
		hidden = true
	}

	// Check if we have a flag value or a base type. Also save the stringer
	// method for later since we might need it for the viper default value.
	setter, stringer := getSetterStringer(field)
	doc, ok := docs[varname]
	if !ok {
		err = fmt.Errorf("No documentation for command line parameter %s (variable %s)", name, varname)
		return
	}

	// Create the cobra flag, either custom, since it implements
	// flagInterface, or a base type.
	if setter.Kind() != reflect.Invalid {
		// Must be first, as some fields with base types might be flags but
		// should be created as flags (say port which is a uint16 really)
		tracer(1, "Will create a varp flag for %s", varname)
		var cflag pflag.Value
		var ok bool
		if field.Kind() == reflect.Ptr {
			cflag, ok = field.Interface().(pflag.Value)
		} else {
			cflag, ok = field.Addr().Interface().(pflag.Value)
		}

		if !ok {
			// Should never happen unless pflag.Value stops implementing
			// flag.Value
			err = fmt.Errorf("Internal error in doBind, cannot cast the flag")
			return
		}
		cmd.PersistentFlags().VarP(cflag, name, short, doc)
	} else {
		switch field.Kind() {
		case reflect.String:
			tracer(1, "Will create a string flag for %s as %s %s",
				varname, name, short)
			cmd.PersistentFlags().StringVarP(field.Addr().Interface().(*string),
				name, short, field.Interface().(string), doc)
		case reflect.Bool:
			tracer(1, "Will create a bool flag for %s", varname)
			cmd.PersistentFlags().BoolVarP(field.Addr().Interface().(*bool),
				name, short, field.Interface().(bool), doc)
		case reflect.Float32:
			tracer(1, "Will create an float32 flag for %s", varname)
			cmd.PersistentFlags().Float32VarP(field.Addr().Interface().(*float32),
				name, short, field.Interface().(float32), doc)
		case reflect.Float64:
			tracer(1, "Will create an float64 flag for %s", varname)
			cmd.PersistentFlags().Float64VarP(field.Addr().Interface().(*float64),
				name, short, field.Interface().(float64), doc)
		case reflect.Int:
			tracer(1, "Will create an int flag for %s", varname)
			cmd.PersistentFlags().IntVarP(field.Addr().Interface().(*int),
				name, short, field.Interface().(int), doc)
		case reflect.Int8:
			tracer(1, "Will create an int8 flag for %s", varname)
			cmd.PersistentFlags().Int8VarP(field.Addr().Interface().(*int8),
				name, short, field.Interface().(int8), doc)
		case reflect.Int16:
			tracer(1, "Will create an int16 flag for %s", varname)
			cmd.PersistentFlags().Int16VarP(field.Addr().Interface().(*int16),
				name, short, field.Interface().(int16), doc)
		case reflect.Int32:
			tracer(1, "Will create an int32 flag for %s", varname)
			cmd.PersistentFlags().Int32VarP(field.Addr().Interface().(*int32),
				name, short, field.Interface().(int32), doc)
		case reflect.Int64:
			tracer(1, "Will create an int64 flag for %s", varname)
			cmd.PersistentFlags().Int64VarP(field.Addr().Interface().(*int64),
				name, short, field.Interface().(int64), doc)
		case reflect.Uint:
			tracer(1, "Will create an uint flag for %s", varname)
			cmd.PersistentFlags().UintVarP(field.Addr().Interface().(*uint),
				name, short, field.Interface().(uint), doc)
		case reflect.Uint8:
			tracer(1, "Will create an uint8 flag for %s", varname)
			cmd.PersistentFlags().Uint8VarP(field.Addr().Interface().(*uint8),
				name, short, field.Interface().(uint8), doc)
		case reflect.Uint16:
			tracer(1, "Will create an uint16 flag for %s", varname)
			cmd.PersistentFlags().Uint16VarP(field.Addr().Interface().(*uint16),
				name, short, field.Interface().(uint16), doc)
		case reflect.Uint32:
			tracer(1, "Will create an uint32 flag for %s", varname)
			cmd.PersistentFlags().Uint32VarP(field.Addr().Interface().(*uint32),
				name, short, field.Interface().(uint32), doc)
		case reflect.Uint64:
			tracer(1, "Will create an uint64 flag for %s", varname)
			cmd.PersistentFlags().Uint64VarP(field.Addr().Interface().(*uint64),
				name, short, field.Interface().(uint64), doc)
		default:
			return fmt.Errorf("Cannot create a flag for %s/%v, unsupported type",
				varname, field.Kind())
		}
	}

	if hidden {
		if err = cmd.PersistentFlags().MarkHidden(name); err != nil {
			// Should not happen
			return
		}
	}

	// Bind the cobra flag to Viper
	if vipername != "" {
		f := cmd.PersistentFlags().Lookup(name)
		if f == nil {
			// Should not happen if the previous code did its job.
			return fmt.Errorf("Internal error, could not find cmdline flag %s in %s",
				name, varname)
		}

		if err = v.BindPFlag(vipername, f); err != nil {
			// Should not happen
			return
		}

		if stringer.Kind() != reflect.Invalid {
			rv := stringer.Call([]reflect.Value{})[0].String()
			v.SetDefault(vipername, rv)
			tracer(1, "Setting default for %s to %s", vipername, rv)
		} else {
			v.SetDefault(vipername, field.Interface())
			tracer(1, "Setting default for %s to %v", vipername,
				field.Interface())
		}
	} else {
		tracer(1, "No viper binding %s %s", viperenv, vipername)
		return
	}

	// Set up the viper env mapping if required, cannot use viper for
	// non-cobra and only env parameters unfortunately :(
	if viperenv != "" {
		ename := appname + "_" + viperenv
		tracer(1, "Binding %s %s", ename, vipername)
		if err = v.BindEnv(vipername, ename); err != nil {
			// Should not hopefully happen
			return
		}

		env[ename] = os.Getenv(ename)
	} else {
		tracer(1, "Skipping as viperenv is empty for %s", vipername)
	}
	return nil
}

// parseTags returns the various parts of our tag
func parseTags(x reflect.StructField) (string, string, string, error) {
	name := x.Name
	tag := x.Tag.Get("greenery")

	if tag == "" {
		return "", "", "", fmt.Errorf("could not find a greenery tag for %s", name)
	}

	tags := strings.Split(tag, sepTag)
	if len(tags) != 3 {
		return "", "", "", fmt.Errorf("Invalid tag for %s, found %d parts instead of 3 in %s", name, len(tags), tag)
	}

	vipername := strings.TrimSpace(tags[1])
	if vipername != "" {
		if strings.HasPrefix(vipername, sepKeyParts) {
			vipername = strings.TrimLeft(vipername, sepKeyParts)
		} else {
			if !strings.Contains(vipername, sepKeyParts) {
				return "", "", "", fmt.Errorf("Invalid config file tag for '%s', no %s present in '%s'", name, sepKeyParts, vipername)
			}
		}

		if strings.Count(vipername, sepKeyParts) > 1 {
			return "", "", "", fmt.Errorf("Invalid config file tag for '%s', more than two %s present in '%s%s'", name, sepKeyParts, sepKeyParts, vipername)
		}
	}

	return strings.TrimSpace(tags[0]), vipername, strings.TrimSpace(tags[2]), nil
}
