package sirius

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type CreateDocument struct {
	Type          string `json:"assuranceType"`
	RequestedDate string `json:"requestedDate"`
	RequestedBy   int    `json:"requestedBy"`
}

type TestFile struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

type AddDocumentRequest struct {
	Type          string   `json:"type"`
	CaseRecNumber string   `json:"caseRecNumber"`
	ParentUuid    string   `json:"parentUuid"`
	Metadata      string   `json:"metadata"`
	File          TestFile `json:"file"`
}

func (c *Client) AddDocument(ctx Context, file multipart.File, documentType string, direction string, date string, notes string) error {
	var body bytes.Buffer
	//w := multipart.NewWriter(&body)

	testFile := TestFile{
		Name:   "api.json",
		Source: "eyJjbGllbnRTdGF0dXMiOlt7ImhhbmRsZSI6Ik9QRU4iLCJsYWJlbCI6Ik9wZW4ifSx7ImhhbmRsZSI6IkFDVElWRSIsImxhYmVsIjoiQWN0aXZlIn0seyJoYW5kbGUiOiJJTkFDVElWRSIsImxhYmVsIjoiSW5hY3RpdmUifSx7ImhhbmRsZSI6IkRFQ0VBU0VEIiwibGFiZWwiOiJEZWNlYXNlZCJ9LHsiaGFuZGxlIjoiQ0xPU0VEIiwibGFiZWwiOiJDbG9zZWQifSx7ImhhbmRsZSI6IkRVUExJQ0FURSIsImxhYmVsIjoiRHVwbGljYXRlIn0seyJoYW5kbGUiOiJSRUdBSU5FRF9DQVBBQ0lUWSIsImxhYmVsIjoiUmVnYWluZWQgY2FwYWNpdHkifSx7ImhhbmRsZSI6IkRFQVRIX05PVElGSUVEIiwibGFiZWwiOiJEZWF0aCBub3RpZmllZCJ9LHsiaGFuZGxlIjoiREVBVEhfQ09ORklSTUVEIiwibGFiZWwiOiJEZWF0aCBjb25maXJtZWQifV0sImNsaWVudEFjY29tbW9kYXRpb24iOlt7ImhhbmRsZSI6IkNBUkVcL05VUlNJTkdcL1JFU0lERU5USUFMIEhPTUUgKFBSSVZBVEVcL0xBXC9SRUdJU1RFUkVEKSIsImxhYmVsIjoiQ2FyZVwvTnVyc2luZ1wvUmVzaWRlbnRpYWwgSG9tZSAoUHJpdmF0ZVwvTEFcL1JlZ2lzdGVyZWQpIn0seyJoYW5kbGUiOiJDT1VOQ0lMIFJFTlRFRCIsImxhYmVsIjoiQ291bmNpbCBSZW50ZWQifSx7ImhhbmRsZSI6IkZBTUlMWSBNRU1CRVJcL0ZSSUVORCdTIEhPTUUiLCJsYWJlbCI6IkZhbWlseSBNZW1iZXJcL0ZyaWVuZCdzIEhvbWUgKGluY2x1ZGluZyBzcG91c2VcL2NpdmlsIHBhcnRuZXIpIn0seyJoYW5kbGUiOiJIRUFMVEggU0VSVklDRSBQQVRJRU5UIiwibGFiZWwiOiJIZWFsdGggU2VydmljZSBQYXRpZW50In0seyJoYW5kbGUiOiJIT1NURUwiLCJsYWJlbCI6Ikhvc3RlbCJ9LHsiaGFuZGxlIjoiSE9URUwiLCJsYWJlbCI6IkhvdGVsIn0seyJoYW5kbGUiOiJIT1VTRSBBU1NPQ0lBVElPTiIsImxhYmVsIjoiSG91c2UgQXNzb2NpYXRpb24ifSx7ImhhbmRsZSI6IkhPVVNJTkcgQVNTT0NJQVRJT04iLCJsYWJlbCI6IkhvdXNpbmcgQXNzb2NpYXRpb24ifSx7ImhhbmRsZSI6IkxBIE5VUlNJTkcgSE9NRSIsImxhYmVsIjoiTEEgTnVyc2luZyBIb21lIn0seyJoYW5kbGUiOiJMQSBOVVJTSU5HIEhPTUUgT1IgUkVTSURFTlRJQUwgSE9NRSIsImxhYmVsIjoiTEEgTnVyc2luZyBIb21lIG9yIFJlc2lkZW50aWFsIEhvbWUifSx7ImhhbmRsZSI6IlNIQVJFRCBMSVZFUyBDQVJFUiIsImxhYmVsIjoiTGl2aW5nIHdpdGggc2hhcmVkIGxpdmVzIGNhcmVyIn0seyJoYW5kbGUiOiJMQSBQQVJUIDMgQUNDT01NT0RBVElPTiIsImxhYmVsIjoiTG9jYWwgQXV0aG9yaXR5IFBhcnQgMyBBY2NvbW1vZGF0aW9uIn0seyJoYW5kbGUiOiJOSFMgQUNDT01NT0RBVElPTiIsImxhYmVsIjoiTkhTIEFjY29tbW9kYXRpb24gZS5nLiBob3NwaXRhbCBvciBob3N0ZWwifSx7ImhhbmRsZSI6Ik5PIEFDQ09NTU9EQVRJT04gVFlQRSIsImxhYmVsIjoiTm8gQWNjb21tb2RhdGlvbiBUeXBlIn0seyJoYW5kbGUiOiJOTyBGSVhFRCBBRERSRVNTIiwibGFiZWwiOiJObyBGaXhlZCBBZGRyZXNzIn0seyJoYW5kbGUiOiJPVEhFUiIsImxhYmVsIjoiT3RoZXIifSx7ImhhbmRsZSI6Ik9XTiBIT01FIiwibGFiZWwiOiJPd24gSG9tZSJ9LHsiaGFuZGxlIjoiUFJJVkFURSBIT1NQSVRBTCIsImxhYmVsIjoiUHJpdmF0ZSBIb3NwaXRhbCJ9LHsiaGFuZGxlIjoiUFJJVkFURSBOVVJTSU5HIEhPTUUiLCJsYWJlbCI6IlByaXZhdGUgTnVyc2luZyBIb21lIn0seyJoYW5kbGUiOiJQUklWQVRFIFJFTlRFRCIsImxhYmVsIjoiUHJpdmF0ZSBSZW50ZWQgKGkuZS4gTm90IENvdW5jaWwpIn0seyJoYW5kbGUiOiJQUklWQVRFIFJFU0lERU5USUFMIEhPTUUiLCJsYWJlbCI6IlByaXZhdGUgUmVzaWRlbnRpYWwgSG9tZSJ9LHsiaGFuZGxlIjoiUkVHSVNURVJFRCBDQVJFIEhPTUUiLCJsYWJlbCI6IlJlZ2lzdGVyZWQgQ2FyZSBIb21lIn0seyJoYW5kbGUiOiJSRVNJREVOVElBTCBFRFVDQVRJT04iLCJsYWJlbCI6IlJlc2lkZW50aWFsIEVkdWNhdGlvbiJ9LHsiaGFuZGxlIjoiU0VDVVJFIEhPU1BJVEFMIiwibGFiZWwiOiJTZWN1cmUgSG9zcGl0YWwifSx7ImhhbmRsZSI6IlNVUEVSVklTRUQgU0hFTFRFUkVEIiwibGFiZWwiOiJTdXBlcnZpc2VkIFNoZWx0ZXJlZCBBY2NvbW1vZGF0aW9uIn0seyJoYW5kbGUiOiJTVVBQT1JURUQgSE9VU0lORyIsImxhYmVsIjoiU3VwcG9ydGVkIEhvdXNpbmcifSx7ImhhbmRsZSI6IlNVUFBPUlRFRCBMSVZJTkciLCJsYWJlbCI6IlN1cHBvcnRlZCBMaXZpbmcifV19",
		Type:   "application/json",
	}

	requestBody := AddDocumentRequest{
		Type:          "Case forum",
		CaseRecNumber: "15674795",
		ParentUuid:    "",
		Metadata:      "12345678",
		File:          testFile,
	}

	//file, err := os.Open("./temp-files/" + filename)
	//if err != nil {
	//	panic(err)
	//}

	err := json.NewEncoder(&body).Encode(requestBody)

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, "/api/public/v1/documents", &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, resp.Body)

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{Errors: v.ValidationErrors}
		}

		return newStatusError(resp)
	}

	return nil
}
