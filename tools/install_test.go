package tools_test

import (
	"testing"

	"get.porter.sh/magefiles/tools"
	"github.com/carolynvs/magex/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureKind(t *testing.T) {
	tools.EnsureKind()
	found, err := pkg.IsCommandAvailable("kind", tools.DefaultKindVersion, "--version")
	require.NoError(t, err)
	assert.True(t, found)
}
