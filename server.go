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
	"time"

	"github.com/labstack/echo"
)

type responseBody []byte

type Route struct {
	HttpMethod string
	Path       string
}

type responseRule struct {
	header            http.Header
	body              responseBody
	statusCode        int
	timeout           time.Duration
	sendCorruptedBody bool
}

func NewMultiRouteServer(routeResponseOptions map[Route][]ResponseRuleOption) (*httptest.Server, map[Route]*RequestRecorder) {
	requestRecorder := make(map[Route]*RequestRecorder)
	e := createEcho()

	routeResponseRules := createRouteResponseRules(routeResponseOptions)
	for route, responseRule := range routeResponseRules {
		routeRequestRecorder := NewRequestRecorder()
		requestRecorder[route] = routeRequestRecorder
		e.Add(route.HttpMethod, route.Path, spyHandler(routeRequestRecorder, responseRule))
	}

	return httptest.NewServer(e), requestRecorder
}

func NewServer(httpMethod, path string, responseRuleOptions ...ResponseRuleOption) (*httptest.Server, *RequestRecorder) {
	requestRecorder := NewRequestRecorder()

	e := createEcho()
	responseRule := createResponseRule(responseRuleOptions)

	e.Add(httpMethod, path, spyHandler(requestRecorder, responseRule))
	return httptest.NewServer(e), requestRecorder
}

func createEcho() *echo.Echo {
	e := echo.New()
	e.Binder = &RequestRecorderBinder{}
	return e
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
	responseRule := &responseRule{statusCode: http.StatusOK}

	for _, responseRuleOption := range responseRuleOptions {
		responseRuleOption(responseRule)
	}

	return *responseRule
}

type RequestRecorderBinder struct{}

func (r *RequestRecorderBinder) Bind(requestRecorder interface{}, ctx echo.Context) error {
	recorder := requestRecorder.(*RequestRecorder)
	return recorder.saveContext(ctx)
}

func spyHandler(requestRecorder *RequestRecorder, res responseRule) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestRecorder.isRequestReceived = true

		if res.sendCorruptedBody {
			// Forces client to read empty buffer and BOOM!
			ctx.Response().Header().Set("Content-Length", "1")
			return nil
		}

		if res.timeout != 0 {
			time.Sleep(res.timeout)
		}

		if err := ctx.Bind(requestRecorder); err != nil {
			return err
		}

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
