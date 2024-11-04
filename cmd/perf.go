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
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"sync"
	"time"

	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/SonicCloudOrg/sonic-android-supply/src/entity"
	"github.com/SonicCloudOrg/sonic-android-supply/src/perfmonUtil"
	"github.com/spf13/cobra"
)

var (
	serial2      string
	isFormat2    bool
	isJson2      bool
	perfOptions2 entity.PerfOption
	pid2         int
	packageName2 string
	refreshTime2 int
	isForce2     bool
)

func init() {
	perfCmd.Flags().StringVarP(&serial2, "serial", "s", "", "device serial (default first device)")
	perfCmd.Flags().IntVarP(&pid2, "pid", "d", -1, "get PID data")
	perfCmd.Flags().StringVarP(&packageName2, "package", "p", "", "app package name")
	perfCmd.Flags().BoolVar(&perfOptions2.SystemCPU, "sys-cpu", false, "get system cpu data")
	perfCmd.Flags().BoolVar(&perfOptions2.SystemMem, "sys-mem", false, "get system memory data")
	//perfmonCmd.Flags().BoolVar(&sysDisk, "sys-disk", false, "get system disk data")
	perfCmd.Flags().BoolVar(&perfOptions2.SystemNetWorking, "sys-network", false, "get system networking data")
	//perfmonCmd.Flags().BoolVar(&perfOptions2.SystemGPU, "gpu", false, "get gpu data")
	perfCmd.Flags().BoolVar(&perfOptions2.ProcFPS, "proc-fps", false, "get fps data")
	perfCmd.Flags().BoolVar(&perfOptions2.ProcThreads, "proc-threads", false, "get process threads")
	//perfmonCmd.Flags().BoolVar(&, "proc-network", false, "get process network data")
	perfCmd.Flags().BoolVar(&perfOptions2.ProcCPU, "proc-cpu", false, "get process cpu data")
	perfCmd.Flags().BoolVar(&perfOptions2.ProcMem, "proc-mem", false, "get process mem data")
	perfCmd.Flags().BoolVar(&isForce2, "force-out", false, "force update pid perf data(applicable to applications being restarted by kill)")
	perfCmd.Flags().IntVarP(&refreshTime2, "refresh", "r", 1000, "data refresh time (millisecond)")
	perfCmd.Flags().BoolVarP(&isFormat2, "format", "f", false, "convert to JSON string and format")
	perfCmd.Flags().BoolVarP(&isJson2, "json", "j", false, "convert to JSON string")
	rootCmd.AddCommand(perfCmd)
}

var perfCmd = &cobra.Command{
	Use:   "perf",
	Short: "Get device performance",
	Long:  "Get device performance",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if serial2 == "" {
			output, err := exec.Command("adb", "devices", "-l").CombinedOutput()
			if err != nil {
				log.Panic(err)
			}
			re := regexp.MustCompile(`(?m)^([^\s]+)\s+device\s+(.+)$`)
			matches := re.FindAllStringSubmatch(string(output), -1)
			if len(matches) == 0 {
				log.Panic("no devices connected")
			}
			for _, m := range matches {
				serial2 = m[1]
				break
			}
		}

		device := adb.NewClient("").DeviceWithSerial2(serial2)
		// data := perfmonUtil.GetSystemMem2(device)
		// fmt.Println(data.ToFormat())
		pidStr := ""
		if isForce {
			perfmonUtil.IsForce = true
			if packageName == "" {
				fmt.Println("please enter packageName.")
				os.Exit(0)
			}
		}
		if pid != -1 && packageName == "" {
			pidStr = fmt.Sprintf("%d", pid)
			packageName, err = perfmonUtil.GetNameOnPid(device, pidStr)
			if err != nil {
				packageName = ""
			}
		}
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)

		if (pid == -1 && packageName == "") &&
			!perfOptions2.SystemCPU &&
			!perfOptions2.SystemGPU &&
			!perfOptions2.SystemNetWorking &&
			!perfOptions2.SystemMem {
			perfOptions2.SystemCPU = true
			perfOptions2.SystemMem = true
			perfOptions2.SystemGPU = true
			perfOptions2.SystemNetWorking = true
		}
		if (pid != -1 || packageName != "") &&
			!perfOptions2.ProcMem &&
			!perfOptions2.ProcCPU &&
			!perfOptions2.ProcThreads &&
			!perfOptions2.ProcFPS {
			perfOptions2.ProcMem = true
			perfOptions2.ProcCPU = true
			perfOptions2.ProcThreads = true
			perfOptions2.ProcFPS = true
		}
		perfmonUtil.PackageName = packageName
		perfmonUtil.Pid = pidStr

		exitCtx, exitChancel := context.WithCancel(context.Background())

		perfmonUtil.UpdatePIDAndPackageCurrentActivity(device, exitCtx)

		perfmonUtil.IntervalTime = float64(refreshTime) / 1000

		timer := time.Tick(time.Duration(int(perfmonUtil.IntervalTime * float64(time.Second))))
		var mu sync.Mutex
		var wg sync.WaitGroup
		for {
			select {
			case <-sig:
				exitChancel()
				os.Exit(0)
			case <-timer:
				go func() {
					perfData := &entity.PerfmonData{
						TimeStamp: time.Now().Unix(),
						System: &entity.SystemInfo{
							CPU:         make(map[string]*entity.SystemCPUInfo),
							MemInfo:     &entity.SystemMemInfo{},
							NetworkInfo: make(map[string]*entity.SystemNetworkInfo),
						},
						Process: &entity.ProcessInfo{
							CPUInfo:    &entity.ProcCpuInfo{},
							MemInfo:    &entity.ProcMemInfo{},
							FPSInfo:    &entity.ProcFPSInfo{},
							ThreadInfo: &entity.ProcTreadsInfo{},
						},
					}

					wg.Add(2)
					go func() {
						defer wg.Done()
						data := perfmonUtil.GetSystemCPU2(device, perfOptions2)
						if data != nil {
							mu.Lock()
							defer mu.Unlock()
							perfData.System.CPU = data.CPU
						}
					}()
					go func() {
						defer wg.Done()
						data := perfmonUtil.GetSystemMem2(device, perfOptions2)
						if data != nil {
							mu.Lock()
							defer mu.Unlock()
							perfData.System.MemInfo = data.MemInfo
						}
					}()
					wg.Wait()
					fmt.Println(perfData.ToFormat())
					// wait all goroutine done to do something
				}()

			}
		}
	},
}
