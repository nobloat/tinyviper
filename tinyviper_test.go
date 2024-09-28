package tinyviper

import (
	"errors"
	"testing"
)

type Config struct {
	UserConfig struct {
		Email    string `env:"MY_APP_EMAIL"`
		Password string `env:"MY_APP_PASSWORD"`
	}
	Endpoint string `env:"MY_APP_ENDPOINT"`
	AppUrl   string `env:"MY_APP_URL"`
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
	cfg := Config{}
	err := LoadFromResolver(&cfg, testEnvResolver{})
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
	if err.Error() != "missing config variables: MY_APP_ENDPOINT,MY_APP_URL" {
		t.Error("Expected error, got wrong one: " + err.Error())
	}
}

func TestConfigNew(t *testing.T) {
	cfg := Config{}
	err := LoadFromResolver(&cfg, EnvResolver{}, NewEnvFileResolver(".env.sample"))
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
