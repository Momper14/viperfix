// Package viperfix provides functions to fix the missing of hirachical getting data when getting multiple values
package viperfix

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var deli = "."

// KeyDelimiter sets the Delimiter for keys
func KeyDelimiter(delimiter string) {
	deli = delimiter
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} {
	return GetStringMapFrom(viper.GetViper(), key)
}

// GetStringMapFrom uses given viperinstance and returns the
// value associated with the key as a map of interfaces.
func GetStringMapFrom(v *viper.Viper, key string) map[string]interface{} {
	var (
		m    = make(map[string]interface{})
		pref = key + deli
		trim = len(pref)
	)

	for _, k := range v.AllKeys() {
		if strings.HasPrefix(k, pref) {
			newKey := k[trim:]

			if !strings.Contains(newKey, deli) {
				m[newKey] = v.Get(k)
				continue
			}

			path := strings.Split(newKey, deli)
			last := len(path) - 1
			tmp := m

			for _, sk := range path[:last] {
				if _, ok := tmp[sk]; !ok {
					tmp[sk] = make(map[string]interface{})
				}

				tmp = tmp[sk].(map[string]interface{})
			}

			tmp[path[last]] = v.Get(k)
		}
	}

	if len(m) == 0 {
		return nil
	}

	return m
}

// UnmarshalKey takes a single key to unmarshals its values into a Struct.
func UnmarshalKey(key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return UnmarshalKeyFrom(viper.GetViper(), key, rawVal, opts...)
}

// UnmarshalKeyFrom uses given viperinstance and takes a single key to unmarshals its values into a Struct.
func UnmarshalKeyFrom(v *viper.Viper, key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return decode(GetStringMapFrom(v, key), defaultDecoderConfig(rawVal, opts...))
}

// Sub returns new Viper instance representing a sub tree of this instance.
// Sub is case-insensitive for a key.
func Sub(key string) *viper.Viper {
	return SubFrom(viper.GetViper(), key)
}

// SubFrom uses given viperinstance and returns new Viper instance representing a sub tree of this instance.
// Sub is case-insensitive for a key.
func SubFrom(v *viper.Viper, key string) *viper.Viper {
	subv := viper.New()

	data := GetStringMapFrom(v, key)
	if data == nil {
		return nil
	}

	subv.MergeConfigMap(data)
	return subv
}

// A wrapper around mapstructure.Decode that mimics the WeakDecode functionality
func decode(input interface{}, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func defaultDecoderConfig(output interface{}, opts ...viper.DecoderConfigOption) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
