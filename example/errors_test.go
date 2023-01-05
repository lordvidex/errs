package example

import (
	"github.com/lordvidex/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Run("testing registration", func(t *testing.T) {
		assert.True(t, errs.IsRegistered(CustomErr))
		assert.True(t, errs.IsRegistered(SecondErr))
	})

	t.Run("testing tojson", func(t *testing.T) {
		byt, err := CustomErr.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`"This is custom error"`), byt)
	})
}
