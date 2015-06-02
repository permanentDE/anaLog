package filepath

import (
	"os"
	"path/filepath"

	"go.permanent.de/anaLog/v1/config"
)

func GetPath(file string) string {
	if filepath.IsAbs(file) {
		return file
	} else {
		path, err := filepath.Abs(GetWorkspace() + "/" + file)
		if err != nil {
			panic(err)
		}
		return path
	}
}

func GetWorkspace() string {
	var path string
	var err error

	if config.Std.AnaLog.Workspace != "" {
		path = config.Std.AnaLog.Workspace
	} else {
		path = filepath.Dir(os.Args[0])
	}

	path, err = filepath.Abs(path)

	if err != nil {
		panic(err)
	}

	return path
}
