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

func (e EnvFileResolver) Get(key string) string {
	return e.Variables[key]
}
func (e EnvResolver) Get(key string) string {
	return os.Getenv(key)
}

type multiResolver struct {
	resolvers []Resolver
}

func NewEnvFileResolver(filename string) *EnvFileResolver {
	readFile, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer readFile.Close()
	r := &EnvFileResolver{make(map[string]string, 0)}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		text := fileScanner.Text()
		index := strings.Index(text, "=")
		if index == -1 || strings.HasPrefix(text, "#") {
			continue
		}
		r.Variables[text[0:index]] = strings.Trim(text[index+1:], " \"'")
	}
	return r
}

func (m multiResolver) Get(key string) string {
	for _, r := range m.resolvers {
		v := r.Get(key)
		if r.Get(key) != "" {
			return v
		}
	}
	return ""
}

func LoadFromResolver[T any](cfg *T, resolver ...Resolver) error {
	res := multiResolver{resolvers: resolver}
	missing := make([]string, 0)
	missing, err := refelectStruct(cfg, res, missing)
	if err != nil {
		return err
	}
	if len(missing) > 0 {
		return errors.New("missing config variables: " + strings.Join(missing, ","))
	}
	return nil
}

func refelectStruct(object any, resolver Resolver, missing []string) ([]string, error) {
	v := reflect.ValueOf(object)
	if v.Elem().Kind() != reflect.Struct {
		return missing, errors.New("type must be a struct")
	}
	e := v.Elem()
	t := e.Type()
	for i := 0; i < t.NumField(); i++ {
		ef := e.Field(i)
		tf := t.Field(i)
		envName := tf.Tag.Get("env")
		if envName != "" {
			if tf.Type != reflect.TypeOf("") {
				return missing, errors.New("env annotated field must have type string")
			}
			if !ef.CanSet() {
				return missing, errors.New("env field must be public")
			}
			value := resolver.Get(envName)
			if value != "" {
				ef.SetString(resolver.Get(envName))
			} else if ef.String() == "" {
				missing = append(missing, envName)
			}
		} else if t.Kind() == reflect.Struct {
			m, err := refelectStruct(ef.Addr().Interface(), resolver, missing)
			if err != nil {
				return missing, err
			}
			missing = append(missing, m...)
		}
	}
	return missing, nil
}
