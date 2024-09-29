package util

import (
	"food_delivery/request"
	"net/http"
	"net/http/httptest"
	"strings"
)

type TestCaseRegisterValidation struct {
	Name         string
	Req          *request.RegisterRequest
	WantError    bool
	WantErrorMsg string
}

type TestCaseCartHandlerGetCart struct {
	TestName    string
	Request     *Request
	HandlerFunc http.HandlerFunc
	Want        *ExpectedResponse // Expected Response
}

type ExpectedResponse struct {
	StatusCode int
	BodyPart   string
}

type Request struct {
	Method      string
	Url         string
	AccessToken string
}

func PreparHandlerTestCases(test *TestCaseCartHandlerGetCart) (request *http.Request, recorder *httptest.ResponseRecorder) {
	request = httptest.NewRequest(test.Request.Method, test.Request.Url, strings.NewReader(""))

	if test.Request.AccessToken != "" {
		request.Header.Set("Authorization", "Bearer "+test.Request.AccessToken)
	}

	return request, httptest.NewRecorder()
}
