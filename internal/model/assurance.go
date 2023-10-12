package model

type Assurance struct {
	Id                 int              `json:"id"`
	Type               AssuranceType    `json:"assuranceType"`
	RequestedDate      string           `json:"requestedDate"`
	RequestedBy        User             `json:"requestedBy"`
	CommissionedDate   string           `json:"commissionedDate"`
	ReportDueDate      string           `json:"reportDueDate"`
	ReportReceivedDate string           `json:"reportReceivedDate"`
	VisitOutcome       VisitOutcomeType `json:"assuranceVisitOutcome"`
	PdrOutcome         PdrOutcomeType   `json:"pdrOutcome"`
	ReportReviewDate   string           `json:"reportReviewDate"`
	ReportMarkedAs     RagRatingType    `json:"reportMarkedAs"`
	Note               string           `json:"note"`
	VisitorAllocated   string           `json:"visitorAllocated"`
	ReviewedBy         User             `json:"reviewedBy"`
	DeputyId           int
}

type AssuranceType struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type VisitOutcomeType struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type PdrOutcomeType struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type RagRatingType struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}
