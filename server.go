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
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
)

type responseBody []byte

type Route struct {
	httpMethod string
	path       string
}

type response struct {
	header     http.Header
	body       responseBody
	statusCode int
}

type responseRule struct {
	header     http.Header
	body       responseBody
	statusCode int
}

func NewMultiRouteServer(routeResponseOptions map[Route][]ResponseRuleOption) (*httptest.Server, *RequestRecorder) {
	requestRecorder := NewRequestRecorder()
	e := echo.New()

	routeResponseRules := createRouteResponseRules(routeResponseOptions)
	for route, responseRule := range routeResponseRules {
		e.Add(route.httpMethod, route.path, spyHandler(requestRecorder, response{responseRule.header, responseRule.body, responseRule.statusCode}))
	}
	return httptest.NewServer(e), requestRecorder
}

func NewServer(httpMethod, path string, responseRuleOptions ...ResponseRuleOption) (*httptest.Server, *RequestRecorder) {
	requestRecorder := NewRequestRecorder()
	e := echo.New()

	responseRule := createResponseRule(responseRuleOptions)

	e.Add(httpMethod, path, spyHandler(requestRecorder, response{responseRule.header, responseRule.body, responseRule.statusCode}))
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

func spyHandler(requestRecorder *RequestRecorder, res response) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestRecorder.saveContext(ctx)

		for key, values := range res.header {
			for _, value := range values {
				ctx.Response().Header().Add(key, value)
			}
		}

		if res.body == nil {
			return ctx.NoContent(res.statusCode)
		}

		ctx.Response().WriteHeader(res.statusCode)
		_, err := ctx.Response().Write(res.body)
		if err != nil {
			return err
		}

		return nil
	}
}
