/*
 *   sonic-android-supply  Supply of ADB.
 *   Copyright (C) 2022  SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published
 *   by the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"github.com/SonicCloudOrg/sonic-android-supply/pkg/utils"
	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share the connected adb device in the network",
	Long:  "Share the connected adb device in the network",
	Run: func(cmd *cobra.Command, args []string) {
		if translatePort == 0 {
			port, err := utils.FindAvailablePort(6174)
			if err != nil {
				log.Panic(err)
			}
			translatePort = port
		}
		if serial == "" {
			output, err := exec.Command("adb", "devices", "-l").CombinedOutput()
			if err != nil {
				log.Panic(err)
			}
			re := regexp.MustCompile(`(?m)^([^\s]+)\s+device\s+(.+)$`)
			matches := re.FindAllStringSubmatch(string(output), -1)
			for _, m := range matches {
				serial = m[1]
				break
			}
			log.Panic("no devices connected")
		}
		adbd := adb.NewADBDaemon2(serial)
		fmt.Printf("Connect with: adb connect %s:%d\n", utils.GetHostIP(), translatePort)
		err := adbd.ListenAndServe(fmt.Sprintf(":%d", translatePort))
		if err != nil {
			log.Panic(err)
		}
	},
}

var translatePort int

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().IntVarP(&translatePort, "translate-port", "p", 0, "translating proxy port")
	shareCmd.Flags().StringVarP(&serial, "serial", "s", "", "device serial")
}
