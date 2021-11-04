package smtp

import (
	"github.com/cookienyancloud/back/pkg/cache"
	"github.com/cookienyancloud/back/pkg/email"
	//"github.com/cookienyancloud/back/pkg/logger"
)


type emailInfo struct {
	Email     string            `json:"email"`
	Variables map[string]string `json:"variables"`
}

type Client struct {
	cache cache.Cache
}

func NewClient(cache cache.Cache) *Client {
	return &Client{cache}
}

func (c *Client) AddEmailToList(input email.AddEmailInput) error {
	//TODO:mail?
	return nil
}

func (c *Client) getToken() (string, error) {
	//TODO:mail?
	return "", nil
}

func (c *Client) authenticate() (string, error) {
	//TODO:mail?
	return "", nil
}
