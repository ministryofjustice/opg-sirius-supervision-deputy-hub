package sirius
//every function has a main function

//MAP data structure

//create a map structure?

func (c *Client) MapTemplates(ctx Context) map[string]string {
	m := make(map[string]string)
	m["kate"] = "hi"
	m["mish"] = "hello"
	return m
}