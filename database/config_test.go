package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyDefaultValue(t *testing.T) {
	defaultConnPool := &DBConnPool{
		MaxIdleConns:    defaultMaxIdleConns,
		ConnMaxIdleTime: defaultConnMaxIdleTime,
		MaxOpenConns:    defaultMaxOpenConns,
		ConnMaxLifetime: defaultConnMaxLifetime,
	}
	cases := map[string]struct {
		cfg         DBConfig
		expected    options
		expectedCfg DBConfig
	}{
		"1. Default": {
			cfg: DBConfig{},
			expected: options{
				ResolverPolicy: RandomPolicy,
			},
			expectedCfg: DBConfig{
				LogLevel:       LogLevelError,
				ResolverPolicy: DefaultPolicyName,
				ConnPool:       defaultConnPool,
			},
		},
		"2. Configured": {
			cfg: DBConfig{
				LogLevel:       LogLevelInfo,
				ResolverPolicy: RoundRobinPolicyName,
			},
			expectedCfg: DBConfig{
				LogLevel:       LogLevelInfo,
				ResolverPolicy: RoundRobinPolicyName,
				ConnPool:       defaultConnPool,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			cfg := tc.cfg
			opt := cfg.applyDefaultValue()

			assert.NotNil(t, opt.ResolverPolicy)
			assert.EqualValues(t, tc.expectedCfg, cfg)
		})
	}
}
