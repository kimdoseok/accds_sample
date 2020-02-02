package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestProcess(t *testing.T) {
	var p Patient
	p.Process("patient.txt", "patient.json")
	jbytes, err := ioutil.ReadFile("patient.json")

	err = json.Unmarshal(jbytes, &p)
	//t.Log(p)
	if err != nil {
		t.Error("There is an error on unmarshalling.")
	}
	if p.PatientId != "AAAAAAAAA" {
		t.Error("Patient ID was not tokenized.")
	}
	if p.MRN != "AAAAAAAAA" {
		t.Error("MRN was not tokenized.")
	}
	if p.PVM != "E43" {
		t.Error("PVM has wrong value.")
	}
	if !p.POA {
		t.Error("POA should return true.")
	}
	if !p.Secondary {
		t.Error("Secondary should return true.")
	}
	if p.ChangeDRGFrom != 699 {
		t.Error("PVM has wrong value.")
	}
	if p.ChangeDRGTo != 698 {
		t.Error("PVM has wrong value.")
	}
	if p.ChangeRWFrom != 1.0327 {
		t.Error("PVM has wrong value.")
	}
	if p.ChangeRWTo != 1.6186 {
		t.Error("PVM has wrong value.")
	}
	if p.DOSFrom != "11/1" {
		t.Error("DOS FROM has wrong value.")
	}
	if p.DOSTo != "11/20" {
		t.Error("DOS TO has wrong value.")
	}
	if len(p.Body) != 598 {
		t.Error("The length of body sould be 598.")
	}
	if p.PIDCount != 2 {
		t.Error("PatientID counting is not 2.")
	}
	if p.MRNCount != 2 {
		t.Error("MRN counting counting is not 2.")
	}

}
