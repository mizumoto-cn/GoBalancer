// This file tests the ip_hash balancer
package balancer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test dataset
var (
	param0 = "http://127.0.0.1:1011"
	param1 = "http://127.0.0.1:1012"
	param2 = "http://127.0.0.1:1013"
)

var (
	ih_1, _   = NewIPHash([]string{param0, param1, param2})
	ih_2, _   = NewIPHash([]string{param0, param1})
	ih_nil, _ = NewIPHash([]string{})
)

// Test IPHash::AddHost
func TestIPHash_AddHost(t *testing.T) {
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
			param2,
			expect{ih_1, nil},
		},
		{
			"test-2",
			ih_2,
			param2,
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
			param1,
			expect{&IPHash{hosts: []string{param0, param2}}, nil},
		},
		{
			"test-2",
			ih_2,
			param2,
			expect{&IPHash{hosts: []string{param0, param1}}, nil},
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
			expect{param0, nil},
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
			assert.Equal(t, c.expect.err, err)
			assert.Equal(t, c.expect.reply, reply)
		})
	}
}

// Test 11 failed to be checked
