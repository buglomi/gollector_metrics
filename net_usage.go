package gollector_metrics

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const NET_BASE_PATH = "/sys/class/net" // base path of the network device statistics files.

var file_pattern = filepath.Join(NET_BASE_PATH, "%s/statistics/")

var file_map = map[string]string{
	"rx_bytes":   "Received (Bytes)",
	"tx_bytes":   "Transmitted (Bytes)",
	"tx_errors":  "Transmission Errors",
	"rx_errors":  "Reception Errors",
	"rx_packets": "Received (Packets)",
	"tx_packets": "Transmitted (Packets)",
}

type NetUsage struct {
	LastMetrics map[string]map[string]uint64
	RWMutex     sync.RWMutex
}

func (nu *NetUsage) readFile(base_path string, metric string) (uint64, error) {
	out, err := ioutil.ReadFile(filepath.Join(base_path, metric))

	if err != nil {
		return 0, err
	}

	out_i, err := strconv.ParseUint(strings.Split(string(out), "\n")[0], 10, 64)

	return out_i, err
}

/*
Collect statistics on network usage. This call returns zeroes on its first
invocation and then returns the difference on each successive poll. The device
must be a valid network device, such as "eth0" or "p2p0".
*/
func (nu *NetUsage) Metric(device string) map[string]uint64 {
	new_metrics := false

	if nu.LastMetrics == nil {
		nu.RWMutex.Lock()
		nu.LastMetrics = make(map[string]map[string]uint64)
		nu.RWMutex.Unlock()
		new_metrics = true
	}

	if nu.LastMetrics[device] == nil {
		nu.RWMutex.Lock()
		nu.LastMetrics[device] = make(map[string]uint64)
		nu.RWMutex.Unlock()
		new_metrics = true
	}

	metrics := make(map[string]uint64)
	difference := make(map[string]uint64)

	base_path := fmt.Sprintf(file_pattern, device)

	for fn, metric := range file_map {
		result, err := nu.readFile(base_path, fn)
		if err == nil {
			metrics[metric] = result
		} else {
			metrics[metric] = 0
		}
	}

	for metric, value := range metrics {
		if new_metrics {
			difference[metric] = 0
			nu.RWMutex.Lock()
			nu.LastMetrics[device][metric] = value
			nu.RWMutex.Unlock()
		} else {
			nu.RWMutex.RLock()
			difference[metric] = value - nu.LastMetrics[device][metric]
			nu.RWMutex.RUnlock()
			nu.RWMutex.Lock()
			nu.LastMetrics[device][metric] = value
			nu.RWMutex.Unlock()
		}

	}

	return difference
}
