package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
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
}

func (c *Client) GetGCMIssues(ctx Context, deputyId int, params GcmIssuesParams) ([]GcmIssue, error) {
	var v []GcmIssue

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/gcm-issues?&filter=%s&sort=%s", deputyId, params.IssueStatus, params.Sort), nil)

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
