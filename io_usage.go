package gollector_metrics

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const DISKSTATS_FILE = "/proc/diskstats"

const (
	DEVICE_DISK uint = iota
	DEVICE_DM        = iota
)

const (
	LINE_ID           uint = 0
	LINE_DEVICE       uint = 2
	LINE_FIRST_METRIC uint = 3
)

var device_to_diskstat_id = map[uint]uint{
	DEVICE_DISK: 8,
	DEVICE_DM:   252,
}

var metric_names = []string{
	"reads issued",
	"reads merged",
	"sectors read",
	"time reading (ms)",
	"writes completed",
	"writes merged",
	"sectors written",
	"time writing (ms)",
	"iops in progress",
	"io time (ms)",
	"weighted io time (ms)",
}

var memLastMetrics map[string]map[string]uint64
var memRWMutex sync.RWMutex

func getDeviceType(device_name string) uint {
	byte_dn := []byte(device_name)

	matched, _ := regexp.Match("^dm-", byte_dn)

	if matched {
		return device_to_diskstat_id[DEVICE_DM]
	}

	return device_to_diskstat_id[DEVICE_DISK]
}

func initLastMetrics(device string) (new_metrics bool) {
	new_metrics = false

	if memLastMetrics == nil {
		memRWMutex.Lock()
		memLastMetrics = make(map[string]map[string]uint64)
		memRWMutex.Unlock()
		new_metrics = true
	}

	if memLastMetrics[device] == nil {
		memRWMutex.Lock()
		memLastMetrics[device] = make(map[string]uint64)
		memRWMutex.Unlock()
		new_metrics = true
	}

	return new_metrics
}

func writeMetric(device string, metric string, value uint64) {
	memRWMutex.Lock()
	memLastMetrics[device][metric] = value
	memRWMutex.Unlock()
}

func readMetric(device string, metric string) (value uint64) {
	memRWMutex.RLock()
	value = memLastMetrics[device][metric]
	memRWMutex.RUnlock()

	return value
}

func getDiskMetrics(device string, device_type uint) (retval map[string]uint64, err error) {
	out, err := ioutil.ReadFile(DISKSTATS_FILE)

	if err != nil {
		return retval, err
	}

	lines := strings.Split(string(out), "\n")
	re, _ := regexp.Compile("[ \t]+")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := re.Split(line, -1)
		parts = parts[1:]

		device_type_parsed, err := strconv.ParseUint(parts[LINE_ID], 10, 64)

		if err != nil {
			return retval, err
		} else if uint(device_type_parsed) == device_type && parts[LINE_DEVICE] == device {
			retval = make(map[string]uint64)

			for i, key := range metric_names {
				retval[key], err = strconv.ParseUint(parts[LINE_FIRST_METRIC+uint(i)], 10, 64)

				if err != nil {
					return make(map[string]uint64), err
				}
			}
		}
	}

	return retval, err
}

/*
Calculate I/O usage. This call produces a difference per each call; that is,
the first call will return zeroes, the next call will return the metrics
collected since the last run, and this will continue for successive calls.

The device used must be a *real* disk. No tmpfs, procfs, etc.
*/
func IOUsage(device string) (map[string]uint64, error) {
	difference := make(map[string]uint64)
	device_type := getDeviceType(device)
	new_metrics := initLastMetrics(device)
	metrics, err := getDiskMetrics(device, device_type)

	/*if new_metrics {*/
	/*log.Log("debug", "New metrics, sending zeroes")*/
	/*}*/

	if err != nil {
		return nil, err
	}

	for metric, value := range metrics {
		if new_metrics {
			difference[metric] = 0
		} else {
			difference[metric] = value - readMetric(device, metric)
			if int64(value-readMetric(device, metric)) < 0 {
				difference[metric] = 0
			}
		}

		writeMetric(device, metric, value)
	}

	return difference, nil
}
