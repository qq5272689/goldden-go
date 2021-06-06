package base_dir

import (
	"os"
	"path/filepath"
)

func GetBaseDir() string {
	selfpath, _ := filepath.Abs(os.Args[0])
	Basedir, _ := filepath.Split(selfpath)
	return Basedir
}
