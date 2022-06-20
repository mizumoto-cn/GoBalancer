// this file tests the random balancer
package balancer

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test Random::AddHost
func TestRandom_AddHost(t *testing.T) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []struct {
		name     string
		balancer Balancer
		args     string
		expect   Balancer
	}{
		{
			"test-1",
			&Random{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1013"}, rand: rand},
			"http://127.0.0.1:1013",
			&Random{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1013"}, rand: rand},
		},
		{
			"test-2",
			&Random{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"}, rand: rand},
			"http://127.0.0.1:1012",
			&Random{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"}, rand: rand},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.balancer.AddHost(c.args)
			assert.Equal(t, c.expect, c.balancer)
		})
	}
}

// Test Random::RemoveHost
func TestRandom_RemoveHost(t *testing.T) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []struct {
		name     string
		balancer Balancer
		args     string
		expect   Balancer
	}{
		{
			"test-1",
			&Random{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1013"}, rand: rand},
			"http://127.0.0.1:1013",
			&Random{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.01:1012"}, rand: rand},
		},
		{
			"test-2",
			&Random{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"}, rand: rand},
			"http://127.0.0.1:1013",
			&Random{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"}, rand: rand},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.balancer.RemoveHost(c.args)
			assert.Equal(t, c.expect, c.balancer)
		})
	}
}

// Test Random::BalanceHost
func TestRandom_BalanceHost(t *testing.T) {
	type expect struct {
		reply string
		err   error
	}
	b1, _ := NewRandom([]string{"http://127.0.0.1:1011"})
	b2, _ := NewRandom([]string{})
	cases := []struct {
		name     string
		balancer Balancer
		args     string
		expect   expect
	}{
		{
			"test-1",
			b1,
			"",
			expect{reply: "http://127.0.0.1:1011", err: nil},
		},
		{
			"test-2",
			b2,
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
