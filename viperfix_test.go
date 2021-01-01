package viperfix_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/momper14/viperfix"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	prefix = "TEST"
)

var yamlExample = []byte(`
log:
  filename: logs/latest.log
  compress: true
  level: info
  timestampformat: 02-01-2006 15:04:05
  max:
    size: 50
    backups: 5
`)

type config struct {
	Filename        string
	Compress        bool
	Level           string
	Timestampformat string
	Max             struct {
		Size    int
		Backups int
		Age     int
	}
}

// init viper
func initYaml() {
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(yamlExample))
}

// define config
func initDefaults() {
	viper.SetDefault("log.filename", "logs/latest.log")
	viper.SetDefault("log.compress", true)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.max.backups", 5)
	viper.SetDefault("log.max.age", 31)
}

func initEnv() {
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.BindEnv("log.level")

	os.Setenv(prefix+"_LOG_LEVEL", "debug")
}

func initCombined() {
	initDefaults()
	initYaml()
	initEnv()

	viper.AutomaticEnv()
	viper.SetTypeByDefaultValue(true)
	os.Setenv(prefix+"_LOG_MAX_BACKUPS", "3")
}

// define config
func initSingleVal() {
	viper.SetDefault("log", "fail")
}

func TestGetStringMapDefault(t *testing.T) {
	viper.Reset()
	initDefaults()
	var expected = map[string]interface{}{
		"filename": "logs/latest.log",
		"compress": true,
		"level":    "info",
		"max": map[string]interface{}{
			"backups": 5,
			"age":     31,
		},
	}

	actual := viperfix.GetStringMap("log")

	assert.Equal(t, expected, actual)
}

func TestGetStringMapConf(t *testing.T) {
	viper.Reset()
	initYaml()
	var expected = map[string]interface{}{
		"filename":        "logs/latest.log",
		"compress":        true,
		"level":           "info",
		"timestampformat": "02-01-2006 15:04:05",
		"max": map[string]interface{}{
			"backups": 5,
			"size":    50,
		},
	}

	actual := viperfix.GetStringMap("log")

	assert.Equal(t, expected, actual)
}

func TestGetStringMapEnv(t *testing.T) {
	viper.Reset()
	initEnv()
	var expected = map[string]interface{}{
		"level": "debug",
	}

	actual := viperfix.GetStringMap("log")

	assert.Equal(t, expected, actual)
}

func TestGetStringMapCombined(t *testing.T) {
	viper.Reset()
	initCombined()
	var expected = map[string]interface{}{
		"filename":        "logs/latest.log",
		"compress":        true,
		"level":           "debug",
		"timestampformat": "02-01-2006 15:04:05",
		"max": map[string]interface{}{
			"backups": 3,
			"age":     31,
			"size":    50,
		},
	}

	actual := viperfix.GetStringMap("log")

	assert.Equal(t, expected, actual)
}

func TestGetStringMapEmpty(t *testing.T) {
	viper.Reset()

	actual := viperfix.GetStringMap("log")

	assert.Nil(t, actual)
}

func TestGetStringMapSingleVal(t *testing.T) {
	viper.Reset()
	initSingleVal()

	actual := viperfix.GetStringMap("log")

	assert.Nil(t, actual)
}

func TestSubDefault(t *testing.T) {
	viper.Reset()
	initDefaults()
	var expected = map[string]interface{}{
		"filename": "logs/latest.log",
		"compress": true,
		"level":    "info",
		"max": map[string]interface{}{
			"backups": 5,
			"age":     31,
		},
	}

	actual := viperfix.Sub("log")

	assert.Equal(t, expected, actual.AllSettings())
}

func TestSubConf(t *testing.T) {
	viper.Reset()
	initYaml()
	var expected = map[string]interface{}{
		"filename":        "logs/latest.log",
		"compress":        true,
		"level":           "info",
		"timestampformat": "02-01-2006 15:04:05",
		"max": map[string]interface{}{
			"backups": 5,
			"size":    50,
		},
	}

	actual := viperfix.Sub("log")

	assert.Equal(t, expected, actual.AllSettings())
}

func TestSubEnv(t *testing.T) {
	viper.Reset()
	initEnv()
	var expected = map[string]interface{}{
		"level": "debug",
	}

	actual := viperfix.Sub("log")

	assert.Equal(t, expected, actual.AllSettings())
}

func TestSubCombined(t *testing.T) {
	viper.Reset()
	initCombined()
	var expected = map[string]interface{}{
		"filename":        "logs/latest.log",
		"compress":        true,
		"level":           "debug",
		"timestampformat": "02-01-2006 15:04:05",
		"max": map[string]interface{}{
			"backups": 3,
			"age":     31,
			"size":    50,
		},
	}

	actual := viperfix.Sub("log")

	assert.Equal(t, expected, actual.AllSettings())
}

func TestSubEmpty(t *testing.T) {
	viper.Reset()

	actual := viperfix.Sub("log")

	assert.Nil(t, actual)
}

func TestSubSingleVal(t *testing.T) {
	viper.Reset()
	initSingleVal()

	actual := viperfix.Sub("log")

	assert.Nil(t, actual)
}

func TestUnmarshalKeyDefault(t *testing.T) {
	viper.Reset()
	initDefaults()

	var expected = config{
		Filename: "logs/latest.log",
		Compress: true,
		Level:    "info",
		Max: struct {
			Size    int
			Backups int
			Age     int
		}{
			Backups: 5,
			Age:     31,
		},
	}

	var actual config
	err := viperfix.UnmarshalKey("log", &actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnmarshalKeyConf(t *testing.T) {
	viper.Reset()
	initYaml()

	var expected = config{
		Filename:        "logs/latest.log",
		Compress:        true,
		Level:           "info",
		Timestampformat: "02-01-2006 15:04:05",
		Max: struct {
			Size    int
			Backups int
			Age     int
		}{
			Size:    50,
			Backups: 5,
		},
	}

	var actual config
	err := viperfix.UnmarshalKey("log", &actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnmarshalKeyEnv(t *testing.T) {
	viper.Reset()
	initEnv()

	var expected = config{
		Level: "debug",
	}

	var actual config
	err := viperfix.UnmarshalKey("log", &actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnmarshalKeyCombined(t *testing.T) {
	viper.Reset()
	initCombined()

	var expected = config{
		Filename:        "logs/latest.log",
		Compress:        true,
		Level:           "debug",
		Timestampformat: "02-01-2006 15:04:05",
		Max: struct {
			Size    int
			Backups int
			Age     int
		}{
			Size:    50,
			Age:     31,
			Backups: 3,
		},
	}

	var actual config
	err := viperfix.UnmarshalKey("log", &actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnmarshalKeyEmpty(t *testing.T) {
	viper.Reset()

	var actual config
	err := viperfix.UnmarshalKey("log", &actual)

	assert.Nil(t, err)
	assert.Empty(t, actual)
}

func TestUnmarshalKeySingleVal(t *testing.T) {
	viper.Reset()
	initSingleVal()

	var actual config
	err := viperfix.UnmarshalKey("log", &actual)

	assert.Nil(t, err)
	assert.Empty(t, actual)
}

func TestKeyDelimiter(t *testing.T) {
	viperfix.KeyDelimiter("-")
}
