package service_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SantiagoBedoya/otp-api/internal/repository/mysql"
	"github.com/SantiagoBedoya/otp-api/internal/service"
)

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening stub db connection: %v", err)
	}
	defer db.Close()

	repo := mysql.NewUserRepository(db)
	svc := service.NewUserService(repo)

	mock.ExpectPrepare("SELECT id, first_name, last_name, email FROM users")
	mock.ExpectQuery("SELECT id, first_name, last_name, email FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email"}).
			AddRow(1, "Santiago", "Bedoya", "santiago@google.com"))

	users, err := svc.GetUsers()
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}
	if len(users) != 1 {
		t.Errorf("expected len shoud 1, got %d", len(users))
	}
}
