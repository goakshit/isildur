package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/goakshit/isildur/core/domain"
	"github.com/goakshit/isildur/core/ports"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SubscriptionsServiceTestSuite struct {
	suite.Suite
	productsRepo      *ports.MockProductsRepository
	subscriptionsRepo *ports.MockSubscriptionsRepository
	service           *SubscriptionService
}

func TestSubscriptionsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionsServiceTestSuite))
}

func (ts *SubscriptionsServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(ts.T())
	ts.productsRepo = ports.NewMockProductsRepository(ctrl)
	ts.subscriptionsRepo = ports.NewMockSubscriptionsRepository(ctrl)
	ts.service = NewSubscriptionService(ts.subscriptionsRepo, ts.productsRepo)
}

func (ts *SubscriptionsServiceTestSuite) TestSubscriptionService_Create() {
	ctx := context.Background()
	productID := uuid.New()
	product := domain.Product{
		ID:             productID,
		Name:           "YOGA 1",
		Description:    "BASIC YOGA",
		MonthlyPrice:   5,
		InstructorName: "A. Dhar",
	}

	type GetByIDMock struct {
		timesToCall int
		retProd     domain.Product
		retErr      error
	}

	type CreateMock struct {
		timesToCall int
		retErr      error
	}

	tc := []struct {
		Name              string
		err               error
		ID                uuid.UUID
		DurationInMonths  int8
		startDate         time.Time
		responsePayload   domain.Product
		getByID           GetByIDMock
		createSubsription CreateMock
	}{
		{
			Name:             "Create Subscription Success",
			ID:               productID,
			DurationInMonths: 3,
			startDate:        time.Now(),
			err:              nil,
			responsePayload:  product,
			getByID: GetByIDMock{
				timesToCall: 1,
				retProd:     product,
				retErr:      nil,
			},
			createSubsription: CreateMock{
				timesToCall: 1,
				retErr:      nil,
			},
		},
		{
			Name:             "Create Subscription failed: product not found",
			ID:               productID,
			DurationInMonths: 3,
			startDate:        time.Now(),
			err:              domain.ErrProductNotfound,
			responsePayload:  product,
			getByID: GetByIDMock{
				timesToCall: 1,
				retProd:     domain.Product{},
				retErr:      domain.ErrProductNotfound,
			},
			createSubsription: CreateMock{
				timesToCall: 0,
			},
		},
		{
			Name:             "Create Subscription failed: invalid product id",
			ID:               uuid.Nil,
			DurationInMonths: 3,
			startDate:        time.Now(),
			err:              domain.ErrProductIDIsInvalid,
			responsePayload:  product,
			getByID: GetByIDMock{
				timesToCall: 1,
				retProd:     domain.Product{},
				retErr:      domain.ErrProductIDIsInvalid,
			},
			createSubsription: CreateMock{
				timesToCall: 0,
			},
		},
		{
			Name:             "Create Subscription failed",
			ID:               productID,
			DurationInMonths: 3,
			startDate:        time.Now(),
			err:              errors.New("something went wrong"),
			responsePayload:  product,
			getByID: GetByIDMock{
				timesToCall: 1,
				retProd:     product,
				retErr:      nil,
			},
			createSubsription: CreateMock{
				timesToCall: 1,
				retErr:      errors.New("something went wrong"),
			},
		},
	}

	for _, tt := range tc {
		ts.Run(tt.Name, func() {
			ts.productsRepo.EXPECT().
				GetByID(gomock.Any(), tt.ID).
				Times(tt.getByID.timesToCall).
				Return(tt.getByID.retProd, tt.getByID.retErr)

			ts.subscriptionsRepo.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Times(tt.createSubsription.timesToCall).
				Return(tt.createSubsription.retErr)
			err := ts.service.CreateSubscription(ctx, tt.ID, tt.DurationInMonths, tt.startDate)
			if tt.err != nil {
				ts.Assert().NotNil(err)
				ts.Assert().EqualError(err, tt.err.Error())
			} else {
				ts.Assert().Nil(err)
			}
		})
	}
}
