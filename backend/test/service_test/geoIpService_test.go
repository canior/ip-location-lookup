package service_test

import (
	"github.com/stretchr/testify/assert"
	"ip-lookup-app/entity"
	"ip-lookup-app/service"
	"os"
	"testing"
)

func TestGeoIpService(t *testing.T) {
	os.Setenv("APP_ENV", "testing")
	service := service.NewGeoIpService()
	var ipLocation entity.IpLocation
	ip := "216.209.131.154"
	service.GetLocation("216.209.131.154", &ipLocation)
	assert.Equal(t, ip, ipLocation.Ip)
}
