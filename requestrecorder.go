// Copyright 2020 StreetByters Community
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aduket

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/clbanning/mxj"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type RequestRecorder struct {
	Body        Body
	Header      http.Header
	Data        []byte
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

	return assert.ObjectsAreEqualValues(expectedRecorderBody, b), nil
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

func (r RequestRecorder) AssertStringBodyEqual(t *testing.T, expectedBody string) bool {
	return assert.Equal(t, expectedBody, string(r.Data))
}

func (r RequestRecorder) AssertJSONBodyEqual(t *testing.T, expectedBody interface{}) bool {
	isEqual, err := r.Body.IsJSONEqual(expectedBody)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	return assert.True(t, isEqual)
}

func (r RequestRecorder) AssertXMLBodyEqual(t *testing.T, expectedXMLBody interface{}) bool {
	isEqual, err := r.Body.IsXMLEqual(expectedXMLBody)
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

func (r RequestRecorder) AssertHeaderEqual(t *testing.T, expectedHeader http.Header) bool {
	for key, value := range expectedHeader {
		actualValue, contains := r.Header[key]
		if !assert.True(t, contains) {
			return false
		}

		if !assert.Equal(t, value, actualValue) {
			return false
		}
	}

	return true
}

func (r *RequestRecorder) saveContext(ctx echo.Context) error {
	if ctx.Request().Header.Get(echo.HeaderContentType) == echo.MIMEApplicationXML {
		r.bindXML(ctx.Request().Body)
	} else if err := ctx.Bind(&r.Body); err != nil {
		data, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			return err
		}
		r.setData(data)
	}

	r.setParams(ctx.ParamNames(), ctx.ParamValues())
	r.setQueryParams(ctx.QueryParams())
	r.setFormParams(ctx.Request().Form)
	r.setHeader(ctx.Request().Header)

	return nil

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

func (r *RequestRecorder) setData(b []byte) {
	r.Data = b
}

func (r *RequestRecorder) setHeader(header http.Header) {
	r.Header = header.Clone()
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
