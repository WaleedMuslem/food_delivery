package handler

import (
	"database/sql"
	"fmt"
	"food_delivery/config"
	"food_delivery/repository"
	"food_delivery/service"
	"food_delivery/test/util"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const UserID uint = 1

type CartHandlerTestSuite struct {
	suite.Suite
	accessToken string
	CartHandler *CartHandler
}

func (s *CartHandlerTestSuite) SetupSuite() {

	cartRepo := repository.NewCartRepositoryFake()

	cfg := &config.Config{
		AccessSecret:          "access",
		AccessLifetimeminutes: 15,
	}

	dsn := "postgres://" + cfg.DbUsername + ":" + cfg.DbPassword + "@localhost/" + cfg.DbName + "?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	tokenService := service.NewTokenService(cfg, db)
	s.CartHandler = NewCartController(tokenService, cartRepo)
	s.accessToken, _ = tokenService.GenerateAccessToken(UserID)
	// s.accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTUsImV4cCI6MTcyNzY0MDgyMH0.Xzkzhhs2G7f3fc3dlIje3XuJDNjCcAK71yRKpdOFlHU"

}

func TestUSerHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CartHandlerTestSuite))
}

func (s *CartHandlerTestSuite) TestWalkCartHandlerGetCart() {
	handlerfunc := s.CartHandler.GetCart
	cases := []*util.TestCaseCartHandlerGetCart{
		{
			TestName: "Susccefuly got cart",
			Request: &util.Request{
				Method:      http.MethodGet,
				Url:         "/cart/getCart",
				AccessToken: s.accessToken,
			},
			HandlerFunc: handlerfunc,
			Want: &util.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "pizza",
			},
		},
	}

	for _, testCase := range cases {
		s.T().Run(testCase.TestName, func(t *testing.T) {
			request, recorder := util.PreparHandlerTestCases(testCase)
			testCase.HandlerFunc(recorder, request)

			fmt.Println("here is what we want ", testCase.Want.BodyPart)
			fmt.Println("here is what recorder ", recorder.Body.String())
			fmt.Println("here is acess ", testCase.Request.AccessToken)

			assert.Contains(t, recorder.Body.String(), testCase.Want.BodyPart)

			if assert.Equal(t, testCase.Want.StatusCode, recorder.Code) {
				if recorder.Code == http.StatusOK {
					util.AssertUserProfileResponse(t, recorder)
				}
			}
		})

	}

}
