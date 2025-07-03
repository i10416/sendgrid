package sendgrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPAddressStruct(t *testing.T) {
	ip := IPAddress{
		IP:         "192.168.1.1",
		Pools:      []string{"pool1", "pool2"},
		Warmup:     true,
		StartDate:  1609459200,
		Subusers:   []string{"subuser1"},
		Rdns:       "example.com",
		AssignedAt: 1609459300,
	}

	assert.Equal(t, "192.168.1.1", ip.IP)
	assert.Len(t, ip.Pools, 2)
	assert.Equal(t, "pool1", ip.Pools[0])
	assert.Equal(t, "pool2", ip.Pools[1])
	assert.True(t, ip.Warmup)
	assert.Equal(t, int64(1609459200), ip.StartDate)
	assert.Len(t, ip.Subusers, 1)
	assert.Equal(t, "subuser1", ip.Subusers[0])
	assert.Equal(t, "example.com", ip.Rdns)
	assert.Equal(t, int64(1609459300), ip.AssignedAt)
}

func TestIPPoolStruct(t *testing.T) {
	pool := IPPool{
		Name: "test-pool",
	}

	assert.Equal(t, "test-pool", pool.Name)
}

func TestIPWarmupStatusStruct(t *testing.T) {
	status := IPWarmupStatus{
		IP:     "192.168.1.1",
		Warmup: true,
	}

	assert.Equal(t, "192.168.1.1", status.IP)
	assert.True(t, status.Warmup)
}

func TestInputAddIPToPoolStruct(t *testing.T) {
	input := InputAddIPToPool{
		IP: "192.168.1.1",
	}

	assert.Equal(t, "192.168.1.1", input.IP)
}

func TestInputAssignIPToSubuserStruct(t *testing.T) {
	input := InputAssignIPToSubuser{
		IPs: []string{"192.168.1.1", "192.168.1.2"},
	}

	assert.Len(t, input.IPs, 2)
	assert.Equal(t, "192.168.1.1", input.IPs[0])
	assert.Equal(t, "192.168.1.2", input.IPs[1])
}

func TestOutputAssignedIPsStruct(t *testing.T) {
	output := OutputAssignedIPs{
		IPs: []string{"192.168.1.1", "192.168.1.2"},
	}

	assert.Len(t, output.IPs, 2)
	assert.Equal(t, "192.168.1.1", output.IPs[0])
	assert.Equal(t, "192.168.1.2", output.IPs[1])
}
