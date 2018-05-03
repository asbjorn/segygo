package segygo

import (
	"flag"
	"fmt"
	"io"
	"os"
	"testing"
)

const filename = "/tmp/testSegy.segy"
const filenameOut = "/tmp/testSegyOut.segy"

func TestMain(m *testing.M) {
	flag.Parse()
	createTestFile()
	os.Exit(m.Run())
}

func check(err error) {
	if err != nil {
		fmt.Println("Error : %s", err.Error())
		os.Exit(1)
	}
}

func createTestFile() {
	srcFile, err := os.Open("data/qcSOTXCor.segy.f0.k1")
	defer srcFile.Close()
	check(err)

	destFile, err := os.Create(filename) // creates if file doesn't exist
	defer destFile.Close()
	check(err)

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	check(err)

	err = destFile.Sync()
	check(err)
}

// Test cases under here

func TestCreateSegy(t *testing.T) {
	s, err := CreateFile(filenameOut)
	if err != nil {
		t.Errorf("Unable to create new segy file: %s", s.Filename)
	}

	s.Header.Jobid = 10
	s.Header.Hns = 100
}

func TestReadSegy(t *testing.T) {
	s, err := OpenFile(filename)
	fmt.Println("SEG-Y file: ", filename)
	fmt.Println("# traces: ", s.NrTraces)
	hdr := s.Header
	if err != nil {
		fmt.Println("Error reading segy file..")
		t.Errorf("Error reading SEG-Y file: %s", s.Filename)
	} else {
		fmt.Println("Job id: ", hdr.Jobid)
		fmt.Println("Line no: ", hdr.Lino)
		fmt.Println("Reno: ", hdr.Reno)
		fmt.Println("Ntrpr: ", hdr.Ntrpr)
		fmt.Println("Nart: ", hdr.Nart)
		fmt.Println("Hdt: ", hdr.Hdt)
		fmt.Println("Dto: ", hdr.Dto)
		fmt.Println("Hns: ", hdr.Hns)
		fmt.Println("Nso: ", hdr.Nso)
		fmt.Println("Format: ", hdr.Format)
		fmt.Println("Fold: ", hdr.Fold)
	}

	//fmt.Println("SEG-Y binary header printout")
	//m := s.GetHeader()

	//for k,v := range m {
	//	fmt.Printf("key = %s, value = %d\n", k, v)
	//}

	//var firstTrace Trace
	//firstTrace, err = s.ReadTrace()
	//firstTrace, err = s.ReadTrace()
}

func TestReadFirstTrace(t *testing.T) {
	s, err := OpenFile(filename)
	fmt.Println("SEG-Y file: ", filename)
	fmt.Println("# traces: ", s.NrTraces)
	trace, err := s.ReadTrace()

	if err != nil {
		t.Errorf("Exception occured: %v", err)
	}

	fmt.Println("Got trace")
	if len(trace.Data) == 0 {
		t.Fatal("trace.Data got 0 samples!")
	} else {
		fmt.Printf("Trace has %d samples..\n", len(trace.Data))
	}
}
