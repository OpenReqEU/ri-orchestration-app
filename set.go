package main

import (
	"strings"
)

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

func (s *set) Remove(value string) {
	delete(s.m, value)
}

func (s *set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

func (s *set) String() string {
	var vals []string
	for val := range s.m {
		vals = append(vals, val+"|"+s.m[val]+"\t")
	}

	return strings.Join(vals, " ")
}
