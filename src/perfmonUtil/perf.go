package perfmonUtil

import (
	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/SonicCloudOrg/sonic-android-supply/src/entity"
)

func GetSystemCPU2(device *adb.Device, perfOptions entity.PerfOption) *entity.SystemInfo {
	if !perfOptions.SystemCPU {
		return nil
	}
	systemInfo := &entity.SystemInfo{}
	err := getCPU(device, systemInfo)
	if err != nil {
		systemInfo.Error = append(systemInfo.Error, err.Error())
	}
	return systemInfo
}

func GetSystemMem2(device *adb.Device, perfOptions entity.PerfOption) *entity.SystemInfo {
	if !perfOptions.SystemMem {
		return nil
	}
	systemInfo := &entity.SystemInfo{
		MemInfo: &entity.SystemMemInfo{},
	}
	err := getMemInfo(device, systemInfo)
	if err != nil {
		systemInfo.Error = append(systemInfo.Error, err.Error())
	}
	return systemInfo
}

func GetSystemNetwork2(device *adb.Device, perfOptions entity.PerfOption) *entity.SystemInfo {
	if !perfOptions.SystemNetWorking {
		return nil
	}
	systemInfo := &entity.SystemInfo{}
	err := getInterfaces(device, systemInfo)
	if err != nil {
		systemInfo.Error = append(systemInfo.Error, err.Error())
	}
	err = getInterfaceInfo(device, systemInfo)
	if err != nil {
		systemInfo.Error = append(systemInfo.Error, err.Error())
	}
	return systemInfo
}

func GetProcCpu2(device *adb.Device, perfOptions entity.PerfOption) *entity.ProcessInfo {
	if !perfOptions.ProcCPU {
		return nil
	}
	return getProcCpu(device)
}

func GetProcMem2(device *adb.Device, perfOptions entity.PerfOption) *entity.ProcessInfo {
	if !perfOptions.ProcMem {
		return nil
	}
	return getProcMem(device)
}

func GetProcFPS2(device *adb.Device, perfOptions entity.PerfOption) *entity.ProcessInfo {
	if !perfOptions.ProcFPS {
		return nil
	}
	return getFPS(device)
}

func GetProcThreads2(device *adb.Device, perfOptions entity.PerfOption) *entity.ProcessInfo {
	if !perfOptions.ProcThreads {
		return nil
	}
	return getThreads(device)
}
