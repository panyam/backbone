package goclient

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	// "os"
	// "os/exec"
	// "strconv"
)

type Authenticator interface {
	AuthenticateRequest(request *http.Request) (*http.Request, error)
}

type LoginAuthenticator struct {
	Username  string
	Password  string
	sessionId string
}

func (auth *LoginAuthenticator) AuthenticateRequest(req *http.Request) (*http.Request, error) {
	var err error
	if auth.sessionId == "" {
		// then login
	} else {
		log.Println("Logging in...")
		// settion session ID cookie on the request
		params := map[string]string{"password": auth.Password}
		url := fmt.Sprintf("/users/%s/login/", auth.Username)
		req, err = MakeRequest("POST", url, params, nil)
		resp, err := SendRequest(req, nil)
		log.Println("Login Response: ", resp, err)
	}
	return req, err
}

type SignedUrlAuthenticator struct {
	AccessToken string
	SecretKey   string
}

func (auth *SignedUrlAuthenticator) AuthenticateRequest(req *http.Request) (*http.Request, error) {
	// sign the URL with the given params
	baseUrl := req.URL.Host + "/" + req.URL.Path
	returnUrl := baseUrl + "?"
	sigString := auth.SecretKey + baseUrl

	req.ParseForm()
	params, err := url.ParseQuery(req.URL.RawQuery)
	params["access_token"] = append(make([]string, 0, 0), auth.AccessToken)

	// create an array to sort the params first for normalisation
	mk := make([]string, len(params))
	i := 0
	for key := range params {
		mk[i] = key
		i++
	}
	sort.Strings(mk)

	for i := range mk {
		values := params[mk[i]]
		sort.Strings(values)
		joinedValues := strings.Join(values, ",")
		sigString += mk[i] + joinedValues
		returnUrl += mk[i] + "=" + joinedValues + "&"
	}

	hasher := sha1.New()
	io.WriteString(hasher, sigString)
	sig := hex.EncodeToString(hasher.Sum(nil))

	returnUrl += "sig=" + sig
	req.URL, err = url.Parse(baseUrl + returnUrl)
	return req, err
}

type DebugAuthenticator struct {
	Username string
}

func (auth *DebugAuthenticator) AuthenticateRequest(req *http.Request) (*http.Request, error) {
	// Append the username param to a request
	return req, nil
}
