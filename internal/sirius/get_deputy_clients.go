package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type apiOrder struct {
	OrderStatus struct {
		Label string `json:"label"`
	}
	LatestSupervisionLevel struct {
		SupervisionLevel struct {
			Label string `json:"label"`
		}
	}
	OrderDate string `json:"orderDate"`
}

type apiOrders []apiOrder

type apiReport struct {
	DueDate        string `json:"dueDate"`
	RevisedDueDate string `json:"revisedDueDate"`
	Status         struct {
		Label string `json:"label"`
	} `json:"status"`
}

type reportReturned struct {
	DueDate        string
	RevisedDueDate string
	StatusLabel    string
}

type apiLatestCompletedVisit struct {
	VisitCompletedDate  string
	VisitReportMarkedAs struct {
		Label string `json:"label"`
	} `json:"visitReportMarkedAs"`
	VisitUrgency struct {
		Label string `json:"label"`
	} `json:"visitUrgency"`
}

type latestCompletedVisit struct {
	VisitCompletedDate  string
	VisitReportMarkedAs string
	VisitUrgency        string
	RagRatingLowerCase  string
}

type apiClient struct {
	ClientId            int    `json:"id"`
	Firstname           string `json:"firstname"`
	Surname             string `json:"surname"`
	CourtRef            string `json:"caseRecNumber"`
	RiskScore           int    `json:"riskScore"`
	ClientAccommodation struct {
		Label string `json:"label"`
	}
	Orders               apiOrders               `json:"orders"`
	OldestReport         apiReport               `json:"oldestNonLodgedAnnualReport"`
	LatestCompletedVisit apiLatestCompletedVisit `json:"latestCompletedVisit"`
	HasActiveREMWarning  bool                    `json:"HasActiveREMWarning"`
}

type Order struct {
	OrderStatus      string
	SupervisionLevel string
	OrderDate        time.Time
}

type Orders []Order

type DeputyClient struct {
	ClientId             int
	Firstname            string
	Surname              string
	CourtRef             string
	RiskScore            int
	AccommodationType    string
	OrderStatus          string
	SupervisionLevel     string
	OldestReport         reportReturned
	LatestCompletedVisit latestCompletedVisit
	HasActiveREMWarning  bool
}

type DeputyClientDetails []DeputyClient

type AriaSorting struct {
	SurnameAriaSort   string
	ReportDueAriaSort string
	CRECAriaSort      string
}

type Page struct {
	PageCurrent int `json:"current"`
	PageTotal   int `json:"total"`
}

type Metadata struct {
	TotalActiveClients int `json:"totalActiveClients"`
}

type ApiClientList struct {
	Clients      []apiClient `json:"clients"`
	Pages        Page        `json:"pages"`
	Metadata     Metadata    `json:"metadata"`
	TotalClients int         `json:"total"`
}

type ClientList struct {
	Clients      DeputyClientDetails
	Pages        Page
	TotalClients int
	Metadata     Metadata
}

type ClientListParams struct {
	DeputyId           int
	Limit              int
	Search             int
	DeputyType         string
	ColumnBeingSorted  string
	SortOrder          string
	OrderStatuses      []string
	AccommodationTypes []string
}

func (c *Client) GetDeputyClients(ctx Context, params ClientListParams) (ClientList, AriaSorting, error) {
	var clientList ClientList
	var apiClientList ApiClientList

	url := fmt.Sprintf("/api/v1/deputies/%s/%d/clients?&limit=%d&page=%d", strings.ToLower(params.DeputyType), params.DeputyId, params.Limit, params.Search)

	filter := params.CreateFilter()
	if filter != "" {
		url = fmt.Sprintf("%s&filter=%s", url, filter)
	}

	req, err := c.newRequest(ctx, http.MethodGet, url, nil)

	if err != nil {
		return clientList, AriaSorting{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return clientList, AriaSorting{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return clientList, AriaSorting{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return clientList, AriaSorting{}, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&apiClientList); err != nil {
		return clientList, AriaSorting{}, err
	}

	var clients DeputyClientDetails
	for _, t := range apiClientList.Clients {
		orders := restructureOrders(t.Orders)
		if len(orders) > 0 {
			var client = DeputyClient{
				ClientId:          t.ClientId,
				Firstname:         t.Firstname,
				Surname:           t.Surname,
				CourtRef:          t.CourtRef,
				RiskScore:         t.RiskScore,
				AccommodationType: t.ClientAccommodation.Label,
				OrderStatus:       getOrderStatus(orders),
				SupervisionLevel:  getMostRecentSupervisionLevel(orders),
				OldestReport: reportReturned{
					t.OldestReport.DueDate,
					t.OldestReport.RevisedDueDate,
					t.OldestReport.Status.Label,
				},
				LatestCompletedVisit: latestCompletedVisit{
					FormatDateTime(IsoDateTimeZone, t.LatestCompletedVisit.VisitCompletedDate, SiriusDate),
					t.LatestCompletedVisit.VisitReportMarkedAs.Label,
					t.LatestCompletedVisit.VisitUrgency.Label,
					strings.ToLower(t.LatestCompletedVisit.VisitReportMarkedAs.Label),
				},
				HasActiveREMWarning: t.HasActiveREMWarning,
			}
			clients = append(clients, client)
		}
	}

	clientList.Clients = clients
	clientList.Pages = apiClientList.Pages
	clientList.TotalClients = apiClientList.TotalClients
	clientList.Metadata = apiClientList.Metadata

	var aria AriaSorting
	aria.SurnameAriaSort = changeSortButtonDirection(params.SortOrder, params.ColumnBeingSorted, "surname")
	aria.ReportDueAriaSort = changeSortButtonDirection(params.SortOrder, params.ColumnBeingSorted, "reportdue")
	aria.CRECAriaSort = changeSortButtonDirection(params.SortOrder, params.ColumnBeingSorted, "crec")

	switch params.ColumnBeingSorted {
	case "reportdue":
		reportDueScoreSort(clients, params.SortOrder)
	case "crec":
		crecScoreSort(clients, params.SortOrder)
	default:
		alphabeticalSort(clients, params.SortOrder)
	}

	return clientList, aria, err
}

func (p ClientListParams) CreateFilter() string {
	var filter string
	for _, s := range p.OrderStatuses {
		filter += "order-status:" + s + ","
	}
	return strings.TrimRight(filter, ",")
}

/*
		GetOrderStatus returns the status of the oldest active order for a client.
	  If there isn’t one, the status of the oldest order is returned.
*/
func getOrderStatus(orders Orders) string {
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].OrderDate.Before(orders[j].OrderDate)
	})

	for _, o := range orders {
		if o.OrderStatus == "Active" {
			return o.OrderStatus
		}
	}
	return orders[0].OrderStatus
}

func getMostRecentSupervisionLevel(orders Orders) string {
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].OrderDate.After(orders[j].OrderDate)
	})
	return orders[0].SupervisionLevel
}

func restructureOrders(apiOrders apiOrders) Orders {
	orders := make(Orders, len(apiOrders))

	for i, t := range apiOrders {
		// reformatting order date to yyyy-dd-mm
		reformattedDate := formatDate(t.OrderDate)

		var supervisionLevel string
		if t.LatestSupervisionLevel.SupervisionLevel.Label != "" {
			supervisionLevel = t.LatestSupervisionLevel.SupervisionLevel.Label
		} else {
			supervisionLevel = ""
		}

		orders[i] = Order{
			OrderStatus:      t.OrderStatus.Label,
			SupervisionLevel: supervisionLevel,
			OrderDate:        reformattedDate,
		}
	}

	updatedOrders := removeOpenStatusOrders(orders)
	return updatedOrders
}

func formatDate(dateString string) time.Time {
	dateTime, _ := time.Parse("02/01/2006", dateString)
	return dateTime
}

func removeOpenStatusOrders(orders Orders) Orders {
	/* An order is open when it's with the Allocations team,
	and so not yet supervised by the PA team */

	var updatedOrders Orders
	for _, o := range orders {
		if o.OrderStatus != "Open" {
			updatedOrders = append(updatedOrders, o)
		}
	}
	return updatedOrders
}

func alphabeticalSort(clients DeputyClientDetails, sortOrder string) DeputyClientDetails {
	if len(clients) > 1 {
		sort.Slice(clients, func(i, j int) bool {
			if sortOrder == "asc" {
				return clients[i].Surname < clients[j].Surname
			} else {
				return clients[i].Surname > clients[j].Surname
			}
		})
	}
	return clients
}

func crecScoreSort(clients DeputyClientDetails, sortOrder string) DeputyClientDetails {
	sort.Slice(clients, func(i, j int) bool {
		if sortOrder == "asc" {
			return clients[i].RiskScore < clients[j].RiskScore
		} else {
			return clients[i].RiskScore > clients[j].RiskScore
		}
	})
	return clients
}

func setDueDateForSort(dueDate, revisedDueDate string) string {
	if revisedDueDate != "" {
		return revisedDueDate
	} else if dueDate != "" {
		return dueDate
	} else {
		return "12/12/9999"
	}
}

func reportDueScoreSort(clients DeputyClientDetails, sortOrder string) DeputyClientDetails {
	sort.Slice(clients, func(i, j int) bool {
		x := setDueDateForSort(clients[i].OldestReport.DueDate, clients[i].OldestReport.RevisedDueDate)
		y := setDueDateForSort(clients[j].OldestReport.DueDate, clients[j].OldestReport.RevisedDueDate)
		dateTimeI := formatDate(x)
		dateTimeJ := formatDate(y)

		if sortOrder == "asc" {
			return dateTimeI.Before(dateTimeJ)
		} else {
			return dateTimeJ.Before(dateTimeI)
		}
	})
	return clients
}

func changeSortButtonDirection(sortOrder string, columnBeingSorted string, functionCalling string) string {
	if functionCalling == columnBeingSorted {
		if sortOrder == "asc" {
			return "ascending"
		} else if sortOrder == "desc" {
			return "descending"
		}
		return "none"
	} else {
		return "none"
	}

}

func (s AriaSorting) GetHTMLSortDirection(direction string) string {
	if direction == "ascending" {
		return "desc"
	}
	return "asc"
}
