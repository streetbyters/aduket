// Copyright 2019 StreetByters Community
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
	"io/ioutil"
	"net/http/httptest"

	"github.com/labstack/echo"
)

type responseBody []byte

func JSONResponse(j interface{}) responseBody {
	jsonBytes, _ := json.Marshal(j)
	return jsonBytes
}

func XMLResponse(x interface{}) responseBody {
	xmlBytes, _ := xml.Marshal(x)
	return xmlBytes
}

func StringResponse(s string) responseBody {
	return []byte(s)
}

func ByteResponse(b []byte) responseBody {
	return b
}

func NoResponse() responseBody {
	return nil
}

func NewServer(httpMethod, path string, statusCode int, response responseBody) (*httptest.Server, *RequestRecorder) {
	requestRecorder := NewRequestRecorder()
	e := createEcho(requestRecorder, httpMethod, path, statusCode, response)
	return httptest.NewServer(e), requestRecorder
}

func createEcho(requestRecorder *RequestRecorder, httpMethod, path string, statusCode int, body responseBody) *echo.Echo {
	e := echo.New()

	e.Add(httpMethod, path, func(ctx echo.Context) error {

		if ctx.Request().Header.Get(echo.HeaderContentType) == echo.MIMEApplicationXML {
			requestRecorder.bindXML(ctx.Request().Body)
		} else if err := ctx.Bind(&requestRecorder.Body); err != nil {
			data, err := ioutil.ReadAll(ctx.Request().Body)
			if err != nil {
				return err
			}
			requestRecorder.setData(data)
		}

		requestRecorder.setParams(ctx.ParamNames(), ctx.ParamValues())
		requestRecorder.setQueryParams(ctx.QueryParams())
		requestRecorder.setFormParams(ctx.Request().Form)

		if body == nil {
			return ctx.NoContent(statusCode)
		}

		ctx.Response().WriteHeader(statusCode)
		_, err := ctx.Response().Write(body)
		if err != nil {
			return err
		}

		return nil
	})

	return e
}
