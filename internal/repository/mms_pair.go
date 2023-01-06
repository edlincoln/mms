package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/edlincoln/mms/internal/model"
	"github.com/edlincoln/mms/internal/utils"
	"github.com/edlincoln/mms/pkg/logger"
)

const (
	queryInsertMmsPair = "INSERT INTO mms_pairs (pair, time_stamp, mms_20, mms_50, mms_200) VALUES($1, $2, $3, $4, $5);"
	queryFindByResult  = "SELECT id, pair, time_stamp, mms_20, mms_50, mms_200 FROM mms_pairs WHERE pair=$1 and time_stamp between $2 and $3;"
	queryCountByResult = "SELECT count(*) FROM mms_pairs WHERE pair=$1;"
)

//go:generate mockery --name=MmsPairRepository --output=../../repository --filename=mms_pair_repository_mock.go --structname=MockMmsPairRepository --inpackage=true
type MmsPairRepository interface {
	Save(ctx context.Context, mmsPair model.MmsPairs) error
	BulkSave(ctx context.Context, mmsPairs []model.MmsPairs) error
	FindByPairAndTimestampRange(ctx context.Context, pair string, from, to time.Time) ([]model.MmsPairs, error)
	Count(ctx context.Context, pair string) (int, error)
}

type MmsPairRepositoryImpl struct {
	manager DatabaseManager
}

func NewMmsPairRepository(manager *DatabaseManager) MmsPairRepository {
	return MmsPairRepositoryImpl{manager: *manager}
}

func (c MmsPairRepositoryImpl) Save(ctx context.Context, mmsPair model.MmsPairs) error {
	stmt, err := c.manager.GetClient().Prepare(queryInsertMmsPair)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return utils.InternalServerError
	}
	defer stmt.Close()

	_, saveErr := stmt.ExecContext(ctx, mmsPair.Pair, mmsPair.Timestamp, mmsPair.Mms20, mmsPair.Mms50, mmsPair.Mms200)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return utils.InternalServerError
	}

	return nil
}

func (c MmsPairRepositoryImpl) BulkSave(ctx context.Context, mmsPairs []model.MmsPairs) error {
	tx, err := c.manager.GetClient().BeginTx(ctx, nil)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return utils.InternalServerError
	}

	for _, mmsPair := range mmsPairs {
		_, saveErr := tx.ExecContext(ctx, queryInsertMmsPair, mmsPair.Pair, mmsPair.Timestamp, mmsPair.Mms20, mmsPair.Mms50, mmsPair.Mms200)
		if saveErr != nil {
			logger.Error("error when trying to save", saveErr)
			tx.Rollback()
			return utils.InternalServerError
		}
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return utils.InternalServerError
	}

	return nil
}

func (c MmsPairRepositoryImpl) FindByPairAndTimestampRange(ctx context.Context, pair string, from, to time.Time) ([]model.MmsPairs, error) {
	mmsPair := make([]model.MmsPairs, 0)
	rows, err := c.manager.GetClient().QueryContext(ctx, queryFindByResult, pair, from, to)
	if err != nil {
		logger.Error("error when trying to prepare statement", err)
		return mmsPair, utils.InternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		mms := new(model.MmsPairs)
		if err := rows.Scan(&mms.Id, &mms.Pair, &mms.Timestamp, &mms.Mms20, &mms.Mms50, &mms.Mms200); err != nil {
			if err != sql.ErrNoRows {
				logger.Error("error when trying to query", err)
				return mmsPair, utils.InternalServerError
			}
			return mmsPair, nil
		}
		mmsPair = append(mmsPair, *mms)
	}

	return mmsPair, nil
}

func (c MmsPairRepositoryImpl) Count(ctx context.Context, pair string) (int, error) {
	var count int
	result := c.manager.GetClient().QueryRowContext(ctx, queryCountByResult, pair)
	if getErr := result.Scan(&count); getErr != nil {
		if getErr != sql.ErrNoRows {
			logger.Error("error when trying to get count", getErr)
			return 0, utils.InternalServerError
		}
		return 0, nil
	}
	return count, nil
}
