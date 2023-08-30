package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	humanize "github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
)

type ResourceUsed struct {
	Cputime  int    `xml:"cput"`
	Walltime int    `xml:"walltime"`
	Memory   Memory `xml:"mem"`
}

type ResourceRequested struct {
	Nodes    string `xml:"nodes"`
	Walltime int    `xml:"walltime"`
	Memory   Memory `xml:"mem"`
	Ncpus    int    `xml:"ncpus"`
}

type Memory uint64

func (mem *Memory) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	v, err := humanize.ParseBytes(s)
	if err != nil {
		return err
	}

	*mem = Memory(v)
	return nil
}

func (mem Memory) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(fmt.Sprintf("%d", mem), start)
}

type JobInfo struct {
	ID        JobId             `xml:"Job_Id"`
	User      string            `xml:"euser"`
	Group     string            `xml:"egroup"`
	Used      ResourceUsed      `xml:"resources_used"`
	Requested ResourceRequested `xml:"Resource_List"`
}

// custom job id XML parser
type JobId string

func (id *JobId) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	*id = JobId(strings.Split(s, ".")[0])

	return nil
}

func (id JobId) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(string(id), start)
}

type Jobs struct {
	JobInfo []JobInfo `xml:"Jobinfo"`
}

var removeNonUTF = func(r rune) rune {
	if r == utf8.RuneError {
		return -1
	}
	return r
}

var removeIllegalCharacter = func(r rune) rune {
	if !isInCharacterRange(r) {
		return -1
	}
	return r
}

var isInCharacterRange = func(r rune) (inrange bool) {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xD7FF ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}

// RemoveNonUTF8Bytes removes bytes that isn't UTF-8 encoded
func RemoveNonUTF8Bytes(data []byte) []byte {
	return bytes.Map(removeNonUTF, data)
}

// RemoveIllegalBytes removes bytes that contains illegal character
func RemoveIllegalBytes(data []byte) []byte {
	return bytes.Map(removeIllegalCharacter, data)
}

func ParseJobLog(logfile string) (*Jobs, error) {

	data, err := os.ReadFile(logfile)
	if err != nil {
		return nil, err
	}

	data = RemoveIllegalBytes(data)

	xmldata := append(
		append([]byte("<Jobs>"), data...),
		[]byte("</Jobs>")...,
	)

	jobs := &Jobs{}

	err = xml.Unmarshal(xmldata, &jobs)

	if err != nil {
		return nil, err
	}

	log.Debugf("[%s] total number of jobs: %d\n", logfile, len(jobs.JobInfo))

	return jobs, nil
}
