package genapi

import (
	"os"
	. "testing"

	"github.com/levenlabs/golib/testutil"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallerStub(t *T) {
	method := "Test.Test"
	args := map[string]interface{}{
		testutil.RandStr(): testutil.RandStr(),
		testutil.RandStr(): testutil.RandStr(),
	}

	stub := CallerStub(func(method string, args interface{}) (interface{}, error) {
		return map[string]interface{}{
			"method": method,
			"args":   args,
		}, nil
	})
	res1 := map[string]interface{}{}
	assert.Nil(t, stub.Call(nil, &res1, method, args))
	assert.Equal(t, map[string]interface{}{
		"method": method,
		"args":   args,
	}, res1)

	type testType struct {
		method string
		args   map[string]interface{}
	}

	stub = CallerStub(func(method string, args interface{}) (interface{}, error) {
		return testType{method, args.(map[string]interface{})}, nil
	})
	res2 := testType{}
	assert.Nil(t, stub.Call(nil, &res2, method, args))
	assert.Equal(t, testType{method, args}, res2)

	assert.Nil(t, stub.Call(nil, nil, method, args))

	stub = CallerStub(func(method string, args interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Nil(t, stub.Call(nil, nil, method, args))
}

func TestInit(t *T) {
	i := 0
	var g *GenAPI
	g = &GenAPI{
		Init: func(g2 *GenAPI) {
			assert.Equal(t, g, g2)
			assert.Equal(t, 0, i)
			assert.Equal(t, TestMode, g2.Mode)
			i++
		},
	}
	g.AppendInit(func(g2 *GenAPI) {
		assert.Equal(t, g, g2)
		assert.Equal(t, 1, i)
		i++
	})
	g.AppendInit(func(g2 *GenAPI) {
		assert.Equal(t, g, g2)
		assert.Equal(t, 2, i)
		i++
	})
	g.TestMode()
	assert.Equal(t, 3, i)

	assert.Panics(t, func() {
		g.AppendInit(func(g2 *GenAPI) {})
	})
}

func TestSRVClientPreprocess(t *T) {
	dc := testutil.RandStr()
	os.Setenv("DATACENTER", dc)
	g := &GenAPI{}
	g.TestMode()
	dcHash := g.getDCHash()
	m := new(dns.Msg)
	m.Answer = []dns.RR{
		// The correct server with the local DC
		dns.RR(&dns.SRV{
			Target:   dcHash + "-" + testutil.RandStr(),
			Port:     uint16(80),
			Priority: uint16(5),
		}),
		// A server with a different DC
		dns.RR(&dns.SRV{
			Target:   testutil.RandStr() + "-" + testutil.RandStr(),
			Port:     uint16(80),
			Priority: uint16(5),
		}),
		// A server with no DC
		dns.RR(&dns.SRV{
			Target:   testutil.RandStr(),
			Port:     uint16(80),
			Priority: uint16(5),
		}),
	}

	g.srvClientPreprocess(m)

	correctPri := []uint16{uint16(4), uint16(5), uint16(5)}
	for i := range m.Answer {
		ansSRV, ok := m.Answer[i].(*dns.SRV)
		require.True(t, ok)
		assert.Equal(t, ansSRV.Priority, correctPri[i])
	}
}
