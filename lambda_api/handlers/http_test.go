package handlers

import "testing"

func Test_functionNameFromURL(t *testing.T) {
	f := functionNameFromURL("/2015-03-31/functions/test-multiprocessing/invocations")
	if f != "test-multiprocessing" {
		t.Error("incorrect parsing for valid URL")
	}

	f = functionNameFromURL("/2015-03-31/functions/test-multiprocessing/invocation")
	if f != "" {
		t.Error("incorrect parsing for invalid URL")
	}
}
