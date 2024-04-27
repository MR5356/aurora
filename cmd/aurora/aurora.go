package main

import (
	"github.com/MR5356/aurora/pkg/config"
	"github.com/MR5356/aurora/pkg/server"
	"github.com/MR5356/aurora/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewAuroraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "aurora",
		Short:   "aurora",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.New(
				config.WithDebug(true),
				config.WithDatabase("sqlite", "db.sqlite"),
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

	return cmd
}

func main() {
	if err := NewAuroraCommand().Execute(); err != nil {
		logrus.Fatal(err)
	}
}
