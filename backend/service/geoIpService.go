package service

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"ip-lookup-app/entity"
	"ip-lookup-app/util"
	"net"
)

type GeoIpService struct {
	file   string
	locale string
}

func NewGeoIpService() *GeoIpService {
	return &GeoIpService{
		file:   util.GetEnv("MAXMIND_DB"),
		locale: util.GetEnv("LOCALE"),
	}
}

func (s *GeoIpService) GetLocation(ip string, ipLocation *entity.IpLocation) error {
	db, err := geoip2.Open(s.file)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	netIP := net.ParseIP(ip)
	if netIP == nil {
		return fmt.Errorf("invalid IP address: %v", err)
	}

	record, err := db.City(netIP)
	if err != nil {
		return fmt.Errorf("failed to lookup IP location: %v", err)
	}

	// Create a new Location struct with the data from the database
	ipLocation.Ip = ip
	ipLocation.City = record.City.Names[s.locale]
	ipLocation.Timezone = record.Location.TimeZone
	ipLocation.AccuracyRadius = record.Location.AccuracyRadius
	ipLocation.PostalCode = record.Postal.Code

	return nil
}
