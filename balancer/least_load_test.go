//Tests least load balancer
package balancer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Test LeastLoad::BalanceHost
func TestLeastLoad_BalanceHost(t *testing.T) {

	expect, err := Build(LeastLoadBalancer, []string{test_param0,
		test_param1, test_param2, test_param3})
	assert.NoError(t, err)
	assert.NotNil(t, expect)
	heap := expect.(*LeastLoad).heap
	assert.NotNil(t, heap)
	assert.Equal(t, uint(4), heap.Num())
	// assert.Equal(t, test_param0, heap.MinimumValue().Tag().(string))
	assert.Equal(t, "Total number: 4, Root Size: 4, Index size: 4,\nCurrent minimum: key(0.000000), tag(http://127.0.0.1:1011), value(&{http://127.0.0.1:1011 0}),\nHeap detail:\n< 0.000000 0.000000 0.000000 0.000000 > \n",
		heap.String())
	boo, err := NewLeastLoad([]string{test_param0,
		test_param1, test_param2, test_param3})
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(expect, boo))

	assert.Equal(t, float64(0), heap.GetValue(test_param0).Key())
	assert.Equal(t, float64(0), heap.GetValue(test_param1).Key())
	assert.Equal(t, float64(0), heap.GetValue(test_param2).Key())
	assert.Equal(t, float64(0), heap.GetValue(test_param3).Key())

	err = expect.RemoveHost(test_param3)
	assert.NoError(t, err) //

	err = expect.Inc(test_param0)
	assert.NoError(t, err) //

	err = expect.Inc(test_param1)
	assert.NoError(t, err) //

	err = expect.Inc(test_param1)
	assert.NoError(t, err) //

	err = expect.Inc(test_param3)
	assert.Error(t, err)

	err = expect.Done(test_param3)
	assert.Error(t, err)

	err = expect.Done(test_param1)
	assert.NoError(t, err) //

	ll, _ := NewLeastLoad([]string{test_param1})
	ll.RemoveHost(test_param3)
	ll.AddHost(test_param0)
	ll.AddHost(test_param1)
	ll.AddHost(test_param2)
	ll.Inc(test_param0)
	ll.Inc(test_param1)
	ll.Inc(test_param1)
	ll.Done(test_param1)
	ll_host, err := ll.BalanceHost("")
	assert.NoError(t, err)
	expect_host, err := expect.BalanceHost("")
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(expect_host, ll_host))
}
