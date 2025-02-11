package main

import "github.com/Sanchir01/avito-testovoe/internal/app"

func main() {
	env, err := app.NewEnv()
	if err != nil {
		panic(err)
	}
	env.Lg.Info("hello")
}
