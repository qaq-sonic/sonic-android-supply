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
