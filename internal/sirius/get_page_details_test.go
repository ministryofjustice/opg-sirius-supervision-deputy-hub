package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPreviousPageNumber(t *testing.T) {
	assert.Equal(t, GetPreviousPageNumber(0), 1)
	assert.Equal(t, GetPreviousPageNumber(1), 1)
	assert.Equal(t, GetPreviousPageNumber(2), 1)
	assert.Equal(t, GetPreviousPageNumber(3), 2)
	assert.Equal(t, GetPreviousPageNumber(5), 4)
}

func SetUpGetNextPageNumber(pageCurrent int, pageTotal int, totalClients int) ClientList {
	clientList := ClientList{
		Pages: Page{
			PageCurrent: pageCurrent,
			PageTotal:   pageTotal,
		},
		TotalClients: totalClients,
	}
	return clientList
}

func TestGetPageDetailsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	clientList := SetUpGetNextPageNumber(1, 1, 1)

	expectedResponse := PageDetails{
		ListOfPages:       []int{1},
		PreviousPage:      1,
		NextPage:          1,
		LimitedPagination: []int{1},
		FirstPage:         1,
		LastPage:          1,
		StoredClientLimit: 25,
		ShowingUpperLimit: 1,
		ShowingLowerLimit: 1,
		LastFilter:        "",
	}

	pageDetails := client.GetPageDetails(getContext(nil), clientList, 1, 25)

	assert.Equal(t, expectedResponse, pageDetails)
}

func TestGetPageDetailsReturnedWithNoClients(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	clientList := SetUpGetNextPageNumber(1, 0, 0)

	expectedResponse := PageDetails{
		ListOfPages:       []int(nil),
		PreviousPage:      1,
		NextPage:          0,
		LimitedPagination: []int{0},
		FirstPage:         0,
		LastPage:          0,
		StoredClientLimit: 25,
		ShowingUpperLimit: 0,
		ShowingLowerLimit: 0,
		LastFilter:        "",
	}

	pageDetails := client.GetPageDetails(getContext(nil), clientList, 1, 25)

	assert.Equal(t, expectedResponse, pageDetails)
}

func TestGetNextPageNumber(t *testing.T) {
	clientList := SetUpGetNextPageNumber(1, 5, 0)

	assert.Equal(t, GetNextPageNumber(clientList, 0), 2)
	assert.Equal(t, GetNextPageNumber(clientList, 2), 3)
	assert.Equal(t, GetNextPageNumber(clientList, 15), 5)
}

func TestGetShowingLowerLimitNumberAlwaysReturns1IfOnly1Page(t *testing.T) {
	clientList := SetUpGetNextPageNumber(1, 0, 13)

	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 25), 1)
	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 50), 1)
	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 100), 1)
}

func TestGetShowingLowerLimitNumberAlwaysReturns0If0Clients(t *testing.T) {
	clientList := SetUpGetNextPageNumber(1, 0, 0)

	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 25), 0)
	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 50), 0)
	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 100), 0)
}

func TestGetShowingLowerLimitNumberCanIncrementOnPages(t *testing.T) {
	clientList := SetUpGetNextPageNumber(2, 0, 100)

	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 25), 26)
	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 50), 51)
	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 100), 101)
}

func TestGetShowingLowerLimitNumberCanIncrementOnManyPages(t *testing.T) {
	clientList := SetUpGetNextPageNumber(5, 0, 5000)

	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 25), 101)
	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 50), 201)
	assert.Equal(t, GetShowingLowerLimitNumber(clientList, 100), 401)
}

func TestGetShowingUpperLimitNumberWillReturnTotalClientsIfOnFinalPage(t *testing.T) {
	clientList := SetUpGetNextPageNumber(1, 0, 10)

	assert.Equal(t, GetShowingUpperLimitNumber(clientList, 25), 10)
	assert.Equal(t, GetShowingUpperLimitNumber(clientList, 50), 10)
	assert.Equal(t, GetShowingUpperLimitNumber(clientList, 100), 10)
}

func TestGetShowingUpperLimitNumberWillIncrementOnManyPages(t *testing.T) {
	clientList := SetUpGetNextPageNumber(1, 0, 1000)

	assert.Equal(t, GetShowingUpperLimitNumber(clientList, 25), 25)
	assert.Equal(t, GetShowingUpperLimitNumber(clientList, 50), 50)
	assert.Equal(t, GetShowingUpperLimitNumber(clientList, 100), 100)
}

func TestGetPaginationLimitsWhenOnFirstPage(t *testing.T) {
	clientList := SetUpGetNextPageNumber(1, 0, 25)
	clientDetails := PageDetails{
		ListOfPages: []int{1, 2, 3, 4},
		LastPage:    4,
	}
	assert.Equal(t, GetPaginationLimits(clientList, clientDetails), []int{1, 2, 3})
}

func TestGetPaginationLimitsWhenOnSecondToLastPage(t *testing.T) {
	clientList := SetUpGetNextPageNumber(3, 0, 25)
	clientDetails := PageDetails{
		ListOfPages: []int{1, 2, 3, 4},
		LastPage:    4,
	}
	assert.Equal(t, GetPaginationLimits(clientList, clientDetails), []int{1, 2, 3, 4})
}

func TestGetPaginationLimitsWhenOnLastPage(t *testing.T) {
	clientList := SetUpGetNextPageNumber(4, 0, 25)
	clientDetails := PageDetails{
		ListOfPages: []int{1, 2, 3, 4},
		LastPage:    4,
	}
	assert.Equal(t, GetPaginationLimits(clientList, clientDetails), []int{2, 3, 4})
}
