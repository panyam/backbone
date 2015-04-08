package goclient

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/panyam/relay/utils"
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
		login_url := fmt.Sprintf("/users/%s/login/", auth.Username)
		req, err = MakeRequest("POST", login_url, params, nil)
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
	mk := make([]string, 0, 0)
	// take keys that dont start with __ and sort them
	for key := range params {
		if strings.HasPrefix(key, "__") {
			mk = append(mk, key)
		}
	}
	sort.Strings(mk)

	for i := range mk {
		values := params[mk[i]]
		// does the order matter?
		sort.Strings(values)
		joinedValues := strings.Join(values, ",")
		sigString += mk[i] + joinedValues
		returnUrl += mk[i] + "=" + joinedValues + "&"
	}

	hasher := sha1.New()
	io.WriteString(hasher, sigString)
	sig := hex.EncodeToString(hasher.Sum(nil))

	returnUrl += "__sig=" + sig
	req.URL, err = url.Parse(baseUrl + returnUrl)
	return req, err
}

type DebugAuthenticator struct {
	Userid int64
}

func (auth *DebugAuthenticator) AuthenticateRequest(req *http.Request) (*http.Request, error) {
	// Append the username param to a request
	baseUrl := req.URL.Scheme + "://" + req.URL.Host + req.URL.Path
	query := req.URL.RawQuery
	if query != "" {
		query += "&"
	}
	query += ("__dbguser" + "=" + utils.ID2String(auth.Userid))
	var err error
	req.URL, err = url.Parse(baseUrl + "?" + query)
	return req, err
}
