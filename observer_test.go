package main

import "testing"

func TestUpdateApp(t *testing.T) {
	induceServerError = false
	updateApp("eu.openreq")
}
