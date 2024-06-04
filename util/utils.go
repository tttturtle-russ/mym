package util

import (
	"Mym/consts"
	"github.com/zcalusic/sysinfo"
	"sort"
)

func StringIn(target string, arr []string) bool {
	sort.Strings(arr)
	i := sort.SearchStrings(arr, target)
	return i < len(arr) && arr[i] == target
}

func IntIn(target int, arr []int) bool {
	sort.Ints(arr)
	i := sort.SearchInts(arr, target)
	return i < len(arr) && arr[i] == target
}

func IsValidIPVersion(ipVersion consts.IPVersion) bool {
	return IntIn(int(ipVersion), []int{4, 6})
}

func IsValidProtocol(protocol consts.Protocol) bool {
	return StringIn(string(protocol), []string{"http", "https"})
}

func GetOSRelease() string {
	var sysInfo sysinfo.SysInfo
	sysInfo.GetSysInfo()
	return sysInfo.OS.Name
}
