package models_test

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skrewby/yapper/database"
	"github.com/skrewby/yapper/models"
	"github.com/skrewby/yapper/utils"
)

func TestUsers_CreateUser(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *pgxpool.Pool
		// Named input parameters for target function.
		email        string
		display_name string
		hash         string
		wantErr      bool
	}{
		{
			name:         "Create user normally",
			db:           db,
			email:        "bob@email.com",
			display_name: "Bob Smith",
			hash:         "1234",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := models.NewUsersModel(tt.db)
			gotErr := u.CreateUser(tt.email, tt.display_name, tt.hash)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateUser() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateUser() succeeded unexpectedly")
			}
		})
	}
}

func TestUsers_UpdateUser(t *testing.T) {
	m := models.NewUsersModel(db)
	err := m.CreateUser("UUSER", "Update User A", "cdefg")
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	user, err := m.GetUserByEmail("UUSER")
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	user.Name = "Update User B"
	user.Active = utils.Pointer(false)
	err = m.UpdateUser(user)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	user, err = m.GetUserByEmail("UUSER")
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if user.Name != "Update User B" || *user.Active {
		t.Errorf("User was not updated correctly")
	}
}

func TestUsers_GetAllUsers(t *testing.T) {
	m := models.NewUsersModel(db)
	query := `DELETE FROM users`
	err := database.Run(db, query, pgx.NamedArgs{})
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	usr, err := m.GetAllUsers()
	if usr != nil {
		t.Errorf("Expected GetAllUsers to return empty")
	}
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	m.CreateUser("GAU01", "Get All Users A", "asdf")
	m.CreateUser("GAU02", "Get All Users B", "qwer")
	m.CreateUser("GAU03", "Get All Users C", "zxcv")
	users, err := m.GetAllUsers()
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	matches := 0
	for _, user := range users {
		if user.Email == "GAU01" || user.Email == "GAU02" || user.Email == "GAU03" {
			matches += 1
		}
	}
	if matches != 3 {
		t.Errorf("Number of users does not match, expected 3 but found %d", matches)
	}
}
