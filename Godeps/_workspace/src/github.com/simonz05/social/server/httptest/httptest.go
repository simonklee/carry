package httptest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type Request http.Request
type Response http.Response

func NewRequest(method string, uri string, body interface{}, params url.Values) (*Request, error) {
	method = strings.ToUpper(method)

	if body != nil && (method != "POST" && method != "PUT") {
		return nil, fmt.Errorf("%s method does not accept body", method)
	}

	var buf io.Reader

	if body != nil {
		b, ok := body.([]byte)
		if ok {
			buf = bytes.NewBuffer(b)
		} else {
			body, err := toURL(body)

			if err != nil {
				return nil, err
			}

			buf = strings.NewReader(body.Encode())
		}
	}

	req, err := http.NewRequest(method, joinURL(uri, params), buf)

	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	rr := Request(*req)
	return &rr, err
}

func (req *Request) Do() (*Response, error) {
	rr, err := http.DefaultClient.Do((*http.Request)(req))

	if err != nil {
		return nil, err
	}

	return (*Response)(rr), nil
}

func (res *Response) ToJSON(i interface{}) error {
	defer res.Body.Close()

	if c := res.Header.Get("Content-Type"); !strings.Contains(c, "application/json") {
		return fmt.Errorf("Unexpected Content-Type, got %s", c)
	}

	reader := bufio.NewReader(res.Body)
	buf, _ := ioutil.ReadAll(reader)
	err := json.Unmarshal(buf, i)
	//fmt.Printf("%s\n", buf)
	//err := json.NewDecoder(res.Body).Decode(v)
	return err
}

type DataErrResponse struct {
	Error map[string]string
}

func (res *Response) ToErr() (*DataErrResponse, error) {
	v := &DataErrResponse{}
	err := res.ToJSON(v)
	return v, err
}

func toURL(query interface{}) (url.Values, error) {
	switch vv := query.(type) {
	case url.Values:
		return query.(url.Values), nil
	case map[string]string:
		val := make(url.Values, len(vv))
		for k, v := range vv {
			val.Add(k, v)
		}
		return val, nil
	default:
		s := reflect.ValueOf(query)
		t := reflect.TypeOf(query)
		val := make(url.Values, s.NumField())

		for i := 0; i < s.NumField(); i++ {
			val.Add(strings.ToLower(t.Field(i).Name), fmt.Sprintf("%v", s.Field(i).Interface()))
		}
		return val, nil
	}
}

func joinURL(endpoint string, args url.Values) string {
	var params string

	if args != nil && len(args) > 0 {
		params = "?" + args.Encode()
	}

	return endpoint + params
}

// type HTTPTest struct {
// 	Method     string
// 	URL        string
// 	Body       interface{}
// 	Args       url.Values
// 	StatusCode int
// }
//
// func NewHTTPTest(method string, url string, body interface{}, args url.Values, statusCode int) *HTTPTest {
// 	return &HTTPTest{method, url, body, args, statusCode}
// }
//
// func (h *HTTPTest) String() string {
// 	var args string
// 	var body string
//
// 	if h.Args != nil {
// 		args = "?" + h.Args.Encode()
// 	}
//
// 	if h.Body != nil {
// 		body = fmt.Sprintf(" %s", h.Body)
// 	}
//
// 	return fmt.Sprintf("%s\t%s%s%s %d", h.Method, h.URL, args, body, h.StatusCode)
// }
//
// func (v *HTTPTest) Test(t *testing.T) *http.Response {
// 	if v.Body != nil && (v.Method == "GET" || v.Method == "DELETE") {
// 		t.Fatalf("%s unexcepted body - ERR", v)
// 	}
//
// 	body, err := toURL(v.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	uri := joinURL(v.URL, v.Args)
// 	req, _ := http.NewRequest(v.Method, uri, strings.NewReader(body.Encode()))
//
// 	if v.Method == "POST" || v.Method == "PUT" {
// 		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	}
//
// 	client := &http.Client{}
// 	r, err := client.Do(req)
//
// 	if err != nil {
// 		t.Fatalf("%s: err sending HTTP request: %s", v, err)
// 	}
//
// 	if r.StatusCode != v.StatusCode {
// 		if r.StatusCode == 400 {
// 			errRes := ReadErrResponse(t, r)
//
// 			for k, v := range errRes.Error {
// 				fmt.Println(k, v)
// 			}
// 		}
//
// 		t.Fatalf("%s: expected status code %d, got %d", v, v.StatusCode, r.StatusCode)
// 	}
//
// 	fmt.Printf("%s - OK\n", v)
// 	return r
// }
//
// func ParseResponse(res *http.Response, v interface{}) error {
// 	defer res.Body.Close()
//
// 	if c := res.Header.Get("Content-Type"); !strings.Contains(c, "application/json") {
// 		return fmt.Errorf("Unexpected Content-Type, got %s", c)
// 	}
//
// 	reader := bufio.NewReader(res.Body)
// 	buf, _ := ioutil.ReadAll(reader)
// 	err := json.Unmarshal(buf, v)
// 	//fmt.Printf("%s\n", buf)
// 	//err := json.NewDecoder(res.Body).Decode(v)
// 	return err
// }
//
// func ReadErrResponse(t *testing.T, res *http.Response) *DataErrResponse {
// 	v := new(DataErrResponse)
// 	ParseResponse(res, v)
// 	return v
// }
