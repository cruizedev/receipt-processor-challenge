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

	cases := []struct {
		request string
		score   int
	}{
		{
			score: 28,
			request: `
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
	`,
		},
		{
			score: 109,
			request: `
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
`,
		},
	}

	for i, c := range cases {
		suite.Run(fmt.Sprintf("case %d", i), func() {
			b := []byte(c.request)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/process", bytes.NewReader(b))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			require.Equal(http.StatusOK, w.Code)

			var id response.ID
			require.NoError(json.NewDecoder(w.Body).Decode(&id))

			suite.T().Logf("sending a process request and having %s", id)

			// figure out the status of receipt by calling points endpoint.
			{
				w := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/points", id.ID), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				suite.engine.ServeHTTP(w, req)
				for w.Code == http.StatusAccepted {
					w = httptest.NewRecorder()
					suite.engine.ServeHTTP(w, req)

					suite.T().Logf("watiing for status to change %d", w.Code)
				}

				require.Equal(http.StatusOK, w.Code)

				var score response.Points
				require.NoError(json.NewDecoder(w.Body).Decode(&score))

				require.Equal(c.score, score.Points)
			}
		})
	}
}

func TestReceiptSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ReceiptSuite))
}
