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

func (ts *SubscriptionsServiceTestSuite) TestSubscriptionService_FetchSubscription() {
	ctx := context.Background()
	productID := uuid.New()
	subscription := domain.Subscription{
		ID:               uuid.New(),
		ProductID:        productID,
		DurationInMonths: 3,
		Tax:              0.77,
		TotalCost:        10.77,
		Status:           domain.SubscriptionStatusInactive,
		StartDate:        time.Now().AddDate(0, 0, 1),
		EndDate:          time.Now().AddDate(0, 3, 1),
	}

	type GetByIDMock struct {
		timesToCall int
		retSub      domain.Subscription
		retErr      error
	}

	tc := []struct {
		Name            string
		err             error
		ID              uuid.UUID
		responsePayload domain.Subscription
		getByID         GetByIDMock
	}{
		{
			Name:            "Fetch subscription success",
			ID:              productID,
			err:             nil,
			responsePayload: subscription,
			getByID: GetByIDMock{
				timesToCall: 1,
				retSub:      subscription,
				retErr:      nil,
			},
		},
		{
			Name:            "Fetch subscription: invalid product id",
			ID:              uuid.Nil,
			err:             domain.ErrSubscriptionIDIsInvalid,
			responsePayload: domain.Subscription{},
			getByID: GetByIDMock{
				timesToCall: 0,
			},
		},
		{
			Name:            "Fetch subscription: product doesn't exist",
			ID:              productID,
			err:             domain.ErrSubscriptionIDIsInvalid,
			responsePayload: domain.Subscription{},
			getByID: GetByIDMock{
				timesToCall: 1,
				retSub:      domain.Subscription{},
				retErr:      domain.ErrSubscriptionIDIsInvalid,
			},
		},
	}

	for _, tt := range tc {
		ts.Run(tt.Name, func() {
			ts.subscriptionsRepo.EXPECT().
				GetByID(gomock.Any(), gomock.Any()).
				Times(tt.getByID.timesToCall).
				Return(tt.getByID.retSub, tt.getByID.retErr)
			gotSubscription, err := ts.service.FetchSubscription(ctx, tt.ID)
			if tt.err != nil {
				ts.Assert().NotNil(err)
				ts.Assert().EqualError(err, tt.err.Error())
			} else {
				ts.Assert().Nil(err)
				ts.Assert().EqualValues(subscription, gotSubscription)
			}
		})
	}
}

func (ts *SubscriptionsServiceTestSuite) TestSubscriptionService_UpdateSubscriptionStatus() {
	ctx := context.Background()
	productID := uuid.New()

	type patchMock struct {
		timesToCall int
		retErr      error
	}

	tc := []struct {
		Name      string
		err       error
		Status    domain.SubscriptionStatus
		ID        uuid.UUID
		patchMock patchMock
	}{
		{
			Name:   "Update subscription status success",
			ID:     productID,
			Status: domain.SubscriptionStatusActive,
			err:    nil,
			patchMock: patchMock{
				timesToCall: 1,
				retErr:      nil,
			},
		},
	}

	for _, tt := range tc {
		ts.Run(tt.Name, func() {
			ts.subscriptionsRepo.EXPECT().
				Patch(gomock.Any(), tt.ID, map[string]interface{}{
					"status": tt.Status,
				}).
				Times(tt.patchMock.timesToCall).
				Return(tt.patchMock.retErr)
			err := ts.service.UpdateSubscriptionStatus(ctx, tt.ID, tt.Status)
			if tt.err != nil {
				ts.Assert().NotNil(err)
				ts.Assert().EqualError(err, tt.err.Error())
			} else {
				ts.Assert().Nil(err)
			}
		})
	}
}
