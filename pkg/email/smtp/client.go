package smtp

import (
	//"bytes"
	//"encoding/json"
	//"errors"
	//"fmt"
	//"io/ioutil"
	//"net/http"
	"time"

	"github.com/cookienyancloud/back/pkg/cache"
	"github.com/cookienyancloud/back/pkg/email"
	//"github.com/cookienyancloud/back/pkg/logger"
)

const (
	cacheTTL = int64(time.Hour) // SendPulse access tokens are valid for 1 hour
)

type authResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type addToListRequest struct {
	Emails []emailInfo `json:"emails"`
}

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
