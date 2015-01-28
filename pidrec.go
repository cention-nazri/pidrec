package pidrec

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type PidFile struct {
	path string
}

// Remove deletes the file held by the PidFile.
func (pf *PidFile) Remove() {
	if pf == nil {
		return
	}
	if len(pf.path) == 0 {
		return
	}
	err := os.Remove(pf.path)
	if err != nil {
		log.Fatal(err)
	}
}

// MustWriteTo writes the pid of the calling process to the given file.
// It panic on any error or if the file already exists.
func MustWriteTo(pidFile string) *PidFile {
	if len(pidFile) == 0 {
		log.Fatal("Error: Filename is empty")
	}
	fi, err := os.Stat(pidFile)
	if err != nil {
		if err.(*os.PathError).Err.Error() != "no such file or directory" {
			log.Fatal(err)
		}
	} else if fi != nil {
		log.Fatal("Pidfile already exist. Is another instance running?: ", pidFile)
	}

	err = ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return &PidFile{pidFile}
}
