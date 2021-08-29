package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"syscall"
)

func main() {
	ctx := context.Background()
	cli := NewHttpCli(false)

	resp, err := cli.doRequest(ctx, "http://localhost:1789/test/context", http.MethodGet, "", "")
	if err != nil {
		fmt.Println(err)
		if e, ok := err.(*url.Error); ok {
			if e == io.EOF {
				fmt.Println("server down")
			}
			err = e.Err
			if op, ok := err.(*net.OpError); ok {
				err = op.Err
			}
			if sys, ok := err.(*os.SyscallError); ok {
				err = sys.Err
				fmt.Println("server refuse connection")

			} else if err != syscall.EPERM {
				fmt.Printf("WriteMsgUnix failed with %v, want EPERM\n ", err)
			}
		}
	}

	fmt.Println(resp.Status, resp.StatusCode)
}

type HttpCli struct {
	Client *http.Client
}

func NewHttpCli(skipVerify bool) *HttpCli {
	defaultHttpCli := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
		// A Timeout of zero means no timeout. default 0
		//Timeout: 0,
	}
	if skipVerify {
		defaultHttpCli = &http.Client{
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}
	return &HttpCli{Client: defaultHttpCli}
}

func (c *HttpCli) client() *http.Client {
	return c.Client
}

func (c *HttpCli) doRequest(ctx context.Context, path, method string, in, out interface{}) (*http.Response, error) {
	var buf bytes.Buffer

	// marshal the input payload into json format and copy
	// to an io.ReadCloser.
	if in != nil {
		json.NewEncoder(&buf).Encode(in)
	}

	endpoint := path
	req, err := http.NewRequest(method, endpoint, &buf)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	res, err := c.client().Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return res, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}
	if res.StatusCode > 299 {
		// if the response body includes an error message
		// we should return the error string.
		if len(body) != 0 {
			return res, errors.New(
				string(body),
			)
		}
		// if the response body is empty we should return
		// the default status code text.
		return res, errors.New(
			http.StatusText(res.StatusCode),
		)
	}
	if out == nil {
		return res, nil
	}
	fmt.Println(string(body))
	return res, json.Unmarshal(body, out)
}

/*
http.Client的Timeout
The timeout includes connection time, any
	// redirects, and reading the response body. The timer remains
	// running after Get, Head, Post, or Do return and will
	// interrupt reading of the Response.Body
*/
