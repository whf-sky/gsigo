package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func NewCurl(url string) *Curl {
	return &Curl{
		url:url,
		query: map[string]string{},
		cookie: map[string]string{},
		header: map[string]string{},
	}
}

type Curl struct {
	url string
	query map[string]string
	cookie map[string]string
	header map[string]string
	body io.Reader
	req *http.Request
	res *http.Response
}

func (c *Curl) Url(url string) *Curl {
	c.url = url
	return c
}

func  (c *Curl) Query(query map[string]string) *Curl {
	c.query = query
	return c
}

func (c *Curl) SetQuery(name string, value string) *Curl {
	c.query[name] = value
	return c
}

func (c *Curl) Body(body io.Reader) *Curl {
	c.body = body
	return c
}

func (c *Curl) Cookie(cookie map[string]string) *Curl {
	c.cookie = cookie
	return c
}

func (c *Curl) SetCookie(name string, value string) *Curl {
	c.cookie[name] = value
	return c
}

func (c *Curl) Header(header map[string]string) *Curl {
	c.header = header
	return c
}

func (c *Curl) SetHeader(name string, value string) *Curl {
	c.header[name] = value
	return c
}

func (c *Curl) SetPostForm() *Curl {
	c.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	return c
}

func (c *Curl) close() {
	_ = c.res.Body.Close()
}

func (c *Curl) getResponse() *http.Response {
	return c.res
}

func (c *Curl) GetBody(method string) ([]byte, error) {
	if err := c.do(method); err != nil {
		return nil, err
	}
	defer c.close()
	body, err := ioutil.ReadAll(c.res.Body)
	return body, err
}

func (c *Curl) GetHeader() http.Header {
	return c.res.Header
}

func (c *Curl) Get() ([]byte, error) {
	return c.GetBody(http.MethodGet)
}

func (c *Curl) Post() ([]byte, error) {
	return c.GetBody(http.MethodPost)
}

func  (c *Curl) Options() ([]byte, error) {
	return c.GetBody(http.MethodOptions)
}

func (c *Curl) Head() (http.Header, error) {
	if err := c.do(http.MethodHead); err != nil {
		return nil, err
	}
	defer c.close()
	return c.res.Header ,nil
}

func (c *Curl) Put() ([]byte, error) {
	return c.GetBody(http.MethodPut)
}

func (c *Curl) Delete() ([]byte, error) {
	return c.GetBody(http.MethodDelete)
}

func (c *Curl) Patch() ([]byte, error) {
	return c.GetBody(http.MethodPatch)
}

func (c *Curl)  Connect() ([]byte, error) {
	return c.GetBody(http.MethodConnect)
}

func (c *Curl) do(method string) error {
	var err error
	client := &http.Client{}
	//url
	if len(c.query) > 0 {
		uparse, _ := url.Parse(c.url)
		qv := uparse.Query()
		for k,v := range c.query {
			qv.Set(k, v)
		}
		query := qv.Encode()
		c.url = uparse.Scheme+"://"+uparse.Host+uparse.Path+"?"+query
	}

	//new request
	c.req, err = http.NewRequest(method, c.url, c.body)
	if err != nil {
		return err
	}

	//add cookie
	if len(c.cookie) > 0 {
		for name, value := range c.cookie  {
			c.req.AddCookie(&http.Cookie{Name:name, Value:value})
		}
	}

	//add header
	if len(c.header) > 0 {
		for name, value := range c.header  {
			c.req.Header.Add(name, value)
		}
	}

	c.res, err = client.Do(c.req)
	if err != nil {
		return  err
	}
	return nil
}