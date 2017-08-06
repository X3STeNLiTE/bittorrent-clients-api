package utorrent

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"booboo-bear.com/utils/nums"

	"github.com/PuerkitoBio/goquery"
)

//Client ...
type Client struct {
	BaseURL   string
	Username  string
	Password  string
	Token     string
	Authorize string
	client    *http.Client
}

//NewClient ...
func NewClient(baseURL string, username string, password string) (ut *Client) {
	if strings.HasSuffix(baseURL, "/") == false {
		baseURL += "/"
	}
	userInfo := url.UserPassword(username, password)
	ut = &Client{
		BaseURL:   baseURL,
		Username:  username,
		Password:  password,
		Authorize: base64.URLEncoding.EncodeToString([]byte(userInfo.String())),
		client:    new(http.Client),
	}
	ut.client.Jar, _ = cookiejar.New(nil)
	ut.login()
	return
}

func (ut *Client) login() error {
	var err error

	tokenURL, _ := url.Parse(ut.BaseURL + "token.html")
	req := http.Request{
		Method: "GET",
		Header: map[string][]string{
			"Authorization":   {"Basic " + ut.Authorize},
			"Accept-Encoding": {"gzip"},
		},
		URL: tokenURL,
	}

	log.Printf("Connecting to ... %v\n", tokenURL)

	res, err := ut.client.Do(&req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ClientError{res.StatusCode, res.Status}
	}
	log.Println("Response code:", res.Status)

	var rc = res.Body
	if res.Header.Get("Content-Encoding") == "gzip" {
		rc, err = gzip.NewReader(res.Body)
		if err != nil {
			log.Panicln(err)
		}
	}
	doc, err := goquery.NewDocumentFromReader(rc)
	if err != nil {
		log.Panicln(err.Error())
	}
	divToken := doc.Find("div")

	//Keep states
	ut.Token = divToken.Text()

	log.Printf("Received Token: %s", ut.Token)
	log.Println("Logged in...")

	return nil
}

func (ut *Client) newRequest(params url.Values, file *os.File) (*http.Request, error) {
	params.Add("t", strconv.FormatInt(time.Now().Unix(), 10))
	params.Add("token", ut.Token)

	var reqURL = fmt.Sprintf("%s?%s", ut.BaseURL, params.Encode())
	var req *http.Request
	var err error

	var headers = http.Header{
		"Authorization":   {"Basic " + ut.Authorize},
		"Accept":          {"application/json"},
		"Accept-Encoding": {"gzip"},
		"Connection":      {"keep-alive"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:53.0) Gecko/20100101 Firefox/53.0"},
	}

	if file != nil {
		body := &bytes.Buffer{}
		w := multipart.NewWriter(body)
		part, err := w.CreateFormFile("torrent_file", file.Name())
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(part, file); err != nil {
			return nil, err
		}
		if err = w.Close(); err != nil {
			return nil, err
		}
		req, err = http.NewRequest(http.MethodPost, reqURL, body)
		headers.Add("Content-Type", w.FormDataContentType())
	} else {
		req, err = http.NewRequest(http.MethodGet, reqURL, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header = headers

	log.Println("Request URL:", reqURL)

	return req, nil
}

func (ut *Client) action(params url.Values) ([]byte, error) {
	tStart := time.Now()

	req, err := ut.newRequest(params, nil)
	if err != nil {
		return nil, err
	}
	res, err := ut.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	log.Printf("[%s] Connected", time.Since(tStart))

	return ioutil.ReadAll(res.Body)
}

//ListFiles ...
func (ut *Client) ListFiles() (*ListFileResponse, error) {
	data, err := ut.action(map[string][]string{
		"list":   {"1"},
		"getmsg": {"1"},
	})
	if err != nil {
		return nil, err
	}
	return unmarshalListFileResponse(data)
}

//ActionResult ...
type ActionResult struct {
	Build int64 `json:"build"`
}

//Separator chuck of 25 element
func (ut *Client) doHashAction(action string, hash []string) error {
	if hLength := len(hash); hLength > 0 {
		offset := 0
		length := 0
		rem := hLength - offset
		for rem > 0 {
			length = nums.MinInt(rem, 25)
			if _, err := ut.action(map[string][]string{
				"action": {action},
				"hash":   hash[offset : offset+length],
			}); err != nil {
				return err
			}
			offset += length
			rem = hLength - offset
			log.Println("OK")
		}
	}
	return nil
}

//Start ...
func (ut *Client) Start(hash []string) error {
	return ut.doHashAction("start", hash)
}

//Stop ...
func (ut *Client) Stop(hash []string) error {
	return ut.doHashAction("stop", hash)
}

//Pause ...
func (ut *Client) Pause(hash []string) error {
	return ut.doHashAction("pause", hash)
}

//ForceStart ...
func (ut *Client) ForceStart(hash []string) error {
	return ut.doHashAction("forcestart", hash)
}

//Unpause ...
func (ut *Client) Unpause(hash []string) error {
	return ut.doHashAction("unpause", hash)
}

//Recheck ...
func (ut *Client) Recheck(hash []string) error {
	return ut.doHashAction("recheck", hash)
}

//Remove ...
func (ut *Client) Remove(hash []string) error {
	return ut.doHashAction("remove", hash)
}

//RemoveData ...
func (ut *Client) RemoveData(hash []string) error {
	return ut.doHashAction("removedata", hash)
}

func (ut *Client) handleActionResponse(resp *http.Response) (*ActionResult, error) {
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, errors.New("bad username or password")
		case http.StatusBadRequest:
			return nil, errors.New("utorrent api returned status code : 400")
		}
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := &ActionResult{}
	err = json.Unmarshal(respData, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//AddURL ..
func (ut *Client) AddURL(torrenturl string) (*ActionResult, error) {
	qs := url.Values{
		"action":       {"add-url"},
		"path":         {""},
		"download_dir": {"0"},
	}
	req, err := ut.newRequest(qs, nil)
	if err != nil {
		return nil, err
	}
	res, err := ut.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ut.handleActionResponse(res)
}

//AddFile ..
func (ut *Client) AddFile(torrentFile *os.File) (*ActionResult, error) {
	qs := url.Values{
		"action":       {"add-file"},
		"path":         {""},
		"download_dir": {"0"},
	}
	req, err := ut.newRequest(qs, torrentFile)
	if err != nil {
		return nil, err
	}
	res, err := ut.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ut.handleActionResponse(res)
}
