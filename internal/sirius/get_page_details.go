package sirius

type PageDetails struct {
	ListOfPages       []int
	PreviousPage      int
	NextPage          int
	LimitedPagination []int
	FirstPage         int
	LastPage          int
	StoredClientLimit int
	ShowingUpperLimit int
	ShowingLowerLimit int
	LastFilter        string
}

func (c *Client) GetPageDetails(ctx Context, clientList ClientList, search int, displayClientLimit int) PageDetails {
	var k PageDetails

	PageDetails := k

	for i := 1; i < clientList.Pages.PageTotal+1; i++ {
		PageDetails.ListOfPages = append(PageDetails.ListOfPages, i)
	}

	PageDetails.PreviousPage = GetPreviousPageNumber(search)

	PageDetails.NextPage = GetNextPageNumber(clientList, search)

	PageDetails.StoredClientLimit = displayClientLimit

	PageDetails.ShowingUpperLimit = GetShowingUpperLimitNumber(clientList, displayClientLimit)

	PageDetails.ShowingLowerLimit = GetShowingLowerLimitNumber(clientList, displayClientLimit)

	if len(PageDetails.ListOfPages) != 0 {
		PageDetails.FirstPage = PageDetails.ListOfPages[0]
		PageDetails.LastPage = PageDetails.ListOfPages[len(PageDetails.ListOfPages)-1]
		PageDetails.LimitedPagination = GetPaginationLimits(clientList, PageDetails)
	} else {
		PageDetails.FirstPage = 0
		PageDetails.LastPage = 0
		PageDetails.LimitedPagination = []int{0}
	}

	return PageDetails
}

func GetPreviousPageNumber(search int) int {
	if search <= 1 {
		return 1
	} else {
		return search - 1
	}
}

func GetNextPageNumber(clientList ClientList, search int) int {
	if search < clientList.Pages.PageTotal {
		if search == 0 {
			return search + 2
		} else {
			return search + 1
		}
	} else {
		return clientList.Pages.PageTotal
	}
}

func GetShowingLowerLimitNumber(clientList ClientList, displayClientLimit int) int {
	if clientList.Pages.PageCurrent == 1 && clientList.TotalClients != 0 {
		return 1
	} else if clientList.Pages.PageCurrent == 1 && clientList.TotalClients == 0 {
		return 0
	} else {
		previousPageNumber := clientList.Pages.PageCurrent - 1
		return previousPageNumber*displayClientLimit + 1
	}
}

func GetShowingUpperLimitNumber(clientList ClientList, displayClientLimit int) int {
	if clientList.Pages.PageCurrent*displayClientLimit > clientList.TotalClients {
		return clientList.TotalClients
	} else {
		return clientList.Pages.PageCurrent * displayClientLimit
	}
}

func GetPaginationLimits(clientList ClientList, clientDetails PageDetails) []int {
	var twoBeforeCurrentPage int
	var twoAfterCurrentPage int
	if clientList.Pages.PageCurrent > 2 {
		twoBeforeCurrentPage = clientList.Pages.PageCurrent - 3
	} else {
		twoBeforeCurrentPage = 0
	}
	if clientList.Pages.PageCurrent+2 <= clientDetails.LastPage {
		twoAfterCurrentPage = clientList.Pages.PageCurrent + 2
	} else if clientList.Pages.PageCurrent+1 <= clientDetails.LastPage {
		twoAfterCurrentPage = clientList.Pages.PageCurrent + 1
	} else {
		twoAfterCurrentPage = clientList.Pages.PageCurrent
	}
	return clientDetails.ListOfPages[twoBeforeCurrentPage:twoAfterCurrentPage]
}
