package geo

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type Geolocation struct {
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
}

var zero = Geolocation{}

func (x Geolocation) IsZero() bool {
	return x.Lat == 0 && x.Lng == 0
}

func ParseGeoLocation(latLng string) (Geolocation, error) {
	if latLng == "" {
		return zero, nil
	}
	x := strings.Split(latLng, ",")
	if len(x) != 2 {
		return zero, errors.New("location parse error")
	}
	lat, err := strconv.ParseFloat(x[0], 64)
	if err != nil {
		return zero, errors.WithMessage(err, "location.lat parse error")
	}
	lng, err := strconv.ParseFloat(x[1], 64)
	if err != nil {
		return zero, errors.WithMessage(err, "location.lng parse error")
	}
	return Geolocation{
		Lat: lat,
		Lng: lng,
	}, nil
}
