package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExecutiveCaseManager struct {
	EcmId     int    `json:"id"`
	EcmName   string `json:"displayName"`
	IsDefault bool
}

type ExecutiveCaseManagerOutgoing struct {
	EcmId int `json:"ecmId"`
}

type DeputyType struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type Firm struct {
	FirmName   string `json:"firmName"`
	FirmId     int    `json:"id"`
	FirmNumber int    `json:"firmNumber"`
}

type DeputySubType struct {
	SubType string `json:"handle"`
}

type HandleLabel struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type DeputyImportantInformation struct {
	Id                        int         `json:"id"`
	AnnualBillingInvoice      HandleLabel `json:"annualBillingInvoice"`
	APAD                      HandleLabel `json:"apad"`
	BankCharges               HandleLabel `json:"bankCharges"`
	Complaints                HandleLabel `json:"complaints"`
	IndependentVisitorCharges HandleLabel `json:"independentVisitorCharges"`
	MonthlySpreadsheet        HandleLabel `json:"monthlySpreadsheet"`
	PanelDeputy               bool        `json:"panelDeputy"`
	ReportSystem              HandleLabel `json:"reportSystemType"`
	OtherImportantInformation string      `json:"otherImportantInformation"`
}

type DeputyDetails struct {
	ID                               int                        `json:"id"`
	DeputyFirstName                  string                     `json:"firstname"`
	DeputySurname                    string                     `json:"surname"`
	DeputyCasrecId                   int                        `json:"deputyCasrecId"`
	DisplayName                      string                     `json:"displayName"`
	CanDelete                        bool                       `json:"canDelete"`
	DeputyNumber                     int                        `json:"deputyNumber"`
	DeputySubType                    DeputySubType              `json:"deputySubType"`
	DeputyStatus                     string                     `json:"deputyStatus"`
	DeputyImportantInformation       DeputyImportantInformation `json:"deputyImportantInformation"`
	OrganisationName                 string                     `json:"organisationName"`
	OrganisationTeamOrDepartmentName string                     `json:"organisationTeamOrDepartmentName"`
	Email                            string                     `json:"email"`
	PhoneNumber                      string                     `json:"phoneNumber"`
	AddressLine1                     string                     `json:"addressLine1"`
	AddressLine2                     string                     `json:"addressLine2"`
	AddressLine3                     string                     `json:"addressLine3"`
	Town                             string                     `json:"town"`
	County                           string                     `json:"county"`
	Postcode                         string                     `json:"postcode"`
	ExecutiveCaseManager             ExecutiveCaseManager       `json:"executiveCaseManager"`
	DeputyType                       DeputyType                 `json:"deputyType"`
	Firm                             Firm                       `json:"firm"`
}

func (c *Client) GetDeputyDetails(ctx Context, defaultPATeam int, defaultPROTeam int, deputyId int) (DeputyDetails, error) {
	var v DeputyDetails

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(SupervisionAPIPath + "/v1/deputies/%d", deputyId), nil)
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

	err = json.NewDecoder(resp.Body).Decode(&v)

	if v.DeputyType.Handle == "PA" && v.ExecutiveCaseManager.EcmId == 0 {
		v.ExecutiveCaseManager.EcmId = defaultPATeam
		v.ExecutiveCaseManager.EcmName = "Public Authority Deputy Team"
	}
	if v.ExecutiveCaseManager.EcmId == defaultPATeam {
		v.ExecutiveCaseManager.IsDefault = true
	}

	if v.DeputyType.Handle == "PRO" && v.ExecutiveCaseManager.EcmId == 0 {
		v.ExecutiveCaseManager.EcmId = defaultPROTeam
		v.ExecutiveCaseManager.EcmName = "Professional deputy team - New deputy order"
	}
	if v.ExecutiveCaseManager.EcmId == defaultPROTeam {
		v.ExecutiveCaseManager.IsDefault = true
	}

	return v, err
}
