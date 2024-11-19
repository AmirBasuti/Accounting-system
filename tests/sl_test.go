package repository_test

import (
	"testing"

	"AccountingSystem/internal/database"
	"AccountingSystem/internal/models"
	"AccountingSystem/internal/repository"
)

func TestSLRepo_Create(t *testing.T) {
	// Setup
	repo := &repository.SLRepo{DB: database.DB}
	code, title := repository.GenerateRandomCodeAndTitle()
	expectedIsDetail := repository.GenerateRandomBool()

	sl := &models.SL{
		Code:     code,
		Title:    title,
		IsDetail: expectedIsDetail,
	}

	// Test
	if err := repo.Create(sl); err != nil {
		t.Fatalf("failed to create SL: %v", err)
	}

	// Cleanup
	if err := repo.Delete(sl.ID, sl.Version); err != nil {
		t.Fatalf("can't delete SL. expected no error, got %v", err)
	}
}

func TestSLRepo_GetByID(t *testing.T) {
	repo := &repository.SLRepo{DB: database.DB}
	code, title := repository.GenerateRandomCodeAndTitle()
	isDetail := repository.GenerateRandomBool()
	sl := &models.SL{
		Code:     code,
		Title:    title,
		IsDetail: isDetail,
	}
	repo.Create(sl)

	fetchedSL, err := repo.GetByID(sl.ID)
	if err != nil {
		t.Fatalf("can't get SL by ID. expected no error, got %v", err)
	}
	if fetchedSL.Code != code {
		t.Fatalf("expected code to be %v, got %v", code, fetchedSL.Code)
	}
	// Cleanup
	if err := repo.Delete(sl.ID, sl.Version); err != nil {
		t.Fatalf("cant delete sl. expected no error, got %v", err)
	}
}

func TestSLRepo_Update(t *testing.T) {
	repo := &repository.SLRepo{DB: database.DB}
	code, title := repository.GenerateRandomCodeAndTitle()
	isDetail := repository.GenerateRandomBool()
	sl := &models.SL{Code: code, Title: title, IsDetail: isDetail}
	repo.Create(sl)

	code, title = repository.GenerateRandomCodeAndTitle()
	sl.Title = title
	if err := repo.Update(sl); err != nil {
		t.Fatalf("failed to update SL: %v", err)
	}
	updatedSL, err := repo.GetByID(sl.ID)
	if err != nil {
		t.Fatalf("can't get SL by ID. expected no error, got %v", err)
	}
	if updatedSL.Title != title {
		t.Fatalf("expected title to be %v, got %v", title, updatedSL.Title)
	}
	sl.Code = code
	if err := repo.Update(sl); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedSL, err = repo.GetByID(sl.ID); err != nil {
		// Cleanup
		t.Fatalf("can't get SL by ID. expected no error, got %v", err)
	}
	if updatedSL.Code != code {
		t.Fatalf("expected code to be %v, got %v", code, updatedSL.Code)
	}
	// Cleanup
	if err := repo.Delete(sl.ID, sl.Version); err != nil {
		t.Fatalf("cant delete sl. expected no error, got %v", err)
	}
}

func TestSLRepo_Delete(t *testing.T) {
	repo := &repository.SLRepo{DB: database.DB}
	code, title := repository.GenerateRandomCodeAndTitle()
	boolian := repository.GenerateRandomBool()
	sl := &models.SL{Code: code, Title: title, IsDetail: boolian}
	repo.Create(sl)

	if err := repo.Delete(sl.ID, sl.Version); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	_, err := repo.GetByID(sl.ID)
	if err == nil {
		t.Fatalf("delete SL did not work. expected error, got nil")
	}

}

func TestSLRepo_Update_VersionMismatch(t *testing.T) {
	repo := &repository.SLRepo{DB: database.DB}
	code, title := repository.GenerateRandomCodeAndTitle()
	boolian := repository.GenerateRandomBool()
	sl := &models.SL{Code: code, Title: title, IsDetail: boolian}
	repo.Create(sl)

	sl.Version = 999 // Simulate version mismatch
	if err := repo.Update(sl); err == nil || err.Error() != "version mismatch: the record has been updated by another process" {
		t.Fatalf("expected version mismatch error when deleting SL, got %v", err)
	}
	// Cleanup
	if err := repo.Delete(sl.ID, 1); err != nil {
		t.Fatalf("cant delete sl. expected no error, got %v", err)
	}
}

func TestSLRepo_Delete_VersionMismatch(t *testing.T) {
	repo := &repository.SLRepo{DB: database.DB}
	code, title := repository.GenerateRandomCodeAndTitle()
	boolian := repository.GenerateRandomBool()
	sl := &models.SL{Code: code, Title: title, IsDetail: boolian}
	repo.Create(sl)
	if err := repo.Delete(sl.ID, 999); err == nil || err.Error() != "version mismatch: the record has been updated by another process" {
		t.Fatalf("expected version mismatch error, got %v", err)
	}
	// Cleanup
	if err := repo.Delete(sl.ID, sl.Version); err != nil {
		t.Fatalf("cant delete sl. expected no error, got %v", err)
	}
}
