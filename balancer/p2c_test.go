// this file tests the p2c balancer
package balancer

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test PowerOfTwoChoices::AddHost
func TestPowerOfTwoChoices_AddHost(t *testing.T) {
	type expect struct {
		balancer Balancer
		err      error
	}
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	h0 := &host_info{test_param0, 0}
	h1 := &host_info{test_param1, 0}
	h2 := &host_info{test_param2, 0}
	cases := []struct {
		name   string
		b      Balancer
		args   string
		expect expect
	}{
		{
			"test-1",
			&PowerOfTwoChoices{hosts: []*host_info{h0, h1},
				rand: rand, loadMap: map[string]*host_info{test_param0: h0, test_param1: h1}},
			test_param0,
			expect{&PowerOfTwoChoices{hosts: []*host_info{h0, h1},
				rand: rand, loadMap: map[string]*host_info{test_param0: h0, test_param1: h1}}, ErrHostAlreadyExists},
		},
		{
			"test-2",
			&PowerOfTwoChoices{hosts: []*host_info{h0, h1},
				rand: rand, loadMap: map[string]*host_info{test_param0: h0, test_param1: h1}},
			test_param2,
			expect{&PowerOfTwoChoices{hosts: []*host_info{h0, h1, h2},
				rand: rand, loadMap: map[string]*host_info{test_param0: h0, test_param1: h1, test_param2: h2}}, nil},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.b.AddHost(c.args)
			assert.Equal(t, c.expect.balancer, c.b)
			assert.Equal(t, c.expect.err, err)
		})
	}
}

// Test PowerOfTwoChoices::RemoveHost
func TestPowerOfTwoChoices_RemoveHost(t *testing.T) {
	type expect struct {
		balancer Balancer
		err      error
	}
	h0 := &host_info{test_param0, 0}
	h1 := &host_info{test_param1, 0}
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []struct {
		name     string
		balancer Balancer
		args     string
		expect   expect
	}{
		{
			"test-1",
			&PowerOfTwoChoices{hosts: []*host_info{h0, h1}, rand: rand,
				loadMap: map[string]*host_info{test_param0: h0, test_param1: h1}},
			test_param0,
			expect{&PowerOfTwoChoices{hosts: []*host_info{h1},
				rand: rand, loadMap: map[string]*host_info{test_param1: h1}}, nil},
		},
		{
			"test-2",
			&PowerOfTwoChoices{hosts: []*host_info{h0, h1}, rand: rand, loadMap: map[string]*host_info{test_param0: h0, test_param1: h1}},
			test_param2,
			expect{&PowerOfTwoChoices{hosts: []*host_info{h0, h1},
				rand: rand, loadMap: map[string]*host_info{test_param0: h0, test_param1: h1}}, ErrHostNotFound},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.balancer.RemoveHost(c.args)
			assert.Equal(t, c.expect.balancer, c.balancer)
			assert.Equal(t, c.expect.err, err)
		})
	}
}

// Test PowerOfTwoChoices::BalanceHost
func TestPowerOfTwoChoices_BalanceHost(t *testing.T) {
	type expect struct {
		reply string
		err   error
	}
	p2c_4, _ := NewPowerOfTwoChoices([]string{test_param0, test_param1, test_param2, test_param3})
	//p2c_nil, _ := NewPowerOfTwoChoices([]string{})
	cases := []struct {
		name     string
		balancer Balancer
		args     string
		expect   expect
	}{
		{
			"test-1",
			p2c_4,
			"key",
			expect{test_param0, nil},
		},
		// {
		// 	"test-2",
		// 	p2c_nil,
		// 	"key",
		// 	expect{"", ErrHostNotFound},
		// },
		// ------> cannot do that as the balancer is nil
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			//		if c.name == "test-1" {
			c.balancer.Inc("http://127.0.0.1:1013")
			c.balancer.Inc("http://127.0.0.1:1013")
			c.balancer.Inc("http://127.0.0.1:1")
			c.balancer.Done("http://127.0.0.1:1")
			c.balancer.Done("http://127.0.0.1:1013")
			//		}
			host, err := c.balancer.BalanceHost(c.args)
			assert.Equal(t, c.expect.reply, host)
			assert.Equal(t, c.expect.err, err)
		})
	}
}
