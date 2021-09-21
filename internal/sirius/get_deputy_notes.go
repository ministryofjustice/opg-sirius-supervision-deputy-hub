package sirius

//func (c *Client) GetDeputyNotes(ctx Context, deputyId int) (error) {
//	//var v
//	//
//	//req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/timeline/%d", deputyId), nil)
//	//
//	//if err != nil {
//	//	return v, err
//	//}
//	//
//	//resp, err := c.http.Do(req)
//	//if err != nil {
//	//	return v, err
//	//}
//	//
//	//defer resp.Body.Close()
//	//
//	//if resp.StatusCode == http.StatusUnauthorized {
//	//	return v, ErrUnauthorized
//	//}
//	//
//	//if resp.StatusCode != http.StatusOK {
//	//	return v, newStatusError(resp)
//	//}
//	//err = json.NewDecoder(resp.Body).Decode(&v)
//
//	//DeputyEvents := EditDeputyEvents(v)
//
//	return err
//
//}