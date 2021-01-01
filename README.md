# Viperfix

- [Viperfix](#viperfix)
  - [What is Viperfix?](#what-is-viperfix)
  - [Why Viperfix?](#why-viperfix)
    - [Example](#example)
  - [How to use](#how-to-use)
    - [Example](#example-1)

## What is Viperfix?

Viperfix is a small tool with an aproach to fix the problems of [Viper](https://github.com/spf13/viper) getting multiple values at once with (env)overrides and defaults.

## Why Viperfix?

I like the idea of Viper. But as I started to get blocks of config into my structs, i noticed a big problem for me. It only uses my config-file and ignores defaults and overrides from env, which I need to work with docker.

### Example

defaults

```go
viper.SetDefault("log.filename", "logs/latest.log")
viper.SetDefault("log.compress", true)
viper.SetDefault("log.level", "info")
viper.SetDefault("log.max.size", 50)
viper.SetDefault("log.max.backups", 5)
viper.SetDefault("log.max.age", 31)
```

config.yml

```yaml
log:
  filename: logs/latest.log
  compress: true
  level: info
  max:
    size: 50
    backups: 5
```

expected from viper.GetStringMap("log")

```go
map[compress:true filename:logs/latest.log level:info max:map[age:31 backups:5 size:50]]
```

but its actual

```go
map[compress:true filename:logs/latest.log level:info max:map[backups:5 size:50]]
```

There are a lot of issues at [Viper](https://github.com/spf13/viper/issues?q=is%3Aissue+is%3Aopen+sub) about Sub(), GetStringMap and Unmarshal() but I didn't find a real solution.

## How to use

First, there is a function to set the key delimter (didn't found a way to get it from viper).

```go
func KeyDelimiter(delimiter string)
```

And for now, there are 3 Functions.

```go
func GetStringMap(key string) map[string]interface{}
func Sub(key string) *viper.Viper
func UnmarshalKey(key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error
```

They are like the viper functions with equal name but use all sources hirachical to get the values.

They use the default viper singleton instance but there are also versions of each function which takes a viper instance to work with.

```go
func GetStringMap(key string) map[string]interface{}

func GetStringMapFrom(v *viper.Viper, key string) map[string]interface{}
```

### Example

for example

```go
viper.SetDefault("log.max.size", 50)
viper.SetDefault("log.max.backups", 5)
viper.SetDefault("log.max.age", 31)

viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
viper.AutomaticEnv()
viper.SetTypeByDefaultValue(true)

os.Setenv("LOG.MAX.AGE", "15")

fmt.Println(viperfix.GetStringMap("log.max"))
```

prints

```go
map[age:31 backups:5 size:50]
```
