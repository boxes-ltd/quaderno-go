package quaderno

import (
	"reflect"
	"runtime/debug"
	"strings"
)

func getVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	pkgPath := reflect.TypeOf(Client{}).PkgPath()

	if info.Main.Path != "" && strings.HasPrefix(pkgPath, info.Main.Path) {
		return info.Main.Version
	}

	var matchVersion string
	var matchLen int

	for _, dep := range info.Deps {
		if strings.HasPrefix(pkgPath, dep.Path) {
			if len(dep.Path) > matchLen {
				matchLen = len(dep.Path)
				matchVersion = dep.Version
			}
		}
	}

	if matchVersion != "" {
		return matchVersion
	}

	return ""
}
