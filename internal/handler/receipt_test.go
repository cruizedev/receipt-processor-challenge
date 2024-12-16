package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cruizedev/receipt-processor-challenge/internal/handler"
	"github.com/cruizedev/receipt-processor-challenge/internal/response"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type ReceiptSuite struct {
	suite.Suite

	engine *echo.Echo
}

func (suite *ReceiptSuite) SetupSuite() {
	suite.engine = echo.New()

	r := handler.NewReceipt()
	r.Register(suite.engine.Group(""))
}

func (suite *ReceiptSuite) TestCount() {
	require := suite.Require()

	b := []byte(`
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },
    {
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },
    {
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },
    {
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },
    {
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
	`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/process", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.engine.ServeHTTP(w, req)
	require.Equal(http.StatusOK, w.Code)

	var id response.ID
	require.NoError(json.NewDecoder(w.Body).Decode(&id))

	// figure out the status of receipt by calling points endpoint.
	{
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/points", id.ID), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		suite.engine.ServeHTTP(w, req)
		for w.Code == http.StatusAccepted {
			suite.engine.ServeHTTP(w, req)
		}

		require.Equal(http.StatusOK, w.Code)

		var score response.Points
		require.NoError(json.NewDecoder(w.Body).Decode(&score))

		require.Equal(28, score.Points)
	}
}

func TestReceiptSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ReceiptSuite))
}
