package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"sort"
	"strings"
)

type Order struct {
	OrderStatus            model.RefData `json:"orderStatus"`
	LatestSupervisionLevel latestSupervisionLevel
	OrderDate              string `json:"orderDate"`
	CaseSubType            string `json:"casesubtype"`
}

type latestSupervisionLevel struct {
	AppliesFrom      string        `json:"appliesFrom"`
	SupervisionLevel model.RefData `json:"supervisionLevel"`
}

type Report struct {
	DueDate        string        `json:"dueDate"`
	RevisedDueDate string        `json:"revisedDueDate"`
	Status         model.RefData `json:"status"`
}

type DeputyClient struct {
	ClientId               int                        `json:"id"`
	Firstname              string                     `json:"firstname"`
	Surname                string                     `json:"surname"`
	CourtRef               string                     `json:"caseRecNumber"`
	RiskScore              int                        `json:"riskScore"`
	ClientAccommodation    model.RefData              `json:"clientAccommodation"`
	Orders                 []Order                    `json:"orders"`
	OldestReport           Report                     `json:"oldestNonLodgedAnnualReport"`
	LatestCompletedVisit   model.LatestCompletedVisit `json:"latestCompletedVisit"`
	HasActiveREMWarning    bool                       `json:"HasActiveREMWarning"`
	SupervisionLevel       string
	OrderStatus            string
	ActivePfaOrderMadeDate string
	HasActiveHWOrder       bool
}

type Page struct {
	PageCurrent int `json:"current"`
	PageTotal   int `json:"total"`
}

type Metadata struct {
	TotalActiveClients int `json:"totalActiveClients"`
}

type ClientList struct {
	Clients      []DeputyClient
	Pages        Page
	TotalClients int
	Metadata     Metadata
}

type ClientListParams struct {
	DeputyId           int
	Limit              int
	Search             int
	DeputyType         string
	Sort               string
	OrderStatuses      []string
	AccommodationTypes []string
	SupervisionLevels  []string
}

func (c *Client) GetDeputyClients(ctx Context, params ClientListParams) (ClientList, error) {
	var clientList ClientList

	url := fmt.Sprintf("/api/v1/deputies/%s/%d/clients?&limit=%d&page=%d&sort=%s", strings.ToLower(params.DeputyType), params.DeputyId, params.Limit, params.Search, params.Sort)

	filter := params.CreateFilter()

	if filter != "" {
		url = fmt.Sprintf("%s&filter=%s", url, filter)
	}

	req, err := c.newRequest(ctx, http.MethodGet, url, nil)

	if err != nil {
		return clientList, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		return clientList, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return clientList, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return clientList, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&clientList); err != nil {
		return clientList, err
	}

	var clients []DeputyClient

	for _, t := range clientList.Clients {
		var client = DeputyClient{
			ClientId:               t.ClientId,
			Firstname:              t.Firstname,
			Surname:                t.Surname,
			CourtRef:               t.CourtRef,
			RiskScore:              t.RiskScore,
			ClientAccommodation:    t.ClientAccommodation,
			OrderStatus:            getOrderStatus(t.Orders),
			ActivePfaOrderMadeDate: getActivePfaOrderMadeDate(t.Orders),
			HasActiveHWOrder:       hasHwOrder(t.Orders),
			SupervisionLevel:       getMostRecentSupervisionLevel(t.Orders),
			OldestReport:           t.OldestReport,
			LatestCompletedVisit: model.LatestCompletedVisit{
				VisitCompletedDate:  FormatDateTime(IsoDateTimeZone, t.LatestCompletedVisit.VisitCompletedDate, SiriusDate),
				VisitReportMarkedAs: t.LatestCompletedVisit.VisitReportMarkedAs,
				VisitUrgency:        t.LatestCompletedVisit.VisitUrgency,
				RagRatingLowerCase:  strings.ToLower(t.LatestCompletedVisit.VisitReportMarkedAs.Label),
			},
			HasActiveREMWarning: t.HasActiveREMWarning,
		}
		clients = append(clients, client)
	}

	clientList.Clients = clients
	clientList.TotalClients = clientList.Metadata.TotalActiveClients

	return clientList, err
}

func (p ClientListParams) CreateFilter() string {
	var filter string
	for _, s := range p.OrderStatuses {
		filter += "order-status:" + s + ","
	}
	for _, k := range p.AccommodationTypes {
		filter += "accommodation:" + strings.Replace(k, " ", "%20", -1) + ","
	}
	for _, s := range p.SupervisionLevels {
		filter += "supervision-level:" + s + ","
	}
	return strings.TrimRight(filter, ",")
}

/*
GetOrderStatus returns the status of the oldest active order for a client and when that order was mdae.

	If there isn’t one, the status of the oldest order is returned.
*/
func getOrderStatus(orders []Order) string {
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].OrderDate == "" {
			orders[i].OrderDate = "31/12/9999"
		}

		if orders[j].OrderDate == "" {
			orders[j].OrderDate = "31/12/9999"
		}

		iDate := model.NewDate(orders[i].OrderDate)
		jDate := model.NewDate(orders[j].OrderDate)

		return iDate.Before(jDate)
	})

	for _, o := range orders {
		if o.OrderStatus.Label == "Active" {
			return o.OrderStatus.Label
		}
	}

	for _, o := range orders {
		if o.OrderStatus.Label != "Open" {
			return o.OrderStatus.Label
		}
	}
	return orders[0].OrderStatus.Label
}

func getActivePfaOrderMadeDate(orders []Order) string {
	for _, o := range orders {
		if o.CaseSubType == "pfa" && o.OrderStatus.Label == "Active" {
			return o.OrderDate
		}
	}
	return ""
}

func hasHwOrder(orders []Order) bool {
	for _, o := range orders {
		if o.CaseSubType == "hw" && o.OrderStatus.Label == "Active" {
			return true
		}
	}
	return false
}

func getMostRecentSupervisionLevel(orders []Order) string {
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].LatestSupervisionLevel.AppliesFrom == "" {
			orders[i].LatestSupervisionLevel.AppliesFrom = "01/01/0001"

		}
		if orders[j].LatestSupervisionLevel.AppliesFrom == "" {
			orders[j].LatestSupervisionLevel.AppliesFrom = "01/01/0001"
		}

		iDate := model.NewDate(orders[i].LatestSupervisionLevel.AppliesFrom)
		jDate := model.NewDate(orders[j].LatestSupervisionLevel.AppliesFrom)

		return iDate.After(jDate)
	})
	return orders[0].LatestSupervisionLevel.SupervisionLevel.Label
}
