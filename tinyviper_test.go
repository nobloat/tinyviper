package tinyviper

import (
	"errors"
	"os"
	"testing"
)

type Config struct {
	UserConfig struct {
		Email    string `env:"MY_APP_EMAIL"`
		Password string `env:"MY_APP_PASSWORD"`
	}
	Endpoint  string `env:"MY_APP_ENDPOINT"`
	AppUrl    string `env:"MY_APP_URL"`
	Undefined string `env:"MY_UNDEFINED"`
	Optional  string `env:"MY_APP_OPTIONAL,omitempty"`
}

type Config2 struct {
	UserConfig struct {
		Email    string `env:"MY_APP_EMAIL"`
		Password string `env:"MY_APP_PASSWORD"`
	}
	Endpoint  string `env:"MY_APP_ENDPOINT"`
	AppUrl    string `env:"MY_APP_URL"`
	Undefined string `env:"MY_UNDEFINED"`
	Optional  string `env:"MY_APP_OPTIONAL"`
}

type testEnvResolver struct{}

func (t testEnvResolver) Get(key string) string {
	switch key {
	case "MY_APP_EMAIL":
		return "someemail@someprovider.org"
	case "MY_APP_PASSWORD":
		return "somepassword@someprovider.org"
	default:
		return ""
	}
}

func TestConfigErrors(t *testing.T) {
	cfg := Config{
		Undefined: "foo",
	}
	err := LoadFromResolver(&cfg, testEnvResolver{})
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
	if err.Error() != "missing config variables: MY_APP_ENDPOINT,MY_APP_URL" {
		t.Error("Expected error, got wrong one: " + err.Error())
	}
}

func TestConfigNew(t *testing.T) {
	cfg := Config{
		Undefined: "foo",
	}
	err := LoadFromResolver(&cfg, NewEnvResolver(), NewEnvFileResolver(".env.sample"))
	if err != nil {
		t.Error(err)
	}
	if cfg.UserConfig.Email != "someemail@someprovider.org" {
		t.Error(errors.New("unexpected email"))
	}
	if cfg.UserConfig.Password != "password2" {
		t.Error(errors.New("unexpected password"))
	}
	if cfg.Endpoint != "some-endpoint" {
		t.Error(errors.New("unexpected endpoint"))
	}
	if cfg.AppUrl != "https://some.exmpale.org/foo?token=bar" {
		t.Error(errors.New("unexpected url"))
	}
}

func TestConfigNewWithoutFile(t *testing.T) {
	cfg := Config{
		Undefined: "foo",
	}
	err := LoadFromResolver(&cfg, EnvResolver{}, NewEnvFileResolver(".env.sample2"))
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
	if err.Error() != "missing config variables: MY_APP_EMAIL,MY_APP_PASSWORD,MY_APP_ENDPOINT,MY_APP_URL" {
		t.Error("Expected error, got wrong one: " + err.Error())
	}
}

func TestConfigOverride(t *testing.T) {
	_ = os.Setenv("MY_APP_EMAIL", "someemail2@someprovider.org")

	cfg := Config{
		Undefined: "foo",
	}
	err := LoadFromResolver(&cfg, NewEnvResolver(), NewEnvFileResolver(".env.sample"))
	if err != nil {
		t.Error(err)
	}

	if cfg.UserConfig.Email != "someemail2@someprovider.org" {
		t.Error(errors.New("unexpected email"))
	}

	if cfg.Undefined != "foo" {
		t.Error(errors.New("unexpected app url"))
	}

	if cfg.UserConfig.Password != "password2" {
		t.Error(errors.New("unexpected password"))
	}

	if cfg.Optional != "" {
		t.Error(errors.New("unexpected optional"))
	}
}

func TestConfigNoOmitMissingVariable(t *testing.T) {
	cfg := Config2{
		Undefined: "foo",
	}

	err := LoadFromResolver(&cfg, NewEnvResolver(), NewEnvFileResolver(".env.sample"))
	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	if err.Error() != "missing config variables: MY_APP_OPTIONAL" {
		t.Error("Expected error, got wrong one: " + err.Error())
	}
}

func TestConfigOmitVariableDefined(t *testing.T) {
	cfg := Config{
		Undefined: "foo",
	}

	err := LoadFromResolver(&cfg, NewEnvResolver(), NewEnvFileResolver(".env.sample3"))
	if err != nil {
		t.Error(err)
	}

	if cfg.Optional != "optional" {
		t.Error(errors.New("unexpected optional"))
	}
}
