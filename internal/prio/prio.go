// Package prio
package prio

import "fmt"

type Priority int

const (
	PRIORITYHIGH Priority = iota + 1
	PRIORITYLOW
)

func PriorityFromString(prioString string) (Priority, error) {
	switch prioString {
	case "high":
		return PRIORITYHIGH, nil
	case "low":
		return PRIORITYLOW, nil
	}

	return 0, fmt.Errorf("unknown priority: %v", prioString)
}
