package greenery

import (
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// Unfortunately viper seems to ignores all overflow/underflow errors for
// ints, so have to use cast directly. Also viper does not expose float32 or
// uints either, using cast directly helps us do that. Note that cast still
// needs some massaging for overflows.

// withErr formats an error message with the common error leading text.
func withErr(vipername string, err error) error {
	return errors.WithMessage(err, fmt.Sprintf("Cannot convert flag value %s", vipername))
}

// checkIntOverflow checks if the specified int would overflow an 8/16/32 int,
// the passed in error is to make the caller code simpler
func checkIntOverflow(err error, vipername string, which int, vc, v interface{}) error {
	if err != nil {
		return withErr(vipername, err)
	}

	val64 := cast.ToInt64(v)
	valc := cast.ToInt64(vc)

	// Need to catch the overflow cases, because unfortunately cast.ToInt will
	// not do so (for int8 if passing 256 it will return -1)
	if val64 != valc {
		return withErr(vipername, fmt.Errorf("%v would overflow an int%d", v, which))
	}
	return nil
}

// checkUIntOverflow checks if the specified uint would overflow an 8/16/32
// int, the passed in error is to make the caller code simpler
func checkUIntOverflow(err error, vipername string, which int, vc, v interface{}) error {
	if err != nil {
		return withErr(vipername, err)
	}

	val64 := cast.ToUint64(v)
	valc := cast.ToUint64(vc)

	// Need to catch the overflow cases, because unfortunately cast.ToUintx
	// seems not to do so (for int8 if passing 256 it will return 0)
	if val64 != valc {
		return withErr(vipername, fmt.Errorf("%v would overflow an uint%d", v, which))
	}
	return nil
}

func getString(vipername string, v interface{}) (string, error) {
	val, err := cast.ToStringE(v)
	if err != nil {
		// Unlikely to fail
		err = withErr(vipername, err)
	}
	return val, err
}

func getStringSlice(vipername string, v interface{}) ([]string, error) {
	val, err := cast.ToStringSliceE(v)
	if err != nil {
		// Note it will convert anything to strings, so unlikely this will
		// fail ever.
		err = withErr(vipername, err)
	}
	return val, err
}

func getBool(vipername string, v interface{}) (bool, error) {
	val, err := cast.ToBoolE(v)
	if err != nil {
		err = withErr(vipername, err)
	}
	return val, err
}

func getFloat32(vipername string, v interface{}) (float32, error) {
	val, err := cast.ToFloat32E(v)
	if err != nil {
		err = withErr(vipername, err)
	}

	// Unfortunately cast.ToFloat32E just returns +Inf here, assume users
	// want to know it would overflow.
	if math.Abs(cast.ToFloat64(v)) > math.MaxFloat32 {
		err = withErr(vipername, fmt.Errorf("%v would overflow a float32 (max %v)", v, math.MaxFloat32))
	}

	return val, err
}

func getFloat64(vipername string, v interface{}) (float64, error) {
	val, err := cast.ToFloat64E(v)
	if err != nil {
		err = withErr(vipername, err)
	}
	return val, err
}

func getInt(vipername string, v interface{}) (int, error) {
	val, err := cast.ToIntE(v)
	if err != nil {
		err = withErr(vipername, err)
	}
	return val, err
}

func getInt8(vipername string, v interface{}) (int8, error) {
	val, err := cast.ToInt8E(v)
	return val, checkIntOverflow(err, vipername, 8, val, v)
}

func getInt16(vipername string, v interface{}) (int16, error) {
	val, err := cast.ToInt16E(v)
	return val, checkIntOverflow(err, vipername, 16, val, v)
}

func getInt32(vipername string, v interface{}) (int32, error) {
	val, err := cast.ToInt32E(v)
	return val, checkIntOverflow(err, vipername, 32, val, v)
}

func getInt64(vipername string, v interface{}) (int64, error) {
	val, err := cast.ToInt64E(v)
	if err != nil {
		err = withErr(vipername, err)
	}
	return val, err
}

func getIntSlice(vipername string, v interface{}) ([]int, error) {
	val, err := cast.ToIntSliceE(v)
	if err != nil {
		err = withErr(vipername, err)
	}
	return val, err
}

func getUint(vipername string, v interface{}) (uint, error) {
	val, err := cast.ToUintE(v)
	if err != nil {
		err = withErr(vipername, err)
	}
	return val, err
}

func getUint8(vipername string, v interface{}) (uint8, error) {
	val, err := cast.ToUint8E(v)
	return val, checkUIntOverflow(err, vipername, 8, val, v)
}

func getUint16(vipername string, v interface{}) (uint16, error) {
	val, err := cast.ToUint16E(v)
	return val, checkUIntOverflow(err, vipername, 16, val, v)
}

func getUint32(vipername string, v interface{}) (uint32, error) {
	val, err := cast.ToUint32E(v)
	return val, checkUIntOverflow(err, vipername, 32, val, v)
}

func getUint64(vipername string, v interface{}) (uint64, error) {
	val, err := cast.ToUint64E(v)
	if err != nil {
		err = withErr(vipername, err)
	}
	return val, err
}

// setField will reflect set a specific field with the relevant value v (viper
// value as opposed to the passed value if viper is not nil). Vipername is
// always passed to make error messages nicer.
func setField(cfg Config, field reflect.Value, vp *viper.Viper, v interface{}, vipername string) (err error) {
	if vp != nil {
		v = vp.Get(vipername)
	}

	var setter reflect.Value

	switch field.Kind() {
	case reflect.String:
		var vs string
		if vs, err = getString(vipername, v); err != nil {
			// Unlikely to fail
			return
		}
		cfg.Tracef("Will assign string: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Bool:
		var vs bool
		if vs, err = getBool(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign bool: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Float32:
		var vs float32
		if vs, err = getFloat32(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign float32: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Float64:
		var vs float64
		if vs, err = getFloat64(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign float64: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Uint:
		var vs uint
		if vs, err = getUint(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign uint: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Uint8:
		var vs uint8
		if vs, err = getUint8(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign uint8: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Uint16:
		var vs uint16
		if vs, err = getUint16(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign uint16: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Uint32:
		var vs uint32
		if vs, err = getUint32(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign uint32: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Uint64:
		var vs uint64
		if vs, err = getUint64(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign uint64: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Int:
		var vs int
		if vs, err = getInt(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign int: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Int8:
		var vs int8
		if vs, err = getInt8(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign int8: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Int16:
		var vs int16
		if vs, err = getInt16(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign int16: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Int32:
		var vs int32
		if vs, err = getInt32(vipername, v); err != nil {
			return
		}
		cfg.Tracef("Will assign int32: %v", vs)
		field.Set(reflect.ValueOf(vs))
	case reflect.Int64:
		// time.Duration serializes/deserializes to int64 apparently, so
		// special case it.
		var vs int64
		if vs, err = getInt64(vipername, v); err != nil {
			return
		}
		if field.Type() == reflect.TypeOf(time.Second) {
			cfg.Tracef("Will assign time.Duration: %v", vs)
			field.Set(reflect.ValueOf(time.Duration(vs)))
		} else {
			cfg.Tracef("Will assign int64: %v", vs)
			field.Set(reflect.ValueOf(vs))
		}
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(v) {
			// time.Time is unmarshaled as a struct in TOML
			cfg.Tracef("Same struct type, assigning as-is: %v", v)
			field.Set(reflect.ValueOf(v))
		} else {
			// but when it's an environment variable of course it's a
			// string. Doubtful users would want to use this, much easier to
			// simply create a Flag type with Set, but just in case don't
			// hardcode time.Time and go through unmarshal instead.
			if reflect.PtrTo(field.Type()).Implements(unmarshalInterface) {
				setter = field.Addr().MethodByName("UnmarshalText")
			} else {
				// Should not happen
				return fmt.Errorf(
					"Unsupported struct type %v for %s, custom flag/environment requiring custom deserialization?", field.Type(), vipername)
			}
		}
	case reflect.Ptr:
		if field.Type() == reflect.TypeOf(v) {
			// TODO: Have not been able to exercise this
			cfg.Tracef("Same ptr type, assigning as-is: %v", v)
			field.Set(reflect.ValueOf(v))
		} else {
			if field.Elem().Type() == reflect.TypeOf(v) {
				cfg.Tracef("We have a ptr to the same type, assigning as-is: %v", v)
				// Ok, we have a ptr to the same in our config, just assign it
				field.Elem().Set(reflect.ValueOf(v))
			} else {
				if field.Type().Implements(unmarshalInterface) {
					setter = field.MethodByName("UnmarshalText")
				} else {
					// Should not happen
					return fmt.Errorf(
						"Unsupported struct type %v for %s, custom flag/environment requiring custom deserialization?", field.Type(), vipername)
				}
			}
		}
	case reflect.Slice:
		switch field.Type().Elem().Kind() {
		case reflect.Int:
			var vs []int
			if vs, err = getIntSlice(vipername, v); err != nil {
				return
			}
			field.Set(reflect.ValueOf(vs))
		case reflect.String:
			var vs []string
			if vs, err = getStringSlice(vipername, v); err != nil {
				// Unlikely to fail
				return
			}
			field.Set(reflect.ValueOf(vs))
		default:
			// Should not happen as binding.go does check for this
			return fmt.Errorf(
				"Only int and string slices are supported, field %s with value %v", vipername, vp.Get(vipername))
		}
	default:
		// Should not happen, but just in case
		return fmt.Errorf(
			"Unsupported type %s for %s, custom flag requiring custom deserialization?: %v", field.Kind(), vipername, vp.Get(vipername))
	}

	if setter.Kind() != reflect.Invalid {
		cfg.Trace("Have a custom unmarshaler")

		// If it's a string, need to convert to []byte
		var vv reflect.Value
		if vvs, ok := v.(string); ok {
			vv = reflect.ValueOf([]byte(vvs))
		} else {
			// TODO: Haven't been able to exercise this, always strings
			if vvb, ok := v.([]byte); ok {
				vv = reflect.ValueOf(vvb)
			} else {
				// Shouldn't happen
				return fmt.Errorf(
					"Unsupported value %v for %s cannot unmarshal", v, vipername)
			}
		}

		p := []reflect.Value{vv}
		cfg.Tracef("Value is a struct, field isn't but has unmarshal, call it to set %v in %s", v, vipername)
		rv := setter.Call(p)
		rve := rv[len(rv)-1]
		if !rve.IsNil() {
			if rerr, ok := rve.Interface().(error); ok {
				return errors.WithMessage(rerr, "Cannot convert flag value via unmarshaling")
			}

			// Should not happen
			return fmt.Errorf("Internal error setting a value, wrong error type returned by UnmarshalText %v", rve.Interface())
		}
	}

	return
}
