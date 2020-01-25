package echolizer

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ffmt.v1"
)

type RequestRecorder struct {
	Body        Body
	Params      map[string]string
	QueryParams url.Values
	FormParams  url.Values
}

type Body map[string]interface{}

func (b Body) IsEqual(body interface{}) (bool, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return false, err
	}

	expectedRecorderBody := Body{}
	if err := json.Unmarshal(bodyJSON, &expectedRecorderBody); err != nil {
		return false, err
	}

	return assert.ObjectsAreEqual(expectedRecorderBody, b), nil
}

func NewRequestRecorder() *RequestRecorder {
	requestRecorder := &RequestRecorder{}
	requestRecorder.Body = make(Body)
	requestRecorder.Params = make(map[string]string)
	return requestRecorder
}

func (r RequestRecorder) AssertBodyEqual(t *testing.T, expectedBody interface{}) bool {

	contentType := getTagKey(expectedBody)
	ffmt.Puts(contentType)

	isEqual, err := r.Body.IsEqual(expectedBody)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	return assert.True(t, isEqual)
}

func (r RequestRecorder) AssertParamEqual(t *testing.T, paramName, paramValue string) bool {
	return assert.Equal(t, r.Params[paramName], paramValue)
}

func (r RequestRecorder) AssertQueryParamEqual(t *testing.T, queryParamName string, queryParamValues []string) bool {
	return assert.Equal(t, r.QueryParams[queryParamName], queryParamValues)
}

func (r RequestRecorder) AssertFormParamEqual(t *testing.T, formParamName string, formValues []string) bool {
	return assert.Equal(t, r.FormParams[formParamName], formValues)
}

func (r *RequestRecorder) setQueryParams(queryParams url.Values) {
	r.QueryParams = queryParams
}

func (r *RequestRecorder) setParams(paramNames, paramValues []string) {
	for index, name := range paramNames {
		r.Params[name] = paramValues[index]
	}
}

func (r *RequestRecorder) setFormParams(formParams url.Values) {
	r.FormParams = formParams
}

func NewEcholizer(httpMethod, path string, statusCode int) (*httptest.Server, *RequestRecorder) {
	requestRecorder := NewRequestRecorder()
	e := createEcho(requestRecorder, httpMethod, path, statusCode, nil)
	return httptest.NewServer(e), requestRecorder
}

func NewEcholizerWithResponse(httpMethod, path string, statusCode int, response interface{}) (*httptest.Server, *RequestRecorder) {
	requestRecorder := NewRequestRecorder()
	e := createEcho(requestRecorder, httpMethod, path, statusCode, response)
	return httptest.NewServer(e), requestRecorder
}

func createEcho(requestRecorder *RequestRecorder, httpMethod, path string, statusCode int, response interface{}) *echo.Echo {
	e := echo.New()

	e.Add(httpMethod, path, func(ctx echo.Context) error {
		ctx.Bind(&requestRecorder.Body)

		ffmt.Puts(requestRecorder.Body)

		requestRecorder.setParams(ctx.ParamNames(), ctx.ParamValues())
		requestRecorder.setQueryParams(ctx.QueryParams())
		requestRecorder.setFormParams(ctx.Request().Form)

		if response == nil {
			return ctx.NoContent(statusCode)
		}

		return ctx.JSON(statusCode, response)
	})

	return e
}

func getTagKey(i interface{}) string {
	return strings.Split(string(reflect.TypeOf(i).Field(0).Tag), ":")[0]
}
