package gollector_metrics

/*
#include <unistd.h>
int get_hz(void) {
  return sysconf(_SC_CLK_TCK);
}
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func getJiffies() (jiffies int64, cpus int64, err error) {
	content, err := ioutil.ReadFile("/proc/stat")

	if err != nil {
		return 0, 0, fmt.Errorf("While processing the cpu_usage package: %s", err.Error())
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if strings.Index(line, "cpu ") == 0 {
			/* cpu with no number is the aggregate of all of them -- this is what we
			 * want to parse
			 */
			parts := strings.Split(line, " ")

			/* 2 - 11 are the time aggregates */
			for x := 2; x <= 11; x++ {

				/* 5 is the idle time, which we don't want */
				if x == 5 {
					continue
				}

				/* integer all the things */
				part, err := strconv.Atoi(parts[x])

				if err != nil {
					return 0, 0, fmt.Errorf("Could not convert integer from string while processing cpu_usage: %s, Error: %s", parts[x], err.Error())
				}

				jiffies += int64(part)
			}

		} else if strings.Index(line, "cpu") == 0 {
			/* cpu with a number is the specific time -- cheat and use this for the
			 * processor count since we've already read it
			 */
			cpus++
		}
	}

	return jiffies, cpus, nil
}

func getJiffyDiff() (int64, int64, error) {
	time1, cpus, err := getJiffies()

	if err != nil {
		return 0, 0, nil
	}

	time.Sleep(1 * time.Second)
	time2, _, err := getJiffies()

	if err != nil {
		return 0, 0, nil
	}

	return time2 - time1, cpus, nil
}

/*
 Obtain the CPUUsage() at the current point. More accurately, it returns a
 [2]float and error based on two sets of jiffies collection, which is
 calculated over a second's time and divided by the kernel's HZ value. The
 float values respectively are the current CPU cores in use and the number of
 total cores in the system

 Note that because a second is required to gather an accurate value, this
 invokes time.Sleep and will block any current goroutine while it collects
 these values.
*/
func CPUUsage() ([2]float64, error) {
	diff, cpus, err := getJiffyDiff()

	if err != nil {
		return [2]float64{0, 0}, err
	}
	return [2]float64{(float64(diff) / float64(C.get_hz())), float64(cpus)}, nil
}
