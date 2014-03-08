package gollector_metrics

/*
#include <stdlib.h>
*/
import "C"

/*
This calls the getloadavg(2) call and returns a 3 element tuple of C.double,
corresponding to the 1 minute, 5 minute, and 15 minute load averages
respectively.
*/
func LoadAverage() [3]C.double {
	var loadavg [3]C.double

	C.getloadavg(&loadavg[0], 3)

	return loadavg
}
