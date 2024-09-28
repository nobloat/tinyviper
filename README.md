# tiny-viper WIP

![ci workflow](https://github.com/nobloat/tiny-viper/actions/workflows/ci.yml/badge.svg)

A minimalistic approach to [spf13/viper](https://github.com/spf13/viper).

## Features
- Read `ENV` variables into a `struct`
- Read a `.env` file into a `struct`
- `< 110` source lines of code
- [No dependencies](go.mod)

Only string fields are supported. 

## Usage

```go
package main

import (
	"github.com/nobloat/tinyviper"
	"fmt"
)

type Config struct {
	UserConfig struct {
		Email    string `env:"MY_APP_EMAIL"`
		Password string `env:"MY_APP_PASSWORD"`
        someOtherProperty string
	}
	Endpoint string `env:"MY_APP_ENDPOINT"`
}

func main() {
  cfg := Config{}
  cfg, err := tinyviper.LoadFromResolver(&cfg, tinyviper.EnvResolver{}, tinyviper.NewEnvFileResolver(".env.sample"))
  if err != nil {
    panic(err)
  }

  fmt.Println("%+v", cfg)
}
```

