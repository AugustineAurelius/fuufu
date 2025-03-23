package geo_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/AugustineAurelius/fuufu/pkg/geo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_URL(t *testing.T) {
	ur := &url.URL{
		Scheme:   "http",
		Host:     "ip-api.com",
		RawQuery: "fields=24903679",
		Path:     "json",
	}
	ur = ur.JoinPath("24.48.0.1")

	assert.Equal(t, "http://ip-api.com/json/24.48.0.1?fields=24903679", ur.String())
}

func Test_GetLocation(t *testing.T) {

	expected := &geo.GeoIPResponse{
		Query:         "24.48.0.1",
		Status:        "success",
		Continent:     "North America",
		ContinentCode: "NA",
		Country:       "Canada",
		CountryCode:   "CA",
		Region:        "QC",
		RegionName:    "Quebec",
		City:          "Montreal",
		District:      "",
		Zip:           "H1K",
		Lat:           45.6085,
		Lon:           -73.5493,
		Timezone:      "America/Toronto",
		Isp:           "Le Groupe Videotron Ltee",
		Org:           "Videotron Ltee",
		As:            "AS5769 Videotron Ltee",
		Asname:        "VIDEOTRON",
		Mobile:        false,
		Proxy:         false,
		Hosting:       false,
	}

	s := geo.NewGeoIPService()
	res, err := s.GetLocation("24.48.0.1")
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func Test_GetLocationRaw(t *testing.T) {

	expected := &geo.GeoIPResponse{
		Query:         "24.48.0.1",
		Status:        "success",
		Continent:     "North America",
		ContinentCode: "NA",
		Country:       "Canada",
		CountryCode:   "CA",
		Region:        "QC",
		RegionName:    "Quebec",
		City:          "Montreal",
		District:      "",
		Zip:           "H1K",
		Lat:           45.6085,
		Lon:           -73.5493,
		Timezone:      "America/Toronto",
		Isp:           "Le Groupe Videotron Ltee",
		Org:           "Videotron Ltee",
		As:            "AS5769 Videotron Ltee",
		Asname:        "VIDEOTRON",
		Mobile:        false,
		Proxy:         false,
		Hosting:       false,
	}

	s := geo.NewGeoIPService()
	res, err := s.GetLocationRaw("24.48.0.1")
	require.NoError(t, err)
	var response geo.GeoIPResponse
	err = json.Unmarshal(res, &response)
	require.NoError(t, err)
	assert.Equal(t, expected, &response)
}
