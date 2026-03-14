package handlers

import (
	"main.go/models"
	"main.go/services"
	"main.go/utils"
)

// Function variables wrapping model/service/utils calls so tests can replace them.
var (
	modelRegisterUser     func(email, password string) (int64, error)        = models.RegisterUser
	modelAuthenticateUser func(email, password string) (*models.User, error) = models.AuthenticateUser
	generateToken         func(userID int64, email string) (string, error)   = utils.GenerateToken

	modelCreateTest      func(test models.TestRequest, userID int64) (int64, error) = models.CreateTest
	modelGetAllTests     func(userID int64) ([]models.Test, error)                  = models.GetAllTests
	modelUpdateTest      func(test models.Test, testId int64) error                 = models.UpdateTest
	modelDeleteTest      func(testID int64) error                                   = models.DeleteTest
	modelCheckTestExists func(testID int64) (bool, error)                           = models.CheckTestIdExists
	modelGetTestResult   func(testRunID int64) (models.TestRun, error)              = models.GetTestRunResult

	serviceStartTestRun func(testID int64, concurrency int) (int64, string, error) = services.StartTestRun
	serviceRunJobs      func(testID int64, concurrency int, testRunID int64)       = services.RunJobs
)
