// +build windows
package functions

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"unicode/utf16"
	"syscall"
	"strconv"

	"golang.org/x/sys/windows/registry"
)

func isElevated() bool {
	currentUser, _ := user.Current()
	return currentUser.Uid == "0"
}
func getArchitecture() string {
	return runtime.GOARCH
}
func getProcessName() string {
	name, err := os.Executable()
	if err != nil {
		return ""
	} else {
		return name
	}
}
func getDomain() string {
	const format = windows.ComputerNameDnsDomain
	n := uint32(64)
	b := make([]uint16, n)
	err := windows.GetComputerNameEx(format, &b[0], &n)
	if err == nil {
		return syscall.UTF16ToString(b[:n])
	}
	return ""
}
func getStringFromBytes(data [65]byte) string {
	stringData := make([]byte, 0, 0)
	for i := range data {
		if data[i] == 0 {
			return string(stringData[:])
		} else {
			stringData = append(stringData, data[i])
		}
	}
	return string(stringData[:])
}
func getOS() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()

	cv, _, err := k.GetStringValue("CurrentVersion")
	if err != nil {
		cv = ""
	}
	pn , _, err := k.GetStringValue("ProductName")
	if err != nil {
		pn = ""
	}
	maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		maj = 0
	}
	min, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
	if err != nil {
		min = 0
	}
	cb, _, err := k.GetStringValue("CurrentBuild")
	if err != nil {
		cb = ""
	}

	return string(pn) + "\n" +
	getHostname() + "\n" +
	cv + "\n" +
	pn + " " + strconv.Itoa(int(maj)) + "." + strconv.Itoa(int(min)) + "." + cb + " " + getArchitecture() + "\n" +
	getArchitecture()
}
func getUser() string {
	currentUser, err := user.Current()
	if err != nil {
		return ""
	} else {
		return currentUser.Username
	}
}
func getPID() int {
	return os.Getpid()
}
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	} else {
		return hostname
	}
}

// Helper function to convert DWORD byte counts to
// human readable sizes.
func UINT32ByteCountDecimal(b uint32) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint32(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float32(b)/float32(div), "kMGTPE"[exp])
}

// Helper function to convert LARGE_INTEGER byte
//  counts to human readable sizes.
func UINT64ByteCountDecimal(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// Helper function to build a string from a WCHAR string
func UTF16ToString(s []uint16) []string {
	var results []string
	cut := 0
	for i, v := range s {
		if v == 0 {
			if i-cut > 0 {
				results = append(results, string(utf16.Decode(s[cut:i])))
			}
			cut = i + 1
		}
	}
	if cut < len(s) {
		results = append(results, string(utf16.Decode(s[cut:])))
	}
	return results
}
