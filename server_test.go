package aduket

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Book struct {
	ISBN string `xml:"isbn"`
	Name string `xml:"name"`
}

type ExpectedResponse struct {
	statusCode int
	header     http.Header
	body       []byte
}

func TestServerResponse(t *testing.T) {
	tests := []struct {
		route               Route
		responseRuleOptions []ResponseRuleOption
		expectedResponse    ExpectedResponse
	}{
		{
			route: Route{HttpMethod: http.MethodGet, Path: "/user"},
			responseRuleOptions: []ResponseRuleOption{
				StatusCode(http.StatusOK),
				JSONBody(User{ID: 123, Name: "kalt"}),
				Header(http.Header{"Content-Type": []string{"application/json"}}),
			},
			expectedResponse: ExpectedResponse{
				statusCode: http.StatusOK,
				header:     http.Header{"Content-Type": []string{"application/json"}},
				body:       jsonMarshal(User{ID: 123, Name: "kalt"}),
			},
		},
		{
			route: Route{HttpMethod: http.MethodPost, Path: "/user"},
			responseRuleOptions: []ResponseRuleOption{
				StatusCode(http.StatusCreated),
				XMLBody(Book{ISBN: "223-123", Name: "n0 n4m3"}),
				Header(http.Header{"Content-Type": []string{"application/xml"}}),
			},
			expectedResponse: ExpectedResponse{
				statusCode: http.StatusCreated,
				header:     http.Header{"Content-Type": []string{"application/xml"}},
				body:       xmlMarshal(Book{ISBN: "223-123", Name: "n0 n4m3"}),
			},
		},
		{
			route: Route{HttpMethod: http.MethodGet, Path: "/hi"},
			responseRuleOptions: []ResponseRuleOption{
				StatusCode(http.StatusOK),
				StringBody("Hello"),
			},
			expectedResponse: ExpectedResponse{
				statusCode: http.StatusOK,
				header:     http.Header{},
				body:       []byte("Hello"),
			},
		},
		{
			route: Route{HttpMethod: http.MethodGet, Path: "/community/best"},
			responseRuleOptions: []ResponseRuleOption{
				StatusCode(http.StatusTeapot),
				ByteBody([]byte{'S', 'T', 'R', 'E', 'E', 'T', ' ', 'B', 'Y', 'T', 'E', 'R', 'S'}),
			},
			expectedResponse: ExpectedResponse{
				statusCode: http.StatusTeapot,
				header:     http.Header{},
				body:       []byte{'S', 'T', 'R', 'E', 'E', 'T', ' ', 'B', 'Y', 'T', 'E', 'R', 'S'},
			},
		},
	}

	for _, test := range tests {
		server, _ := NewServer(test.route.HttpMethod, test.route.Path, test.responseRuleOptions...)
		defer server.Close()
		testRouteResponse(t, server.URL, test.route, test.expectedResponse)
	}
}

type bodyAssertFunc func(interface{}, Body) (bool, error)

func TestServerRequestRecorderBody(t *testing.T) {
	tests := []struct {
		route          Route
		request        func(string) *http.Request
		bodyAssertFunc bodyAssertFunc
		expectedBody   interface{}
	}{
		{
			route: Route{http.MethodPost, "/user"},
			request: func(url string) *http.Request {
				return newJSONRequest(http.MethodPost, url+"/user", User{ID: 133, Name: "Ken"})
			},
			bodyAssertFunc: isJSONEqual,
			expectedBody:   User{ID: 133, Name: "Ken"},
		},
		{
			route: Route{http.MethodPost, "/book"},
			request: func(url string) *http.Request {
				return newXMLRequest(http.MethodPost, url+"/book", Book{ISBN: "123-321-123", Name: "SICP"})
			},
			bodyAssertFunc: isXMLEqual,
			expectedBody:   Book{ISBN: "123-321-123", Name: "SICP"},
		},
	}

	for _, test := range tests {
		server, requestRecorder := NewServer(test.route.HttpMethod, test.route.Path)
		defer server.Close()

		_, err := http.DefaultClient.Do(test.request(server.URL))
		assert.Nil(t, err)

		testRouteRequestRecorderBody(t, test.expectedBody, requestRecorder, test.bodyAssertFunc)
	}
}

func TestServerWithTimeout(t *testing.T) {
	server, _ := NewServer(http.MethodGet, "/user", Timeout(20*time.Millisecond))
	defer server.Close()

	req := newJSONRequest(http.MethodGet, server.URL+"/user", http.NoBody)

	client := http.Client{
		Timeout: 5 * time.Millisecond,
	}

	_, err := client.Do(req)
	assert.NotNil(t, err)
}

func TestMultiRouteServerResponse(t *testing.T) {
	tests := []struct {
		routeResponseRuleOptions map[Route][]ResponseRuleOption
		expectedRouteResponses   map[Route]ExpectedResponse
	}{
		{
			routeResponseRuleOptions: map[Route][]ResponseRuleOption{
				{http.MethodGet, "/user"}: {
					StatusCode(http.StatusOK),
					JSONBody(User{ID: 123, Name: "kalt"}),
					Header(http.Header{"Content-Type": []string{"application/json"}}),
				},
				{http.MethodGet, "/book"}: {
					StatusCode(http.StatusTeapot),
					Header(http.Header{"Content-Type": []string{"application/xml"}}),
					XMLBody(Book{ISBN: "9780262510875", Name: "Structure and Interpretation of Computer Programs"}),
				},
			},
			expectedRouteResponses: map[Route]ExpectedResponse{
				{http.MethodGet, "/user"}: {
					statusCode: http.StatusOK,
					header:     http.Header{"Content-Type": []string{"application/json"}},
					body:       jsonMarshal(User{ID: 123, Name: "kalt"}),
				},
				{http.MethodGet, "/book"}: {
					statusCode: http.StatusTeapot,
					header:     http.Header{"Content-Type": []string{"application/xml"}},
					body:       xmlMarshal(Book{ISBN: "9780262510875", Name: "Structure and Interpretation of Computer Programs"}),
				},
			},
		},
	}

	for _, test := range tests {
		server, _ := NewMultiRouteServer(test.routeResponseRuleOptions)
		defer server.Close()

		for route := range test.routeResponseRuleOptions {
			expectedResponse := test.expectedRouteResponses[route]
			testRouteResponse(t, server.URL, route, expectedResponse)
		}
	}
}

func TestMultiRouteServerRequestRecorderBody(t *testing.T) {
	tests := []struct {
		routes                             []Route
		requests                           func(string) []*http.Request
		routeBodyAssertFunc                map[Route]bodyAssertFunc
		expectedRouteRequestRecorderBodies map[Route]interface{}
	}{
		{
			routes: []Route{
				{http.MethodPost, "/user"},
				{http.MethodPost, "/book"},
			},
			requests: func(url string) []*http.Request {
				reqs := []*http.Request{
					newJSONRequest(http.MethodPost, url+"/user", User{ID: 1222, Name: "nonono"}),
					newXMLRequest(http.MethodPost, url+"/book", Book{ISBN: "123-321-123", Name: "SICP"}),
				}
				return reqs
			},
			routeBodyAssertFunc: map[Route]bodyAssertFunc{
				{http.MethodPost, "/user"}: isJSONEqual,
				{http.MethodPost, "/book"}: isXMLEqual,
			},
			expectedRouteRequestRecorderBodies: map[Route]interface{}{
				{http.MethodPost, "/user"}: User{ID: 1222, Name: "nonono"},
				{http.MethodPost, "/book"}: Book{ISBN: "123-321-123", Name: "SICP"},
			},
		},
	}

	for _, test := range tests {
		routeResponseRules := make(map[Route][]ResponseRuleOption)
		for _, route := range test.routes {
			routeResponseRules[route] = []ResponseRuleOption{}
		}

		server, requestRecorder := NewMultiRouteServer(routeResponseRules)
		defer server.Close()

		for _, request := range test.requests(server.URL) {
			_, err := http.DefaultClient.Do(request)
			assert.Nil(t, err)
		}

		for _, route := range test.routes {
			bodyAssertFunc := test.routeBodyAssertFunc[route]
			expectedBody := test.expectedRouteRequestRecorderBodies[route]

			testRouteRequestRecorderBody(t, expectedBody, requestRecorder[route], bodyAssertFunc)
		}
	}
}

func testRouteResponse(t *testing.T, serverURL string, route Route, expectedResponse ExpectedResponse) {
	request, err := http.NewRequest(route.HttpMethod, serverURL+route.Path, http.NoBody)
	assert.Nil(t, err)

	response, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse.statusCode, response.StatusCode)

	actualBody, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse.body, actualBody)

	assert.True(t, isHeaderContains(expectedResponse.header, response.Header))
}

func testRouteRequestRecorderBody(t *testing.T, expectedBody interface{}, requestRecorder *RequestRecorder, bodyAssertFunc bodyAssertFunc) {
	isBodyEqual, err := bodyAssertFunc(expectedBody, requestRecorder.Body)
	assert.Nil(t, err)
	assert.True(t, isBodyEqual)
}

func jsonMarshal(j interface{}) []byte {
	m, _ := json.Marshal(j)
	return m
}

func xmlMarshal(x interface{}) []byte {
	m, _ := xml.Marshal(x)
	return m
}
