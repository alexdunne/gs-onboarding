package gateway

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexdunne/gs-onboarding/internal/gateway/hackernews"
	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type itemsResponse struct {
	Items []models.Item `json:"items"`
}

func TestGetAllItems(t *testing.T) {
	type testcase struct {
		name               string
		hn                 *hackernews.Mock
		expectMocks        func(t *testing.T, hn *hackernews.Mock)
		expectedStatusCode int
		expectedItems      int
	}

	tests := []testcase{
		{
			name: "no items",
			hn:   &hackernews.Mock{},
			expectMocks: func(t *testing.T, hn *hackernews.Mock) {
				hn.On("FetchAll", mock.Anything).Return([]models.Item{}, nil)

			},
			expectedStatusCode: 200,
			expectedItems:      0,
		},
		{
			name: "one item",
			hn:   &hackernews.Mock{},
			expectMocks: func(t *testing.T, hn *hackernews.Mock) {
				hn.On("FetchAll", mock.Anything).Return([]models.Item{{ID: 1}}, nil)

			},
			expectedStatusCode: 200,
			expectedItems:      1,
		},
		{
			name: "two items",
			hn:   &hackernews.Mock{},
			expectMocks: func(t *testing.T, hn *hackernews.Mock) {
				hn.On("FetchAll", mock.Anything).Return([]models.Item{{ID: 1}, {ID: 2}}, nil)

			},
			expectedStatusCode: 200,
			expectedItems:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectMocks != nil {
				tt.expectMocks(t, tt.hn)
			}

			context, res := setUpRequest(http.MethodPost, "/all")

			h := Handler{
				HNClient: tt.hn,
			}

			err := h.GetAllItems(context)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, res.Code)

			var resBody itemsResponse
			err = json.Unmarshal(res.Body.Bytes(), &resBody)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedItems, len(resBody.Items))
		})
	}
}

func TestGetStories(t *testing.T) {
	type testcase struct {
		name               string
		hn                 *hackernews.Mock
		expectMocks        func(t *testing.T, hn *hackernews.Mock)
		expectedStatusCode int
		expectedItems      int
	}

	tests := []testcase{
		{
			name: "no items",
			hn:   &hackernews.Mock{},
			expectMocks: func(t *testing.T, hn *hackernews.Mock) {
				hn.On("FetchStories", mock.Anything).Return([]models.Item{}, nil)

			},
			expectedStatusCode: 200,
			expectedItems:      0,
		},
		{
			name: "one item",
			hn:   &hackernews.Mock{},
			expectMocks: func(t *testing.T, hn *hackernews.Mock) {
				hn.On("FetchStories", mock.Anything).Return([]models.Item{{ID: 1}}, nil)

			},
			expectedStatusCode: 200,
			expectedItems:      1,
		},
		{
			name: "two items",
			hn:   &hackernews.Mock{},
			expectMocks: func(t *testing.T, hn *hackernews.Mock) {
				hn.On("FetchStories", mock.Anything).Return([]models.Item{{ID: 1}, {ID: 2}}, nil)

			},
			expectedStatusCode: 200,
			expectedItems:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectMocks != nil {
				tt.expectMocks(t, tt.hn)
			}

			context, res := setUpRequest(http.MethodPost, "/all")

			h := Handler{
				HNClient: tt.hn,
			}

			err := h.GetStories(context)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, res.Code)

			var resBody itemsResponse
			err = json.Unmarshal(res.Body.Bytes(), &resBody)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedItems, len(resBody.Items))
		})
	}
}

func TestGetJobs(t *testing.T) {
	type testcase struct {
		name               string
		hn                 *hackernews.Mock
		expectMocks        func(t *testing.T, hn *hackernews.Mock)
		expectedStatusCode int
		expectedItems      int
	}

	tests := []testcase{
		{
			name: "no items",
			hn:   &hackernews.Mock{},
			expectMocks: func(t *testing.T, hn *hackernews.Mock) {
				hn.On("FetchJobs", mock.Anything).Return([]models.Item{}, nil)

			},
			expectedStatusCode: 200,
			expectedItems:      0,
		},
		{
			name: "one item",
			hn:   &hackernews.Mock{},
			expectMocks: func(t *testing.T, hn *hackernews.Mock) {
				hn.On("FetchJobs", mock.Anything).Return([]models.Item{{ID: 1}}, nil)

			},
			expectedStatusCode: 200,
			expectedItems:      1,
		},
		{
			name: "two items",
			hn:   &hackernews.Mock{},
			expectMocks: func(t *testing.T, hn *hackernews.Mock) {
				hn.On("FetchJobs", mock.Anything).Return([]models.Item{{ID: 1}, {ID: 2}}, nil)

			},
			expectedStatusCode: 200,
			expectedItems:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectMocks != nil {
				tt.expectMocks(t, tt.hn)
			}

			context, res := setUpRequest(http.MethodPost, "/all")

			h := Handler{
				HNClient: tt.hn,
			}

			err := h.GetJobs(context)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, res.Code)

			var resBody itemsResponse
			err = json.Unmarshal(res.Body.Bytes(), &resBody)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedItems, len(resBody.Items))
		})
	}
}

func setUpRequest(method string, endpoint string) (echo.Context, *httptest.ResponseRecorder) {
	router := echo.New()

	request := httptest.NewRequest(method, endpoint, nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()

	context := router.NewContext(request, response)

	return context, response
}
