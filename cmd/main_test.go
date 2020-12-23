package main

import (
	"testing"
)

// export OOPS_ENV_FILE=${PWD}/tests/env
func TestLoadConfigNoTLS(t *testing.T) {
	loadConfig()

	if conf.TLS {
		t.Fatal("Expected to not use TLS")
	}

	// TODO add more checks
}

func TestGetCSSBox(t *testing.T) {
	// rice.MustFindBox() panics if it can't find the provided dir
	getRiceBox("../css")
}
