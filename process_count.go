package gollector_metrics

func ProcessCount(command string) (int, error) {
	pids, err := GetPids(command)

	if err != nil {
		return 0, err
	}

	return len(pids), nil
}
