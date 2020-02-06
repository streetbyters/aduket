<p align="center">
    <b>Aduket</b>
    <br>
    <img height="200px" src="assets/aduket.gif">

</p>

<p align="center">
    <i>Straight-forward HTTP client testing, assertions included!</i>
</p>

Simple `httptest.Server` wrapper with a little request recorder spice on it. No special DSL, no complex API to learn. Just create a server and fire your request like an **Hadouken** then assert them.


### Todo:
* Add example usages
* Add docs
* Add response headers to NewServer
* Add request header assertions
* Add multiple request assertion logic
* Extract Request().Body to requestRecorder.Body binding logic to CustomBinder
* Add NewServerWithTimeout for testing API timeouts
* http.RoundTripper interface can be implemented to mock arbitrary URLs
* A Builder can be written to NewServer for ease of use
