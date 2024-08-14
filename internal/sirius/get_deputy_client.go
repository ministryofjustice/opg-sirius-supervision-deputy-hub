package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ClientWithOrderDeputy struct {
	ClientId  int    `json:"id"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
	CourtRef  string `json:"caseRecNumber"`
	Cases     []struct {
		Deputies []struct {
			Deputy struct {
				Id int `json:"id"`
			} `json:"deputy"`
		} `json:"deputies"`
	} `json:"cases"`
}

func (c *Client) GetDeputyClient(ctx Context, caseRecNumber string, deputyId int) (ClientWithOrderDeputy, error) {
	var v ClientWithOrderDeputy

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/clients/caserec/%s", caseRecNumber), nil)

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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ClientWithOrderDeputy{}, ValidationError{Errors: v.ValidationErrors}
		}

		return ClientWithOrderDeputy{}, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, err
	}

	linked := checkIfClientLinkedToDeputy(v, deputyId)

	if !linked {
		validationErrors := ValidationErrors{
			"deputy": {
				"deputyClientLink": "Case number does not belong to this deputy",
			},
		}

		return ClientWithOrderDeputy{}, ValidationError{Errors: validationErrors}
	}
	return v, nil
}

func checkIfClientLinkedToDeputy(client ClientWithOrderDeputy, deputyId int) bool {
	for i := 0; i < len(client.Cases); {
		deputiesForOrder := client.Cases[i].Deputies
		for j := 0; j < len(deputiesForOrder); {
			if deputiesForOrder[j].Deputy.Id == deputyId {
				return true
			}
			j++
		}
		i++
	}
	return false
}
