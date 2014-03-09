package gollector_metrics

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

/*
 * Get the process ids for a given process name, full path is required.
 *
 * Returns strings because it's easier for most of the things we'll use this
 * for.
 */
func GetPids(process string) ([]string, error) {
	pids := []string{}

	dir, err := os.Open("/proc")

	if err != nil {
		return nil, fmt.Errorf("Could not open /proc for reading: " + err.Error())
	}

	defer dir.Close()

	proc_files, err := dir.Readdirnames(0)

	if err != nil {
		return nil, fmt.Errorf("Could not read directory names from /proc: " + err.Error())
	}

	all_pids := []string{}
	// XXX totally cheating here -- the only all-numeric filenames in this dir
	// will be pid directories. This should be faster than 4 bajillion stat
	// calls (that I'd have to do this to anyway).
	for _, fn := range proc_files {
		_, err := strconv.Atoi(fn)
		if err == nil {
			all_pids = append(all_pids, fn)
		}
	}

	for _, pid := range all_pids {
		path := "/proc/" + pid + "/cmdline"
		file, err := os.Open(path)

		if err != nil {
			return nil, fmt.Errorf("Could not open " + path + ": " + err.Error())
		}

		defer file.Close()

		cmdline, err := ioutil.ReadAll(file)

		if err != nil {
			return nil, fmt.Errorf("Could not read from " + path + ": " + err.Error())
		}

		cmdline_parts := strings.Split(string(cmdline), "\x00")

		if len(cmdline_parts) > 1 {
			cmdline_parts = cmdline_parts[0 : len(cmdline_parts)-1]
		}

		string_cmd := strings.Join(cmdline_parts, " ")

		if string_cmd == process {
			pids = append(pids, pid)
		}
	}

	return pids, nil
}
