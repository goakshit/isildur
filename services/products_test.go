package services

import (
	"context"
	"testing"

	"github.com/goakshit/isildur/core/domain"
	"github.com/goakshit/isildur/core/ports"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ProductsServiceTestSuite struct {
	suite.Suite
	productsRepo *ports.MockProductsRepository
	service      *ProductsService
}

func TestProductsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductsServiceTestSuite))
}

func (ts *ProductsServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(ts.T())
	ts.productsRepo = ports.NewMockProductsRepository(ctrl)
	ts.service = NewProductsService(ts.productsRepo)
}

func (ts *ProductsServiceTestSuite) TestProductsService_FetchProduct() {
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

	tc := []struct {
		Name            string
		ID              uuid.UUID
		err             error
		responsePayload domain.Product
		getByID         GetByIDMock
	}{
		{
			Name:            "Fetch Product Success",
			ID:              productID,
			err:             nil,
			responsePayload: product,
			getByID: GetByIDMock{
				timesToCall: 1,
				retProd:     product,
				retErr:      nil,
			},
		},
		{
			Name:            "Fetch Product: Not found",
			ID:              productID,
			err:             domain.ErrProductNotfound,
			responsePayload: domain.Product{},
			getByID: GetByIDMock{
				timesToCall: 1,
				retProd:     domain.Product{},
				retErr:      domain.ErrProductNotfound,
			},
		},
		{
			Name:            "Fetch Product: Invalid UUID",
			ID:              uuid.Nil,
			err:             domain.ErrProductIDIsInvalid,
			responsePayload: product,
			getByID: GetByIDMock{
				timesToCall: 0,
			},
		},
	}

	for _, tt := range tc {
		ts.Run(tt.Name, func() {
			ts.productsRepo.EXPECT().
				GetByID(gomock.Any(), tt.ID).
				Times(tt.getByID.timesToCall).
				Return(tt.getByID.retProd, tt.getByID.retErr)

			gotProduct, err := ts.service.FetchProduct(ctx, tt.ID)
			if tt.err != nil {
				ts.Assert().NotNil(err)
				ts.Assert().EqualError(err, tt.err.Error())
			} else {
				ts.Assert().Nil(err)
				ts.Assert().EqualValues(tt.responsePayload, gotProduct)
			}
		})
	}
}

func (ts *ProductsServiceTestSuite) TestProductsService_FetchAllProducts() {
	ctx := context.Background()
	productID := uuid.New()
	product := domain.Product{
		ID:             productID,
		Name:           "YOGA 1",
		Description:    "BASIC YOGA",
		MonthlyPrice:   5,
		InstructorName: "A. Dhar",
	}

	type GetAllMock struct {
		timesToCall int
		retProd     []domain.Product
		retErr      error
	}

	tc := []struct {
		Name            string
		err             error
		responsePayload []domain.Product
		getByAll        GetAllMock
	}{
		{
			Name:            "Fetch Products Success",
			err:             nil,
			responsePayload: []domain.Product{product},
			getByAll: GetAllMock{
				timesToCall: 1,
				retProd:     []domain.Product{product},
				retErr:      nil,
			},
		},
		{
			Name:            "No products available",
			err:             nil,
			responsePayload: []domain.Product{},
			getByAll: GetAllMock{
				timesToCall: 1,
				retProd:     []domain.Product{},
				retErr:      nil,
			},
		},
	}

	for _, tt := range tc {
		ts.Run(tt.Name, func() {
			ts.productsRepo.EXPECT().
				GetAll(gomock.Any()).
				Times(tt.getByAll.timesToCall).
				Return(tt.getByAll.retProd, tt.getByAll.retErr)

			gotProducts, err := ts.service.FetchAllProducts(ctx)
			if tt.err != nil {
				ts.Assert().NotNil(err)
				ts.Assert().EqualError(err, tt.err.Error())
			} else {
				ts.Assert().Nil(err)
				ts.Assert().EqualValues(tt.responsePayload, gotProducts)
			}
		})
	}
}
