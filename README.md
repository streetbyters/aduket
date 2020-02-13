<h1 align="center">Aduket</h1>
<p align="center">
    <img height="200px" src="assets/aduket.gif">
</p>
<p align="center">
    <i>Straight-forward HTTP client testing, assertions included!</i>
</p>

Simple `httptest.Server` wrapper with a little request recorder spice on it. No special DSL, no complex API to learn. Just create a server and fire your request like an **Hadouken** then assert it.

## TODO
 - [x] Add support for accessing to each RequestRecorder of Route of MultiRouteServer
 - [ ] Extract Request().Body to requestRecorder.Body binding logic to CustomBinder
 - [ ] Add badges to README.md
 - [ ] Add NewServerWithTimeout for testing API timeouts
 - [ ] http.RoundTripper interface can be implemented to mock arbitrary URLs
 - [ ] Add multiple request assertion logic
 - [ ] Add example usages
 - [ ] Add docs
 - [ ] Add meaningful error messages to assertions
 - [x] Add response headers to NewServer
 - [x] Add request header assertions
 - [x] Add multiple route registration to server

## LICENSE
Copyright 2020 StreetByters Community

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.