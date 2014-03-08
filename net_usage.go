package gollector_metrics

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var net_base_path = "/sys/class/net"
var file_pattern = filepath.Join(net_base_path, "%s/statistics/")

var file_map = map[string]string{
	"rx_bytes":   "Received (Bytes)",
	"tx_bytes":   "Transmitted (Bytes)",
	"tx_errors":  "Transmission Errors",
	"rx_errors":  "Reception Errors",
	"rx_packets": "Received (Packets)",
	"tx_packets": "Transmitted (Packets)",
}

var netLastMetrics map[string]map[string]uint64
var netRWMutex sync.RWMutex

func readFile(base_path string, metric string) (uint64, error) {
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
func NetUsage(device string) map[string]uint64 {
	new_metrics := false

	if netLastMetrics == nil {
		netRWMutex.Lock()
		netLastMetrics = make(map[string]map[string]uint64)
		netRWMutex.Unlock()
		new_metrics = true
	}

	if netLastMetrics[device] == nil {
		netRWMutex.Lock()
		netLastMetrics[device] = make(map[string]uint64)
		netRWMutex.Unlock()
		new_metrics = true
	}

	metrics := make(map[string]uint64)
	difference := make(map[string]uint64)

	base_path := fmt.Sprintf(file_pattern, device)

	for fn, metric := range file_map {
		result, err := readFile(base_path, fn)
		if err == nil {
			metrics[metric] = result
		} else {
			metrics[metric] = 0
		}
	}

	for metric, value := range metrics {
		if new_metrics {
			difference[metric] = 0
			netRWMutex.Lock()
			netLastMetrics[device][metric] = value
			netRWMutex.Unlock()
		} else {
			netRWMutex.RLock()
			difference[metric] = value - netLastMetrics[device][metric]
			netRWMutex.RUnlock()
			netRWMutex.Lock()
			netLastMetrics[device][metric] = value
			netRWMutex.Unlock()
		}

	}

	return difference
}
