package gomagnet

import "os"

var (
	apiHost      = "https://bot.magnet.chat"
	clientID     string
	clientSecret string
	username     string
	password     string
)

func init() {
	// Set environment variables configuration.
	if os.Getenv("MAGNET_API_HOST") != "" {
		apiHost = os.Getenv("MAGNET_API_HOST")
	}
	if os.Getenv("MAGNET_CLIENT_ID") != "" {
		clientID = os.Getenv("MAGNET_CLIENT_ID")
	}
	if os.Getenv("MAGNET_CLIENT_SECRET") != "" {
		clientSecret = os.Getenv("MAGNET_CLIENT_SECRET")
	}
	if os.Getenv("MAGNET_USERNAME") != "" {
		username = os.Getenv("MAGNET_USERNAME")
	}
	if os.Getenv("MAGNET_PASSWORD") != "" {
		password = os.Getenv("MAGNET_PASSWORD")
	}
}
