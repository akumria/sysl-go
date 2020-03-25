package config

import (
	"testing"
	"time"

	"github.com/alecthomas/assert"
)

type TestMyConfig struct {
	Server TestServer `yaml:"server"`
}

type TestServer struct {
	AdminServer TestAdminServer `yaml:"adminServer"`
}

type TestAdminServer struct {
	ContextTimeout time.Duration `yaml:"contextTimeout"`
	HTTP           TestHTTP      `yaml:"http"`
}

type TestHTTP struct {
	BasePath     string `yaml:"basePath"`
	ReadTimeout  string `yaml:"readTimeout"`
	WriteTimeout string `yaml:"writeTimeout"`
}

type TestDownstreamConfig struct {
	ContextTimeout time.Duration        `yaml:"contextTimeout"`
	Fenergo        CommonDownstreamData `yaml:"fenergo"`
	Qas            CommonDownstreamData `yaml:"qas"`
}

func TestSReadConfig(t *testing.T) {
	t.Parallel()

	lib := LibraryConfig{}
	gen := GenCodeConfig{}
	gen.Downstream = &TestDownstreamConfig{}
	myConfig := TestMyConfig{}
	err := ReadConfig("testdata/config.yaml", &lib, &gen, &myConfig)

	assert.Nil(t, err)

	assert.Equal(t, 2*time.Second, myConfig.Server.AdminServer.ContextTimeout)
	assert.Equal(t, "/admintest", myConfig.Server.AdminServer.HTTP.BasePath)

	assert.False(t, lib.Log.ReportCaller)

	assert.Equal(t, 8080, gen.Upstream.HTTP.Common.Port)
	assert.Equal(t, 8081, gen.Upstream.GRPC.Port)
	assert.Equal(t, 120*time.Second, gen.Downstream.(*TestDownstreamConfig).Fenergo.ClientTimeout)
	assert.Equal(t, "https://stubs.bah.apps.x.gcpnp.anz", gen.Downstream.(*TestDownstreamConfig).Qas.ServiceURL)
}
