package gomagnet

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
)

// ClientConfig object used for client creation
type ClientConfig struct {
	ClientID        string
	ClientSecret    string
	Username        string
	Password        string
	AccessToken     string
	RefreshToken    string
	TokenExpiration *time.Time
	TokenType       string
	Debug           bool
}

// NewClientConfig constructs a ClientConfig object with the environment variables set as default
func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Username:     username,
		Password:     password,
	}
}

// Client object
type Client struct {
	config *ClientConfig

	Transport http.RoundTripper
}

// New constructor from configuration
func New(cc *ClientConfig) *Client {
	return &Client{
		config: cc,
	}
}

// AuthRequest data
type AuthRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

// AuthResponse data
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// GetAccessToken exchanges secrets for an access_token.
func (c *Client) GetAccessToken() (AuthResponse, error) {
	var r AuthResponse

	var rd AuthRequest
	rd.GrantType = "password"
	rd.ClientID = c.config.ClientID
	rd.ClientSecret = c.config.ClientSecret
	rd.Username = c.config.Username
	rd.Password = c.config.Password
	data, err := json.Marshal(rd)
	if err != nil {
		return r, err
	}

	// Parse URL
	base, err := url.Parse(apiHost)
	if err != nil {
		return r, err
	}
	base.Path = path.Join(base.Path, "oauth", "token")
	uri := base.String()

	if c.config.Debug {
		log.Printf("Sending auth request to %s with payload: %s", uri, data)
	}

	// Create request
	authRequest, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	if err != nil {
		return r, err
	}
	authRequest.Header.Add("Content-Type", "application/json")

	// Perform request
	client := &http.Client{Transport: c.Transport}
	resp, err := client.Do(authRequest)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	if c.config.Debug {
		log.Printf("Got auth response: %s", body)
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return r, err
	}

	// Set current client credentials
	c.config.AccessToken = r.AccessToken
	c.config.RefreshToken = r.RefreshToken
	texp := time.Now().Add(time.Second * time.Duration(r.ExpiresIn))
	c.config.TokenExpiration = &texp
	c.config.TokenType = r.TokenType

	return r, nil
}

// Request executes any MagnetChat API method using the current client configuration
func (c *Client) Request(method, endpoint string, params url.Values, data, response interface{}) error {
	// Parse URL
	base, err := url.Parse(apiHost)
	if err != nil {
		return err
	}
	base.Path = path.Join(base.Path, "api", endpoint)
	uri := base.String()

	// Parse params
	if params == nil {
		params = url.Values{}
	}
	encodedParams := params.Encode()
	if encodedParams != "" {
		uri += "?" + params.Encode()
	}

	// Debug
	if c.config.Debug {
		log.Printf("NEW %s REQUEST TO %s with payload: %+v", method, uri, data)
	}

	// Create request
	var req *http.Request
	if data != nil {
		var eData []byte
		eData, err = json.Marshal(data)
		if err != nil {
			return err
		}
		req, err = http.NewRequest(method, uri, bytes.NewBuffer(eData))
		if err != nil {
			return err
		}
	} else {
		req, err = http.NewRequest(method, uri, nil)
		if err != nil {
			return err
		}
	}

	// Set Auth
	req.Header.Set("Authorization", c.config.TokenType+" "+c.config.AccessToken)

	// Perform request
	client := &http.Client{Transport: c.Transport}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Debug
	if c.config.Debug {
		log.Printf("RESPONSE FROM %s: %s", base.String(), string(body))
	}

	// Unmarshal into response
	if len(body) > 0 {
		err = json.Unmarshal(body, response)
	}

	return err
}

// GetContacts method
func (c *Client) GetContacts(params url.Values) ([]Contact, error) {
	r := ContactResponse{}
	err := c.Request("GET", "contact", params, nil, &r)
	return r.Data, err
}
