package main

import (
	"os"

	"github.com/adamgoose/ghpatch/cmd"
	"github.com/adamgoose/ghpatch/lib"
	"github.com/defval/di"
	"github.com/google/go-github/v56/github"
	"github.com/spf13/viper"
)

func main() {
	if err := lib.App.Apply(
		di.Provide(func() (*github.Client, error) {
			c := github.NewClient(nil).
				WithAuthToken(viper.GetString("github_token"))

			return c, nil
		}),
	); err != nil {
		panic(err)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.SetEnvPrefix("ghpatch")
	viper.AutomaticEnv()
}
