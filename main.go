// accds_sample project main.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Patient struct {
	PatientId     string  `json:"patientid"`
	MRN           string  `json:"mrn"`
	PVM           string  `json:"pvm"`
	POA           bool    `json:"poa"`
	Secondary     bool    `json:"secondary"`
	ChangeDRGFrom int64   `json:"changedrgfrom"`
	ChangeDRGTo   int64   `json:"changedrgto"`
	ChangeRWFrom  float64 `json:"changerwfrom"`
	ChangeRWTo    float64 `json:"changerwto"`
	DOSFrom       string  `json:"dosfrom"`
	DOSTo         string  `json:"dosto"`
	Body          string  `json:"body"`
	PIDCount      int64   `json:"pidcount"`
	MRNCount      int64   `json:"mrncount"`
}

func (p Patient) CheckError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

func (p Patient) Process(sname string, tname string) bool {
	f, err := os.Open(sname)
	if err != nil {
		p.CheckError(err)
		return false
	}
	defer f.Close()

	p.PIDCount = 0
	p.MRNCount = 0
	scanner := bufio.NewScanner(f)
	linenum := 0
	pidprefix := ""
	mrnprefix := ""

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t\n\v'\f\r")
		conline := strings.Join(strings.Fields(line), "")

		if linenum == 0 {
			linenum++
			continue
		}
		linenum++

		if len(conline) == 0 {
			continue
		}

		found := false
		compiled := regexp.MustCompile("(?i)^PATIENTID[\t: ]+(1009[0-9]{5}|1008[0-9]{5})$")
		res := compiled.FindAllStringSubmatch(line, -1)
		if len(res) > 0 {
			//p.PatientId = res[0][1]
			if res[0][1][:4] == "1009" {
				pidprefix = "1009"
				mrnprefix = "1008"
			} else {
				pidprefix = "1008"
				mrnprefix = "1009"
			}
			p.PatientId = "AAAAAAAAA"
			p.PIDCount += 1
			found = true
		}

		compiled = regexp.MustCompile("(?i)^MRN[\t: ]+(1009[0-9]{5}|1008[0-9]{5})$")
		res = compiled.FindAllStringSubmatch(line, -1)
		if len(res) > 0 {
			//p.MRN = res[0][1]
			p.MRN = "AAAAAAAAA"
			p.MRNCount += 1
			found = true
		}

		compiled = regexp.MustCompile("(?i)PVM[\t: ]+(\\w\\d+)")
		res = compiled.FindAllStringSubmatch(line, -1)
		if len(res) > 0 {
			p.PVM = res[0][1]
			found = true
		}

		compiled = regexp.MustCompile("(?i)POA[\t: ]+(YES|NO)")
		res = compiled.FindAllStringSubmatch(line, -1)
		if len(res) > 0 {
			if strings.ToUpper(res[0][1]) == "YES" {
				p.POA = true
			} else {
				p.POA = false
			}
			found = true
		}

		compiled = regexp.MustCompile("(?i)SECONDARY[\t: ]+(YES|NO)")
		res = compiled.FindAllStringSubmatch(line, -1)
		if len(res) > 0 {
			if strings.ToUpper(res[0][1]) == "YES" {
				p.Secondary = true
			} else {
				p.Secondary = false
			}
			found = true
		}

		compiled = regexp.MustCompile("(?i)CHANGE[\t ]+DRG[\t: ]+FROM?[\t ]+(\\d+)[\t ]+TO[\t ]+(\\d+)")
		res = compiled.FindAllStringSubmatch(line, -1)
		if len(res) > 0 {
			n, err := strconv.Atoi(res[0][1])
			p.CheckError(err)
			p.ChangeDRGFrom = int64(n)
			n, err = strconv.Atoi(res[0][2])
			p.CheckError(err)
			p.ChangeDRGTo = int64(n)
			found = true
		}

		compiled = regexp.MustCompile("(?i)CHANGE[\t ]+RW[\t: ]+FROM?[\t ]+(\\d+\\.\\d+)[\t ]+TO[\t ]+(\\d+\\.\\d+)")
		res = compiled.FindAllStringSubmatch(line, -1)
		if len(res) > 0 {
			n, err := strconv.ParseFloat(res[0][1], 64)
			p.CheckError(err)
			p.ChangeRWFrom = float64(n)
			n, err = strconv.ParseFloat(res[0][2], 64)
			p.CheckError(err)
			p.ChangeRWTo = float64(n)
			found = true
		}

		compiled = regexp.MustCompile("(?i)DOS[\t: ]+(\\d+/\\d+)[\t ]+TO[\t ]+(\\d+/\\d+)")
		res = compiled.FindAllStringSubmatch(line, -1)

		if len(res) > 0 {
			p.DOSFrom = res[0][1]
			p.DOSTo = res[0][2]
			found = true
		}

		if !found {
			compiled := regexp.MustCompile(fmt.Sprintf("(%s[0-9]{5})", pidprefix))
			res := compiled.FindAllStringSubmatch(line, -1)
			if len(res) > 0 {
				p.PIDCount += int64(len(res[0]) - 1)
			}
			compiled = regexp.MustCompile(fmt.Sprintf("(%s[0-9]{5})", mrnprefix))
			res = compiled.FindAllStringSubmatch(line, -1)
			if len(res) > 0 {
				p.MRNCount += int64(len(res[0]) - 1)
			}
			compiled = regexp.MustCompile("(1009[0-9]{5}|1008[0-9]{5})")
			p.Body += compiled.ReplaceAllLiteralString(line, "AAAAAAAAA") + "\n"
		}

	}
	if p.CheckError(scanner.Err()) {
		return false
	}

	g, err := os.OpenFile(tname, os.O_RDWR|os.O_CREATE, 0644)
	if p.CheckError(err) {
		return false
	}
	defer g.Close()
	jbytes, err := json.MarshalIndent(p, "", "    ")
	g.Write(jbytes)
	return true
}

func main() {
	t := time.Now()
	var p Patient
	p.Process("patient.txt", "patient.json")

	fmt.Println(time.Since(t))
}
