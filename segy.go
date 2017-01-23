package segygo

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	//"log"
	"os"
	"unsafe"
	"reflect"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("segygo")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

const Version = "0.1"
const SEGY_DESC_HDR_LEN = 3200
const SEGY_BIN_HDR_LEN = 400
const SEGY_TRACE_HDR_LEN = 240

type BinHeader struct {
	Jobid, Lino, Reno int32
	Ntrpr, Nart int16
	Hdt, Dto, Hns, Nso uint16
	Format, Fold, Tsort, Vscode, Hsfs, Hsfe, Hslen, Hstyp, Schn, Hstas, Hstae, Htatyp, Hcorr, Bgrcv, Rcvm, Mfeet, Polyt, Vgpol int16
	Hunass [170]int16 // unassigned
}

type TraceHeader struct {
	Tracel int32
	Tracer int32
	Fldr int32
	Tracf int32
	Ep int32
	CDP int32
	CDPT int32
	Trid int16
	Nvs int16
	Nhs int16
	Duse int16
	Offset int32
	Gelev int32
	Selev int32
	Sdepth int32
	Gdel int32
	Sdel int32
	SwDep int32
	GwDep int32
	Scalel int16
	Scalco int16
	Sx int32
	Sy int32
	Gx int32
	Gy int32
	CoUnit int16
	WeVel int16
	SweVel int16
	Sut, Gut, Sstat, Gstat, Tstat, Laga, Lagb, Delrt, Muts, Mute int16
	Ns, Dt uint16
	Gain, Igc, Igi, Corr, Sfs, Sfe, Slen, Styp, Stas, Stae, Tatyp int16
	Afilf, Afils, NoFilf, NoFils, Lcf, Hcf, Lcs, Hcs, Year, Day int16
	Hour, Minute, Sec, Timbas, Trwf, Grnors, Grnofr, Grnlof, Gaps, Otrav int16
	D1, F1, D2, F2, Ungpow, Unscale float32
	Ntr int32
	Mark, Shortpad int16
	Unass [14]int16 // unassigned short array
}

type Trace struct {
	TraceHeader
	Data []float32
}

type SegyFile struct {
	Filename string
	Header   BinHeader
	NrTraces int64
	file     *os.File
	Position int64
	LogLevel logging.Level
}

func CreateFile(filename string) (SegyFile, error) {
	var s SegyFile
	var binHdr BinHeader
	f, err := os.Create(filename)
	defer f.Close()

	if err != nil {
		return s, err
	}

	s.LogLevel = logging.WARNING

	// Setup proper logging
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1Formatter)
	logging.SetLevel(s.LogLevel, "")

	log.Debugf("Creating SEG-Y file: %s", s.Filename)

	s.Filename = filename
	s.Header = binHdr
	s.NrTraces = 0
	s.file = f
	s.Position = 0

	accum := make([]byte, 3200)
	//r := bytes.NewWriter(accum)
	//binary.Write()
	buff := bytes.NewBuffer(accum)
	if err = binary.Write(buff, binary.BigEndian, &s.Header); err != nil {
		log.Errorf("Error creating buffer to hold binary header for segy file: %s. Msg: %s", s.Filename, err)
		return s, err
	}

	n, err := f.Write(buff.Bytes())
	if err != nil {
		log.Errorf("Error writing binary header to segy file: %s. Msg: %s", s.Filename, err)
		return s, err
	}
	log.Debugf("Wrote %d bytes to file: %s", n, s.Filename)

	return s, err

}

func OpenFile(filename string) (SegyFile, error) {
	var s SegyFile
	var binHdr BinHeader

	s.Filename = filename
	s.LogLevel = logging.WARNING

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return s, err
	}

	// Setup proper logging
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1Formatter)
	logging.SetLevel(s.LogLevel, "")

	accum := []byte{}
	accum = append(accum, b...)

	accum2 := accum[3200:]
	r := bytes.NewReader(accum2)
	log.Debugf("Number of bytes: %d", r.Len())

	if err = binary.Read(r, binary.BigEndian, &binHdr); err != nil {
		log.Errorf("Error reading segy file (bigendian). %s", err)
		return s, err
	}

	// Open and store the os.File object in our struct
	file, err := os.Open(s.Filename)
	s.file = file
	defer file.Close()

	s.Header = binHdr
	s.NrTraces = s.GetNrTraces()

	return s, err
}

func (s *SegyFile) SetVerbose(verbose bool) {

	if verbose {
		s.LogLevel = logging.DEBUG
		logging.SetLevel(s.LogLevel, "")
	} else {
		s.LogLevel = logging.WARNING
		logging.SetLevel(s.LogLevel, "")
	}

}

func (s *SegyFile) GetNrTraces() int64 {
	fi, err := s.file.Stat()
	if err != nil {
		log.Warning("unable to get Stat()")
		log.Fatal(err)
	}
	size := fi.Size()
	nSamples := s.Header.Hns
	txtAndBinarySize := int64(SEGY_DESC_HDR_LEN + SEGY_BIN_HDR_LEN)
	nTraces := ((size - txtAndBinarySize) / (int64(SEGY_TRACE_HDR_LEN) + int64(nSamples)*int64(unsafe.Sizeof(float32(1)))))

	return nTraces
}

func (s *SegyFile) GetNrSamples() int32 {
	return int32(s.Header.Hns)
}

func (s *SegyFile) GetHeader() map[string]interface{} {
	m := make(map[string]interface{})
	v := reflect.ValueOf(s.Header)
	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Name
		val := v.Field(i).Interface()
		log.Debugf("name = %s, value = %d", key, val)
		m[key] = val
	}

	return m
}

func (s *SegyFile) ReadTrace() (Trace, error) {
	trace := Trace{}
	traceBuff := make([]float32, s.GetNrSamples())
	byteBuff := make([]byte, s.GetNrSamples()*4)
	trace.Data = traceBuff

	trcHdrBuff := make([]byte, SEGY_TRACE_HDR_LEN)
	log.Info("trying to read trc hdr")
	bytesRead, err := s.file.Read(trcHdrBuff)
	if err != nil {
		log.Fatal(err)
		return trace, err
	}

	trcHdrReader := bytes.NewReader(trcHdrBuff)
	err = binary.Read(trcHdrReader, binary.BigEndian, &trace.TraceHeader)
	if err != nil {
		log.Fatal(err)
		return trace, err
	}

	log.Info("Reading first trace data")
	bytesRead, err = s.file.Read(byteBuff)
	if err != nil {
		log.Fatal(err)
		return trace, err
	}

	if bytesRead == 0 {
		log.Infof("No bytes read for trace #", s.Position)
	}

	for i := range trace.Data {
		trace.Data[i] = float32(binary.BigEndian.Uint32(byteBuff[i*4 : (i+1)*4]))
	}

	// Then figure out the size of the data, and read it
	return trace, nil
}

//func (s *SegyFile) 
