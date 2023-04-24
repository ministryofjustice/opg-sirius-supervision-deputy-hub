package sirius

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type AssuranceVisit struct {
	Id                  int                 `json:"id"`
	AssuranceType       AssuranceTypes      `json:"assuranceType"`
	RequestedDate       time.Time           `json:"requestedDate"`
	RequestedBy         User                `json:"requestedBy"`
	CommissionedDate    time.Time           `json:"commissionedDate"`
	ReportDueDate       time.Time           `json:"reportDueDate"`
	ReportReceivedDate  time.Time           `json:"reportReceivedDate"`
	VisitOutcome        VisitOutcomeTypes   `json:"assuranceVisitOutcome"`
	PdrOutcome          PdrOutcomeTypes     `json:"pdrOutcome"`
	ReportReviewDate    time.Time           `json:"reportReviewDate"`
	VisitReportMarkedAs VisitRagRatingTypes `json:"assuranceVisitReportMarkedAs"`
	Note                string              `json:"note"`
	VisitorAllocated    string              `json:"visitorAllocated"`
	ReviewedBy          User                `json:"reviewedBy"`
	NullDateForFrontEnd time.Time
}

func (c *Client) GetAssuranceVisitById(ctx Context, deputyId int, visitId int) (AssuranceVisit, error) {
	var v AssuranceVisit

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/assurance-visits/%d", deputyId, visitId), nil)

	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	io.Copy(os.Stdout, resp.Body)

	err = json.NewDecoder(resp.Body).Decode(&v)

	AssuranceVisitFormatted := formatAssuranceVisit(v)

	fmt.Println("sirius date")
	fmt.Println(v.ReportReviewDate)
	fmt.Println(v.ReportDueDate)
	fmt.Println(v.RequestedDate)
	fmt.Println(v.CommissionedDate)
	fmt.Println(v.ReportReceivedDate)

	fmt.Println("parse 1")
	fmt.Println(time.Now().UTC())
	fmt.Println(v.ReportDueDate.UTC())
	//fmt.Println(v.RequestedDate.Local())
	//fmt.Println(v.CommissionedDate.Location())
	fmt.Println(v.ReportReceivedDate)

	return AssuranceVisitFormatted, err
}

func formatAssuranceVisit(v AssuranceVisit) AssuranceVisit {
	nullDate, _ := time.Parse("2006-01-02T15:04:05+00:00", "0001-01-01 00:00:00 +0000 UTC")

	updatedVisit := AssuranceVisit{
		AssuranceType:       v.AssuranceType,
		RequestedDate:       v.RequestedDate,
		Id:                  v.Id,
		RequestedBy:         v.RequestedBy,
		CommissionedDate:    v.CommissionedDate,
		ReportDueDate:       v.ReportDueDate,
		ReportReceivedDate:  v.ReportReceivedDate,
		ReportReviewDate:    v.ReportReviewDate,
		VisitOutcome:        v.VisitOutcome,
		PdrOutcome:          v.PdrOutcome,
		VisitReportMarkedAs: v.VisitReportMarkedAs,
		Note:                v.Note,
		VisitorAllocated:    v.VisitorAllocated,
		ReviewedBy:          v.ReviewedBy,
		NullDateForFrontEnd: nullDate,
	}

	return updatedVisit
}
