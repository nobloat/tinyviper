package tinyviper

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strings"
)

type Resolver interface {
	Get(key string) string
}

type EnvResolver struct{}

type EnvFileResolver struct {
	Variables map[string]string
}

func NewEnvFileResolver(filename string) (*EnvFileResolver, error) {
	readFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	r := &EnvFileResolver{make(map[string]string, 0)}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		text := fileScanner.Text()
		parts := strings.Split(text, "=")
		if len(parts) != 2 || strings.HasPrefix(text, "#") {
			continue
		}
		r.Variables[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return r, readFile.Close()
}

func (e EnvFileResolver) Get(key string) string {
	return e.Variables[key]
}

func (e EnvResolver) Get(key string) string {
	return os.Getenv(key)
}

func NewEnvConfig[T any]() (*T, error) {
	cfg := new(T)
	res := EnvResolver{}
	err := ReflectStruct(cfg, res)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func NewEnvFileConfig[T any](filename string) (*T, error) {
	cfg := new(T)
	res, err := NewEnvFileResolver(filename)
	if err != nil {
		return nil, err
	}
	err = ReflectStruct(cfg, res)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func ReflectStruct(object any, resolver Resolver) error {
	v := reflect.ValueOf(object)
	if v.Elem().Kind() != reflect.Struct {
		return errors.New("type must be a struct")
	}
	e := v.Elem()
	t := e.Type()
	for i := 0; i < t.NumField(); i++ {
		ef := e.Field(i)
		tf := t.Field(i)
		envName := tf.Tag.Get("env")
		if envName != "" {
			if tf.Type != reflect.TypeOf("") {
				return errors.New("env annotated field must have type string")
			}
			if !ef.CanSet() {
				return errors.New("env field must be public")
			}
			value := resolver.Get(envName)
			if value == "" {
				return errors.New("env variable " + envName + " is not set")
			}
			ef.SetString(resolver.Get(envName))
		} else if t.Kind() == reflect.Struct {
			err := ReflectStruct(ef.Addr().Interface(), resolver)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
