---
title:      "Viper Multi Config"
date:       2021-08-18
categories: ["Snippet"]
tags:       ["golang", "dev", "cli"]
url:        "post/viper_multi_config"
---

This post will explain how to overload configurations in Golang using
[viper][viper_repo] and [cobra][cobra_repo] libraries.
The use case may be a niche one, but I find this to be an easy to understand
and pretty clear way to do configuration merge.

## The problem and the expectation

I am deploying most of my application on Kubernetes nowadays (I also moved
my personal infrastructure to it a few weeks ago). This solution come with
a big tooling ecosystem and a really opiniated way to deploy.

Most of the time applications you deploy come with a default configuration
stored in a config map (using Kustomize, Helm, ....).
My problem was that I wanted to package one of my application with a default
configuration, then overload it only for specific entries. Here is an example
to make this a bit more clear.

```toml
[log]
level = "info"
type = "json"

[sql]
host = "localhost"
port = 3306
```

My application come packaged with a configuration file stored in a config map.
Let's imagine I want to only set a different log level for one of my deployment.
I would like to avoid doing this through the command line since some entries may
be tricky to map to a proper option.

The idea would be to keep this config map, and create a second one with only
the entries I want to modify. For example:

```toml
[log]
level = "debug"

[sql]
host = "my.remote.endpoint"
```

In memory my application will merge those two files to create a final configuration
looking like this:

```toml
[log]
level = "debug"
type = "json"

[sql]
host = "my.remote.endpoint"
port = 3306
```

My application will then be called this way:
`myapp --config=config_1.toml --config=config_2.toml`.
It will load config files in order overriding the config entries as it goes
and keeping the base one if they are not modified.

## The implementation

```go
package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Configuration struct {
	Logger        struct {
		Level string `mapstructure:"level"`
		Type  string `mapstructure:"type"`
	} `mapstructure:"log"`
	SQL struct {
		Host string `mapstructure:"database"`
		Port int    `mapstructure:"table"`
	} `mapstructure:"sql"`
}

func init() {
	flags := command.PersistentFlags()
	flags.StringArray("config", []string{}, "Path to configuration files")
}

func config(cmd *cobra.Command) (*Configuration, error) {
	configuration := &Configuration{}
	configs, err := cmd.PersistentFlags().GetStringArray("config")
	if err != nil {
		return nil, err
	}
	for _, config := range configs {
		viper.SetConfigFile(config)
		if err := viper.MergeInConfig(); err != nil {
			return nil, err
		}
	}

	return configuration, viper.Unmarshal(configuration)
}

var (
	command = &cobra.Command{
		Use:   "myapp [flags]",
		RunE: func(cmd *cobra.Command, args []string) error {
			// https://github.com/spf13/cobra/issues/340
			cmd.SilenceUsage = true
			config, err := config(cmd)
			if err != nil {
				return err
			}
			// ...
		},
	}
)
```

Here the magic is happening in the `config` function. We are using the global `viper`
object and we push the entries of our option. Then we call the `MergeInConfig` function
and it will perform what we were looking for.

This is pretty straightforward but a bit hidden in the doc so not really obvious.
Hope it helps someone out there, I find this feature pretty convenient and I
will most probably integrate it in my upcoming developments.


[viper_repo]: https://github.com/spf13/viper
[cobra_repo]: https://github.com/spf13/cobra
