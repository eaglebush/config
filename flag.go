package cfg

import (
	"strconv"
	"strings"
)

// Flag - dynamic flags structure
type Flag struct {
	Key   string
	Value *string
}

// Bool - return a boolean from flag value
func (f Flag) Bool() *bool {

	var ret *bool

	if f.Value == nil {
		return ret
	}

	v := strings.TrimSpace(*f.Value)
	v = strings.ToLower(v)
	ret = new(bool)

	switch v {
	case "1", "on", "yes", "enabled", "true":
		*ret = true
		return ret
	}

	return ret
}

// Int64 - return an int64 from flag value
func (f Flag) Int64() *int64 {
	if f.Value == nil {
		return nil
	}

	v := strings.TrimSpace(*f.Value)
	vi, _ := strconv.ParseInt(v, 0, 64)
	return &vi
}

// Int - return an int from flag value
func (f Flag) Int() *int {
	if f.Value == nil {
		return nil
	}

	v := strings.TrimSpace(*f.Value)
	vi, _ := strconv.Atoi(v)
	return &vi
}

// Float - return an float from flag value
func (f Flag) Float() *float32 {
	if f.Value == nil {
		return nil
	}

	v := strings.TrimSpace(*f.Value)
	vi, _ := strconv.ParseFloat(v, 32)

	ret := new(float32)
	*ret = float32(vi)
	return ret
}

// Float - return an float from flag value
func (f Flag) Float64() *float64 {
	if f.Value == nil {
		return nil
	}

	v := strings.TrimSpace(*f.Value)
	vi, _ := strconv.ParseFloat(v, 32)

	ret := new(float64)
	*ret = float64(vi)
	return ret
}

// String - return a string from flag value
func (f Flag) String() *string {
	if f.Value == nil {
		return nil
	}

	return f.Value
}
