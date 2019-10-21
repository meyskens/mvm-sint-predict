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
	serverAddress string

	rootCmd = &cobra.Command{
		Use:   "sint-client",
		Short: "sint-client is the CLI client for mvm-sint-predict",
		Long:  `sint-client is the CLI client for mvm-sint-predict. It will connect to sint-server.`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func initConfig() {
	rootCmd.PersistentFlags().StringVarP(&serverAddress, "server-address", "s", "localhost:8080", "server to connect to")

	viper.BindPFlag("server-address", rootCmd.PersistentFlags().Lookup("server-address"))
	viper.AutomaticEnv()
}

func main() {
	flag.Parse()
	err := rootCmd.Execute()
	if err != nil {
		glog.Error(err)
	}
}
