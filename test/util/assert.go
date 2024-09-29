package util

import (
	"encoding/json"
	"food_delivery/model"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertRegisterValidationResult(t *testing.T, testCase TestCaseRegisterValidation, gotErr error) {
	t.Helper()

	if testCase.WantError {
		assert.Error(t, gotErr)
		assert.Equal(t, testCase.WantErrorMsg, gotErr.Error())
	} else {
		assert.NoError(t, gotErr)
	}
}

func AssertUserProfileResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	var r model.Cart
	err := json.Unmarshal(recorder.Body.Bytes(), &r)

	if assert.NoError(t, err) {
		assert.Equal(t, model.Cart{
			CartID: 1,
			Items: []model.CartItem{
				{
					CartId:     1,
					Name:       "pizza",
					Image:      "image",
					ProductID:  1,
					Quantity:   1,
					Price:      20,
					TotalPrice: 20,
				},
			}}, r)
	}
}
