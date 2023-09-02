package service_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SantiagoBedoya/otp-api/internal/dto"
	"github.com/SantiagoBedoya/otp-api/internal/repository/mysql"
	"github.com/SantiagoBedoya/otp-api/internal/service"
	mmysql "github.com/go-sql-driver/mysql"
)

func TestSignIn(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening stub db connection: %v", err)
	}
	defer db.Close()

	repo := mysql.NewUserRepository(db)
	svc := service.NewAuthService(repo)

	mock.ExpectPrepare("SELECT id, email, password FROM users WHERE email = ?")
	mock.ExpectQuery("SELECT id, email, password FROM users WHERE email = ?").
		WithArgs("santiago@google.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(1, "santiago@google.com", "hash123"))

	_, err = svc.SignIn(&dto.SignInDto{
		Email:    "santiago@google.com",
		Password: "santiago123",
	})

	if err == nil {
		t.Errorf("Expected error %v, got nil", service.ErrInvalidPassword)
	}

	mock.ExpectPrepare("SELECT id, email, password FROM users WHERE email = ?")
	mock.ExpectQuery("SELECT id, email, password FROM users WHERE email = ?").
		WithArgs("santiago@google.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(1, "santiago@google.com", "$2a$10$DwfxEpjq0gE2YW3OcWUF9eOQJzfdrkvRN8j3O06Olbm2ng6wNcpMK"))

	token, err := svc.SignIn(&dto.SignInDto{
		Email:    "santiago@google.com",
		Password: "santiago123",
	})

	if err != nil {
		t.Errorf("unxpected error: %v", err)
	}

	if len(token) == 0 {
		t.Errorf("token should not be empty")
	}
}

func TestGenerateAccessToken(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := mysql.NewUserRepository(db)
	svc := service.NewAuthService(repo)

	token, err := svc.GenerateAccessToken("1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(token) == 0 {
		t.Error("Token should not be empty")
	}
}

func TestSaveUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := mysql.NewUserRepository(db)
	svc := service.NewAuthService(repo)

	mock.ExpectPrepare("INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)").
		ExpectExec().
		WithArgs("santiago", "bedoya", "santiago@google.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	user, err := svc.SaveUser(&dto.SignUpDto{
		FirstName: "santiago",
		LastName:  "bedoya",
		Email:     "santiago@google.com",
		Password:  "santiago123",
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if user == nil {
		t.Error("The user should not be empty")
	}
	if user != nil && user.ID != 1 {
		t.Errorf("Expected userID '1', got %d", user.ID)
	}

	mock.ExpectPrepare("INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)").
		ExpectExec().
		WithArgs("santiago", "bedoya", "santiago@google.com", sqlmock.AnyArg()).
		WillReturnError(&mmysql.MySQLError{Number: 1062})

	user, err = svc.SaveUser(&dto.SignUpDto{
		FirstName: "santiago",
		LastName:  "bedoya",
		Email:     "santiago@google.com",
		Password:  "santiago123",
	})

	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}

	if err == nil {
		t.Errorf("Expected error %v, got nil", service.ErrEmailInUse)
	}

	if err != nil {
		if !errors.Is(err, service.ErrEmailInUse) {
			t.Errorf("expected error: %v, got %v", service.ErrEmailInUse, err)
		}
	}
}

func TestSetUser2FAValid(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := mysql.NewUserRepository(db)
	svc := service.NewAuthService(repo)

	mock.ExpectPrepare(mysql.SetValidSecretQuery).
		ExpectExec().
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = svc.SetUser2FAValid("1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSetUser2FA(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := mysql.NewUserRepository(db)
	svc := service.NewAuthService(repo)

	mock.ExpectPrepare(mysql.SaveUserSecretQuery).
		ExpectExec().
		WithArgs("secret123", "1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = svc.SetUser2FA("1", "secret123")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := mysql.NewUserRepository(db)
	svc := service.NewAuthService(repo)

	mock.ExpectPrepare(mysql.GetUserByIDQuery).
		ExpectQuery().
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"email", "secret_2fa"}).
			AddRow("santiago@google.com", "secret123"))

	user, err := svc.GetUserByID("1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if user == nil {
		t.Error("Unexpected nil user")
	}

	if user != nil && user.Email != "santiago@google.com" {
		t.Errorf("Expected user email 'santiago@google.com', got %v", user.Email)
	}

	if user != nil && user.Secret2FA != "secret123" {
		t.Errorf("Expected user secret 2FA 'secret123', got %v", user.Secret2FA)
	}
}
