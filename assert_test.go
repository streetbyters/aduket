package aduket

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var dummyRequest = httptest.NewRequest(http.MethodGet, "http://streetbyters.com", http.NoBody)

func TestAssertJSONBodyEqual(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	expectedPayload := User{Name: "noname"}
	request := newJSONRequest(http.MethodPost, "", expectedPayload)

	ctx := echo.New().NewContext(request, nil)

	requestRecorder := NewRequestRecorder()
	requestRecorder.saveContext(ctx)

	// tester := &testing.T{}
	requestRecorder.AssertJSONBodyEqual(t, User{Name: "noname"})
	// assert.True(t, )
	// assert.False(t, tester.Failed())

	// assert.False(t, requestRecorder.AssertJSONBodyEqual(tester, User{Name: "lel"}))
	// assert.True(t, tester.Failed())
}

func TestAssertStringBodyEqual(t *testing.T) {
	expectedPayload := "Hello"
	request := newStringRequest(http.MethodPost, "", expectedPayload)

	ctx := echo.New().NewContext(request, nil)

	requestRecorder := NewRequestRecorder()
	requestRecorder.saveContext(ctx)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertStringBodyEqual(tester, expectedPayload))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertStringBodyEqual(tester, "olleH"))
	assert.True(t, tester.Failed())
}

func TestAssertXMLBodyEqual(t *testing.T) {
	type Book struct {
		Name string `xml:"name"`
	}

	expectedPayload := Book{Name: "noname"}

	request := newXMLRequest(http.MethodPost, "", expectedPayload)

	ctx := echo.New().NewContext(request, nil)

	requestRecorder := NewRequestRecorder()
	requestRecorder.saveContext(ctx)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertXMLBodyEqual(tester, expectedPayload))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertXMLBodyEqual(tester, Book{Name: "lel"}))
	assert.True(t, tester.Failed())
}

func TestAssertParamEqual(t *testing.T) {
	ctx := echo.New().NewContext(dummyRequest, nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues("123")

	requestRecorder := NewRequestRecorder()
	requestRecorder.saveContext(ctx)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertParamEqual(tester, "id", "123"))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertParamEqual(tester, "id", "321"))
	assert.True(t, tester.Failed())
}

func TestAssertQueryParamEqual(t *testing.T) {
	ctx := echo.New().NewContext(dummyRequest, nil)
	ctx.QueryParams().Add("name", "Joe")

	requestRecorder := NewRequestRecorder()
	requestRecorder.saveContext(ctx)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertQueryParamEqual(tester, "name", []string{"Joe"}))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertQueryParamEqual(tester, "name", []string{"Doe"}))
	assert.True(t, tester.Failed())
}

func TestAssertFormParamEqual(t *testing.T) {
	ctx := echo.New().NewContext(dummyRequest, nil)
	ctx.Request().Form = url.Values{"name": []string{"Joe"}}

	requestRecorder := NewRequestRecorder()
	requestRecorder.saveContext(ctx)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertFormParamEqual(tester, "name", []string{"Joe"}))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertFormParamEqual(tester, "name", []string{"Doe"}))
	assert.True(t, tester.Failed())
}

func TestAssertHeaderContains(t *testing.T) {
	ctx := echo.New().NewContext(dummyRequest, nil)
	ctx.Request().Header = http.Header{"Test": []string{"123"}}

	requestRecorder := NewRequestRecorder()
	requestRecorder.saveContext(ctx)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertHeaderContains(tester, http.Header{"Test": []string{"123"}}))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertHeaderContains(tester, http.Header{"Test": []string{"noo"}}))
	assert.False(t, requestRecorder.AssertHeaderContains(tester, http.Header{"West": []string{"123"}}))
}

func TestAssertNoRequest(t *testing.T) {
	requestRecorder := NewRequestRecorder()

	tester := &testing.T{}
	assert.True(t, requestRecorder.AssertNoRequest(tester))
	assert.False(t, tester.Failed())

	ctx := echo.New().NewContext(dummyRequest, nil)
	spyHandler(requestRecorder, responseRule{})(ctx)

	assert.False(t, requestRecorder.AssertNoRequest(tester))
	assert.True(t, tester.Failed())
}

func newStringRequest(method, url, body string) *http.Request {
	request, _ := http.NewRequest(method, url, strings.NewReader(body))
	return request
}

func newJSONRequest(method, url string, body interface{}) *http.Request {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")

	return request
}

func newXMLRequest(method, url string, body interface{}) *http.Request {
	requestBody, _ := xml.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/xml")

	return request
}
