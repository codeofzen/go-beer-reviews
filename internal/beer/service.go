package beer

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound = errors.New("beer not found")
	ErrUnknownMethod = errors.New("unknown method")
)


type RepoBeer struct {
	ID 			string 		`json:"id"`
	Name		string 		`json:"name"`
	CountryISO 	string 		`json:"countryISO"`
	CreatedAt 	time.Time 	`json:"created_at"`
}


type Repository interface {
	GetBeers() ([]RepoBeer, error)	
	GetBeer(id string) (*RepoBeer, error)
	CreateBeer(name, country string) (*RepoBeer, error)
}


type BeerService struct {
	repository	Repository
}

func (b *BeerService) GetBeers() ([]RepoBeer, error) {
	return b.repository.GetBeers()
}

func (b *BeerService) GetBeer(id string) (*RepoBeer, error) {
	return b.repository.GetBeer(id)
}

func (b *BeerService) CreateBeer(name, country string) (*RepoBeer, error) {
	return b.repository.CreateBeer(name, country)
}

func NewHandler(r Repository) (string, http.HandlerFunc, error) {

	srv := BeerService{
		repository: r,
	}

	f := func(w http.ResponseWriter, r *http.Request) {
	
		w.Header().Set("content-type", "application/json")

		results, err := handleRequest(&srv, r)

		if err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		b, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)

	}

	return "/beers", f, nil
}


func handleRequest(srv *BeerService, req *http.Request) ([]RepoBeer, error) {
	switch method := req.Method; method {
		
		case http.MethodGet:
			path := req.URL.Path
			log.Printf("GET - %s\n", path)

			// handle request for specific beer: /beers/{id}
			if res, id := extractEntityID(req.URL.Path); res {
				beer, err := srv.GetBeer(id)
				if err != nil {
					return nil, err
				} else {
					return []RepoBeer{*beer}, nil
				}
			}

			// handle requests for all beers
			beers, err := srv.GetBeers()
			if err != nil {
				return nil, ErrInternalServerError
			}
			return beers, nil

			
		case http.MethodPost:
			path := req.URL.Path
			log.Printf("POST - %s\n", path)
			beer, err := srv.CreateBeer("Fake Beer", "uk")
			if err != nil {
				return nil, ErrInternalServerError
			}
			return []RepoBeer{*beer}, nil
		
		default:
			log.Printf("Unknown method - %s\n", method)
			return nil, ErrUnknownMethod
	}
}


func extractEntityID(path string) (bool, string) {

	res := strings.TrimPrefix(path, "/beers")
	if len(res) == 0 || len(res) == 1 {
		return false, ""
	}

	if res[len(res)-1:] == "/" {
		return true, res[1:len(res)-1] 
	}

	return true, res[1:]
}