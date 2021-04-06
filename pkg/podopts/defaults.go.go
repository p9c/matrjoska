package podopts

import (
	"reflect"
)

// GetDefaultConfig returns a Config struct pristine factory fresh
func GetDefaultConfig() (c *Config) {
	I.Ln("getting default config")
	c = &Config{
		Commands: GetCommands(),
		Map:      GetConfigs(),
	}
	c.RunningCommand = c.Commands[0]
	// I.S(c.Commands[0])
	// I.S(c.Map)
	t := reflect.ValueOf(c)
	t = t.Elem()
	for i := range c.Map {
		tf := t.FieldByName(i)
		if tf.IsValid() && tf.CanSet() && tf.CanAddr() {
			val := reflect.ValueOf(c.Map[i])
			tf.Set(val)
		}
	}
	// I.S(c)
	return
}

