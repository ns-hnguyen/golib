package genapi

import (
	. "testing"

	"github.com/levenlabs/golib/testutil"
	"github.com/mediocregopher/lever"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfigurator struct {
	name string
	val  string
}

func (tc *testConfigurator) Params() []lever.Param {
	return []lever.Param{{Name: tc.name}}
}

func (tc *testConfigurator) WithParams(l *lever.Lever) {
	tc.val, _ = l.ParamStr(tc.name)
}

func TestConfig(t *T) {
	c := Config{Name: "test"}
	a := testConfigurator{name: testutil.RandStr()}
	aVal := testutil.RandStr()
	require.Nil(t, os.SetEnv("test_"+a.name, aVal))

	b := testConfigurator{name: testutil.RandStr()}
	bVal := testutil.RandStr()
	require.Nil(t, os.SetEnv("test_"+b.name, bVal))

	c.Add(a, b)
	assert.Nil(t, c.configurate())

	assert.Equal(t, aVal, a.val)
	assert.Equal(t, bVal, b.val)

	// Reset and add a again, configurate should fail due to duplicate parameter
	a.val = ""
	b.val = ""
	c.Add(a)
	assert.NotNil(t, c.configurate())
}
