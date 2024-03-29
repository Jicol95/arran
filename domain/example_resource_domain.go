package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/jicol-95/arran/dal"
	"github.com/labstack/echo/v4"
)

type ExampleResourceService struct {
	logger echo.Logger
	tm     dal.TransactionManager
	repo   dal.ExampleResourceRepository
}

func (svc *ExampleResourceService) CreateExampleResource(name string) (dal.ExampleResource, error) {
	resource := dal.ExampleResource{ID: uuid.NewString(), Name: name, CreatedAt: time.Now().UTC()}

	tx, _ := svc.tm.BeginTx()

	if err := svc.repo.Insert(resource, tx); err != nil {
		svc.logger.Error("Failed to insert example resource")
		tx.Rollback()
		return dal.ExampleResource{}, err
	}

	if err := tx.Commit(); err != nil {
		svc.logger.Error("Failed to insert example resource")
		return dal.ExampleResource{}, err
	}

	svc.logger.Info("Successfully inserted example resource")

	return resource, nil
}

func (svc *ExampleResourceService) UpdateExampleResource(id string, name string) (dal.ExampleResource, error) {
	tx, _ := svc.tm.BeginTx()

	resource, err := svc.repo.FetchByID(id, true, tx)

	if err != nil {
		svc.logger.Error("Failed to fetch example resource for update")
		tx.Rollback()
		return dal.ExampleResource{}, err
	}

	resource.Name = name

	if err := svc.repo.Update(resource, tx); err != nil {
		svc.logger.Error("Failed to insert example resource")
		tx.Rollback()
		return dal.ExampleResource{}, err
	}

	tx.Commit()

	svc.logger.Info("Successfully inserted example resource")

	return resource, nil
}

func (svc *ExampleResourceService) GetExampleResourceById(id string) (dal.ExampleResource, error) {
	resource, err := svc.repo.FetchByID(id, false, nil)

	if err != nil {
		svc.logger.Error(err)
		return dal.ExampleResource{}, nil
	}

	return resource, nil
}

func (svc *ExampleResourceService) DeleteExmpleResourceById(id string) error {
	err := svc.repo.DeleteByID(id, nil)

	if err != nil {
		svc.logger.Error(err)
	}

	return err
}

func NewExampleResourceService(
	logger echo.Logger,
	tm dal.TransactionManager,
	repo dal.ExampleResourceRepository,
) ExampleResourceService {
	return ExampleResourceService{
		logger: logger,
		tm:     tm,
		repo:   repo,
	}
}
