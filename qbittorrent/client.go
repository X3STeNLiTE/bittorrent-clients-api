package qbittorrent

//Implemented by following API Document
//https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

//Client ... qBittorrent Configure
type Client struct {
	BaseURL   string
	Username  string
	Password  string
	Token     string
	Cookies   http.CookieJar
	Authorize string
	client    *http.Client
	connected bool
}

type httpMethod string

const (
	methodGet  httpMethod = "GET"
	methodPost httpMethod = "POST"
)

//NewClient ...
func NewClient(baseURL string, username string, password string) (qb *Client, err error) {
	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}
	userInfo := url.UserPassword(username, password)
	//Check URL syntax
	_, err = url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	qb = &Client{
		BaseURL:   baseURL,
		Username:  username,
		Password:  password,
		Authorize: base64.URLEncoding.EncodeToString([]byte(userInfo.String())),
		client:    new(http.Client),
	}
	qb.client.Jar, _ = cookiejar.New(nil)
	err = qb.login()
	return
}

func (qb *Client) do(method httpMethod, path string, query url.Values, form url.Values) (data []byte, err error) {
	tStart := time.Now()

	defer log.Println(times.TimeTracks(tStart, "HttpRequest"))

	_url, _ := url.Parse(qb.BaseURL + path)
	if query != nil {
		_url.RawQuery = maps.Join(query, "&")
	}
	_urlStr := _url.String()
	_req, _ := http.NewRequest(string(method), _urlStr, strings.NewReader(form.Encode()))
	if method == methodPost {
		_req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	log.Printf("Connecting to %s ...", _urlStr)
	_res, _err := qb.client.Do(_req)
	if _err != nil {
		return nil, _err
	}

	defer _res.Body.Close()

	log.Println(times.TimeTracks(tStart, "Connected"))
	if _res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s : %s", _urlStr, _res.Status)
	}
	return ioutil.ReadAll(_res.Body)
}

func (qb *Client) login() error {
	if qb.connected {
		return errors.New("Connected")
	}
	form := url.Values{}
	form.Add("username", qb.Username)
	form.Add("password", qb.Password)
	_, err := qb.do(methodPost, "/login", nil, form)
	if err != nil {
		return err
	}
	qb.connected = true
	return nil
}

//Logout ...
func (qb *Client) Logout() error {
	if !qb.connected {
		return errors.New("Not Connected")
	}
	_, err := qb.do(methodGet, "/logout", nil, nil)
	return err
}

//APIVersion ...
func (qb *Client) APIVersion() (string, error) {
	data, err := qb.do(methodGet, "/version/api", nil, nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

//APIMinimumVersion ...
func (qb *Client) APIMinimumVersion() (string, error) {
	data, err := qb.do(methodGet, "/version/api_min", nil, nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

//ListTorrent ...
func (qb *Client) ListTorrent() (torrents *[]Torrent, err error) {
	if !qb.connected {
		return nil, errors.New("Not Connected")
	}
	data, err := qb.do(methodGet, "/query/torrents", nil, nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &torrents)
	return
}

//Common command executor
func (qb *Client) doCommand(cmd string, params url.Values) error {
	qb.do(methodPost, "/command/"+cmd, nil, params)
	return nil
}

//Start ...
//https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation#resume-torrent
func (qb *Client) Start(hash []string) error {
	return qb.doCommand("resume", map[string][]string{
		"hash": hash,
	})
}

//StartAll ...
//https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation#resume-all-torrents
func (qb *Client) StartAll(hash []string) error {
	return qb.doCommand("resumeAll", nil)
}

//Pause ...
//https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation#pause-torrent
func (qb *Client) Pause(hash []string) error {
	return qb.doCommand("pause", map[string][]string{
		"hash": hash,
	})
}

//PauseAll ...
//https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation#pause-all-torrents
func (qb *Client) PauseAll(hash []string) error {
	return qb.doCommand("pauseAll", nil)
}

//Remove ...
//https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation#delete-torrent
func (qb *Client) Remove(hash []string) error {
	return qb.doCommand("delete", map[string][]string{
		"hashs": {strings.Join(hash, "|")},
	})
}

//RemoveData ...
//https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation#delete-torrent-with-downloaded-data
func (qb *Client) RemoveData(hash []string) error {
	return qb.doCommand("deletePerm", map[string][]string{
		"hashs": {strings.Join(hash, "|")},
	})
}

//Recheck ...
//https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation#recheck-torrent
func (qb *Client) Recheck(hash []string) error {
	return qb.doCommand("pause", map[string][]string{
		"hash": hash,
	})
}

//Command ...
func (qb *Client) Command(action string, hash []string) error {
	return nil
}
