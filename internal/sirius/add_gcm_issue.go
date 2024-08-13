package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

type CreateGcmIssue struct {
	ClientCaseRecNumber string        `json:"caseRecNumber"`
	GcmIssueType        model.RefData `json:"gcmIssueType"`
	Notes               string        `json:"notes"`
}

type GcmIssue struct {
	Id     int `json:"id"`
	Client struct {
		Id            int    `json:"id"`
		CaseRecNumber string `json:"caseRecNumber"`
		Firstname     string `json:"firstname"`
		Surname       string `json:"surname"`
	} `json:"client"`
	CreatedDate   string `json:"createdDate"`
	CreatedByUser struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"createdByUser"`
	Notes        string        `json:"notes"`
	GcmIssueType model.RefData `json:"gcmIssueType"`
}

//
//{
//	"id":1,
//	"client":{"id":66,"caseRecNumber":"48217682","firstname":"Hamster","surname":"Person"},
//	"createdDate":"13\/08\/2024",
//	"receivedDate":"08\/08\/2024",
//	"createdByUser":{"id":101,"name":"PROTeam1","phoneNumber":"12345678","displayName":"PROTeam1 User1","deleted":false,"email":"pro1@opgtest.com","firstname":"PROTeam1","surname":"User1","roles":["OPG User","Case Manager"],"locked":false,"suspended":false},
//	"notes":"v",
//	"GCMIssueType":{"handle":"MISSING_INFORMATION","label":"Missing information"}
//},
//{
//	"id":2,
//	"client":{"id":66,"caseRecNumber":"48217682","firstname":"Hamster","surname":"Person"},
//	"createdDate":"13\/08\/2024",
//	"receivedDate":"08\/08\/2024",
//	"createdByUser":{"id":101,"name":"PROTeam1","phoneNumber":"12345678","displayName":"PROTeam1 User1","deleted":false,"email":"pro1@opgtest.com","firstname":"PROTeam1","surname":"User1","roles":["OPG User","Case Manager"],"locked":false,"suspended":false},
//	"notes":"v",
//	"GCMIssueType":{"handle":"MISSING_INFORMATION","label":"Missing information"}
//},
//{
//	"id":3,
//	"client":{"id":66,"caseRecNumber":"48217682","firstname":"Hamster","surname":"Person"},"createdDate":"13\/08\/2024","receivedDate":"07\/09\/2024","createdByUser":{"id":101,"name":"PROTeam1","phoneNumber":"12345678","displayName":"PROTeam1 User1","deleted":false,"email":"pro1@opgtest.com","firstname":"PROTeam1","surname":"User1","roles":["OPG User","Case Manager"],"locked":false,"suspended":false},"notes":"kate","GCMIssueType":{"handle":"MISSING_INFORMATION","label":"Missing information"}},{"id":4,"client":{"id":66,"caseRecNumber":"48217682","firstname":"Hamster","surname":"Person"},"createdDate":"13\/08\/2024","receivedDate":"01\/08\/2024","createdByUser":{"id":101,"name":"PROTeam1","phoneNumber":"12345678","displayName":"PROTeam1 User1","deleted":false,"email":"pro1@opgtest.com","firstname":"PROTeam1","surname":"User1","roles":["OPG User","Case Manager"],"locked":false,"suspended":false},"notes":"new refdata","GCMIssueType":{"handle":"DEPUTY_FEES_INCORRECT","label":"Deputy fees incorrect"}},{"id":5,"client":{"id":66,"caseRecNumber":"48217682","firstname":"Hamster","surname":"Person"},"createdDate":"13\/08\/2024","receivedDate":"03\/08\/2024","createdByUser":{"id":101,"name":"PROTeam1","phoneNumber":"12345678","displayName":"PROTeam1 User1","deleted":false,"email":"pro1@opgtest.com","firstname":"PROTeam1","surname":"User1","roles":["OPG User","Case Manager"],"locked":false,"suspended":false},"notes":"tyr","GCMIssueType":{"handle":"MISSING_INFORMATION","label":"Missing information"}}

func (c *Client) AddGcmIssue(ctx Context, clientCaseRecNumber, notes string, gcmIssueType model.RefData, deputyId int) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(CreateGcmIssue{
		ClientCaseRecNumber: clientCaseRecNumber,
		GcmIssueType:        gcmIssueType,
		Notes:               notes,
	})

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/deputies/%d/gcm-issues", deputyId), &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	//io.Copy(os.Stdout, resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {

		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{Errors: v.ValidationErrors}
		}

		return newStatusError(resp)
	}
	return nil
}
