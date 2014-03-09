package gollector_metrics

/*
#include <unistd.h>
unsigned int get_pgsz(void) {
  return sysconf(_SC_PAGESIZE);
}
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func ProcessMemoryUsage(command string) (uint, error) {
	total := uint(0)
	page_size := uint(C.get_pgsz())

	pids, err := GetPids(command)

	if err != nil {
		return 0, err
	}

	for _, pid := range pids {
		path := "/proc/" + pid + "/statm"
		f, err := os.Open(path)

		if err != nil {
			return 0, fmt.Errorf("Could not open " + path + ": " + err.Error())
		}

		defer f.Close()

		content, err := ioutil.ReadAll(f)

		if err != nil {
			return 0, fmt.Errorf("Could not read from " + path + ": " + err.Error())
		}

		parts := strings.Split(string(content), " ")
		mem, err := strconv.Atoi(parts[1])

		if err != nil {
			return 0, fmt.Errorf("Trouble converting resident size " + parts[1] + " to integer: " + err.Error())
		}

		total += uint(mem) * page_size
	}

	return total, nil
}
