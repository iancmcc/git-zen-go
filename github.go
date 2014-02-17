package main

import (
	"bufio"
	"bytes"
	"code.google.com/p/gopass"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var baseURL, _ = url.Parse("https://api.github.com/")

const (
	tokenfile   = ".zendev.gitauth"
	req         = `{"scopes":["repo"],"note":"git-zen"}`
	mediaTypeV3 = "application/vnd.github.v3+json"
)

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("GitHub Username: ")
	username, _ := reader.ReadString('\n')
	password, _ := gopass.GetPass("GitHub Password: ")
	return username, password
}

func getURL(path string) (result *url.URL, err error) {
	p, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	result = baseURL.ResolveReference(p)
	return result, nil
}

func requestToken(username, password string) (token string, err error) {
	client := &http.Client{}
	authurl, err := getURL("authorizations")
	if err != nil {
		return token, err
	}
	req, err := http.NewRequest("POST", authurl.String(),
		strings.NewReader(req))
	if err != nil {
		return token, err
	}
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}
	rb := make(map[string]interface{})
	json.Unmarshal(body, &rb)
	return rb["token"].(string), err
}

func getOAuthToken() (token string, err error) {
	usr, err := user.Current()
	if err != nil {
		return token, err
	}
	fname := filepath.Join(usr.HomeDir, tokenfile)
	tkn, err := ioutil.ReadFile(fname)
	if err != nil {
		// Assume file didn't exist; get the token
		u, p := credentials()
		token, err = requestToken(u, p)
		if err != nil {
			return token, err
		}
		if err := ioutil.WriteFile(fname, []byte(token), os.ModePerm); err != nil {
			// Maybe log here but don't fail
			return token, err
		}
	}
	return string(tkn), err
}

func githubRequest(method string, url string, body interface{}) (req *http.Request, err error) {
	u, err := getURL(url)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err = http.NewRequest("POST", u.String(), buf)
	if err != nil {
		return nil, err
	}
	token, err := getOAuthToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", mediaTypeV3)
	req.Header.Add("Authorization", "token "+token)
	return req, err
}
