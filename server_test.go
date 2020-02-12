package aduket

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"testing"

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

func TestServer(t *testing.T) {
	serverTests := []struct {
		route               Route
		responseRuleOptions []ResponseRuleOption
		expectedResponse    ExpectedResponse
	}{
		{
			route: Route{httpMethod: http.MethodGet, path: "/user"},
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
			route: Route{httpMethod: http.MethodPost, path: "/user"},
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
			route: Route{httpMethod: http.MethodGet, path: "/hi"},
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
			route: Route{httpMethod: http.MethodGet, path: "/community/best"},
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

	for _, test := range serverTests {
		server, _ := NewServer(test.route.httpMethod, test.route.path, test.responseRuleOptions...)
		defer server.Close()
		testRoute(t, server.URL, test.route, test.expectedResponse)
	}
}

func TestMultiRouteServer(t *testing.T) {

	multiRouteServerTests := []struct {
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

	for _, test := range multiRouteServerTests {
		server, _ := NewMultiRouteServer(test.routeResponseRuleOptions)
		defer server.Close()

		for route := range test.routeResponseRuleOptions {
			expectedResponse := test.expectedRouteResponses[route]
			testRoute(t, server.URL, route, expectedResponse)
		}
	}
}

func testRoute(t *testing.T, serverURL string, route Route, expectedResponse ExpectedResponse) {
	request, err := http.NewRequest(route.httpMethod, serverURL+route.path, http.NoBody)
	assert.Nil(t, err)

	response, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse.statusCode, response.StatusCode)

	actualBody, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse.body, actualBody)

	assertHeaderContains(t, expectedResponse.header, response.Header)
}

func jsonMarshal(j interface{}) []byte {
	m, _ := json.Marshal(j)
	return m
}

func xmlMarshal(x interface{}) []byte {
	m, _ := xml.Marshal(x)
	return m
}

func assertHeaderContains(t *testing.T, expectedHeader, actualHeader http.Header) bool {
	for key, value := range expectedHeader {
		actualValue, contains := actualHeader[key]
		if !assert.True(t, contains) {
			return false
		}
		if !assert.Equal(t, value, actualValue) {
			return false
		}
	}
	return true
}
