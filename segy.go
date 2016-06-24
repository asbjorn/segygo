package segy

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"fmt"
	"log"
	"unsafe"
)

const Version = "0.1"
const SEGY_DESC_HDR_LEN = 3200
const SEGY_BIN_HDR_LEN = 400
const SEGY_TRACE_HDR_LEN = 240

type BinHeader struct {
	Jobid		int32
	Lino		int32
	Reno		int32
	Ntrpr		int16
	Nart		int16
	Hdt			uint16
	Dto			uint16
	Hns			uint16
	Nso			uint16
	Format		int16
	Fold		int16
	Tsort		int16
	Vscode		int16
	Hsfs		int16
	Hsfe		int16
	Hslen		int16
	Hstyp		int16
	Schn		int16
	Hstas		int16
	Hstae		int16
	Htatyp		int16
	Hcorr		int16
	Bgrcv		int16
	Rcvm		int16
	Mfeet		int16
	Polyt		int16
	Vgpol		int16
	Hunass		[170]int16 // unassigned
}

type TraceHeader struct {
	tull		int32
}

type Trace struct {
	TraceHeader
	Data		*byte
}

type SegyFile struct {
	Filename string
	Header	BinHeader
	NrTraces int64
	file 	*os.File
}


func OpenFile(filename string) (SegyFile, error) {
	var s SegyFile;
	var binHdr BinHeader;
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return s, err
	}

	s.Filename = filename

	accum := []byte{}
	accum = append(accum, b...)

	accum2 := accum[3200:]
	r := bytes.NewReader(accum2)
	fmt.Println("Number of bytes:", r.Len())

	if err = binary.Read(r, binary.BigEndian, &binHdr); err != nil {
		fmt.Println("Error reading segy file (bigendian). ", err)
		return s, err
	}

	// Open and store the os.File object in our struct
	file, err := os.Open(s.Filename)
	s.file = file

	s.Header = binHdr
	s.NrTraces = s.getNrTraces()

	return s, err
}

func (s *SegyFile) getNrTraces() int64 {
	fi, err := s.file.Stat()
	if err != nil {
		fmt.Println("unable to get Stat()")
		log.Fatal(err)
	}
	size := fi.Size()
	nSamples := s.Header.Hns
	txtAndBinarySize := int64(SEGY_DESC_HDR_LEN + SEGY_BIN_HDR_LEN)
	nTraces := ((size - txtAndBinarySize) / (int64(SEGY_TRACE_HDR_LEN) + int64(nSamples) * int64(unsafe.Sizeof(float32(1)))))

	return nTraces
}

func (s *SegyFile) getNrSamples() int32 {
	return int32(s.Header.Hns)
}

func (s *SegyFile) ReadTrace() (Trace, error) {
	// First read the TraceHeader
	//data := []byte{}
	//data = append(data, 

	// Then figure out the size of the data, and read it
	return Trace{}, nil
}
