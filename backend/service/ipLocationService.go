package service

import (
	"ip-lookup-app/entity"
	"ip-lookup-app/repository"
)

type IpLocationService struct {
	redisRepo    repository.RedisRepository
	geoIpService GeoIpService
}

func NewIpLocationService(redisRepo *repository.RedisRepository, geoIpService *GeoIpService) *IpLocationService {
	return &IpLocationService{
		redisRepo:    *redisRepo,
		geoIpService: *geoIpService,
	}
}

func (ils *IpLocationService) GetLocation(ip string) (entity.IpLocation, error) {
	// Try to get location from Redis cache
	var ipLocation entity.IpLocation
	err := ils.redisRepo.GetData(ip, &ipLocation)
	if err != nil {
		return ipLocation, err
	}

	// cache miss, try to get location from GeoIP database
	if ipLocation.Ip == "" {
		err := ils.geoIpService.GetLocation(ip, &ipLocation)
		if err != nil {
			return ipLocation, err
		}

		err = ils.redisRepo.WriteData(ip, ipLocation)
		if err != nil {
			return ipLocation, err
		}
	}

	return ipLocation, nil
}

func (ils *IpLocationService) GetBulkLocation(ips []string) ([]entity.IpLocation, error) {
	var locations []entity.IpLocation
	for _, ip := range ips {
		loc, err := ils.GetLocation(ip)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}
