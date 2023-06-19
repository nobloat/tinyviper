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
}

func TestConfig(t *testing.T) {
	cfg, err := NewEnvFileConfig[Config](".env.sample")
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
}
