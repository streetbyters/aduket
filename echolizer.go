package echolizer

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/clbanning/mxj"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type RequestRecorder struct {
	Body        Body
	Params      map[string]string
	QueryParams url.Values
	FormParams  url.Values
}

type Body map[string]interface{}

func (b Body) IsJSONEqual(body interface{}) (bool, error) {
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

func (b Body) IsXMLEqual(body interface{}) (bool, error) {
	bodyXML, err := xml.Marshal(body)
	if err != nil {
		return false, err
	}

	mv, err := mxj.NewMapXml(bodyXML)
	if err != nil {
		return false, err
	}

	expectedRecorderBody := mv.Old()

	return assert.ObjectsAreEqualValues(b, expectedRecorderBody), nil
}

func NewRequestRecorder() *RequestRecorder {
	requestRecorder := &RequestRecorder{}
	requestRecorder.Body = make(Body)
	requestRecorder.Params = make(map[string]string)
	return requestRecorder
}

func (r RequestRecorder) AssertBodyEqual(t *testing.T, expectedBody interface{}) bool {
	bodyType := getTagKey(expectedBody)

	isEqualFunc := r.getIsEqualFunctionByContentType(bodyType)

	isEqual, err := isEqualFunc(expectedBody)
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

func (r RequestRecorder) getIsEqualFunctionByContentType(bodyType string) func(body interface{}) (bool, error) {
	switch bodyType {
	case "xml":
		return r.Body.IsXMLEqual
	case "json":
		return r.Body.IsJSONEqual
	default:
		return func(body interface{}) (bool, error) {
			return false, errors.New("Unsupported body type")
		}
	}
}

func (r *RequestRecorder) bindXML(from io.ReadCloser) error {
	body, err := ioutil.ReadAll(from)
	if err != nil {
		return err
	}

	mv, err := mxj.NewMapXml(body)
	if err != nil {
		return err
	}

	r.Body = mv.Old()

	return nil
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
		if ctx.Request().Header.Get(echo.HeaderContentType) == echo.MIMEApplicationXML {
			requestRecorder.bindXML(ctx.Request().Body)
		} else {
			ctx.Bind(&requestRecorder.Body)
		}

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
