package main

import (
	"github.com/MR5356/aurora/pkg/config"
	"github.com/MR5356/aurora/pkg/server"
	"github.com/MR5356/aurora/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	port            int
	debug           bool
	dbDriver, dbDSN string
)

func NewAuroraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "aurora",
		Short:   "aurora",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.New(
				config.WithDebug(debug),
				config.WithPort(port),
				config.WithDatabase(dbDriver, dbDSN),
			)

			svc, err := server.New(cfg)
			if err != nil {
				logrus.Fatalf("server.New failed, err: %v", err)
			}

			if err := svc.Run(); err != nil {
				logrus.Fatalf("server.Run failed, err: %v", err)
			}
			return nil
		},
	}

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 80, "server port")
	cmd.PersistentFlags().StringVar(&dbDriver, "dbDriver", "sqlite", "database driver")
	cmd.PersistentFlags().StringVar(&dbDSN, "dbDSN", "db.sqlite", "database DSN")

	return cmd
}

func main() {
	if err := NewAuroraCommand().Execute(); err != nil {
		logrus.Fatal(err)
	}
}
