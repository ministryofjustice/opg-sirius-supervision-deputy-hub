package sirius

//func TestAddDocument(t *testing.T) {
//	mockClient := &mocks.MockClient{}
//	client, _ := NewClient(mockClient, "http://localhost:3000")
//
//	json := `{
//		"date": "14/06/2024",
//		"description": "<p>Note content</p>",
//		"direction": {
//			"handle": "INCOMING",
//			"label": "Incoming"
//		},
//		"name": "Test",
//		"type": {
//			"handle": "CASE_FORUM",
//			"label": "Case forum"
//		},
//		"personId": 68,
//		"fileName": "Screenshot 2024-06-19 at 11.16.27.png",
//		"file": {
//			"name": "Screenshot 2024-06-19 at 11.16.27.png",
//			"type": "image/png",
//			"source": "VBORw0KGgoAAAANSUhEUgAABg0AAAMOCA",
//		},
//		"fileSource" : "VBORw0KGgoAAAANSUhEUgAABg0AAAMOCA"
//	}`
//
//	r := io.NopCloser(bytes.NewReader([]byte(json)))
//
//	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
//		return &http.Response{
//			StatusCode: 201,
//			Body:       r,
//		}, nil
//	}
//
//	err := client.AddDocument(getContext(nil), "fake note title", "fake note text", 76, 51, "PA")
//	assert.Nil(t, err)
//}
