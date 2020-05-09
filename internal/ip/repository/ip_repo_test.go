package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
)

func TestIpRepository_CreateFirstIpMask(t *testing.T) {
	db, mock, _ := sqlmock.New()
	ipRepo := NewIpRepository(db)
	customErr := errors.New("smth happend")

	ip := "123"
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO ip \(ip_value,counter\) VALUES \(\$1,\$2\)`).WithArgs(ip, 1).WillReturnError(customErr)
	mock.ExpectCommit()

	tx, _ := db.Begin()

	if err := ipRepo.CreateFirstIpMask(ip); err != customErr {
		tx.Rollback()
		t.Error("error at CreateFirstIpMask", err)
	}
	tx.Commit()
}

func TestIpRepository_SetMaskDataToDefault(t *testing.T) {
	db, mock, _ := sqlmock.New()
	ipRepo := NewIpRepository(db)
	customErr := errors.New("smth happend")

	ip := "123"
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE ip SET counter \= 1 , time \= NULL where ip_value \= \$1`).WithArgs(ip).WillReturnError(customErr)
	mock.ExpectCommit()

	tx, _ := db.Begin()

	if err := ipRepo.SetMaskDataToDefault(ip); err != customErr {
		tx.Rollback()
		t.Error("error at SetMaskDataToDefault", err)
	}
	tx.Commit()
}

func TestIpRepository_UpdateMaskData(t *testing.T) {
	db, mock, _ := sqlmock.New()
	ipRepo := NewIpRepository(db)
	customErr := errors.New("smth happend")

	ip := "123"
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE ip SET counter \= counter \+ 1 where ip_value \= \$1`).WithArgs(ip).WillReturnError(customErr)
	mock.ExpectCommit()

	tx, _ := db.Begin()

	if err := ipRepo.UpdateMaskData(ip); err != customErr {
		tx.Rollback()
		t.Error("error at UpdateMaskData", err)
	}
	tx.Commit()
}

func TestIpRepository_ResetIpCoolDown(t *testing.T) {
	db, mock, _ := sqlmock.New()
	ipRepo := NewIpRepository(db)
	customErr := errors.New("smth happend")

	ip := "123"
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE ip SET counter \= 0 , time \= NULL where ip_value \= \$1`).WithArgs(ip).WillReturnError(customErr)
	mock.ExpectCommit()

	tx, _ := db.Begin()

	if err := ipRepo.ResetIpCoolDown(ip); err != customErr {
		tx.Rollback()
		t.Error("error at ResetIpCoolDown", err)
	}
	tx.Commit()
}

func TestIpRepository_SetMaskCoolDown(t *testing.T) {
	db, mock, _ := sqlmock.New()
	ipRepo := NewIpRepository(db)
	customErr := errors.New("smth happend")

	ip := "123"
	limit := 2
	_time := time.Now()
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE ip SET counter \= \$2 , time \= \$3 where ip_value \= \$1`).WithArgs(ip, limit, _time).WillReturnError(customErr)
	mock.ExpectCommit()

	tx, _ := db.Begin()

	if err := ipRepo.SetMaskCoolDown(ip, limit, _time); err != customErr {
		tx.Rollback()
		t.Error("error at SetMaskCoolDown", err)
	}
	tx.Commit()
}

func TestIpRepository_GetMaskData(t *testing.T) {
	db, mock, _ := sqlmock.New()
	ipRepo := NewIpRepository(db)
	customErr := errors.New("smth happend")

	ip := "123"
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT counter , time FROM ip WHERE ip_value \= \$1`).WithArgs(ip).WillReturnError(customErr)
	mock.ExpectCommit()

	tx, _ := db.Begin()

	if counter, limiter, err := ipRepo.GetMaskData(ip); err != customErr || counter != 0 || limiter != "" {
		tx.Rollback()
		t.Error("error at GetMaskData", err)
	}

	tx.Commit()
	limit := 2
	time := "time"
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT counter , time FROM ip WHERE ip_value \= \$1`).WithArgs(ip).WillReturnRows(sqlmock.NewRows([]string{"counter", "time"}).AddRow(limit, time))
	mock.ExpectCommit()

	tx, _ = db.Begin()

	if counter, limiter, err := ipRepo.GetMaskData(ip); err != nil || counter != limit || limiter != time {
		tx.Rollback()
		t.Error("error at GetMaskData", err)
	}

	tx.Commit()
}
