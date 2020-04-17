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
	k := strings.ToLower(key)

	for i := range c.Flags {
		if k2 := strings.TrimSpace(strings.ToLower(c.Flags[i].Key)); k == k2 {
			return c.Flags[i]
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

// String - return a string from flag value
func (f Flag) String() string {
	return strings.TrimSpace(f.Value)
}
