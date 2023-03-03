package gsrv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	t.Run("bad_config", func(t *testing.T) {
		s, err := New("badaddress")
		require.Error(t, err)
		require.Nil(t, s)
	})

	t.Run("default_configs", func(t *testing.T) {
		address := "0.0.0.0:8080"
		s, err := New(address)
		require.NoError(t, err)
		require.NotNil(t, s)
		require.NotNil(t, s.logger)
		require.Equal(t, s.address, address)
		require.Equal(t, s.timeout, defaultTimeout)
		require.NotNil(t, s.listener)
	})

	t.Run("custom_timeout", func(t *testing.T) {
		timeout := 10 * time.Minute
		s, err := New("0.0.0.0:8081", WithTimeout(timeout))
		require.NoError(t, err)
		require.NotNil(t, s)
		require.Equal(t, s.timeout, timeout)
	})
}
