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

type IOUsage struct {
	LastMetrics map[string]map[string]uint64
	RWMutex     sync.RWMutex
}

func (io *IOUsage) getDeviceType(device_name string) uint {
	byte_dn := []byte(device_name)

	matched, _ := regexp.Match("^dm-", byte_dn)

	if matched {
		return device_to_diskstat_id[DEVICE_DM]
	}

	return device_to_diskstat_id[DEVICE_DISK]
}

func (io *IOUsage) initLastMetrics(device string) (new_metrics bool) {
	new_metrics = false

	if io.LastMetrics == nil {
		io.RWMutex.Lock()
		io.LastMetrics = make(map[string]map[string]uint64)
		io.RWMutex.Unlock()
		new_metrics = true
	}

	if io.LastMetrics[device] == nil {
		io.RWMutex.Lock()
		io.LastMetrics[device] = make(map[string]uint64)
		io.RWMutex.Unlock()
		new_metrics = true
	}

	return new_metrics
}

func (io *IOUsage) writeMetric(device string, metric string, value uint64) {
	io.RWMutex.Lock()
	io.LastMetrics[device][metric] = value
	io.RWMutex.Unlock()
}

func (io *IOUsage) readMetric(device string, metric string) (value uint64) {
	io.RWMutex.RLock()
	value = io.LastMetrics[device][metric]
	io.RWMutex.RUnlock()

	return value
}

func (io *IOUsage) getDiskMetrics(device string, device_type uint) (retval map[string]uint64, err error) {
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
func (io *IOUsage) Metric(device string) (map[string]uint64, error) {
	difference := make(map[string]uint64)
	device_type := io.getDeviceType(device)
	new_metrics := io.initLastMetrics(device)
	metrics, err := io.getDiskMetrics(device, device_type)

	if err != nil {
		return nil, err
	}

	for metric, value := range metrics {
		if new_metrics {
			difference[metric] = 0
		} else {
			difference[metric] = value - io.readMetric(device, metric)
			if int64(value-io.readMetric(device, metric)) < 0 {
				difference[metric] = 0
			}
		}

		io.writeMetric(device, metric, value)
	}

	return difference, nil
}
