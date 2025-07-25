package configs

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

var config *Config

func Global() *Config {
	return config
}

type Config struct {
	Logger      Logger
	Environment string `env:"ENV"`
	ServiceName string `env:"SERVICE_NAME"`
	Server      Server
}

func IsDevelopment() bool {
	return config.Environment == "prod" || config.Environment == "production"
}

type Server struct {
	Port     int `env:"SERVER_PORT"`
	Timeouts Timeouts
}

type Timeouts struct {
	Read        int `env:"SERVER_TIMEOUT_READ"`
	Write       int `env:"SERVER_TIMEOUT_WRITE"`
	Shutdown    int `env:"SERVER_TIMEOUT_SHUTDOWN"`
	ReadHeaders int `env:"SERVER_TIMEOUT_READ_HEADERS"`
	Idle        int `env:"SERVER_TIMEOUT_IDLE"`
}

type Logger struct {
	Level string `env:"LOG_LEVEL"`
}

func LoadConfig() error {
	config = &Config{}
	return populateStruct(config)
}

func populateStruct(s interface{}) error {
	durationType := reflect.TypeOf(time.Duration(0))
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("expected a non-nil pointer to struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("expected a pointer to struct")
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			if err := populateStruct(field.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		envKey := fieldType.Tag.Get("env")
		if envKey == "" {
			continue
		}

		envVal, ok := os.LookupEnv(envKey)
		if !ok {
			continue
		}

		if field.Type() == durationType {
			d, err := time.ParseDuration(envVal)
			if err != nil {
				return fmt.Errorf("invalid duration for %s: %v", envKey, err)
			}
			field.Set(reflect.ValueOf(d))
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envVal)
		case reflect.Int:
			val, err := strconv.Atoi(envVal)
			if err != nil {
				return fmt.Errorf("invalid int for %s: %v", envKey, err)
			}
			field.SetInt(int64(val))
		case reflect.Bool:
			val, err := strconv.ParseBool(envVal)
			if err != nil {
				return fmt.Errorf("invalid bool for %s: %v", envKey, err)
			}
			field.SetBool(val)
		default:
			return fmt.Errorf("unsupported kind %s for field %s", field.Kind(), fieldType.Name)
		}
	}
	return nil
}
