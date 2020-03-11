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
	"fmt"
	"net/http"
	"testing"

	"github.com/clbanning/mxj"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

func (r RequestRecorder) AssertStringBodyEqual(t *testing.T, expectedBody string) bool {
	return assert.Equal(t, expectedBody, string(r.Data))
}

func (r RequestRecorder) AssertJSONBodyEqual(t *testing.T, expectedBody interface{}) bool {
	isEqual, err := isJSONEqual(expectedBody, r.Body)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	if isEqual {
		return true
	}

	failTest(t, "JSON Bodies are not equal!", func() {
		printJSONDiff(expectedBody, r.Body)
	})
	return false
}

func (r RequestRecorder) AssertXMLBodyEqual(t *testing.T, expectedXMLBody interface{}) bool {
	isEqual, err := isXMLEqual(expectedXMLBody, r.Body)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	if isEqual {
		return true
	}

	failTest(t, "XML Bodies are not equal!")
	return false
}

func (r RequestRecorder) AssertParamEqual(t *testing.T, paramName, paramValue string) bool {
	if assert.ObjectsAreEqual(r.Params[paramName], paramValue) {
		return true
	}

	failTest(t, fmt.Sprintf("Param name '%s' is not equal to '%s'", paramName, paramValue), func() {
		color.Red("Actual:\t  %s", r.Params[paramName])
		color.Yellow("Expected: %s", paramValue)
	})
	return false
}

func (r RequestRecorder) AssertQueryParamEqual(t *testing.T, queryParamName string, queryParamValues []string) bool {
	if assert.ObjectsAreEqual(r.QueryParams[queryParamName], queryParamValues) {
		return true
	}

	failTest(t, fmt.Sprintf("QueryParam name '%s' is not equal to %v", queryParamName, queryParamValues), func() {
		color.Red("Actual:\t  %s", r.QueryParams[queryParamName])
		color.Yellow("Expected: %s", queryParamValues)
	})
	return false
}

func (r RequestRecorder) AssertFormParamEqual(t *testing.T, formParamName string, formValues []string) bool {
	if assert.ObjectsAreEqual(r.FormParams[formParamName], formValues) {
		return true
	}

	failTest(t, fmt.Sprintf("FormParam name '%s' is not equal to %v", formParamName, formValues), func() {
		color.Red("Actual:\t  %s", r.FormParams[formParamName])
		color.Yellow("Expected: %s", formValues)
	})
	return false
}

func (r RequestRecorder) AssertHeaderEqual(t *testing.T, expectedHeader http.Header) bool {
	if isHeaderContains(expectedHeader, r.Header) {
		return true
	}

	failTest(t, "HTTP Headers are not equal", func() {
		printJSONDiff(expectedHeader, r.Header)
	})
	return false
}

func isHeaderContains(expectedHeader, actualHeader http.Header) bool {
	for key, value := range expectedHeader {
		actualValue, contains := actualHeader[key]
		if !contains {
			return false
		}

		if !assert.ObjectsAreEqualValues(value, actualValue) {
			return false
		}
	}
	return true
}

func isJSONEqual(expectedBody interface{}, actualBody Body) (bool, error) {
	bodyJSON, err := json.Marshal(expectedBody)
	if err != nil {
		return false, err
	}

	expectedRecorderBody := Body{}
	if err := json.Unmarshal(bodyJSON, &expectedRecorderBody); err != nil {
		return false, err
	}
	return assert.ObjectsAreEqualValues(expectedRecorderBody, actualBody), nil
}

func isXMLEqual(expectedBody interface{}, actualBody Body) (bool, error) {
	bodyXML, err := xml.Marshal(expectedBody)
	if err != nil {
		return false, err
	}

	mv, err := mxj.NewMapXml(bodyXML)
	if err != nil {
		return false, err
	}

	expectedRecorderBody := mv.Old()

	return assert.ObjectsAreEqualValues(expectedRecorderBody, actualBody), nil
}

func failTest(t *testing.T, message string, callback ...func()) {
	color.Red(message)
	if len(callback) > 0 {
		callback[0]()
	}
	t.Fail()
}

func printJSONDiff(expected, actual interface{}) {
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)

	differ := diff.New()
	diff, _ := differ.Compare(expectedJSON, actualJSON)

	var asciiFormatMap map[string]interface{}
	json.Unmarshal(expectedJSON, &asciiFormatMap)

	formatter := formatter.NewAsciiFormatter(asciiFormatMap, formatter.AsciiFormatterConfig{
		Coloring: true,
	})

	diffString, _ := formatter.Format(diff)

	color.Yellow("Difference:")
	fmt.Println(diffString)
}
