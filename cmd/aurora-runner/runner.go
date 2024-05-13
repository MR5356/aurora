package main

import (
	"github.com/MR5356/aurora/pkg/domain/runner"
	"github.com/MR5356/aurora/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	port        int
	host, token string
	debug       bool
)

func NewRunnerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "runner",
		Short:   "runner",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Infof("token: %s", token)
			return runner.Run(&runner.Config{
				Host:  host,
				Port:  port,
				Token: token,
				Debug: debug,
			})
		},
	}

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	cmd.PersistentFlags().IntVar(&port, "port", 8080, "server port")
	cmd.PersistentFlags().StringVar(&host, "host", "localhost", "token")
	cmd.PersistentFlags().StringVarP(&token, "token", "t", "", "token")

	return cmd
}

func main() {
	if err := NewRunnerCommand().Execute(); err != nil {
		logrus.Fatal(err)
	}

}
