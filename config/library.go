package config

import (
	"encoding/base64"
	"time"

	"github.com/anz-bank/sysl-go/common"
	"github.com/anz-bank/sysl-go/jwtauth"
	"github.com/anz-bank/sysl-go/validator"
	"github.com/sirupsen/logrus"
)

// LibraryConfig struct.
type LibraryConfig struct {
	Log            LogConfig             `yaml:"log"`
	Profiling      bool                  `yaml:"profiling"`
	Health         bool                  `yaml:"health"`
	Authentication *AuthenticationConfig `yaml:"authentication"`
}

type AdminConfig struct {
	ContextTimeout time.Duration          `yaml:"contextTimeout" validate:"nonnil"`
	HTTP           CommonHTTPServerConfig `yaml:"http"`
}

// LogConfig struct.
type LogConfig struct {
	Format       string        `yaml:"format" validate:"nonnil,oneof=color json text"`
	Splunk       *SplunkConfig `yaml:"splunk"`
	Level        logrus.Level  `yaml:"level" validate:"nonnil"`
	ReportCaller bool          `yaml:"caller" mapstructure:"caller"`
}

// SplunkConfig struct.
type SplunkConfig struct {
	TokenBase64 common.SensitiveString `yaml:"tokenBase64" validate:"nonnil,base64"`
	Index       string                 `yaml:"index" validate:"nonnil"`
	Target      string                 `yaml:"target" validate:"nonnil,url"`
	Source      string                 `yaml:"source" validate:"nonnil"`
	SourceType  string                 `yaml:"sourceType" validate:"nonnil"`
}

// AuthenticationConfig struct.
type AuthenticationConfig struct {
	JWTAuth *jwtauth.Config `yaml:"jwtauth"`
}

func (s *SplunkConfig) Token() string {
	b, _ := base64.StdEncoding.DecodeString(s.TokenBase64.Value())

	return string(b)
}

func (c *LibraryConfig) Validate() error {
	// existing validation
	if err := validator.Validate(c); err != nil {
		return err
	}

	return nil
}
