package perfmonUtil

import (
	"fmt"
	"testing"

	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
)

func TestGetSystemMem2(t *testing.T) {
	device := adb.NewClient("").DeviceWithSerial2("AXGE022414002023")
	data := GetSystemMem2(device)
	fmt.Println(data.ToFormat())
}
