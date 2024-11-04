package main

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/SonicCloudOrg/sonic-android-supply/pkg/utils"
	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/spf13/cobra"
)

var translatePort int

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().IntVarP(&translatePort, "translate-port", "p", 0, "translating proxy port")
	shareCmd.Flags().StringVarP(&serial, "serial", "s", "", "device serial")
}

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share the connected adb device in the network",
	Long:  "Share the connected adb device in the network",
	RunE: func(cmd *cobra.Command, args []string) error {
		if translatePort == 0 {
			port, err := utils.FindAvailablePort(6174)
			if err != nil {
				return err
			}
			translatePort = port
		}
		if serial == "" {
			output, err := exec.Command("adb", "devices", "-l").CombinedOutput()
			if err != nil {
				return err
			}
			re := regexp.MustCompile(`(?m)^([^\s]+)\s+device\s+(.+)$`)
			matches := re.FindAllStringSubmatch(string(output), -1)
			for _, m := range matches {
				serial = m[1]
				break
			}
			return fmt.Errorf("no devices connected")
		}
		adbd := adb.NewADBDaemon2(serial)
		logger.Info(fmt.Sprintf("Connect with: adb connect %s:%d\n", utils.GetHostIP(), translatePort))
		err := adbd.ListenAndServe(fmt.Sprintf(":%d", translatePort))
		if err != nil {
			return err
		}
		return nil
	},
}
