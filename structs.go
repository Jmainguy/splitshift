package main

import (
	"encoding/xml"
)

// Shift : A shift of work
type Shift struct {
	XMLName   xml.Name `xml:"Shift"`
	AideName  string   `xml:"AideName,omitempty"`
	StartTime string   `xml:"StartTime,omitempty"`
	EndTime   string   `xml:"EndTime,omitempty"`
	Date      string   `xml:"Date,omitempty"`
	Amount    float64  `xml:"Amount,omitempty"`
	Hours     float64  `xml:"Hours,omitempty"`
	Rate      float64  `xml:"Rate,omitempty"`
}
