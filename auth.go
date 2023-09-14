package api_client_go

func (c *Client) GetToken() (string, error) {
	return c.AuthClient.Token()
}
