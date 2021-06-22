package types

import (
	"strconv"
)

func SliceStringToInterface(d []string) []interface{} {
	a := []interface{}{}
	for _, s := range d {
		a = append(a, s)
	}
	return a
}

func SliceStringToInt(d []string) ([]int, error) {
	a := []int{}
	for _, s := range d {
		si, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		a = append(a, si)
	}
	return a, nil
}
