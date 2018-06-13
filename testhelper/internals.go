package testhelper

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
)

var basePType = reflect.TypeOf(&greenery.BaseConfig{})
var baseType = reflect.TypeOf(greenery.BaseConfig{})

func removeFiles(t *testing.T, cleanup *[]string, mtx sync.Locker) {
	if mtx != nil {
		mtx.Lock()
		defer func() {
			mtx.Unlock()
		}()

		for _, f := range *cleanup {
			err := os.Remove(f)
			require.NoError(t, err)
			t.Logf("Removed temp file: %s", f)
		}
		(*cleanup) = []string{}
	}
}

func compareField(t *testing.T, tc *TestCase, f reflect.StructField, v reflect.Value, comparedValues map[string]bool, baseValues map[string]interface{}, comps map[string]CompareFunc) {
	if isUnexported(f.Name) {
		return
	}

	vf := v.FieldByName(f.Name)

	// Did the user have a value or comparer listed for this field?
	if cv, ok := tc.ExpectedValues[f.Name]; ok {
		cf := cv.Compare
		if cf == nil {
			cf = comps[f.Name]
		}

		if cv.Accessor != "" {
			t.Logf(spew.Sdump(cv.Value))
			t.Logf(spew.Sdump(vf.Interface()))
			val := getValue(t, f.Name, cv.Accessor, vf.Interface())
			require.Equal(t, cv.Value, val, "Field %s", f.Name)
		} else {
			if cf != nil {
				var err error
				if cv.DefaultValue {
					err = cf(t, f.Name, baseValues[f.Name], vf.Interface())
				} else {
					err = cf(t, f.Name, cv.Value, vf.Interface())
				}
				require.NoError(t, err, "Custom compare function for %s", f.Name)
			} else {
				require.Equal(t, cv.Value, vf.Interface(), "Field %s", f.Name)
			}
		}

		comparedValues[f.Name] = true
	} else {
		// In this case compare to the "clean cfg" value
		if cf, ok := comps[f.Name]; ok {
			require.NoError(t, cf(t, f.Name, baseValues[f.Name], vf.Interface()))
		} else {
			require.Equal(t, baseValues[f.Name], vf.Interface(), "Field %s", f.Name)
		}
	}
}

func addCleanupFiles(t *testing.T, files []string, cleanup *[]string, mtx sync.Locker) {
	if cleanup != nil {
		if mtx != nil {
			mtx.Lock()
			t.Log("Locked cleanup mtx")
			defer func() {
				t.Log("Unlocked cleanup mtx")
				mtx.Unlock()
			}()

			for _, f := range files {
				t.Logf("Added %s to the cleanup list", f)
				*cleanup = append(*cleanup, f)
			}
		}
	}
}

// TempFileT is a utility function typically used to create temporary
// configuration files to test with.
func TempFileT(t *testing.T, af afero.Fs, prefix, suffix, contents string, cleanup *[]string, mtx sync.Locker, real bool) (name string, err error) {
	var cfgFile afero.File

	// Since ioutil.Tempfile does not unfortunately allow for a suffix, create
	// a normal tempfile, then use that name + the wanted suffix as the temp
	// file, if it exists try again.

	// Seems very unlikely we'd fail more than 100 times here, especially on
	// afero.
	for i := 0; i < 100; i++ {
		cfgFile, err = afero.TempFile(af, "", prefix)
		if suffix == "" {
			// If the user had no suffix, we're done
			break
		}

		var original = cfgFile
		defer func() {
			err = original.Close()
			require.NoError(t, err)
			err = af.Remove(original.Name())
			require.NoError(t, err)
			t.Logf("Suffix asked, removed temp file %s and recreated with the suffix", original.Name())
		}()

		cfgFile, err = af.OpenFile(original.Name()+suffix,
			os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)

		if err != nil {
			if os.IsExist(err) {
				// Bad luck, let's retry
				continue
			} else {
				return
			}
		}

		// If we're here we're done
		break
	}

	tName := cfgFile.Name()
	if real {
		t.Logf("Add temp file to the cleanup list: %s", tName)
		addCleanupFiles(t, []string{tName}, cleanup, mtx)
	}

	_, err = cfgFile.Write([]byte(contents))
	require.NoError(t, err)

	err = cfgFile.Close()
	require.NoError(t, err)

	t.Logf("Written the file contents: %s", contents)
	return cfgFile.Name(), nil
}

func setLogging(t *testing.T, tc *TestCase, cmdLine []string, cleanup *[]string, mtx sync.Locker) (string, []string) {
	logFile, err := TempFileT(t, tc.af, "tlog", ".log", "", cleanup, mtx, tc.RealFilesystem)
	require.NoError(t, err)

	newCmdLine := []string{
		"--log-file",
		logFile,
	}

	return logFile, append(newCmdLine, cmdLine...)
}

func regexpMatches(regex, s string) bool {
	r := regexp.MustCompile(regex)
	return (r.FindStringIndex(s) != nil)
}

func verifyValues(t *testing.T, tc *TestCase, cfg, cleanCfg greenery.Config, cleanValues map[string]interface{}, comps map[string]CompareFunc) {
	comparedValues := map[string]bool{}
	v := reflect.ValueOf(cfg).Elem()

	cleanCfgValue := reflect.ValueOf(cleanCfg).Elem()
	cfgType := reflect.TypeOf(cleanCfgValue.Interface())

	for i := 0; i < cfgType.NumField(); i++ {
		f := cfgType.Field(i)

		if isUnexported(f.Name) {
			continue
		}

		// Base config values validation
		if f.Type == basePType {
			vbase := v.Field(i).Elem()

			for ibase := 0; ibase < baseType.NumField(); ibase++ {
				compareField(t, tc, baseType.Field(ibase), vbase, comparedValues, cleanValues, comps)
			}
		} else {
			// Custom values
			compareField(t, tc, f, v, comparedValues, cleanValues, comps)
		}
	}

	if len(comparedValues) != len(tc.ExpectedValues) {
		errMsg := "The following fields were not found in the configuration: "
		notFound := []string{}
		for candidate := range tc.ExpectedValues {
			if _, ok := comparedValues[candidate]; !ok {
				notFound = append(notFound, candidate)
			}
		}

		require.FailNow(t, errMsg+strings.Join(notFound, ","))
	}
}

func verifyGoldFile(t *testing.T, tf *TestFile, contents []byte) {
	gcontents, err := ioutil.ReadFile(tf.Source)

	if *GoldUpdate {
		match := true
		if tf.Custom != nil {
			match = tf.Custom(t, gcontents, contents)
		} else {
			if len(contents) != len(gcontents) || err != nil {
				match = false
			} else {
				for i := range contents {
					if contents[i] != gcontents[i] {
						match = false
						break
					}
				}
			}
		}

		if !match {
			err = ioutil.WriteFile(tf.Source, contents, 0644)
			require.NoError(t, err, "Cannot update gold file %s", tf.Source)
			t.Logf("Gold file %s updated\n", tf.Source)
		} else {
			t.Logf("Gold file %s does not need updating\n", tf.Source)
		}
	} else {
		require.NoError(t, err)
		if tf.Custom != nil {
			require.True(t, tf.Custom(t, gcontents, contents), "The gold file %s contents differ (custom compare): gold\n%s\n\nfound\n%s", tf.Source, string(gcontents), string(contents))
		} else {
			require.EqualValues(t, gcontents, contents, "The gold file %s contents differ: %s", tf.Source, string(contents))
		}
	}
}

func verifyGoldFiles(t *testing.T, tc *TestCase) {
	for _, k := range tc.GoldFiles {
		location := k.Location
		if location[0] != filepath.Separator {
			var cwd string
			cwd, err := os.Getwd()
			require.NoError(t, err)
			location = filepath.Join(cwd, location)
		}

		contents, err := afero.ReadFile(tc.af, location)
		require.NoError(t, err, "Expected gold file %s was not created", location)
		s, err := tc.af.Stat(location)
		require.NoError(t, err, "Cannot stat expected gold file %s", location)
		require.EqualValues(t, k.Perms, s.Mode(), "Created gold file %s has the wrong mode %v, was expecting %v", location, s.Mode(), k.Perms)

		verifyGoldFile(t, &k, contents)
	}
}

func verifyGoldLog(t *testing.T, tc *TestCase, logFile string) {
	write := true

	var err error
	var actual, expected []byte
	var actualP map[string]interface{}
	var expectedP map[string]interface{}
	var expectedL, actualL []string

	// Meh, don't like using gotos that much, but saves a bunch of
	// nested ifs here to shortcircuit when we know we have to
	// overwrite the gold file because it doesn't match.
	actual, err = afero.ReadFile(tc.af, logFile)
	if *GoldUpdate && err != nil {
		goto mustupdate
	}
	require.NoError(t, err)
	actualL = strings.Split(string(actual), "\n")
	expected, err = ioutil.ReadFile(tc.GoldLog)
	if *GoldUpdate && err != nil {
		goto mustupdate
	}

	require.NoError(t, err)
	expectedL = strings.Split(string(expected), "\n")
	if *GoldUpdate && !assert.EqualValues(t, len(actualL), len(expectedL)) {
		goto mustupdate
	}
	require.EqualValues(t, len(actualL), len(expectedL), "The golden and produced log files are different lengths")

	// For every line compare, overriding ts since it
	// obviously changes every run.
	for i := 0; i < len(actualL)-1; i++ {
		err = json.Unmarshal([]byte(actualL[i]), &actualP)
		if err != nil {
			goto mustupdate
		}
		require.NoError(t, err, "Cannot unmarshal the log entry.")

		err = json.Unmarshal([]byte(expectedL[i]), &expectedP)
		if err != nil {
			goto mustupdate
		}
		require.NoError(t, err, "Cannot unmarshal the log entry.")
		expectedP["ts"] = actualP["ts"]

		if *GoldUpdate && !assert.EqualValues(t, actualP, expectedP) {
			goto mustupdate
		}
		require.EqualValues(t, actualP, expectedP, "The log values differ")
	}
	write = false

mustupdate:
	if *GoldUpdate {
		if write {
			err = ioutil.WriteFile(tc.GoldLog, actual, 0644)
			require.NoError(t, err)
			t.Logf("Log gold file %s updated\n", tc.GoldLog)
		} else {
			t.Logf("Log gold file %s does not need updating\n", tc.GoldLog)
		}
	}
}

func verifyStdOutErr(t *testing.T, tc *TestCase, outS, errS string) {
	if tc.GoldStdOut != nil {
		verifyGoldFile(t, tc.GoldStdOut, []byte(outS))
	} else {
		if tc.OutStdOutRegex == "" {
			require.EqualValues(t, tc.OutStdOut, outS, "stdout does not match what was requested")
		}
	}

	if tc.GoldStdErr != nil {
		verifyGoldFile(t, tc.GoldStdErr, []byte(errS))
	} else {
		if tc.OutStdErrRegex == "" {
			require.EqualValues(t, tc.OutStdErr, errS, "stderr does not match what was requested")
		}
	}

	if tc.OutStdOutRegex != "" {
		require.Regexp(t, tc.OutStdOutRegex, outS, "stdout does not match what was requested via regex")
	}

	if tc.OutStdErrRegex != "" {
		require.Regexp(t, tc.OutStdErrRegex, errS, "stderr does not match what was requested via regex")
	}
}

func getValues(cfg greenery.Config) map[string]interface{} {
	cfgValue := reflect.ValueOf(cfg).Elem()
	cfgType := reflect.TypeOf(cfgValue.Interface())
	values := map[string]interface{}{}
	for i := 0; i < cfgType.NumField(); i++ {
		f := cfgType.Field(i)
		if isUnexported(f.Name) {
			continue
		}

		if f.Type == basePType {
			vbase := cfgValue.Field(i).Elem()

			for ibase := 0; ibase < baseType.NumField(); ibase++ {
				fbase := baseType.Field(ibase)

				if isUnexported(fbase.Name) {
					continue
				}

				values[fbase.Name] = vbase.FieldByName(fbase.Name).Interface()
			}
		} else {
			values[f.Name] = cfgValue.Field(i).Interface()
		}
	}
	return values
}

func precreateFiles(t *testing.T, af afero.Fs, files []TestFile) []string {
	var err error
	var created []string

	for _, k := range files {
		var contents []byte
		if k.Source != "" {
			contents, err = ioutil.ReadFile(k.Source)
			require.NoError(t, err)
		} else {
			contents = k.Contents
		}

		location := k.Location
		if location[0] != filepath.Separator {
			var cwd string
			cwd, err = os.Getwd()
			require.NoError(t, err)
			location = filepath.Join(cwd, location)
		}

		err = afero.WriteFile(af, location, contents, k.Perms)
		require.NoError(t, err, "Cannot write precreate file %s", location)
		created = append(created, location)
	}

	return created
}

func verifyLog(t *testing.T, af afero.Fs, logFile string, loglines []LogLine, all bool) {
	logB, err := afero.ReadFile(af, logFile)
	require.NoError(t, err)
	llen := len(loglines)

	// Short circuit the empty log scenario
	if llen == 0 {
		require.EqualValues(t, 0, len(logB))
	} else {
		logSL := strings.Split(string(logB), "\n")
		logSL = logSL[:len(logSL)-1]

		if all {
			require.EqualValues(t, len(loglines), len(logSL), "If LogAllLines is set, the length of OutLogLines must be the same as the output")
		}

		var li int
		var payload map[string]interface{}
	outer:
		for _, v := range logSL {
			require.NoError(t, json.Unmarshal([]byte(v), &payload), "Cannot unmarshal the log entry, does this logger support structured output? (baseLogger doesn't): %s.", v)

			l := loglines[li]
			if l.Level != "" {
				if !assert.ObjectsAreEqualValues(l.Level, payload["level"]) {
					continue
				}
			}

			if l.Msg != "" {
				if !assert.ObjectsAreEqualValues(l.Msg, payload["msg"]) {
					continue
				}
			}

			if l.MsgRegex != "" {
				ms, ok := payload["msg"].(string)
				assert.True(t, ok, "The 'msg' log field is not a string")
				if !regexpMatches(l.MsgRegex, ms) {
					continue
				}
			}

			for kk, vv := range l.Custom {
				pl, ok := payload[kk]
				if !ok {
					t.Logf("Field %s does not exist in payload %v", kk, payload)
					continue outer
				}
				if !assert.ObjectsAreEqualValues(vv, pl) {
					t.Logf("Field %s value %v did not match %v", kk, vv, payload[kk])
					continue outer
				}
			}

			for kk, vv := range l.CustomRegex {
				ms, ok := payload[kk].(string)
				assert.True(t, ok, "The '%s' log field is not a string", kk)
				if !regexpMatches(vv, ms) {
					t.Logf("Regexp %s did not match %s", vv, ms)
					continue outer
				}
			}

			// If we are here everything matched, onto the next LogLine
			li++
			if li == llen {
				break
			}
		}

		require.EqualValues(t, li, llen, "Not all the requested loglines were matched, we are at LogLine %v", li)
	}

}

func isUnexported(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return !unicode.IsUpper(r)
}

func getValue(t *testing.T, name, accessor string, current interface{}) interface{} {
	cv := reflect.ValueOf(current)
	getter := cv.MethodByName(accessor)

	require.NotEqualf(t, reflect.Invalid, getter.Kind(), "Type %v to compare for variable %s does not have an accessor %s method!", cv.Type(), name, accessor)
	rv := getter.Call([]reflect.Value{})

	require.True(t, len(rv) > 0, "Getter for type %v for variable %s did not return a result", cv.Type(), name)
	return rv[0].Interface()
}

// CompareGetterToGetter is FIXME:DOC
func CompareGetterToGetter(t *testing.T, name string, expected, current interface{}) error {
	evalue := getValue(t, name, "GetTyped", expected)
	cvalue := getValue(t, name, "GetTyped", current)
	require.Equalf(t, evalue, cvalue, "Variable %s", name)
	return nil
}

// CompareValueToGetter is FIXME:DOC
func CompareValueToGetter(t *testing.T, name string, expected, current interface{}) error {
	cvalue := getValue(t, name, "GetTyped", current)
	require.Equalf(t, expected, cvalue, "Variable %s", name)
	return nil
}
