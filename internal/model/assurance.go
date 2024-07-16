package model

type Assurance struct {
	Id                 int       `json:"id"`
	Type               RefData   `json:"assuranceType"`
	RequestedDate      string    `json:"requestedDate"`
	RequestedBy        User      `json:"requestedBy"`
	CommissionedDate   string    `json:"commissionedDate"`
	ReportDueDate      string    `json:"reportDueDate"`
	ReportReceivedDate string    `json:"reportReceivedDate"`
	VisitOutcome       RefData   `json:"assuranceVisitOutcome"`
	PdrOutcome         RefData   `json:"pdrOutcome"`
	ReportReviewDate   string    `json:"reportReviewDate"`
	ReportMarkedAs     RAGRating `json:"reportMarkedAs"`
	Note               string    `json:"note"`
	VisitorAllocated   string    `json:"visitorAllocated"`
	ReviewedBy         User      `json:"reviewedBy"`
	DeputyId           int
}
