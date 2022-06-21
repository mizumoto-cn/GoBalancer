// This file tests the ip_hash balancer
package balancer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test dataset
const (
	test_param0 = "http://127.0.0.1:1011"
	test_param1 = "http://127.0.0.1:1012"
	test_param2 = "http://127.0.0.1:1013"
	test_param3 = "http://127.0.0.1:1014"
)

// Test IPHash::AddHost
func TestIPHash_AddHost(t *testing.T) {
	var (
		ih_1, _ = NewIPHash([]string{test_param0, test_param1, test_param2})
		ih_2, _ = NewIPHash([]string{test_param0, test_param1})
	)
	type expect struct {
		balancer Balancer
		err      error
	}
	cases := []struct {
		name     string
		balancer Balancer
		args     string
		expect   expect
	}{
		{
			"test-1",
			ih_1,
			test_param2,
			expect{ih_1, nil},
		},
		{
			"test-2",
			ih_2,
			test_param2,
			expect{ih_1, nil},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.balancer.AddHost(c.args)
			assert.Equal(t, c.expect.err, err)
			assert.Equal(t, c.expect.balancer, c.balancer)
		})
	}
}

// Test IPHash::RemoveHost
func TestIPHash_RemoveHost(t *testing.T) {
	var (
		ih_1, _ = NewIPHash([]string{test_param0, test_param1, test_param2})
		ih_2, _ = NewIPHash([]string{test_param0, test_param1})
	)
	type expect struct {
		balancer Balancer
		err      error
	}
	cases := []struct {
		name     string
		balancer Balancer
		args     string
		expect   expect
	}{
		{
			"test-1",
			ih_1,
			test_param1,
			expect{&IPHash{hosts: []string{test_param0, test_param2}}, nil},
		},
		{
			"test-2",
			ih_2,
			test_param2,
			expect{&IPHash{hosts: []string{test_param0, test_param1}}, nil},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.balancer.RemoveHost(c.args)
			assert.Equal(t, c.expect.err, err)
			assert.Equal(t, c.expect.balancer, c.balancer)
		})
	}
}

// Test IPHash::BalanceHost
func TestIPHash_BalanceHost(t *testing.T) {
	var (
		ih_1, _   = NewIPHash([]string{test_param0, test_param1, test_param2})
		ih_nil, _ = NewIPHash([]string{})
	)
	type expect struct {
		reply string
		err   error
	}
	cases := []struct {
		name     string
		balancer Balancer
		key      string
		expect   expect
	}{
		{
			"test-1",
			ih_1,
			"192.168.1.1",
			expect{test_param0, nil},
		},
		{
			"test-2",
			ih_nil,
			"192.168.1.1",
			expect{"", ErrHostNotFound},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reply, err := c.balancer.BalanceHost(c.key)
			//assert.Equal(t, ih_1, IPHash{hosts: []string{param0, param1, param2}}.hosts)
			assert.Equal(t, c.expect.err, err)
			assert.Equal(t, c.expect.reply, reply)
		})
	}
}

// Test 11 failed to be checked
