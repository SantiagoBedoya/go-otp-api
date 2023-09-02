package service_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SantiagoBedoya/otp-api/internal/dto"
	"github.com/SantiagoBedoya/otp-api/internal/repository/mysql"
	"github.com/SantiagoBedoya/otp-api/internal/service"
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
