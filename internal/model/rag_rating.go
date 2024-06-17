package model

import "strings"

type RAGRating RefData

type VisitReportMarkedAs struct {
	Name   string
	Colour string
}

func (ragRating RAGRating) GetRiskMarker() VisitReportMarkedAs {
	var markedAs VisitReportMarkedAs
	switch strings.ToLower(ragRating.Handle) {
	case "red":
		markedAs.Name = "High risk"
		markedAs.Colour = "red"
	case "amber":
		markedAs.Name = "Medium risk"
		markedAs.Colour = "orange"
	case "green":
		markedAs.Name = "Low risk"
		markedAs.Colour = "green"
	}
	return markedAs
}
