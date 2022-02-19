package beer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)


var (
	beer1 = RepoBeer{"1111", "Beer One",   "de", time.Date(2021, time.February, 01, 01, 0, 0, 0, time.UTC)}
	beer2 = RepoBeer{"2222", "Beer Two",   "de", time.Date(2021, time.February, 01, 01, 0, 0, 0, time.UTC)}
	beer3 = RepoBeer{"3333", "Beer Three", "de", time.Date(2021, time.February, 01, 01, 0, 0, 0, time.UTC)}
	beer4 = RepoBeer{"4444", "Beer Four",  "de", time.Date(2021, time.February, 01, 01, 0, 0, 0, time.UTC)}
	beer5 = RepoBeer{"5555", "Beer Five",  "de", time.Date(2021, time.February, 01, 01, 0, 0, 0, time.UTC)}
	expectedbeers = []RepoBeer{beer1, beer2, beer3, beer4, beer5}
)

type MockRepository struct {}
func (m *MockRepository) GetBeers() ([]RepoBeer, error)	{
	return expectedbeers, nil
}

func (m *MockRepository) GetBeer(id string) (*RepoBeer, error) {
	for _, b := range expectedbeers {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, ErrNotFound
}

func (m *MockRepository) CreateBeer(name, country string) (*RepoBeer, error) {
	return &RepoBeer{"121212", name, country, time.Now()}, nil
}


func TestAcceptMethods(t *testing.T) {

	var testdata = []struct {
		request *http.Request
		expectedError error
		message string
	}{
		{httptest.NewRequest(http.MethodGet, "/beers", nil), nil, "method GET shall be supported"},
		{httptest.NewRequest(http.MethodPost, "/beers", nil), nil, "method POST shall be supported"},
		{httptest.NewRequest(http.MethodDelete, "/beers", nil), ErrUnknownMethod, "method DELETE shall not not supported"},
		{httptest.NewRequest(http.MethodPut, "/beers", nil), ErrUnknownMethod, "method PUT shall not not supported"},
		{httptest.NewRequest(http.MethodHead, "/beers", nil), ErrUnknownMethod, "method PUT shall not not supported"},
	}

	srv := &BeerService{repository: &MockRepository{}}

	for _, tt := range testdata {
		_, actualErrror := handleRequest(srv, tt.request)
		require.Equal(t, tt.expectedError, actualErrror, tt.message)
	}

}

func TestGetBeers(t *testing.T) {
	srv := &BeerService{repository: &MockRepository{}}
	actualbeers, err := srv.GetBeers()
	require.NoError(t, err, "return a list of valid beers of the required size")
	require.Equal(t, len(actualbeers), len(expectedbeers))

}

func TestGetBeer(t *testing.T) {
	var testdata = []struct {
		id string
		expectedErr error
		message string
	} {
		{"1111", nil, "requesting with an valid id shall be successful"},
		{"9999", ErrNotFound, "requesting with an valid id shall return not-found error"},
	}

	srv := &BeerService{repository: &MockRepository{}}

	for _, tt := range testdata {
		_, actualErr := srv.GetBeer(tt.id)
		require.Equal(t, tt.expectedErr, actualErr)
	}
}

func TestCreateBeer(t *testing.T) {


	var (
		name = "My own beer"
		country = "de"
	)
	
	srv := &BeerService{repository: &MockRepository{}}

	res, err := srv.CreateBeer(name, country)
	require.NoError(t, err, "creating a beer shall be successful")
	require.Equal(t, res.Name, name, "shall return the valid name")
	require.Equal(t, res.Name, name, "shall return the valid country")
	require.Equal(t, time.Now().After(res.CreatedAt), true, "shall have a creation data in the past")

}

func TestGetBeersHTTP(t *testing.T) {

	repository := MockRepository{}
	path, handleFunc, err := NewHandler(&repository)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, path, nil)
	actualResponse := httptest.NewRecorder()

	handleFunc(actualResponse, request)
	require.Equal(t, http.StatusOK, actualResponse.Code)
	actualBeers := []RepoBeer{}
	unmarshalErr := json.Unmarshal(actualResponse.Body.Bytes(), &actualBeers)
	require.NoError(t, unmarshalErr)

	require.Equal(t, len(actualBeers), len(expectedbeers), "shall return the correct number of beers")


}

func TestGetSingleHTTP(t *testing.T) {

	expectedbeer := expectedbeers[1]

	repository := MockRepository{}
	path, handleFunc, err := NewHandler(&repository)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, path + "/" + expectedbeer.ID, nil)
	actualResponse := httptest.NewRecorder()

	handleFunc(actualResponse, request)
	require.Equal(t, http.StatusOK, actualResponse.Code)
	actualBeers := []RepoBeer{}
	unmarshalErr := json.Unmarshal(actualResponse.Body.Bytes(), &actualBeers)
	require.NoError(t, unmarshalErr)
	require.Equal(t, 1, len(actualBeers))

	actualBeer := actualBeers[0]
	require.Equal(t, actualBeer.ID, expectedbeer.ID, "shall return the exect beer requested")


}


func TestExtractEntityID(t *testing.T) {

	var testdata = []struct {
		path string
		expectedID string
		expectedResult bool
		message string
	} {
		{"/beers", "", false, "shall detect that no id is given"},
		{"/beers/", "", false, "shall detect that no id is given"},
		{"/beers/1111", "1111", true, "shall detect valid id"},
		{"/beers/2222/", "2222", true, "shall detect valid id"},
		{"/beers/2222/someprefix", "2222/someprefix", true, "handles invalid urls in a consistent way"},
	}

	for _, tt := range testdata {	
		res, id := extractEntityID(tt.path)
		require.Equal(t, tt.expectedResult, res, tt.message)
		require.Equal(t, tt.expectedID, id, tt.message)
	}

}