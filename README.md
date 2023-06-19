# tiny-viper WIP
A minimalistic approach to [spf13/viper](https://github.com/spf13/viper).

## Features
- Read `ENV` variables into a `struct`
- Read a `.env` file into a `struct`

Only string fields are supported. 

## Usage

```go
type Config struct {
	UserConfig struct {
		Email    string `env:"MY_APP_EMAIL"`
		Password string `env:"MY_APP_PASSWORD"`
    someOtherProperty string
	}
	Endpoint string `env:"MY_APP_ENDPOINT"`
}

func main() {
  //cfg, err := NewEnvConfig[Config]()  //Read from env
  cfg, err := NewEnvFileConfig[Config](".env.sample") //Read from .env file
  if err != nil {
    panic(err)
  }

  fmt.Println("%+v", cfg)
}
```

