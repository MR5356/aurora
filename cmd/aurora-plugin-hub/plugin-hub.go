package main

import (
	"github.com/MR5356/aurora/pkg/domain/runner/pluginhub"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	port  int
	debug bool
)

func NewPluginHubCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin-hub",
		Short: "plugin-hub",
		RunE: func(cmd *cobra.Command, args []string) error {
			ph := pluginhub.New()
			return ph.Run(port)
		},
	}

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	cmd.PersistentFlags().IntVar(&port, "port", 8081, "server port")

	return cmd
}

func main() {
	if err := NewPluginHubCommand().Execute(); err != nil {
		logrus.Fatal(err)
	}
}
