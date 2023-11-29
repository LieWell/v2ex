package httpclient

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

var defaultSetting = Config{
	UserAgent:        "http-client/v0.1",
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,
	Gzip:             true,
	DumpBody:         true,
}

var defaultCookieJar http.CookieJar
var mutex sync.Mutex

// createDefaultCookie creates a global CookieJar to store cookies.
func createDefaultCookie() {
	mutex.Lock()
	defer mutex.Unlock()
	defaultCookieJar, _ = cookiejar.New(nil)
}

// NewRequest return *HttpRequest with specific method
func NewRequest(rawUrl, method string) *Request {
	var resp http.Response
	u, err := url.Parse(rawUrl)
	if err != nil {
		log.Printf("http lib parse RawURL error: %v", err)
	}
	req := http.Request{
		URL:        u,
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	return &Request{
		url:     rawUrl,
		req:     &req,
		params:  map[string][]string{},
		files:   map[string]string{},
		setting: defaultSetting,
		resp:    &resp,
	}
}

func Get(url string) *Request {
	return NewRequest(url, http.MethodGet)
}

func Post(url string) *Request {
	return NewRequest(url, http.MethodPost)
}

func Put(url string) *Request {
	return NewRequest(url, http.MethodPut)
}

func Delete(url string) *Request {
	return NewRequest(url, http.MethodDelete)
}

func Head(url string) *Request {
	return NewRequest(url, "HEAD")
}

// Config is the http.Client setting
type Config struct {
	ShowDebug        bool
	UserAgent        string
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	TLSClientConfig  *tls.Config
	Proxy            func(*http.Request) (*url.URL, error)
	Transport        http.RoundTripper
	CheckRedirect    func(req *http.Request, via []*http.Request) error
	EnableCookie     bool
	Gzip             bool
	DumpBody         bool
	Retries          int // if set to -1 means will retry forever
}

// Request provides more useful methods for requesting one url than http.Request.
type Request struct {
	url     string
	req     *http.Request
	params  map[string][]string
	files   map[string]string
	setting Config
	resp    *http.Response
	body    []byte
	dump    []byte
}

// GetURL return the request url
func (b *Request) GetURL() string {
	return b.url
}

// GetRequest return the request object
func (b *Request) GetRequest() *http.Request {
	return b.req
}

// Setting Change request config
func (b *Request) Setting(setting Config) *Request {
	b.setting = setting
	return b
}

// SetBasicAuth sets the request's Authorization header to use HTTP Basic Authentication with the provided username and password.
func (b *Request) SetBasicAuth(username, password string) *Request {
	b.req.SetBasicAuth(username, password)
	return b
}

// SetEnableCookie sets enable/disable cookiejar
func (b *Request) SetEnableCookie(enable bool) *Request {
	b.setting.EnableCookie = enable
	return b
}

// SetUserAgent sets User-Agent header field
func (b *Request) SetUserAgent(useragent string) *Request {
	b.setting.UserAgent = useragent
	return b
}

// Debug sets show debug or not when executing request.
func (b *Request) Debug(isdebug bool) *Request {
	b.setting.ShowDebug = isdebug
	return b
}

// Retries sets Retries times.
// default is 0 means no retried.
// -1 means retried forever.
// others means retried times.
func (b *Request) Retries(times int) *Request {
	b.setting.Retries = times
	return b
}

// DumpBody setting whether need to Dump the Body.
func (b *Request) DumpBody(isdump bool) *Request {
	b.setting.DumpBody = isdump
	return b
}

// DumpRequest return the DumpRequest
func (b *Request) DumpRequest() []byte {
	return b.dump
}

// SetTimeout sets connect time out and read-write time out for Request.
func (b *Request) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *Request {
	b.setting.ConnectTimeout = connectTimeout
	b.setting.ReadWriteTimeout = readWriteTimeout
	return b
}

// SetTLSClientConfig sets tls connection configurations if visiting https url.
func (b *Request) SetTLSClientConfig(config *tls.Config) *Request {
	b.setting.TLSClientConfig = config
	return b
}

// Header add header item string in request.
func (b *Request) Header(key, value string) *Request {
	b.req.Header.Set(key, value)
	return b
}

// SetHost set the request host
func (b *Request) SetHost(host string) *Request {
	b.req.Host = host
	return b
}

// SetProtocolVersion Set the protocol version for incoming requests.
// Client requests always use HTTP/1.1.
func (b *Request) SetProtocolVersion(vers string) *Request {
	if len(vers) == 0 {
		vers = "HTTP/1.1"
	}

	major, minor, ok := http.ParseHTTPVersion(vers)
	if ok {
		b.req.Proto = vers
		b.req.ProtoMajor = major
		b.req.ProtoMinor = minor
	}

	return b
}

// SetCookie add cookie into request.
func (b *Request) SetCookie(cookie *http.Cookie) *Request {
	b.req.Header.Add("Cookie", cookie.String())
	return b
}

// SetTransport set the setting transport
func (b *Request) SetTransport(transport http.RoundTripper) *Request {
	b.setting.Transport = transport
	return b
}

// SetProxy set the http proxy
// example:
//
//	func(req *http.Request) (*url.URL, error) {
//		u, _ := url.ParseRequestURI("http://127.0.0.1:8080")
//		return u, nil
//	}
func (b *Request) SetProxy(proxy func(*http.Request) (*url.URL, error)) *Request {
	b.setting.Proxy = proxy
	return b
}

// SetCheckRedirect specifies the policy for handling redirects.
//
// If CheckRedirect is nil, the Client uses its default policy,
// which is to stop after 10 consecutive requests.
func (b *Request) SetCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) *Request {
	b.setting.CheckRedirect = redirect
	return b
}

// Param adds query param in to request.
// params build query string as ?key1=value1&key2=value2...
func (b *Request) Param(key, value string) *Request {
	if param, ok := b.params[key]; ok {
		b.params[key] = append(param, value)
	} else {
		b.params[key] = []string{value}
	}
	return b
}

// PostFile add a post file to the request
func (b *Request) PostFile(form, filename string) *Request {
	b.files[form] = filename
	return b
}

// Body adds request raw body.
// it supports string and []byte.
func (b *Request) Body(data interface{}) *Request {
	switch t := data.(type) {
	case string:
		bf := bytes.NewBufferString(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.ContentLength = int64(len(t))
	case []byte:
		bf := bytes.NewBuffer(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.ContentLength = int64(len(t))
	}
	return b
}

// XMLBody adds request raw body encoding by XML.
func (b *Request) XMLBody(obj interface{}) (*Request, error) {
	if b.req.Body == nil && obj != nil {
		byts, err := xml.Marshal(obj)
		if err != nil {
			return b, err
		}
		b.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		b.req.ContentLength = int64(len(byts))
		b.req.Header.Set("Content-Type", "application/xml")
	}
	return b, nil
}

// YAMLBody adds request raw body encoding by YAML.
func (b *Request) YAMLBody(obj interface{}) (*Request, error) {
	if b.req.Body == nil && obj != nil {
		byts, err := yaml.Marshal(obj)
		if err != nil {
			return b, err
		}
		b.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		b.req.ContentLength = int64(len(byts))
		b.req.Header.Set("Content-Type", "application/x+yaml")
	}
	return b, nil
}

// JSONBody adds request raw body encoding by JSON.
func (b *Request) JSONBody(obj interface{}) (*Request, error) {
	if b.req.Body == nil && obj != nil {
		byts, err := json.Marshal(obj)
		if err != nil {
			return b, err
		}
		b.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		b.req.ContentLength = int64(len(byts))
		b.req.Header.Set("Content-Type", "application/json")
	}
	return b, nil
}

func (b *Request) buildURL(paramBody string) {
	// build GET url with query string
	if b.req.Method == http.MethodGet && len(paramBody) > 0 {
		if strings.Contains(b.url, "?") {
			b.url += "&" + paramBody
		} else {
			b.url = b.url + "?" + paramBody
		}
		return
	}

	// build POST/PUT/PATCH url and body
	if (b.req.Method == http.MethodPost || b.req.Method == http.MethodPut || b.req.Method == http.MethodPatch || b.req.Method == http.MethodDelete) && b.req.Body == nil {
		// with files
		if len(b.files) > 0 {
			pr, pw := io.Pipe()
			bodyWriter := multipart.NewWriter(pw)
			go func() {
				for form, filename := range b.files {
					fileWriter, err := bodyWriter.CreateFormFile(form, filename)
					if err != nil {
						log.Printf("http client create file writer error: %v", err)
					}
					fh, err := os.Open(filename)
					if err != nil {
						log.Printf("http client open file [%v] error: %v", filename, err)
					}

					_, err = io.Copy(fileWriter, fh)
					fh.Close()
					if err != nil {
						log.Printf("http client copy file error: %v", err)
					}
				}
				for k, v := range b.params {
					for _, vv := range v {
						bodyWriter.WriteField(k, vv)
					}
				}
				bodyWriter.Close()
				pw.Close()
			}()
			b.Header("Content-Type", bodyWriter.FormDataContentType())
			b.req.Body = ioutil.NopCloser(pr)
			return
		}

		// with params
		if len(paramBody) > 0 {
			b.Header("Content-Type", "application/x-www-form-urlencoded")
			b.Body(paramBody)
		}
	}
}

func (b *Request) getResponse() (*http.Response, error) {
	if b.resp.StatusCode != 0 {
		return b.resp, nil
	}
	resp, err := b.DoRequest()
	if err != nil {
		return nil, err
	}
	b.resp = resp
	return resp, nil
}

// DoRequest will do the client.Do
func (b *Request) DoRequest() (resp *http.Response, err error) {
	var paramBody string
	if len(b.params) > 0 {
		var buf bytes.Buffer
		for k, v := range b.params {
			for _, vv := range v {
				buf.WriteString(url.QueryEscape(k))
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(vv))
				buf.WriteByte('&')
			}
		}
		paramBody = buf.String()
		paramBody = paramBody[0 : len(paramBody)-1]
	}

	b.buildURL(paramBody)
	urlParsed, err := url.Parse(b.url)
	if err != nil {
		return nil, err
	}

	b.req.URL = urlParsed

	trans := b.setting.Transport

	if trans == nil {
		// create default transport
		trans = &http.Transport{
			TLSClientConfig:     b.setting.TLSClientConfig,
			Proxy:               b.setting.Proxy,
			Dial:                TimeoutDialer(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout),
			MaxIdleConnsPerHost: 100,
		}
	} else {
		// if b.transport is *http.Transport then set the config.
		if t, ok := trans.(*http.Transport); ok {
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = b.setting.TLSClientConfig
			}
			if t.Proxy == nil {
				t.Proxy = b.setting.Proxy
			}
			if t.Dial == nil {
				t.Dial = TimeoutDialer(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout)
			}
		}
	}

	var jar http.CookieJar
	if b.setting.EnableCookie {
		if defaultCookieJar == nil {
			createDefaultCookie()
		}
		jar = defaultCookieJar
	}

	client := &http.Client{
		Transport: trans,
		Jar:       jar,
	}

	if b.setting.UserAgent != "" && b.req.Header.Get("User-Agent") == "" {
		b.req.Header.Set("User-Agent", b.setting.UserAgent)
	}

	if b.setting.CheckRedirect != nil {
		client.CheckRedirect = b.setting.CheckRedirect
	}

	if b.setting.ShowDebug {
		dump, err := httputil.DumpRequest(b.req, b.setting.DumpBody)
		if err != nil {
			log.Println(err.Error())
		}
		b.dump = dump
	}
	// retries default value is 0, it will run once.
	// retries equal to -1, it will run forever until success
	// retries is setted, it will retries fixed times.
	for i := 0; b.setting.Retries == -1 || i <= b.setting.Retries; i++ {
		resp, err = client.Do(b.req)
		if err == nil {
			break
		}
	}
	return resp, err
}

// String returns the body string in response.
// it calls Response inner.
func (b *Request) String() (string, error) {
	data, err := b.Bytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Bytes returns the body []byte in response.
// it calls Response inner.
func (b *Request) Bytes() ([]byte, error) {
	if b.body != nil {
		return b.body, nil
	}
	resp, err := b.getResponse()
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	if b.setting.Gzip && resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		b.body, err = ioutil.ReadAll(reader)
		return b.body, err
	}
	b.body, err = ioutil.ReadAll(resp.Body)
	return b.body, err
}

// ToFile saves the body data in response to one file.
// it calls Response inner.
func (b *Request) ToFile(filename string) error {
	resp, err := b.getResponse()
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	err = pathExistAndMkdir(filename)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// Check that the file directory exists, there is no automatically created
func pathExistAndMkdir(filename string) (err error) {
	filename = path.Dir(filename)
	_, err = os.Stat(filename)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(filename, os.ModePerm)
		if err == nil {
			return nil
		}
	}
	return err
}

// ToJSON returns the map that marshals from the body bytes as json in response .
// it calls Response inner.
func (b *Request) ToJSON(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// ToXML returns the map that marshals from the body bytes as xml in response .
// it calls Response inner.
func (b *Request) ToXML(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

// ToYAML returns the map that marshals from the body bytes as yaml in response .
// it calls Response inner.
func (b *Request) ToYAML(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

// Response executes request client gets response manually.
func (b *Request) Response() (*http.Response, error) {
	return b.getResponse()
}

// TimeoutDialer returns functions of connection dialer with timeout config for http.Transport Dial field.
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		err = conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, err
	}
}
