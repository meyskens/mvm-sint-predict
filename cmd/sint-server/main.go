package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "sint-server",
		Short: "TODO",
		Long:  `TODO`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func initConfig() {
	viper.AutomaticEnv()
}

func main() {
	flag.Parse()
	err := rootCmd.Execute()
	if err != nil {
		glog.Error(err)
	}
}
