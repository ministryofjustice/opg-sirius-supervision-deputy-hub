package sirius

import (
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

func SetUpGetNextPageNumber(pageCurrent int, pageTotal int, totalTasks int) ClientList {
	clientList := ClientList{
		Pages: Page{
			PageCurrent: pageCurrent,
			PageTotal:   pageTotal,
		},
		TotalClients: totalTasks,
	}
	return clientList
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

func TestGetShowingLowerLimitNumberAlwaysReturns0If0Tasks(t *testing.T) {
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

func MakeListOfPagesRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
