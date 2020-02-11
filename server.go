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
	"net/http"
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

type route struct {
	httpMethod string
	path       string
}

type routeResponse struct {
	statusCode int
	header     http.Header
	body       responseBody
}

type server struct {
	*httptest.Server
	routes map[route]routeResponse
}

type RouteResponseOption func(*routeResponse)

func StatusCode(statusCode int) RouteResponseOption {
	return func(r *routeResponse) {
		r.statusCode = statusCode
	}
}

func JSONBody(body interface{}) RouteResponseOption {
	return func(r *routeResponse) {
		r.body = JSONResponse(body)
	}
}

func Header(header http.Header) RouteResponseOption {
	return func(r *routeResponse) {
		r.header = header
	}
}

func NewServe(responseRules map[route][]RouteResponseOption) (*server, *RequestRecorder) {

	routes := make(map[route]routeResponse)

	for route, responseOption := range responseRules {
		routeResponse := &routeResponse{}
		for _, opt := range responseOption {
			opt(routeResponse)
		}
		routes[route] = *routeResponse
	}

	requestRecorder := NewRequestRecorder()
	e := echo.New()

	for route, responseRule := range routes {
		e.Add(route.httpMethod, route.path, func(ctx echo.Context) error {
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

			for key, values := range responseRule.header {
				for _, value := range values {
					ctx.Response().Header().Add(key, value)
				}
			}

			if responseRule.body == nil {
				return ctx.NoContent(responseRule.statusCode)
			}

			ctx.Response().WriteHeader(responseRule.statusCode)
			_, err := ctx.Response().Write(responseRule.body)
			if err != nil {
				return err
			}

			return nil
		})
	}
	s := &server{httptest.NewServer(e), routes}

	return s, requestRecorder
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
