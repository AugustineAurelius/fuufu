package geo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/AugustineAurelius/fuufu/pkg/cache"
)

type GeoIPResponse struct {
	Query         string  `json:"query"`
	Status        string  `json:"status"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Timezone      string  `json:"timezone"`
	Isp           string  `json:"isp"`
	Org           string  `json:"org"`
	As            string  `json:"as"`
	Asname        string  `json:"asname"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
	Hosting       bool    `json:"hosting"`
}

type GeoIPService struct {
	URL      *url.URL
	cache    *cache.Weak[string, GeoIPResponse]
	cacheRaw *cache.Weak[string, []byte]
}

func NewGeoIPService() *GeoIPService {
	return &GeoIPService{
		URL: &url.URL{
			Scheme:   "http",
			Host:     "ip-api.com",
			RawQuery: "fields=24903679",
			Path:     "json",
		},
		cache:    cache.NewWeak[string, GeoIPResponse](),
		cacheRaw: cache.NewWeak[string, []byte](),
	}
}

func (s *GeoIPService) GetLocation(ip string) (*GeoIPResponse, error) {
	res, ok := s.cache.Get(ip)
	if ok {
		return res, nil
	}

	url := s.URL.JoinPath(ip)
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GeoIPResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	s.cache.Set(ip, result)

	return &result, nil
}

func (s *GeoIPService) GetLocationRaw(ip string) ([]byte, error) {
	res, ok := s.cacheRaw.Get(ip)
	if ok {
		return *res, nil
	}

	url := s.URL.JoinPath(ip)
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.cacheRaw.Set(ip, result)

	return result, nil
}
