package perfmonUtil

import (
	"fmt"
	"testing"

	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
)

var device *adb.Device

func setupDevice(serial string) {
	device = adb.NewClient("").DeviceWithSerial2(serial)
}

func TestGetFPS(t *testing.T) {
	setupDevice("AXGE022414002023")
	r, _ := getProcessFPSBySurfaceFlinger(device, "com.android.browser")
	fmt.Println(r)
}

func TestGetPackageCurrentActivity(t *testing.T) {
	setupDevice("AXGE022414002023")
	fmt.Println(getPackageCurrentActivity(device, "com.tencent.mm", "19799"))
}

func TestGet(t *testing.T) {
	setupDevice("AXGE022414002023")
	fmt.Println(GetPidOnPackageName(device, "com.tencent.mm"))
}
