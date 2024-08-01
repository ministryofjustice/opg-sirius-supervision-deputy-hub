package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"io"
	"mime/multipart"
	"net/http"
)

type mockTemplate struct {
	executed         bool
	executedTemplate bool
	lastVars         interface{}
	lastW            io.Writer
	error            error
}

func (m *mockTemplate) Execute(w io.Writer, vars any) error {
	m.executed = true
	m.lastVars = vars
	m.lastW = w
	return m.error
}

func (m *mockTemplate) ExecuteTemplate(w io.Writer, name string, vars any) error {
	m.executedTemplate = true
	m.lastVars = vars
	m.lastW = w
	return m.error
}

type mockRoute struct {
	client   ApiClient
	data     any
	executed bool
	lastW    io.Writer
	error
}

func (r *mockRoute) Client() ApiClient {
	return r.client
}

func (r *mockRoute) execute(w http.ResponseWriter, req *http.Request, data any) error {
	r.executed = true
	r.lastW = w
	r.data = data
	return r.error
}

type mockApiClient struct {
	error                 error
	CurrentUserDetails    sirius.UserDetails
	DeputyDetails         sirius.DeputyDetails
	AssignDeputyToFirmErr error
	AddFirm               int
	AddFirmDetailsErr     error
}

func (m mockApiClient) GetUserDetails(sirius.Context) (sirius.UserDetails, error) {
	return m.CurrentUserDetails, m.error
}

func (m mockApiClient) GetDeputyDetails(sirius.Context, int, int, int) (sirius.DeputyDetails, error) {
	return m.DeputyDetails, m.error
}

func (m *mockApiClient) AddFirmDetails(ctx sirius.Context, deputyId sirius.FirmDetails) (int, error) {
	return m.AddFirm, m.AddFirmDetailsErr
}

func (m *mockApiClient) AssignDeputyToFirm(ctx sirius.Context, deputyId int, firmId int) error {
	return m.AssignDeputyToFirmErr
}

func (m mockApiClient) GetAssurances(ctx sirius.Context, deputyId int) ([]model.Assurance, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) AddAssurance(ctx sirius.Context, assuranceType string, requestedDate string, userId, deputyId int) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) UpdateAssurance(ctx sirius.Context, manageAssuranceForm sirius.UpdateAssuranceDetails, deputyId, visitId int) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetVisitors(ctx sirius.Context) ([]model.Visitor, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetRagRatingTypes(ctx sirius.Context) ([]model.RAGRating, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetVisitOutcomeTypes(ctx sirius.Context) ([]model.RefData, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetPdrOutcomeTypes(ctx sirius.Context) ([]model.RefData, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetAssuranceById(ctx sirius.Context, deputyId int, visitId int) (model.Assurance, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDeputyClients(context sirius.Context, params sirius.ClientListParams) (sirius.ClientList, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetAccommodationTypes(context sirius.Context) ([]model.RefData, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetSupervisionLevels(context sirius.Context) ([]model.RefData, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDeputyContacts(context sirius.Context, i int) (sirius.ContactList, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) UpdateDeputyContactDetails(context sirius.Context, i int, details sirius.DeputyContactDetails) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetContactById(ctx sirius.Context, deputyId int, contactId int) (sirius.Contact, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) DeleteContact(context sirius.Context, i int, i2 int) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) AddContact(context sirius.Context, i int, form sirius.ContactForm) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) UpdateContact(context sirius.Context, i int, i2 int, form sirius.ContactForm) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) DeleteDeputy(context sirius.Context, i int) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) EditDeputyTeamDetails(context sirius.Context, details sirius.DeputyDetails) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDeputyTeamMembers(context sirius.Context, i int, details sirius.DeputyDetails) ([]model.TeamMember, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) ChangeECM(context sirius.Context, outgoing sirius.ExecutiveCaseManagerOutgoing, details sirius.DeputyDetails) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDeputyDocuments(ctx sirius.Context, deputyId int, sort string) (sirius.DocumentList, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) AddDocument(ctx sirius.Context, file multipart.File, filename string, documentType string, direction string, date string, notes string, deputyId int) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDocumentDirections(ctx sirius.Context) ([]model.RefData, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDocumentTypes(ctx sirius.Context) ([]model.RefData, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) ReplaceDocument(ctx sirius.Context, file multipart.File, filename, documentType, direction, date, notes string, deputyId, documentId int) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDocumentById(ctx sirius.Context, deputyId, documentId int) (model.Document, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetFirms(context sirius.Context) ([]sirius.FirmForList, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDeputyNotes(context sirius.Context, i int) (sirius.DeputyNoteCollection, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) AddNote(ctx sirius.Context, title, note string, deputyId, userId int, deputyType string) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) UpdateImportantInformation(context sirius.Context, i int, details sirius.ImportantInformationDetails) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDeputyAnnualInvoiceBillingTypes(ctx sirius.Context) ([]model.RefData, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDeputyBooleanTypes(ctx sirius.Context) ([]model.RefData, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDeputyReportSystemTypes(ctx sirius.Context) ([]model.RefData, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetTaskTypesForDeputyType(ctx sirius.Context, deputyType string) ([]model.TaskType, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetTasks(ctx sirius.Context, deputyId int) (sirius.TaskList, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) AddTask(ctx sirius.Context, deputyId int, taskType string, typeName string, dueDate string, notes string, assigneeId int) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetTask(context sirius.Context, i int) (model.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) CompleteTask(context sirius.Context, i int, s string) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) UpdateTask(ctx sirius.Context, deputyId, taskId int, dueDate, notes string, assigneeId int) error {
	//TODO implement me
	panic("implement me")
}

func (m mockApiClient) GetDeputyEvents(context sirius.Context, i int) (sirius.DeputyEvents, error) {
	//TODO implement me
	panic("implement me")
}
