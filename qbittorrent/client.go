package qbittorrent

//Implemented by following API Document
//https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

func (qb *Client) do(method httpMethod, path string, query url.Values, form url.Values) (body io.ReadCloser, err error) {
	urlBuf := new(bytes.Buffer)
	urlBuf.WriteString(qb.BaseURL)
	urlBuf.WriteString(path)
	if query != nil {
		urlBuf.WriteByte('?')
		urlBuf.WriteString(query.Encode())
	}
	req, _ := http.NewRequest(string(method), urlBuf.String(), strings.NewReader(form.Encode()))
	if method == http.MethodPost {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	res, _err := qb.client.Do(req)
	if _err != nil {
		return nil, _err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s : %s", urlBuf.String(), res.Status)
	}
	return res.Body, nil
}

func (qb *Client) login() error {
	if qb.connected {
		return errors.New("Connected")
	}
	form := url.Values{}
	form.Add("username", qb.Username)
	form.Add("password", qb.Password)
	_, err := qb.do(http.MethodPost, "/login", nil, form)
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
	_, err := qb.do(http.MethodGet, "/logout", nil, nil)
	return err
}

func readAllString(body io.ReadCloser) (string, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

//APIVersion ...
func (qb *Client) APIVersion() (string, error) {
	body, err := qb.do(http.MethodGet, "/version/api", nil, nil)
	if err != nil {
		return "", err
	}
	defer body.Close()
	return readAllString(body)
}

//APIMinimumVersion ...
func (qb *Client) APIMinimumVersion() (string, error) {
	body, err := qb.do(http.MethodGet, "/version/api_min", nil, nil)
	if err != nil {
		return "", err
	}
	defer body.Close()
	return readAllString(body)
}

//ListTorrent ...
func (qb *Client) ListTorrent() (torrents *[]Torrent, err error) {
	if !qb.connected {
		return nil, errors.New("Not Connected")
	}
	body, err := qb.do(http.MethodGet, "/query/torrents", nil, nil)
	if err != nil {
		return
	}
	jd := json.NewDecoder(body)
	err = jd.Decode(&torrents)
	if err != nil {
		return nil, err
	}
	return torrents, nil
}

//Common command executor
func (qb *Client) doCommand(cmd string, params url.Values) error {
	qb.do(http.MethodPost, "/command/"+cmd, nil, params)
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

// Recheck ...
// https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation#recheck-torrent
func (qb *Client) Recheck(hash []string) error {
	return qb.doCommand("pause", map[string][]string{
		"hash": hash,
	})
}

// Command ...
func (qb *Client) Command(action string, hash []string) error {
	return nil
}

// UploadFile .
// https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation#upload-torrent-from-disk
func (qb *Client) UploadFile(file *os.File) error {
	mpFile, _ := ioutil.TempFile(os.TempDir(), "qbtmp")
	defer os.Remove(mpFile.Name())

	wr := multipart.NewWriter(mpFile)
	part, err := wr.CreateFormFile("torrents", filepath.Base(file.Name()))
	if err != nil {
		return err
	}
	io.Copy(part, file)

	ulURL := createURL("/command/upload", nil)

	req, err := http.NewRequest(http.MethodPost, ulURL, mpFile)
	if err != nil {
		return err
	}
	res, err := qb.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload torrent: %v", err)
	}
	return nil
}

func createURL(path string, query url.Values) string {
	buf := new(bytes.Buffer)
	buf.WriteString(path)
	if query != nil {
		buf.WriteByte('?')
		buf.WriteString(query.Encode())
	}
	return buf.String()
}
