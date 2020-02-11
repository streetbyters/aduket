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
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
)

type responseBody []byte

type Route struct {
	httpMethod string
	path       string
}

type responseRule struct {
	statusCode int
	header     http.Header
	body       responseBody
}

func NewMultiRouteServer(routeResponseOptions map[Route][]ResponseRuleOption) (*httptest.Server, *RequestRecorder) {
	requestRecorder := NewRequestRecorder()
	e := echo.New()

	routeResponseRules := createRouteResponseRules(routeResponseOptions)
	for route, responseRule := range routeResponseRules {
		e.Add(route.httpMethod, route.path, spyHandler(requestRecorder, responseRule.header, responseRule.body, responseRule.statusCode))
	}

	return httptest.NewServer(e), requestRecorder
}

func NewServer(httpMethod, path string, responseRuleOptions ...ResponseRuleOption) (*httptest.Server, *RequestRecorder) {
	requestRecorder := NewRequestRecorder()
	e := echo.New()

	responseRule := createResponseRule(responseRuleOptions)

	e.Add(httpMethod, path, spyHandler(requestRecorder, responseRule.header, responseRule.body, responseRule.statusCode))
	return httptest.NewServer(e), requestRecorder
}

func createRouteResponseRules(routeResponseOptions map[Route][]ResponseRuleOption) map[Route]responseRule {
	routeResponseRules := make(map[Route]responseRule)
	for route, responseOption := range routeResponseOptions {
		rule := createResponseRule(responseOption)
		routeResponseRules[route] = rule
	}

	return routeResponseRules
}

func createResponseRule(responseRuleOptions []ResponseRuleOption) responseRule {
	responseRule := &responseRule{}
	for _, responseRuleOption := range responseRuleOptions {
		responseRuleOption(responseRule)
	}

	return *responseRule
}

func spyHandler(requestRecorder *RequestRecorder, header http.Header, body responseBody, statusCode int) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		contextToRequestRecorder(ctx, requestRecorder)

		for key, values := range header {
			for _, value := range values {
				ctx.Response().Header().Add(key, value)
			}
		}

		if body == nil {
			return ctx.NoContent(statusCode)
		}

		ctx.Response().WriteHeader(statusCode)
		_, err := ctx.Response().Write(body)
		if err != nil {
			return err
		}

		return nil
	}
}

func contextToRequestRecorder(ctx echo.Context, requestRecorder *RequestRecorder) error {
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
	requestRecorder.setHeader(ctx.Request().Header)

	return nil
}
