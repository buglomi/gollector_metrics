package gollector_metrics

/*
Count the number of processes running with the given commandline. Returns int,
error. Count will be zero on any error.
*/
func ProcessCount(command string) (int, error) {
	pids, err := GetPids(command)

	if err != nil {
		return 0, err
	}

	return len(pids), nil
}
