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
// limitations under the License

package aduket

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (r RequestRecorder) AssertStringBodyEqual(t *testing.T, expectedBody string) bool {
	return assert.Equal(t, expectedBody, string(r.Data))
}

func (r RequestRecorder) AssertJSONBodyEqual(t *testing.T, expectedBody interface{}) bool {
	expectedBody, err := json.Marshal(expectedBody)
	if err != nil {
		t.Error("expected body could not marshaled to json")
	}
	return assert.Equal(t, expectedBody, r.Body)
}

func (r RequestRecorder) AssertXMLBodyEqual(t *testing.T, expectedXMLBody interface{}) bool {
	expectedBody, err := xml.Marshal(expectedXMLBody)
	if err != nil {
		t.Error("expected body could not marshaled to xml")
	}
	return assert.Equal(t, expectedBody, r.Body)
}

func (r RequestRecorder) AssertParamEqual(t *testing.T, paramName, paramValue string) bool {
	return assert.Equal(t, paramValue, r.Params[paramName])
}

func (r RequestRecorder) AssertQueryParamEqual(t *testing.T, queryParamName string, queryParamValues []string) bool {
	return assert.Equal(t, queryParamValues, r.QueryParams[queryParamName])
}

func (r RequestRecorder) AssertFormParamEqual(t *testing.T, formParamName string, formValues []string) bool {
	return assert.Equal(t, formValues, r.FormParams[formParamName])
}

func (r RequestRecorder) AssertHeaderContains(t *testing.T, expectedHeader http.Header) bool {
	return assert.True(t, isHeaderContains(expectedHeader, r.Header))
}

func (r RequestRecorder) AssertNoRequest(t *testing.T) bool {
	return assert.False(t, r.isRequestReceived)
}

func isHeaderContains(expectedHeader, actualHeader http.Header) bool {
	assertionResult := true
	for key, value := range expectedHeader {
		headerValue := actualHeader.Values(key)
		assertionResult = assertionResult && assert.ObjectsAreEqualValues(headerValue, value)
	}
	return assertionResult
}

func isJSONEqual(expectedBody interface{}, actualBody Body) (bool, error) {
	expectedBytes, err := json.Marshal(expectedBody)
	return assert.ObjectsAreEqual(expectedBytes, actualBody), err
}

func isXMLEqual(expectedBody interface{}, actualBody Body) (bool, error) {
	expectedBytes, err := xml.Marshal(expectedBody)
	return assert.ObjectsAreEqualValues(expectedBytes, actualBody), err
}
