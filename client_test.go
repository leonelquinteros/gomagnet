package gomagnet

import (
	"net/url"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	conf := NewClientConfig()
	conf.Debug = true
	c := New(conf)

	token, err := c.GetAccessToken()
	if err != nil {
		t.Fatal(err)
	}
	if c.config.AccessToken != token.AccessToken {
		t.Error("Failed to set AccessToken")
	}
	if c.config.TokenType != token.TokenType {
		t.Error("Failed to set TokenType")
	}
	if c.config.RefreshToken != token.RefreshToken {
		t.Error("Failed to set RefreshToken")
	}
	if c.config.TokenExpiration == nil || !c.config.TokenExpiration.After(time.Now()) {
		t.Error("Failed to set TokenExpiration")
	}

	params := url.Values{}

	contacts, err := c.GetContacts(params)
	if err != nil {
		t.Error(err)
	}

	t.Log(contacts)
}
