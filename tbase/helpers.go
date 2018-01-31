package tbase

import (
	"strconv"
	"strings"
)

func IntList(val string) []int {
	s := StringList(val)
	if s == nil {
		return nil
	}
	x := make([]int, len(s))
	for k, v := range s {
		i, _ := strconv.Atoi(v)
		x[k] = i
	}
	return x
}

func StringList(val string) []string {
	if val == "" {
		return nil
	}
	return strings.Split(val, ",")
}
