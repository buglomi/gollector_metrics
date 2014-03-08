package gollector_metrics

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

/*
MemoryUsage() returns a map[string]int of keys describing the values and an
integer with the value in bytes.

Values Are:

* Total: total memory

* Free: free memory

* Used: used memory

* Swap Total: total swap

* Swap Free: available swap
*/
func MemoryUsage() (map[string]int, error) {
	content, err := ioutil.ReadFile("/proc/meminfo")

	var total, buffers, cached, free, swap_total, swap_free int

	if err != nil {
		return nil, fmt.Errorf("While processing the mem_usage package: %s", err.Error())
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		parts := strings.Split(line, " ")
		id := len(parts) - 2

		switch parts[0] {
		case "MemTotal:":
			total, err = strconv.Atoi(parts[id])
		case "MemFree:":
			free, err = strconv.Atoi(parts[id])
		case "Cached:":
			cached, err = strconv.Atoi(parts[id])
		case "Buffers:":
			buffers, err = strconv.Atoi(parts[id])
		case "SwapTotal:":
			swap_total, err = strconv.Atoi(parts[id])
		case "SwapFree:":
			swap_free, err = strconv.Atoi(parts[id])
		}

		if err != nil {
			return nil, fmt.Errorf("Could not convert integer from string while processing cpu_usage: %s: error: %s", parts[id], err.Error())
		}
	}

	return map[string]int{
		"Total":      total * 1024,
		"Free":       (buffers + cached + free) * 1024,
		"Used":       total*1024 - ((buffers + cached + free) * 1024),
		"Swap Total": swap_total * 1024,
		"Swap Free":  swap_free * 1024,
	}, nil
}
