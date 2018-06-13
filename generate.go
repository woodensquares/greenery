package greenery

import (
	"encoding"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shibukawa/configdir"
	"github.com/woodensquares/greenery/internal/doc"
)

func commentify(s, comment string) string {
	var o string
	if s == "" {
		return ""
	}

	for _, l := range strings.Split(strings.TrimSuffix(s, "\n"), "\n") {
		o = o + comment + l + "\n"
	}
	return o
}

func fileName(dname, name, location string) (wanted string, err error) {
	if name == "" {
		name = dname + ".toml"
	}

	if !strings.Contains(name, string(os.PathSeparator)) {
		switch location {
		case "cwd":
			wanted, err = os.Getwd()
			if err != nil {
				err = errors.WithMessage(err,
					"Cannot access the current directory")
				return
			}
			wanted = path.Join(wanted, dname+".toml")
		case "user":
			wanted = configdir.New(dname, name).QueryFolders(configdir.Global)[0].Path
		case "system":
			wanted = configdir.New(dname, name).QueryFolders(configdir.System)[0].Path
		}
	}

	return
}

type cfgLine struct {
	parent    string
	doc       string
	value     string
	child     string
	skipvalue bool
}

var marshalInterface = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
var unmarshalInterface = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

func serializeHelper(field reflect.Value, extra bool, name string) (string, error) {
	var rv string
	switch field.Kind() {
	case reflect.String:
		rv = "\"" + field.String() + "\""
	case reflect.Bool:
		rv = strconv.FormatBool(field.Bool())
	case reflect.Float32, reflect.Float64:
		rv = fmt.Sprintf("%v", field.Float())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv = fmt.Sprintf("%v", field.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv = fmt.Sprintf("%v", field.Int())
	default:
		if !extra {
			return "", fmt.Errorf(
				"Internal error, unsupported type %s for %s", field.Kind(), name)
		}
	}
	return rv, nil
}

func paramHelper(cfg Config, docs, fallback map[string]string, val reflect.Value, x reflect.StructField) (*cfgLine, error) {
	field := val.FieldByName(x.Name)
	extra := false

	cobra, vipername, _, _ := parseTags(x)
	if vipername == "" {
		return nil, nil
	}

	if strings.HasSuffix(cobra, sepCmdParts+"custom") {
		extra = true
	}

	// Check if we have a flag or not.
	var stringer reflect.Value
	var marshaler reflect.Value
	if field.Addr().Elem().Kind() == reflect.Ptr {
		if field.Type().Implements(flagInterface) {
			stringer = field.Addr().Elem().MethodByName("String")
		}
		if field.Type().Implements(marshalInterface) {
			marshaler = field.Addr().Elem().MethodByName("MarshalText")
		}
	} else {
		if reflect.PtrTo(field.Type()).Implements(flagInterface) {
			stringer = field.Addr().MethodByName("String")
		}
		if reflect.PtrTo(field.Type()).Implements(marshalInterface) {
			marshaler = field.Addr().MethodByName("MarshalText")
		}
	}

	// All viper tags are parent.child. No support for grandparents yet.
	vv := strings.Split(vipername, sepKeyParts)
	var parent, child string
	switch len(vv) {
	case 1:
		parent = ""
		child = vv[0]
	case 2:
		parent = vv[0]
		child = vv[1]
	default:
		// Should not happen, bind should have caught this
		return nil, fmt.Errorf("Internal error, malformed cfgfile tag %s", vipername)
	}

	var rv string
	if marshaler.Kind() != reflect.Invalid {
		res := marshaler.Call([]reflect.Value{})
		rve := res[1]
		if !rve.IsNil() {
			// Should never happen
			return nil, rve.Interface().(error)
		}
		rv = string(res[0].Bytes())
	} else if stringer.Kind() != reflect.Invalid {
		rv = "\"" + stringer.Call([]reflect.Value{})[0].String() + "\""
	} else {
		var err error
		if field.Kind() == reflect.Slice {
			rv = "[ "
			var rve []string
			for ie := 0; ie < field.Len(); ie++ {
				var s string
				if s, err = serializeHelper(field.Index(ie), extra, x.Name); err != nil {
					// Should not happen, bind should have caught this
					return nil, errors.WithMessage(err, "in a slice context only base types are supported")
				}
				rve = append(rve, s)
			}
			rv = rv + strings.Join(rve, ", ") + " ]"
		} else {
			if rv, err = serializeHelper(field, extra, x.Name); err != nil {
				return nil, err
			}
		}
	}

	// If the user did not provide any config file specific docs, use the
	// cmdline docs as a default. Allow not having config docs for specific
	// variables if the user so decides.
	d, ok := docs[x.Name]
	if !ok {
		d, ok = fallback[x.Name]
		if !ok {
			return nil, fmt.Errorf("Config file variable %s has no documentation, set it to empty if it's meant to not have it", x.Name)
		}
	}

	if extra {
		cfg.Tracef("Init: Extra values are special, no default, just doc for %s.%s", parent, child)
		if d != "" {
			return &cfgLine{
				parent:    parent,
				child:     child,
				doc:       d,
				skipvalue: true,
			}, nil
		}

		cfg.Tracef("No doc was provided for the special variable, skipping")
		return nil, nil
	}

	cfg.Tracef("Init: Setting default for %s.%s to %s", parent, child, rv)
	return &cfgLine{
		parent: parent,
		child:  child,
		doc:    d,
		value:  rv,
	}, nil
}

func fileContents(cfg Config) (confText string, err error) {
	// Add everything to the template, we are interested only in viper
	// fields here, so only anything with a viper tag set. All default values
	// are already set by cfg.Initialize()
	t := reflect.TypeOf(cfg).Elem()
	val := reflect.ValueOf(cfg).Elem()

	_, d := cfg.GetDocs()
	docs := d.ConfigFile
	fallback := d.CmdLine
	if t, ok := docs[doc.ConfigHeader]; ok {
		confText = commentify(t, "# ") + "\n"
	}

	cfgVars := make([]*cfgLine, 0)
	for i := 0; i < val.NumField(); i++ {
		x := t.Field(i)

		if x.Type == basePType {
			v2 := val.Field(i).Elem()
			for i2 := 0; i2 < baseType.NumField(); i2++ {
				x2 := baseType.Field(i2)

				var n *cfgLine
				if n, err = paramHelper(cfg, docs, fallback, v2, x2); err != nil {
					// Currently should not happen
					return
				}

				if n != nil {
					cfgVars = append(cfgVars, n)
				}
			}
			continue
		}

		var n *cfgLine
		if n, err = paramHelper(cfg, docs, fallback, val, x); err != nil {
			return
		}

		if n != nil {
			cfgVars = append(cfgVars, n)
		}
	}

	sort.Slice(cfgVars, func(l, r int) bool {
		// Root section variables should always come first no matter what
		switch strings.Compare(cfgVars[l].parent, cfgVars[r].parent) {
		case -1:
			if cfgVars[r].parent == doc.ConfigHeader {
				return false
			}
			return true
		case 1:
			if cfgVars[l].parent == doc.ConfigHeader {
				return true
			}
			return false
		}

		// skipvalue always come last, as they usually are blocks of
		// documentation / sample
		if cfgVars[l].skipvalue && !cfgVars[r].skipvalue {
			return false
		}
		if !cfgVars[l].skipvalue && cfgVars[r].skipvalue {
			return true
		}

		// Sort the rest alphabetically, leaving equal values as-is
		if strings.Compare(cfgVars[l].child, cfgVars[r].child) == 1 {
			return false
		}

		return true
	})

	var cparent string
	seen := map[string]string{}
	seenDoc := map[string]string{}
	for _, v := range cfgVars {
		cand := v.parent + "!" + v.child
		if pValue, next := seen[cand]; next {
			if pValue != v.value {
				err = fmt.Errorf("More than one configuration variable corresponds to config file variable %s with different defaults: %s and %s for example", v.parent+"."+v.child, v.value, pValue)
				return
			}

			if seenDoc[cand] != v.doc {
				err = fmt.Errorf("More than one configuration variable corresponds to config file variable %s with different doc lines: \"%s\" and \"%s\" for example", v.parent+"."+v.child, v.doc, seenDoc[cand])
				return
			}
			continue
		}
		seen[cand] = v.value
		seenDoc[cand] = v.doc

		if cparent != v.parent {
			cparent = v.parent

			confText += "\n"
			if t, ok := docs[cparent+sepKeyParts]; ok {
				confText += commentify(t, "# ")
			}

			confText += "[" + cparent + "]\n"
		}

		if v.skipvalue {
			confText += commentify(v.doc, "")
		} else {
			confText += commentify(v.doc, "# ")
			confText += v.child + " = " + v.value + "\n"
		}
	}

	return
}

// initCfgFile creates a new default config file in the specified location, it
// will refuse to overwrite an already existing configuration file.
func initCfgFile(icfg Config, cfg *BaseConfig) (used string, err error) {
	var wanted string
	if wanted, err = fileName(cfg.s_appName, icfg.GetConfigFile(), cfg.CfgLocation.Value); err != nil {
		return
	}

	used = wanted
	cfgDir := filepath.Dir(wanted)

	// The directory might not exist in the xdg case
	err = cfg.s_fs.MkdirAll(cfgDir, 0700)
	if err != nil {
		return "", errors.WithMessage(err, fmt.Sprintf("Cannot create config directory %s", cfgDir))
	}

	perms := os.O_WRONLY | os.O_CREATE
	if cfg.CfgForce {
		perms = perms | os.O_TRUNC
	} else {
		perms = perms | os.O_EXCL
	}

	f, err := cfg.s_fs.OpenFile(wanted, perms, 0644)
	if err != nil {
		return "", errors.WithMessage(err, fmt.Sprintf("Cannot create config file %s", wanted))
	}
	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	var confText string
	if confText, err = fileContents(icfg); err != nil {
		return
	}

	var tmpl *template.Template
	if tmpl, err = template.New("initconfig").Funcs(template.FuncMap{
		"date": func() string { return time.Now().Format(time.RFC3339) },
	}).Parse(confText); err != nil {
		// Should never happen
		return "", errors.WithMessage(err, fmt.Sprintf("Internal error parsing the template"))
	}

	// Our custom flags, like loglevel, define String() so they will be
	// properly converted to strings without having to do anything else.
	err = tmpl.Execute(f, cfg)
	if err != nil {
		// Should not happen
		return "", errors.WithMessage(err, fmt.Sprintf("Internal error executing the template"))
	}

	return
}
