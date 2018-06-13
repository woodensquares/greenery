package greenery

import (
	"fmt"
	"net"
	"strconv"
)

// --------------------------------------------------------------------------

// IntValue is a flag struct containing an integer with customizeable minimum
// and maximum values.
type IntValue struct {
	Value int
	min   int
	max   int
	name  string
}

// NewIntValue returns a new IntValue flag, its value is set to min. Name is a
// name for this flag, typically the configuration variable name, and is used
// for error messages. Min and max are the values to be used as min/max values
// to validate against.
func NewIntValue(name string, min, max int) *IntValue {
	return &IntValue{
		name:  name,
		min:   min,
		max:   max,
		Value: min,
	}
}

// NewDefaultIntValue returns a new IntValue flag set to the specified "set"
// value. Name is a name for this flag, typically the configuration
// variable name, and is used for error messages. Min and max are the values
// to be used as min/max values to validate against.
func NewDefaultIntValue(name string, set, min, max int) *IntValue {
	v := NewIntValue(name, min, max)
	if err := v.SetInt(set); err != nil {
		panic(err.Error())
	}

	return v
}

// SetInt will set an integer and validate it for correctness
func (i *IntValue) SetInt(d int) (err error) {
	if d < i.min || d > i.max {
		return fmt.Errorf("Invalid value %d for variable %s, should be between %d and %d",
			d, i.name, i.min, i.max)
	}

	i.Value = d
	return
}

// Set will set a value parsing from a string, while validating for correctness
func (i *IntValue) Set(s string) (err error) {
	d, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("Variable %s, %s, cannot be converted to a number", i.name, s)
	}

	return i.SetInt(d)
}

// GetTyped is typically used for tests and returns the flag int value
func (i *IntValue) GetTyped() int {
	return i.Value
}

// Type will return a string describing the type of the flag, it is required
// to fulfill pflag.Value and will be printed in the help messages
func (i *IntValue) Type() string {
	return "int"
}

// String will return a string representation of the flag value
func (i *IntValue) String() string { return strconv.Itoa(i.Value) }

// UnmarshalText is used for TOML configuration file unmarshaling, and will
// set the value in the flag with validation.
func (i *IntValue) UnmarshalText(text []byte) error {
	return i.Set(string(text))
}

// MarshalText is used for TOML configuration file marshaling, it is used when
// generating the config files via config generate.
func (i *IntValue) MarshalText() (text []byte, err error) {
	return []byte(strconv.Itoa(int(i.Value))), nil
}

// --------------------------------------------------------------------------

// CustomStringHandler defines a function used to validate the custom string
// for correctness. This function will be passed, in order, the name of this
// flag, the value the user would like to set, and the arbitrary flag
// data. The function should return the value that will be set in the flag,
// which could be different from what the user sends, and an error if the
// value is not acceptable.
type CustomStringHandler func(string, string, interface{}) (string, error)

// CustomStringValue is a flag struct that has a customizeable validator /
// converter.
type CustomStringValue struct {
	Value   string
	handler CustomStringHandler
	name    string
	data    interface{}
}

// NewCustomStringValue returns a new CustomStringValue flag, the default
// value is set to whatever the validation function will return when passed
// the empty string. Name is a name for this flag, typically the configuration
// variable name, and is used for error messages. Data is arbitrary data that
// will be stored in the flag and passed to the custom validation function on
// validation calls.
func NewCustomStringValue(name string, validate CustomStringHandler, data interface{}) *CustomStringValue {
	return NewDefaultCustomStringValue(name, "", validate, data)
}

// NewDefaultCustomStringValue returns a new CustomStringValue flag with the
// default value being what the validation function will return when passed
// the string specified in set. Name is a name for this flag, typically the
// configuration variable name, and is used for error messages. Data is
// arbitrary data that will be stored in the flag and passed to the custom
// validation function on validation calls. If the set value does not fulfill
// the validate function, this function will panic.
func NewDefaultCustomStringValue(name string, set string, validate CustomStringHandler, data interface{}) *CustomStringValue {
	if validate != nil {
		var err error

		// Users can also get a different default value via their handler
		if set, err = validate(name, set, data); err != nil {
			panic(fmt.Sprintf("While creating variable %s, error %v", name, err.Error()))
		}
	}

	return &CustomStringValue{
		Value:   set,
		handler: validate,
		name:    name,
		data:    data,
	}
}

// Set will set the string value, while validating for correctness by calling
// the custom validation function.
func (c *CustomStringValue) Set(s string) (err error) {
	if s, err = c.handler(c.name, s, c.data); err == nil {
		c.Value = s
	}

	return
}

// GetTyped is typically used for tests and returns the flag string value
func (c *CustomStringValue) GetTyped() string {
	return c.Value
}

// Type will return a string describing the type of the flag, it is required
// to fulfill pflag.Value and will be printed in the help messages
func (c *CustomStringValue) Type() string {
	return "string"
}

// String will return the string flag value
func (c *CustomStringValue) String() string { return c.Value }

// UnmarshalText is used for TOML configuration file unmarshaling, and will
// set the value in the flag with validation.
func (c *CustomStringValue) UnmarshalText(text []byte) error {
	return c.Set(string(text))
}

// MarshalText is used for TOML configuration file marshaling, it is used when
// generating the config files via config generate.
func (c *CustomStringValue) MarshalText() (text []byte, err error) {
	return []byte("\"" + c.Value + "\""), nil
}

// ---------------------------------------------------------------------------

// IPValue is a CustomStringValue flag struct containing, and validating for,
// an IP address (v4 or v6).
type IPValue struct {
	*CustomStringValue
}

// validateIP will validate the IP address, used as a CustomStringValue
// validation function
func validateIP(name, s string, data interface{}) (string, error) {
	if net.ParseIP(s) == nil {
		return "", fmt.Errorf("Invalid value %s, should be an IP address", s)
	}
	return s, nil
}

// NewIPValue returns a new IPValue flag, its value is set to "0.0.0.0". Name
// is a name for this flag, typically the configuration variable name, and is
// used for error messages.
func NewIPValue(name string) *IPValue {
	return &IPValue{
		CustomStringValue: NewDefaultCustomStringValue(name, "0.0.0.0", validateIP, nil),
	}
}

// NewDefaultIPValue returns a new IPValue flag set to the passed set
// value. Name is a name for this flag, typically the configuration variable
// name, and is used for error messages. If the set value is not a valid IP,
// this function will panic.
func NewDefaultIPValue(name, set string) *IPValue {
	return &IPValue{
		CustomStringValue: NewDefaultCustomStringValue(name, set, validateIP, nil),
	}
}

// --------------------------------------------------------------------------

// EnumValue is a CustomStringValue flag struct containing, and validating for,
// a specified set of acceptable string values set on creation.
type EnumValue struct {
	*CustomStringValue
}

// validateEnum will validate the string, used as a CustomStringValue
// validation function
func validateEnum(name, s string, data interface{}) (string, error) {
	enum := data.([]string)

	var valid bool
	// Don't expect enums to be large, if more performance was needed just
	// sort data on creation and bisect.
	for _, v := range enum {
		if v == s {
			valid = true
			break
		}
	}

	if !valid {
		r := ""
		for i, l := range enum {
			if i < (len(enum) - 1) {
				r = r + l + ", "
			} else {
				r = r + l + "."
			}
		}

		return "", fmt.Errorf("Invalid value %s for variable %s, should be one of %s",
			s, name, r)
	}

	return s, nil
}

// NewEnumValue returns a new EnumValue flag, its value is set to the first
// value in the valid enums list. Name is a name for this flag, typically the
// configuration variable name, and is used for error messages. valid is a set
// of valid values for this flag, at least one must be present.
func NewEnumValue(name string, valid ...string) *EnumValue {
	if len(valid) == 0 {
		panic(fmt.Sprintf("Flag %s has no valid values", name))
	}
	return &EnumValue{
		CustomStringValue: NewDefaultCustomStringValue(name, valid[0], validateEnum, valid),
	}
}

// NewDefaultEnumValue returns a new EnumValue flag, its value is set to the
// passed set value. Name is a name for this flag, typically the configuration
// variable name, and is used for error messages. valid is a set of valid
// values for this flag, at least one must be present. If the set value is not
// a valid enum value this function will panic.
func NewDefaultEnumValue(name, set string, valid ...string) *EnumValue {
	if len(valid) == 0 {
		panic(fmt.Sprintf("Flag %s has no valid values", name))
	}
	return &EnumValue{
		CustomStringValue: NewDefaultCustomStringValue(name, set, validateEnum, valid),
	}
}

// --------------------------------------------------------------------------

// PortValue is a flag representing a port value, from 0 to 65535. An IntValue
// flag with these limits could have also been used, either directly or as an
// embedded struct.
type PortValue uint16

// NewPortValue returns a new PortValue flag, its value is set to 0.
func NewPortValue() PortValue {
	return PortValue(0)
}

// SetInt will set an integer and validate it for correctness
func (p *PortValue) SetInt(i int) (err error) {
	if i < 0 || i > 65535 {
		err = fmt.Errorf("Invalid port %v, must be an integer between 0 and 65535", i)
	}

	*p = PortValue(i)
	return
}

// Set will set a value parsing from a string, while validating for correctness
func (p *PortValue) Set(e string) (err error) {
	var i int
	if e == "" {
		i = 0
	} else {
		i, err = strconv.Atoi(e)
	}

	if err != nil {
		return fmt.Errorf("%s, cannot be converted to a port", e)
	}

	return p.SetInt(i)
}

// GetTyped is typically used for tests and returns the flag uint16 value
func (p *PortValue) GetTyped() uint16 {
	return uint16(*p)
}

// Type will return a string describing the type of the flag, it is required
// to fulfill pflag.Value and will be printed in the help messages
func (p *PortValue) Type() string {
	return "uint16"
}

// String will return a string representation of the flag value
func (p *PortValue) String() string { return strconv.Itoa(int(*p)) }

// UnmarshalText is used for TOML configuration file unmarshaling, and will
// set the value in the flag with validation.
func (p *PortValue) UnmarshalText(text []byte) error {
	return p.Set(string(text))
}

// MarshalText is used for TOML configuration file marshaling, it is used when
// generating the config files via config generate.
func (p *PortValue) MarshalText() (text []byte, err error) {
	return []byte(strconv.Itoa(int(*p))), nil
}
