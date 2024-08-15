package sirius

//
//func TestDeputyClientReturned(t *testing.T) {
//	mockClient := &mocks.MockClient{}
//	client, _ := NewClient(mockClient, "http://localhost:3000")
//
//	json := `{
//		"id":66,
//		"caseRecNumber":"43787324",
//		"firstname":"Hamster",
//		"surname":"Person",
//		"cases":[
//			{
//				"id":94,
//				"deputies":[
//					{
//						"deputy":
//							{
//								"id":67
//							}
//					}
//				]
//			},
//			{
//				"id":95,
//				"deputies":[
//					{
//						"deputy":
//							{
//								"id":67
//							}
//					}
//				]
//			}
//		]
//	}`
//
//	r := io.NopCloser(bytes.NewReader([]byte(json)))
//
//	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
//		return &http.Response{
//			StatusCode: 200,
//			Body:       r,
//		}, nil
//	}
//
//	expectedResponse := ClientWithOrderDeputy{
//		ClientId:  66,
//		Firstname: "Hamster",
//		Surname:   "Person",
//		CourtRef:  "43787324",
//		Cases: []Case{
//			{
//				Deputies: []Deputies{
//					{
//						Deputy: OrderDeputy{Id: 67},
//					},
//				},
//			},
//			{
//				Deputies: []Deputies{
//					{
//						Deputy: OrderDeputy{Id: 67},
//					},
//				},
//			},
//		},
//	}
//
//	expectedClient, err := client.GetDeputyClient(getContext(nil), "43787324", 67)
//
//	assert.Equal(t, expectedResponse, expectedClient)
//	assert.Equal(t, nil, err)
//}
//
//func TestDeputyClientReturnsErrorIfDeputyNotLinkedToClient(t *testing.T) {
//	mockClient := &mocks.MockClient{}
//	client, _ := NewClient(mockClient, "http://localhost:3000")
//
//	json := `{
//		"id":66,
//		"caseRecNumber":"43787324",
//		"firstname":"Hamster",
//		"surname":"Person",
//		"cases":[
//			{
//				"id":94,
//				"deputies":[
//					{
//						"deputy":
//							{
//								"id":67
//							}
//					}
//				]
//			}
//		]
//	}`
//
//	r := io.NopCloser(bytes.NewReader([]byte(json)))
//
//	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
//		return &http.Response{
//			StatusCode: 200,
//			Body:       r,
//		}, nil
//	}
//
//	expectedResponse := ValidationError(
//		ValidationError{
//			Message: "",
//			Errors: ValidationErrors{
//				"deputy": map[string]string{"deputyClientLink": "Case number does not belong to this deputy"},
//			},
//		},
//	)
//
//	actualClient, actualError := client.GetDeputyClient(getContext(nil), "43787324", 999)
//	assert.Equal(t, expectedResponse, actualError)
//	assert.Equal(t, ClientWithOrderDeputy{}, actualClient)
//
//}
//
//func TestGetDeputyClientReturnsNewStatusError(t *testing.T) {
//	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusMethodNotAllowed)
//	}))
//	defer svr.Close()
//
//	client, _ := NewClient(http.DefaultClient, svr.URL)
//
//	contact, err := client.GetDeputyClient(getContext(nil), "123456", 76)
//
//	expectedResponse := ClientWithOrderDeputy{}
//
//	assert.Equal(t, expectedResponse, contact)
//	assert.Equal(t, StatusError{
//		Code:   http.StatusMethodNotAllowed,
//		URL:    svr.URL + "/api/v1/clients/caserec/123456",
//		Method: http.MethodGet,
//	}, err)
//}
//
//func TestGetDeputyClientReturnsUnauthorisedClientError(t *testing.T) {
//	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusUnauthorized)
//	}))
//	defer svr.Close()
//
//	client, _ := NewClient(http.DefaultClient, svr.URL)
//
//	contact, err := client.GetDeputyClient(getContext(nil), "123456", 76)
//
//	expectedResponse := ClientWithOrderDeputy{}
//
//	assert.Equal(t, ErrUnauthorized, err)
//	assert.Equal(t, expectedResponse, contact)
//}
//
//func TestValidationErrorReturnedIfClientNotLinkedToDeputy(t *testing.T) {
//	testClient := ClientWithOrderDeputy{
//		ClientId:  123,
//		Firstname: "test",
//		Surname:   "Client",
//		CourtRef:  "12345",
//		Cases: []Case{
//			{
//				[]Deputies{
//					{
//						Deputy: OrderDeputy{
//							77,
//						},
//					},
//				},
//			},
//		},
//	}
//
//	assert.Equal(t, checkIfClientLinkedToDeputy(testClient, 77), true)
//	assert.Equal(t, checkIfClientLinkedToDeputy(testClient, 76), false)
//}
//
//func TestCheckIfClientLinkedToDeputy(t *testing.T) {
//	testClient := ClientWithOrderDeputy{
//		ClientId:  123,
//		Firstname: "test",
//		Surname:   "Client",
//		CourtRef:  "12345",
//		Cases: []Case{
//			{
//				[]Deputies{
//					{
//						Deputy: OrderDeputy{
//							77,
//						},
//					},
//				},
//			},
//		},
//	}
//
//	assert.Equal(t, checkIfClientLinkedToDeputy(testClient, 77), true)
//	assert.Equal(t, checkIfClientLinkedToDeputy(testClient, 76), false)
//}
