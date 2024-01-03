package model

type LatestCompletedVisit struct {
	VisitCompletedDate  string
	VisitReportMarkedAs RefData
	VisitUrgency        RefData
	RagRatingLowerCase  string
}

type RAGRating struct {
	Name   string
	Colour string
}

func (ragColour LatestCompletedVisit) GetRAGRating() RAGRating {
	var rag RAGRating
	switch ragColour.RagRatingLowerCase {
	case "red":
		rag.Name = "High risk"
		rag.Colour = "red"
	case "amber":
		rag.Name = "Medium risk"
		rag.Colour = "orange"
	case "green":
		rag.Name = "Low risk"
		rag.Colour = "green"
	}
	return rag
}
