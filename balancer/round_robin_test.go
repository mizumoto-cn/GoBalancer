// this file tests the round robin balancer.
package balancer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test RoundRobin::AddHost
func TestRoundRobin_AddHost(t *testing.T) {
	rr1, _ := NewRoundRobin([]string{"http://127.0.0.1:1011",
		"http://127.0.0.1:1012", "http://127.0.0.1:1013"})
	rr2, _ := NewRoundRobin([]string{"http://127.0.0.1:1011",
		"http://127.0.0.1:1012"})
	cases := []struct {
		name     string
		balancer Balancer
		args     string
		expect   Balancer
	}{
		{
			"test-1",
			rr1,
			"http://127.0.0.1:1013",
			&RoundRobin{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1013"}, i: 0},
		},
		{
			"test-2",
			rr2,
			"http://127.0.0.1:1014",
			&RoundRobin{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1014"}, i: 0},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.balancer.AddHost(c.args)
			assert.Equal(t, c.expect, c.balancer)
		})
	}
}

// Test RoundRobin::RemoveHost
func TestRoundRobin_RemoveHost(t *testing.T) {
	rr1, _ := NewRoundRobin([]string{"http://127.0.0.1:1011",
		"http://127.0.0.1:1012", "http://127.0.0.1:1013"})
	rr2, _ := NewRoundRobin([]string{"http://127.0.0.1:1011",
		"http://127.0.0.1:1012"})
	cases := []struct {
		name     string
		balancer Balancer
		args     string
		expect   Balancer
	}{
		{
			"test-1",
			rr1,
			"http://127.0.0.1:1013",
			&RoundRobin{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"}, i: 0},
		},
		{
			"test-2",
			rr2,
			"http://127.0.0.1:1014",
			&RoundRobin{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"}, i: 0},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.balancer.RemoveHost(c.args)
			assert.Equal(t, c.expect, c.balancer)
		})
	}
}

// Test RoundRobin::BalanceHost
func TestRoundRobin_BalanceHost(t *testing.T) {
	type expect struct {
		reply string
		err   error
	}
	rr1, _ := NewRoundRobin([]string{"http://127.0.0.1:1011"})
	rr2, _ := NewRoundRobin([]string{})
	cases := []struct {
		name     string
		balancer Balancer
		args     string
		expect   expect
	}{
		{
			"test-1",
			rr1,
			"",
			expect{reply: "http://127.0.0.1:1011", err: nil},
		},
		{
			"test-2",
			rr2,
			"",
			expect{reply: "", err: ErrHostNotFound},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reply, err := c.balancer.BalanceHost(c.args)
			assert.Equal(t, c.expect.reply, reply)
			assert.Equal(t, c.expect.err, err)
		})
	}
}
