package service

import (
	"food_delivery/request"
	"food_delivery/test/util"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RegisterValidationTestSuite struct {
	suite.Suite
	Req *request.RegisterRequest
}

func (suite *RegisterValidationTestSuite) SetupSuite() {
	suite.Req = &request.RegisterRequest{
		Email:     "Ismail@gmail.com",
		FirstName: "Ismail",
		LastName:  "Waleed",
		Password:  "Ismail@123456789",
		Phone:     "156498152136",
	}
}

func (suite *RegisterValidationTestSuite) SetupTest() {

}

func (suite *RegisterValidationTestSuite) TearDownTest() {

}

func (suite *RegisterValidationTestSuite) TearDownSuite() {

}

func TestTokenServiceSuite(t *testing.T) {
	suite.Run(t, new(RegisterValidationTestSuite))
}

func (suite *RegisterValidationTestSuite) TestRegisterValidation() {

	testCases := []util.TestCaseRegisterValidation{
		{
			Name:         "valid input",
			Req:          suite.Req,
			WantError:    false,
			WantErrorMsg: "",
		},
		{
			Name: "invaild Email",
			Req: &request.RegisterRequest{
				FirstName: suite.Req.FirstName,
				LastName:  suite.Req.LastName,
				Email:     "ismailgmaildwa",
				Password:  suite.Req.Password,
				Phone:     suite.Req.Phone,
			},
			WantError:    true,
			WantErrorMsg: "invalid email format",
		},
		{
			Name: "invaild Password",
			Req: &request.RegisterRequest{
				FirstName: suite.Req.FirstName,
				LastName:  suite.Req.LastName,
				Email:     suite.Req.Email,
				Password:  "WrongPassword",
				Phone:     suite.Req.Phone,
			},
			WantError:    true,
			WantErrorMsg: "password must be at least 8 characters long and contain at least one number, one uppercase letter, and one special character",
		},
		{
			Name: "invaild Password",
			Req: &request.RegisterRequest{
				FirstName: suite.Req.FirstName,
				LastName:  suite.Req.LastName,
				Email:     suite.Req.Email,
				Password:  suite.Req.Password,
				Phone:     "123456",
			},
			WantError:    true,
			WantErrorMsg: "phone number must be between 10 and 15 digits",
		},
	}

	for _, testCase := range testCases {
		suite.T().Run(testCase.Name, func(t *testing.T) {
			got := ValidateInput(testCase.Req)

			util.AssertRegisterValidationResult(t, testCase, got)
		})
	}

}

func BenchmarkRegisterValidation(b *testing.B) {

	Req := &request.RegisterRequest{
		Email:     "Ismail@gmail.com",
		FirstName: "Ismail",
		LastName:  "Waleed",
		Password:  "Ismail@123456789",
		Phone:     "156498152136",
	}

	for i := 0; i < b.N; i++ {
		_ = ValidateInput(Req)
	}
}
