package main

import "encoding/json"
import "io/ioutil"
import "testing"

func TestUnmarshalReport(t *testing.T) {
	var report containerVulnerabilityReport

	fileBytes, err := ioutil.ReadFile("./example.json")
	if err != nil {
		t.Fatal(err)
	}

	if err = json.Unmarshal(fileBytes, &report); err != nil {
		t.Fatal(err)
	}

	t.Logf("Report - %+v", report)

}
