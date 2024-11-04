package main


import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"

	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/SonicCloudOrg/sonic-android-supply/src/entity"
	"github.com/SonicCloudOrg/sonic-android-supply/src/perfmonUtil"
	"github.com/SonicCloudOrg/sonic-android-supply/src/util"
	"github.com/spf13/cobra"
)

var perfmonCmd = &cobra.Command{
	Use:   "perfmon",
	Short: "Get device performance",
	Long:  "Get device performance",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if serial == "" {
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
				serial = m[1]
				break
			}
		}

		device := adb.NewClient("").DeviceWithSerial2(serial)
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

		if (pidStr == "" && packageName == "") &&
			!perfOptions.SystemCPU &&
			!perfOptions.SystemGPU &&
			!perfOptions.SystemNetWorking &&
			!perfOptions.SystemMem {
			sysAllParamsSet()
		}
		if (pidStr != "" || packageName != "") &&
			!perfOptions.ProcMem &&
			!perfOptions.ProcCPU &&
			!perfOptions.ProcThreads &&
			!perfOptions.ProcFPS {
			//sysAllParamsSet()
			perfOptions.ProcMem = true
			perfOptions.ProcCPU = true
			perfOptions.ProcThreads = true
			perfOptions.ProcFPS = true
		}
		perfmonUtil.PackageName = packageName
		perfmonUtil.Pid = pidStr

		exitCtx, exitChancel := context.WithCancel(context.Background())

		perfmonUtil.UpdatePIDAndPackageCurrentActivity(device, exitCtx)

		perfmonUtil.IntervalTime = float64(refreshTime) / 1000

		var perfDataChan = make(chan *entity.PerfmonData)
		perfmonUtil.GetSystemCPU(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetSystemMem(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetSystemNetwork(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetProcCpu(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetProcMem(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetProcFPS(device, perfOptions, perfDataChan, exitCtx)
		perfmonUtil.GetProcThreads(device, perfOptions, perfDataChan, exitCtx)

		for {
			select {
			case <-sig:
				exitChancel()
				os.Exit(0)
			case perfData, ok := <-perfDataChan:
				if ok {
					fmt.Println(util.Format(perfData, isFormat, isJson))
				}
			}
		}
	},
}

var (
	serial   string
	isFormat bool
	isJson   bool
	perfOptions entity.PerfOption
	pid         int
	packageName string
	refreshTime int
	isForce     bool
)

func sysAllParamsSet() {
	perfOptions.SystemCPU = true
	perfOptions.SystemMem = true
	perfOptions.SystemGPU = true
	perfOptions.SystemNetWorking = true
}

func init() {
	rootCmd.AddCommand(perfmonCmd)
	perfmonCmd.Flags().StringVarP(&serial, "serial", "s", "", "device serial (default first device)")
	perfmonCmd.Flags().IntVarP(&pid, "pid", "d", -1, "get PID data")
	perfmonCmd.Flags().StringVarP(&packageName, "package", "p", "", "app package name")
	perfmonCmd.Flags().BoolVar(&perfOptions.SystemCPU, "sys-cpu", false, "get system cpu data")
	perfmonCmd.Flags().BoolVar(&perfOptions.SystemMem, "sys-mem", false, "get system memory data")
	//perfmonCmd.Flags().BoolVar(&sysDisk, "sys-disk", false, "get system disk data")
	perfmonCmd.Flags().BoolVar(&perfOptions.SystemNetWorking, "sys-network", false, "get system networking data")
	//perfmonCmd.Flags().BoolVar(&perfOptions.SystemGPU, "gpu", false, "get gpu data")
	perfmonCmd.Flags().BoolVar(&perfOptions.ProcFPS, "proc-fps", false, "get fps data")
	perfmonCmd.Flags().BoolVar(&perfOptions.ProcThreads, "proc-threads", false, "get process threads")
	//perfmonCmd.Flags().BoolVar(&, "proc-network", false, "get process network data")
	perfmonCmd.Flags().BoolVar(&perfOptions.ProcCPU, "proc-cpu", false, "get process cpu data")
	perfmonCmd.Flags().BoolVar(&perfOptions.ProcMem, "proc-mem", false, "get process mem data")
	perfmonCmd.Flags().BoolVar(&isForce, "force-out", false, "force update pid perf data(applicable to applications being restarted by kill)")
	perfmonCmd.Flags().IntVarP(&refreshTime, "refresh", "r", 1000, "data refresh time (millisecond)")
	perfmonCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
	perfmonCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
}
