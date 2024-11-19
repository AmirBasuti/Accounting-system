package repository_test

import (
	"AccountingSystem/internal/database"
	"AccountingSystem/internal/models"
	"AccountingSystem/internal/repository"
	"testing"
)

var createdDLs map[uint]uint

func TestMain(m *testing.M) {
	database.Connect()
	createdDLs = make(map[uint]uint)
	m.Run()
}

func TestDLRepo_Create(t *testing.T) {
	repo := &repository.DLRepo{DB: database.DB}
	code, title := repository.GenerateRandomCodeAndTitle()
	detail := &models.DL{Code: code, Title: title}

	if err := repo.Create(detail); err != nil {
		t.Error(err)
	}

	tests := []struct {
		dl   *models.DL
		want string
	}{
		{&models.DL{Code: code, Title: "test"}, "expected: Code already exists, got: nil"},
		{&models.DL{Code: "test", Title: title}, "expected: Title already exists, got: nil"},
		{&models.DL{Code: code}, "expected: Title can't be null, got: nil"},
		{&models.DL{Title: title}, "expected: Code can't be null, got: nil"},
	}

	for _, tt := range tests {
		if err := repo.Create(tt.dl); err == nil {
			t.Error(tt.want)
		}
	}

	createdDLs[detail.ID] = detail.Version
}

func TestDLRepo_Update(t *testing.T) {
	repo := &repository.DLRepo{DB: database.DB}
	code, title := repository.GenerateRandomCodeAndTitle()
	detail := &models.DL{Code: code, Title: title}

	if err := repo.Create(detail); err != nil {
		t.Error(err, "can't create detail")
	}

	code, title = repository.GenerateRandomCodeAndTitle()
	detail.Title = title
	if err := repo.Update(detail); err != nil {
		t.Error(err, "can't update detail")
	}
	if detail.Title != title {
		t.Errorf("expected: %s, got: %s", title, detail.Title)
	}

	detail.Code = code
	if err := repo.Update(detail); err != nil {
		t.Error(err, "can't update detail")
	}
	if detail.Code != code {
		t.Errorf("expected: %s, got: %s", code, detail.Code)
	}

	createdDLs[detail.ID] = detail.Version
}

func TestDLRepo_GetByID(t *testing.T) {
	repo := &repository.DLRepo{DB: database.DB}
	code, title := repository.GenerateRandomCodeAndTitle()
	detail := &models.DL{Code: code, Title: title}

	if err := repo.Create(detail); err != nil {
		t.Error(err, "can't create detail")
	}

	detail2, err := repo.GetByID(detail.ID)
	if err != nil {
		t.Error(err, "can't get detail")
	}
	if detail.ID != detail2.ID {
		t.Errorf("expected: %d, got: %d", detail.ID, detail2.ID)
	}

	if _, err := repo.GetByID(0); err == nil {
		t.Error("expected: detail not found, got: nil")
	}

	createdDLs[detail.ID] = detail.Version
}

func TestDLRepo_Delete(t *testing.T) {
	repo := &repository.DLRepo{DB: database.DB}
	code, title := repository.GenerateRandomCodeAndTitle()
	detail := &models.DL{Code: code, Title: title}

	if err := repo.Create(detail); err != nil {
		t.Error(err, "can't create detail")
	}

	if err := repo.Delete(detail.ID, detail.Version); err != nil {
		t.Error(err, "can't delete detail")
	}

	if err := repo.Delete(detail.ID, detail.Version); err == nil {
		t.Error("expected: detail not found, got: nil")
	}

	for k, v := range createdDLs {
		if err := repo.Delete(k, v); err != nil {
			t.Error(err, "can't delete detail")
		}
	}
}
