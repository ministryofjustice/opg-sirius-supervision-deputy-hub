package sirius

type addTask struct {
	Title    string `json:"name"`
	Note     string `json:"description"`
	UserId   int    `json:"createdById"`
	NoteType string `json:"noteType"`
}

func (c *Client) AddTask(ctx Context, deputyId int, taskType string, dueDate string, notes string) error {
	//var noteType = getNoteType(deputyType)
	//var body bytes.Buffer
	//err := json.NewEncoder(&body).Encode(addNoteRequest{
	//	Title:    title,
	//	Note:     note,
	//	UserId:   userId,
	//	NoteType: noteType,
	//})
	//if err != nil {
	//	return err
	//}
	//
	//req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/deputy/%d/notes", deputyId), &body)
	//
	//if err != nil {
	//	return err
	//}
	//req.Header.Set("Content-Type", "application/json")
	//
	//resp, err := c.http.Do(req)
	//if err != nil {
	//	return err
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode == http.StatusUnauthorized {
	//	return ErrUnauthorized
	//}
	//
	//if resp.StatusCode != http.StatusCreated {
	//	var v struct {
	//		ValidationErrors ValidationErrors `json:"validation_errors"`
	//	}
	//
	//	if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
	//		return ValidationError{Errors: v.ValidationErrors}
	//	}
	//
	//	return newStatusError(resp)
	//}

	return nil
}
