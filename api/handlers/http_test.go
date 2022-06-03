package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goakshit/isildur/core/domain"
	"github.com/goakshit/isildur/core/ports"
	"github.com/goakshit/isildur/platform/constants"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type HttpTestSuite struct {
	suite.Suite
	prodSvc *ports.MockProductsService
	subsSvc *ports.MockSubscriptionService
}

func TestSubscriptionsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(HttpTestSuite))
}

func (ts *HttpTestSuite) SetupTest() {
	ctrl := gomock.NewController(ts.T())
	ts.prodSvc = ports.NewMockProductsService(ctrl)
	ts.subsSvc = ports.NewMockSubscriptionService(ctrl)
}

func getProduct() domain.Product {
	return domain.Product{
		ID:             uuid.New(),
		Name:           "YOGA 1",
		Description:    "BASIC YOGA",
		MonthlyPrice:   5,
		InstructorName: "A. Dhar",
	}
}

func (ts *HttpTestSuite) TestHttpHandlers_FetchAllProducts() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")

	expectedProducts := []domain.Product{
		getProduct(),
	}

	expectedProductsBytes, _ := json.Marshal(expectedProducts)

	ts.prodSvc.EXPECT().FetchAllProducts(gomock.Any()).
		Times(1).
		Return(expectedProducts, nil)

	hndlr := NewHTTPHandler(ts.subsSvc, ts.prodSvc)
	hndlr.FetchAllProducts(c)
	ts.Assert().EqualValues(http.StatusOK, w.Code)

	data, err := io.ReadAll(w.Result().Body)
	ts.Assert().Nil(err)
	ts.Assert().EqualValues(data, expectedProductsBytes)
}

func (ts *HttpTestSuite) TestHttpHandlers_FetchProductSuccess() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	productID := uuid.New()
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")
	c.AddParam(constants.ProductIDKey, productID.String())

	expectedProduct := getProduct()

	expectedProductBytes, _ := json.Marshal(expectedProduct)

	ts.prodSvc.EXPECT().FetchProduct(gomock.Any(), productID).
		Times(1).
		Return(expectedProduct, nil)

	hndlr := NewHTTPHandler(ts.subsSvc, ts.prodSvc)
	hndlr.FetchProduct(c)
	ts.Assert().EqualValues(http.StatusOK, w.Code)

	data, err := io.ReadAll(w.Result().Body)
	ts.Assert().Nil(err)
	ts.Assert().EqualValues(data, expectedProductBytes)
}

func (ts *HttpTestSuite) TestHttpHandlers_FetchProductNotFound() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	productID := uuid.New()
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")
	c.AddParam(constants.ProductIDKey, productID.String())

	expectedResponse := []byte(`{"status_code":404,"error":"product not found"}`)

	ts.prodSvc.EXPECT().FetchProduct(gomock.Any(), productID).
		Times(1).
		Return(domain.Product{}, domain.ErrProductNotfound)

	hndlr := NewHTTPHandler(ts.subsSvc, ts.prodSvc)
	hndlr.FetchProduct(c)
	ts.Assert().EqualValues(http.StatusNotFound, w.Code)

	data, err := io.ReadAll(w.Result().Body)
	ts.Assert().Nil(err)
	ts.Assert().EqualValues(data, expectedResponse)
}

func (ts *HttpTestSuite) TestHttpHandlers_UpdateSubscription() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	productID, subsID := uuid.New(), uuid.New()
	subscription := domain.Subscription{
		ID:               uuid.New(),
		ProductID:        productID,
		DurationInMonths: 3,
		Tax:              0.77,
		TotalCost:        10.77,
		Status:           domain.SubscriptionStatusInactive,
		StartDate:        time.Now(),
		EndDate:          time.Now().AddDate(0, 3, 0),
	}

	type fetchSubscriptionMock struct {
		timesToCall int
		retSubs     domain.Subscription
		retErr      error
	}

	type updateSubscriptionStatusMock struct {
		timesToCall int
		retErr      error
	}

	tt := []struct {
		name             string
		subscriptionID   uuid.UUID
		status           domain.SubscriptionStatus
		expectedResponse []byte
		fsmock           fetchSubscriptionMock
		usmock           updateSubscriptionStatusMock
	}{
		{
			name:             "Update Subscription status success",
			subscriptionID:   subsID,
			status:           domain.SubscriptionStatusActive,
			expectedResponse: []byte(`{"message":"Successfully updated the subscription status.","status_code":200}`),
			fsmock: fetchSubscriptionMock{
				timesToCall: 1,
				retSubs:     subscription,
				retErr:      nil,
			},
			usmock: updateSubscriptionStatusMock{
				timesToCall: 1,
				retErr:      nil,
			},
		},
	}

	for _, tc := range tt {
		ts.Run(tc.name, func() {
			r := &http.Request{
				Header: make(http.Header),
				URL: &url.URL{
					RawQuery: "status=" + tc.status.String(),
				},
			}
			r.Method = "PATCH"
			r.Header.Set("Content-Type", "application/json")
			c.AddParam(constants.SubscriptionIDKey, subsID.String())
			c.Request = r
			ts.subsSvc.EXPECT().FetchSubscription(gomock.Any(), tc.subscriptionID).
				Times(tc.fsmock.timesToCall).
				Return(tc.fsmock.retSubs, tc.fsmock.retErr)

			ts.subsSvc.EXPECT().UpdateSubscriptionStatus(gomock.Any(), tc.subscriptionID, tc.status).
				Times(tc.usmock.timesToCall).
				Return(tc.usmock.retErr)

			hndlr := NewHTTPHandler(ts.subsSvc, ts.prodSvc)
			hndlr.UpdateSubscriptionStatus(c)

			data, err := io.ReadAll(w.Result().Body)
			ts.Assert().Nil(err)
			ts.Assert().EqualValues(data, tc.expectedResponse)
		})
	}
}
