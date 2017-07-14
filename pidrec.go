package pidrec

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// PidFile represents a file that contains a process id.
type PidFile struct {
	path string
	pid  int
}

// Remove deletes the file associated with the PidFile. The file is not deleted
// if the process id recorded in the file is different than the one associated
// with PidFile.
func (pf *PidFile) Remove() error {
	if pf == nil {
		return nil
	}
	if len(pf.path) == 0 {
		return nil
	}
	pid, err := getPid(pf.path)
	if err != nil {
		return err
	}
	if pid != pf.pid {
		return fmt.Errorf("%s: expecting %d, found %d", pf.path, pf.pid, pid)
	}

	return os.Remove(pf.path)
}

func getPid(file string) (int, error) {
	f, err := os.Open(file)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	i := 1
	var str string
	for s.Scan() {
		if i > 1 {
			return -1, fmt.Errorf("%s has spurious content", file)
		}
		str = s.Text()
		i++
	}
	str = strings.Trim(str, " \t\n")
	pid, err := strconv.Atoi(str)
	if err != nil {
		return -1, fmt.Errorf("%s: content not a number '%v'", file, err)
	}
	return pid, nil
}

// MustWriteTo writes the pid of the calling process to the given file.
// It panics on any error or if the file already exists.
func MustWriteTo(pidFile string) *PidFile {
	if len(pidFile) == 0 {
		panic("empty filename")
	}
	fi, err := os.Stat(pidFile)
	if err != nil {
		if err.(*os.PathError).Err.Error() != "no such file or directory" {
			panic(err)
		}
	} else if fi != nil {
		panic(fmt.Sprintf("%s already exist, is another instance running?", pidFile))
	}

	pid := os.Getpid()
	err = ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d\n", pid)), 0644)
	if err != nil {
		panic(err)
	}
	return &PidFile{
		path: pidFile,
		pid:  pid,
	}
}
