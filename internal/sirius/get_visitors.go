package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Visitors []Visitor

type Visitor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetVisitors(ctx Context) (Visitors, error) {
	var v Visitors

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/visitors", nil)

	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("error")
		fmt.Println(v)
		fmt.Println("resp")
		fmt.Println(resp)
		return v, newStatusError(resp)
	}
	err = json.NewDecoder(resp.Body).Decode(&v)

	return v, err
}
