package ghnc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	SIGNUP_URL        = "https://github.com/join?source=login"
	NAMECHECK_API_URL = "https://github.com/signup_check/username"
)

type GHClient struct {
	HttpClient *http.Client
	AuthToken  string
}

// UsernameAvailable checks if a username is available.
func (g *GHClient) UsernameAvailable(username string) (available bool, reason string, err error) {
	// Make the request
	form := url.Values{}
	form.Add("value", username)
	form.Add("authenticity_token", g.AuthToken)

	req, err := http.NewRequest("POST", NAMECHECK_API_URL, strings.NewReader(form.Encode()))
	if err != nil {
		return false, "", err
	}

	// Send the request
	resp, err := g.HttpClient.Do(req)
	if err != nil {
		return false, "", err
	}

	// check the status code
	if resp.StatusCode == 200 {
		// 200 == username available,
		return true, "", nil
	} else if resp.StatusCode == 422 {
		// 422 == username not available, with reason as response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, "", err
		}

		return false, string(body), nil
	}

	return false, "", errors.New(fmt.Sprintf("Unknown response code %v", resp.StatusCode))
}

// GetGHClient gets the required cookies and authenticy token and returns a GHClient.
func GetGHClient() (*GHClient, error) {
	// Setup the HTTP client with cookiejar
	cj, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{Jar: cj}

	// Get cookies from the page
	resp, err := client.Get(SIGNUP_URL)
	if err != nil {
		return nil, err
	}

	// Get the authenticity token
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	var authToken string
	doc.Find("#user_login").Each(func(i int, sel *goquery.Selection) {
		token, exists := sel.Attr("data-autocheck-authenticity-token")
		if exists {
			authToken = token
		}

	})

	// Fail if failed to find the auth token
	if authToken == "" {
		return nil, errors.New("Failed to find data-autocheck-authenticity-token")
	}

	return &GHClient{HttpClient: &client, AuthToken: authToken}, nil
}
