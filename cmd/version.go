package main

import (
	"fmt"

	"github.com/SonicCloudOrg/sonic-android-supply/pkg/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version code of sas",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(version.VERSION)
		return nil
	},
}
