package variables

// App name and version
var AppName string
var AppVersion string

// CMD
var DiskUtil string
var HdiUtil string

func init() {
	// App name and version
	AppName = "MacOS Ramdisk Creator"
	AppVersion = "1.0.0"

	// Set commands
	DiskUtil = "diskutil"
	HdiUtil = "hdiutil"
}
