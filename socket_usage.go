package gollector_metrics

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

/*
Types of sockets for SocketUsage()
*/
var SOCK_TYPES = []string{
	"tcp",
	"tcp6",
	"udp",
	"udp6",
	"udplite",
	"udplite6",
	"unix",
}

/*
Given a socket type in SOCK_TYPES, yields the count of sockets for that type.
Returns 0, error for any error condition.
*/
func SocketUsage(sock_type string) (int, error) {
	found_sock_type := false

	for _, val := range SOCK_TYPES {
		if sock_type == val {
			found_sock_type = true
			break
		}
	}

	if !found_sock_type {
		return 0, fmt.Errorf("Invalid socket type: " + sock_type)
	}

	f, err := os.Open("/proc/self/net/" + sock_type)

	if err != nil {
		return 0, fmt.Errorf("Could not open socket information for " + sock_type + ": " + err.Error())
	}

	defer f.Close()

	content, err := ioutil.ReadAll(f)

	if err != nil {
		return 0, fmt.Errorf("Trouble reading socket information for type " + sock_type + ": " + err.Error())
	}

	lines := strings.Split(string(content), "\n")
	return len(lines) - 1, nil // there's a one line header in these files
}
