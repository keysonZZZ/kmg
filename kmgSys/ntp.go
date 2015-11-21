package kmgSys

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgPlatform"
)

//目前只支持 Ubuntu
// TODO data race?
func SyncTime() {
	if !kmgPlatform.IsLinux(){
		return
	}
	if !kmgCmd.Exist("ntpdate") {
		kmgCmd.MustRun("apt-get install -y ntpdate")
	}
	kmgCmd.MustRun("ntpdate -u pool.ntp.org")
}
