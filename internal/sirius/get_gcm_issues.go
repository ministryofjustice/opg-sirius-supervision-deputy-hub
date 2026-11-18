package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

type GcmClient struct {
	Id            int    `json:"id"`
	CaseRecNumber string `json:"caseRecNumber"`
	Firstname     string `json:"firstname"`
	Surname       string `json:"surname"`
}

type UserInformation struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type GcmIssue struct {
	Id            int             `json:"id"`
	Client        GcmClient       `json:"client"`
	CreatedDate   string          `json:"createdDate"`
	CreatedByUser UserInformation `json:"createdByUser"`
	Notes         string          `json:"notes"`
	GcmIssueType  model.RefData   `json:"gcmIssueType"`
	ClosedOnDate  string          `json:"closedOnDate"`
	ClosedByUser  UserInformation `json:"closedByUser"`
}

type GcmIssuesParams struct {
	IssueStatus string
	Sort        string
	Limit       int
	Page        int
}

type GcmIssueList struct {
	GcmIssues []GcmIssue `json:"gcmIssues"`
	Pages     Page       `json:"pages"`
	Total     int        `json:"total"`
	Limit     int        `json:"limit"`
}

func (c *Client) GetGCMIssues(ctx Context, deputyId int, params GcmIssuesParams) (GcmIssueList, error) {
	var v GcmIssueList

	url := fmt.Sprintf(
		SupervisionAPIPath+"/v1/deputies/%d/gcm-issues?filter=%s&sort=%s&limit=%d&page=%d",
		deputyId, params.IssueStatus, params.Sort, params.Limit, params.Page,
	)

	req, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	return v, err
}
