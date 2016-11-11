package segygo

import (
	"testing"
)


func TestCreateSegy(t *testing.T) {
	var filename = "/tmp/testSegy.segy"

	s, err := CreateFile(filename)
	if err != nil {
		t.Errorf("Unable to create new segy file: %s", s.Filename)
	}

	s.Header.Jobid = 10
	s.Header.Hns = 100
}
