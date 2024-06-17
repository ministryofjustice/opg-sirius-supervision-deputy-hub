package model

type LatestCompletedVisit struct {
	VisitCompletedDate  string
	VisitReportMarkedAs RAGRating
	VisitUrgency        RefData
	RagRatingLowerCase  string
}
