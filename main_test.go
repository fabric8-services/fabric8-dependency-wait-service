package main

import (
	"runtime"
	"testing"
)

func Test_isIn(t *testing.T) {
	httpGood := []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}
	if !isIn(httpGood, 200) {
		t.Errorf("Expected 200 to be in httpGood")
	}
	if isIn(httpGood, 20) {
		t.Errorf("Expected 20 to not be in httpGood")
	}
}

func Test_isInPath(t *testing.T) {

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if err := isInPath("date"); err != nil {
			t.Errorf("Expected date to be in path.")
		}
	}

	if err := isInPath("dateabcdefgh"); err == nil {
		t.Errorf("Expected random cmd to not be in path.")
	}
}

func Test_splitPostgresURL(t *testing.T) {
	var h, p, hExp, pExp string
	var err error
	h, p, err = splitPostgresURL("postgres://127.0.0.1:230")
	if err != nil {
		t.Errorf("Expected error to be nil. Got %v.\n", err)
	}
	hExp = "127.0.0.1"
	if h != hExp {
		t.Errorf("Expected host to be %s. Got %s.\n", hExp, h)
	}
	pExp = "230"
	if p != pExp {
		t.Errorf("Expected port to be %s. Got %s.\n", pExp, p)
	}

	h, p, err = splitPostgresURL("postgres://127.0.0.1:230/")
	if err != nil {
		t.Errorf("Expected error to be %v. Got %v.\n", err)
	}
	hExp = "127.0.0.1"
	if h != hExp {
		t.Errorf("Expected host to be %s. Got %s.\n", hExp, h)
	}
	pExp = "230"
	if p != pExp {
		t.Errorf("Expected port to be %s. Got %s.\n", pExp, p)
	}

	h, p, err = splitPostgresURL("postgres://localhost/")
	if err != nil {
		t.Errorf("Expected error to be %v. Got %v.\n", err)
	}
	hExp = "localhost"
	if h != hExp {
		t.Errorf("Expected host to be %s. Got %s.\n", hExp, h)
	}
	pExp = ""
	if p != pExp {
		t.Errorf("Expected port to be %s. Got %s.\n", pExp, p)
	}

	h, p, err = splitPostgresURL("1postgres://localhost/")
	if err == nil {
		t.Errorf("Expected error to be non nil. Got %v.\n", err)
	}
}
