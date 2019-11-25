package main

/*
 * set implementation
 */
var exists = struct{}{}

type set struct {
	m map[string]string
}

// NewSet is a custom implementation for imitating a set in golang
func NewSet() *set {
	s := &set{}
	s.m = make(map[string]string)
	return s
}

func (s *set) Add(value string, interval string) {
	s.m[value] = interval
}
