package cfg

import (
	"strconv"
	"strings"
)

// Flag - dynamic flags structure
type Flag struct {
	Key   string
	Value string
}

// Flag - get a flag value
func (c *Configuration) Flag(key string) Flag {

	if c.Flags == nil {
		return Flag{}
	}

	flgs := *c.Flags

	for i := range flgs {
		if strings.EqualFold(key, flgs[i].Key) {
			return flgs[i]
		}
	}

	return Flag{}
}

// Bool - return a boolean from flag value
func (f Flag) Bool() bool {
	v := strings.TrimSpace(f.Value)
	v = strings.ToLower(v)
	switch v {
	case "1", "on", "yes", "enabled", "true":
		return true
	}
	return false
}

// Int64 - return an int64 from flag value
func (f Flag) Int64() int64 {
	v := strings.TrimSpace(f.Value)
	vi, _ := strconv.ParseInt(v, 0, 64)
	return vi
}

// Int - return an int from flag value
func (f Flag) Int() int {
	v := strings.TrimSpace(f.Value)
	vi, _ := strconv.Atoi(v)
	return vi
}

// Float - return an float from flag value
func (f Flag) Float() float32 {
	v := strings.TrimSpace(f.Value)
	vi, _ := strconv.ParseFloat(v, 32)
	return float32(vi)
}

// Float - return an float from flag value
func (f Flag) Float64() float64 {
	v := strings.TrimSpace(f.Value)
	vi, _ := strconv.ParseFloat(v, 64)
	return vi
}

// String - return a string from flag value
func (f Flag) String() string {
	return f.Value
}
