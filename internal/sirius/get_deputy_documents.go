package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

//type ApiDocument struct {
//	Id               int    `json:"id"`
//	Name             string `json:"name"`
//	JobTitle         string `json:"jobTitle"`
//	Email            string `json:"email"`
//	PhoneNumber      string `json:"phoneNumber"`
//	OtherPhoneNumber string `json:"otherPhoneNumber"`
//	Notes            string `json:"notes"`
//	IsMainContact    bool   `json:"isMainContact"`
//	IsNamedDeputy    bool   `json:"isNamedDeputy"`
//}

func (c *Client) GetDeputyDocuments(ctx Context, deputyId int) (*[]model.Document, error) {
	var documentList model.Response

	deputyId = 81
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/persons/%d/documents", deputyId), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&documentList); err != nil {
		return nil, err
	}

	return &documentList.Documents, err
}
