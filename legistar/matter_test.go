package legistar

import (
	"fmt"
	"testing"
)

func TestMatterText_SimplifiedText(t *testing.T) {
	type testCase struct {
		Plain, Expected string
	}
	tests := []testCase{
		{
			Plain: `
Res. No. 730
 
..Title
Resolution supporting the mission and growth of the Climate Museum.
..Body
 
By Council Members Gennaro, Hanif, Cab�n, Richardson Jordan, Nurse, Riley and Williams
 
	Whereas, According to the National Aeronautics and Space Administration, Earth&#39;s climate is changing at a rate not seen in the past 10,000 years: the global temperature is increasing, oceans are warming, sea levels are rising, ice and snow levels are decreasing, and the frequency of extreme weather events is increasing; and`,
			Expected: `Whereas, According to the National Aeronautics and Space Administration, Earth&#39;s climate is changing at a rate not seen in the past 10,000 years: the global temperature is increasing, oceans are warming, sea levels are rising, ice and snow levels are decreasing, and the frequency of extreme weather events is increasing; and`,
		},
	}
	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			mt := &MatterText{Plain: tc.Plain}
			if got := mt.SimplifiedText(); got != tc.Expected {
				t.Errorf("got %q expected %q", got, tc.Expected)
			}
		})

	}
}

func TestMatterText_SimplifiedRTF(t *testing.T) {
	type testCase struct {
		RTF, Expected string
	}
	tests := []testCase{
		{
			RTF:      "{\\rtf1\\fbidis\\ansi\\ansicpg1252\\deff0\\deflang1033\\deflangfe1033{\\fonttbl{\\f0\\froman\\fprq2\\fcharset0 Times New Roman;}}\r\n\\viewkind4\\uc1\\pard\\ltrpar\\hyphpar0\\qc\\f0\\fs24 Int. No. 1-A\\par\r\n\\pard\\ltrpar\\hyphpar0\\par\r\n\\pard\\ltrpar\\hyphpar0\\qj By Council Members A, B and C (by request of the D)\\par\r\n\\par\r\n\\v ..Title\\par\r\n\\v0 A Local Law to amend the New York city charter and the administrative code, in relation to x, y and z.\\par\r\n\\pard\\ltrpar\\hyphpar0\\v ..Body\\par\r\n\\v0\\par\r\n\\ul Be it enacted by the Council as follows:\\par\r\n",
			Expected: "{\\rtf1\\fbidis\\ansi\\ansicpg1252\\deff0\\deflang1033\\deflangfe1033{\\fonttbl{\\f0\\froman\\fprq2\\fcharset0 Times New Roman;}}\n\\viewkind4\\uc1\n\\v0\\par\n\\ul Be it enacted by the Council as follows:\\par\n",
		},
	}
	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			mt := &MatterText{RTF: tc.RTF}
			if got := mt.SimplifiedRTF(); got != tc.Expected {
				t.Errorf("got %q expected %q", got, tc.Expected)
			}
		})

	}
}
