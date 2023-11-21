package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/jicol-95/arran/dao"
	"github.com/labstack/echo/v4"
)

type ExampleResourceService struct {
	logger echo.Logger
	tm     dao.TransactionManager
	repo   dao.ExampleResourceRepository
}

func (svc *ExampleResourceService) CreateExampleResource(name string) (dao.ExampleResource, error) {
	resource := dao.ExampleResource{ID: uuid.NewString(), Name: name, CreatedAt: time.Now().UTC()}

	tx, _ := svc.tm.BeginTx()

	if err := svc.repo.Insert(resource, tx); err != nil {
		svc.logger.Error("Failed to insert example resource")
		tx.Rollback()
		return dao.ExampleResource{}, err
	}

	if err := tx.Commit(); err != nil {
		svc.logger.Error("Failed to insert example resource")
		return dao.ExampleResource{}, nil
	}

	svc.logger.Info("Successfully inserted example resource")

	return resource, nil
}

func (svc *ExampleResourceService) GetExampleResourceById(id string) (dao.ExampleResource, error) {
	resource, err := svc.repo.FetchByID(id, false, nil)

	if err != nil {
		svc.logger.Error(err)
		return dao.ExampleResource{}, nil
	}

	return resource, nil
}

func (svc *ExampleResourceService) GetAllExampleResources() ([]dao.ExampleResource, error) {
	resources, err := svc.repo.FetchAll(nil)

	if err != nil {
		svc.logger.Error(err)
		return nil, err
	}

	return resources, nil
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
	tm dao.TransactionManager,
	repo dao.ExampleResourceRepository,
) ExampleResourceService {
	return ExampleResourceService{
		logger: logger,
		tm:     tm,
		repo:   repo,
	}
}