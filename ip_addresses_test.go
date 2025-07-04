package sendgrid

import (
	"context"
	"net/http"
	"net/url"
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

func TestGetIPAddresses(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"ip":"192.168.1.1","pools":["pool1"],"warmup":true,"start_date":1609459200,"subusers":["subuser1"],"rdns":"example.com","assigned_at":1609459300}]`))
	})

	ctx := context.Background()
	ips, err := client.GetIPAddresses(ctx)

	assert.NoError(t, err)
	assert.Len(t, ips, 1)
	assert.Equal(t, "192.168.1.1", ips[0].IP)
	assert.Equal(t, []string{"pool1"}, ips[0].Pools)
	assert.True(t, ips[0].Warmup)
	assert.Equal(t, int64(1609459200), ips[0].StartDate)
	assert.Equal(t, []string{"subuser1"}, ips[0].Subusers)
	assert.Equal(t, "example.com", ips[0].Rdns)
	assert.Equal(t, int64(1609459300), ips[0].AssignedAt)
}

func TestGetAssignedIPAddresses(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/assigned", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"ip":"192.168.1.1","pools":["pool1"],"warmup":false}]`))
	})

	ctx := context.Background()
	ips, err := client.GetAssignedIPAddresses(ctx)

	assert.NoError(t, err)
	assert.Len(t, ips, 1)
	assert.Equal(t, "192.168.1.1", ips[0].IP)
	assert.Equal(t, []string{"pool1"}, ips[0].Pools)
	assert.False(t, ips[0].Warmup)
}

func TestGetRemainingIPCount(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/remaining", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"remaining":5,"total":10}`))
	})

	ctx := context.Background()
	result, err := client.GetRemainingIPCount(ctx)

	assert.NoError(t, err)
	assert.Equal(t, float64(5), result["remaining"])
	assert.Equal(t, float64(10), result["total"])
}

func TestGetIPAddress(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/192.168.1.1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ip":"192.168.1.1","pools":["pool1"],"warmup":true,"rdns":"example.com"}`))
	})

	ctx := context.Background()
	ip, err := client.GetIPAddress(ctx, "192.168.1.1")

	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", ip.IP)
	assert.Equal(t, []string{"pool1"}, ip.Pools)
	assert.True(t, ip.Warmup)
	assert.Equal(t, "example.com", ip.Rdns)
}

func TestGetIPPools(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"name":"pool1"},{"name":"pool2"}]`))
	})

	ctx := context.Background()
	pools, err := client.GetIPPools(ctx)

	assert.NoError(t, err)
	assert.Len(t, pools, 2)
	assert.Equal(t, "pool1", pools[0].Name)
	assert.Equal(t, "pool2", pools[1].Name)
}

func TestCreateIPPool(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"name":"new-pool"}`))
	})

	ctx := context.Background()
	pool, err := client.CreateIPPool(ctx, "new-pool")

	assert.NoError(t, err)
	assert.Equal(t, "new-pool", pool.Name)
}

func TestGetIPPool(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/test-pool", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"test-pool"}`))
	})

	ctx := context.Background()
	pool, err := client.GetIPPool(ctx, "test-pool")

	assert.NoError(t, err)
	assert.Equal(t, "test-pool", pool.Name)
}

func TestUpdateIPPool(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/old-pool", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"new-pool"}`))
	})

	ctx := context.Background()
	pool, err := client.UpdateIPPool(ctx, "old-pool", "new-pool")

	assert.NoError(t, err)
	assert.Equal(t, "new-pool", pool.Name)
}

func TestDeleteIPPool(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/test-pool", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	err := client.DeleteIPPool(ctx, "test-pool")

	assert.NoError(t, err)
}

func TestAddIPToPool(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/test-pool/ips", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
	})

	ctx := context.Background()
	err := client.AddIPToPool(ctx, "test-pool", "192.168.1.1")

	assert.NoError(t, err)
}

func TestRemoveIPFromPool(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/test-pool/ips/192.168.1.1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	err := client.RemoveIPFromPool(ctx, "test-pool", "192.168.1.1")

	assert.NoError(t, err)
}

func TestStartIPWarmup(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/warmup/192.168.1.1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ip":"192.168.1.1","warmup":true}`))
	})

	ctx := context.Background()
	status, err := client.StartIPWarmup(ctx, "192.168.1.1")

	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", status.IP)
	assert.True(t, status.Warmup)
}

func TestStopIPWarmup(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/warmup/192.168.1.1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ip":"192.168.1.1","warmup":false}`))
	})

	ctx := context.Background()
	status, err := client.StopIPWarmup(ctx, "192.168.1.1")

	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", status.IP)
	assert.False(t, status.Warmup)
}

func TestGetIPWarmupStatus(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/warmup/192.168.1.1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ip":"192.168.1.1","warmup":true}`))
	})

	ctx := context.Background()
	status, err := client.GetIPWarmupStatus(ctx, "192.168.1.1")

	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", status.IP)
	assert.True(t, status.Warmup)
}

func TestGetAllIPWarmupStatus(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/warmup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"ip":"192.168.1.1","warmup":true},{"ip":"192.168.1.2","warmup":false}]`))
	})

	ctx := context.Background()
	statuses, err := client.GetAllIPWarmupStatus(ctx)

	assert.NoError(t, err)
	assert.Len(t, statuses, 2)
	assert.Equal(t, "192.168.1.1", statuses[0].IP)
	assert.True(t, statuses[0].Warmup)
	assert.Equal(t, "192.168.1.2", statuses[1].IP)
	assert.False(t, statuses[1].Warmup)
}

// Error cases for IP addresses
func TestGetIPAddresses_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "Internal server error"}`))
	})

	ctx := context.Background()
	_, err := client.GetIPAddresses(ctx)

	assert.Error(t, err)
}

func TestGetAssignedIPAddresses_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/assigned", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "Internal server error"}`))
	})

	ctx := context.Background()
	_, err := client.GetAssignedIPAddresses(ctx)

	assert.Error(t, err)
}

func TestGetRemainingIPCount_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/remaining", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "Internal server error"}`))
	})

	ctx := context.Background()
	_, err := client.GetRemainingIPCount(ctx)

	assert.Error(t, err)
}

func TestGetIPAddress_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/192.168.1.1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "IP address not found"}`))
	})

	ctx := context.Background()
	_, err := client.GetIPAddress(ctx, "192.168.1.1")

	assert.Error(t, err)
}

func TestGetIPPools_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "Internal server error"}`))
	})

	ctx := context.Background()
	_, err := client.GetIPPools(ctx)

	assert.Error(t, err)
}

func TestCreateIPPool_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Invalid pool name"}`))
	})

	ctx := context.Background()
	_, err := client.CreateIPPool(ctx, "invalid-pool")

	assert.Error(t, err)
}

func TestGetIPPool_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/test-pool", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "Pool not found"}`))
	})

	ctx := context.Background()
	_, err := client.GetIPPool(ctx, "test-pool")

	assert.Error(t, err)
}

func TestUpdateIPPool_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/old-pool", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "Pool not found"}`))
	})

	ctx := context.Background()
	_, err := client.UpdateIPPool(ctx, "old-pool", "new-pool")

	assert.Error(t, err)
}

func TestDeleteIPPool_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/test-pool", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "Pool not found"}`))
	})

	ctx := context.Background()
	err := client.DeleteIPPool(ctx, "test-pool")

	assert.Error(t, err)
}

func TestAddIPToPool_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/test-pool/ips", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Invalid IP address"}`))
	})

	ctx := context.Background()
	err := client.AddIPToPool(ctx, "test-pool", "invalid-ip")

	assert.Error(t, err)
}

func TestRemoveIPFromPool_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/test-pool/ips/192.168.1.1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "IP not found in pool"}`))
	})

	ctx := context.Background()
	err := client.RemoveIPFromPool(ctx, "test-pool", "192.168.1.1")

	assert.Error(t, err)
}

func TestStartIPWarmup_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/warmup/192.168.1.1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "IP already in warmup"}`))
	})

	ctx := context.Background()
	_, err := client.StartIPWarmup(ctx, "192.168.1.1")

	assert.Error(t, err)
}

func TestStopIPWarmup_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/warmup/192.168.1.1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "IP not in warmup"}`))
	})

	ctx := context.Background()
	_, err := client.StopIPWarmup(ctx, "192.168.1.1")

	assert.Error(t, err)
}

func TestGetIPWarmupStatus_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/warmup/192.168.1.1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "IP not found"}`))
	})

	ctx := context.Background()
	_, err := client.GetIPWarmupStatus(ctx, "192.168.1.1")

	assert.Error(t, err)
}

func TestGetAllIPWarmupStatus_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/warmup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "Internal server error"}`))
	})

	ctx := context.Background()
	_, err := client.GetAllIPWarmupStatus(ctx)

	assert.Error(t, err)
}

// URL escaping tests
func TestGetIPAddress_URLEscape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ip":"192.168.1.1/24","pools":["pool1"]}`))
	})

	ctx := context.Background()
	ip, err := client.GetIPAddress(ctx, "192.168.1.1/24")

	assert.NoError(t, err)
	assert.NotNil(t, ip)
	assert.Equal(t, "192.168.1.1/24", ip.IP)
}

func TestGetIPPool_URLEscape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"pool name"}`))
	})

	ctx := context.Background()
	pool, err := client.GetIPPool(ctx, "pool name")

	assert.NoError(t, err)
	assert.NotNil(t, pool)
	assert.Equal(t, "pool name", pool.Name)
}

func TestUpdateIPPool_URLEscape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"new pool"}`))
	})

	ctx := context.Background()
	pool, err := client.UpdateIPPool(ctx, "old pool", "new pool")

	assert.NoError(t, err)
	assert.NotNil(t, pool)
	assert.Equal(t, "new pool", pool.Name)
}

func TestDeleteIPPool_URLEscape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	err := client.DeleteIPPool(ctx, "test pool")

	assert.NoError(t, err)
}

func TestAddIPToPool_URLEscape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
	})

	ctx := context.Background()
	err := client.AddIPToPool(ctx, "test pool", "192.168.1.1")

	assert.NoError(t, err)
}

func TestRemoveIPFromPool_URLEscape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/pools/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	err := client.RemoveIPFromPool(ctx, "test pool", "192.168.1.1/24")

	assert.NoError(t, err)
}

func TestStartIPWarmup_URLEscape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/warmup/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ip":"192.168.1.1/24","warmup":true}`))
	})

	ctx := context.Background()
	status, err := client.StartIPWarmup(ctx, "192.168.1.1/24")

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "192.168.1.1/24", status.IP)
	assert.True(t, status.Warmup)
}

func TestStopIPWarmup_URLEscape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/warmup/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ip":"192.168.1.1/24","warmup":false}`))
	})

	ctx := context.Background()
	status, err := client.StopIPWarmup(ctx, "192.168.1.1/24")

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "192.168.1.1/24", status.IP)
	assert.False(t, status.Warmup)
}

func TestGetIPWarmupStatus_URLEscape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/ips/warmup/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ip":"192.168.1.1/24","warmup":true}`))
	})

	ctx := context.Background()
	status, err := client.GetIPWarmupStatus(ctx, "192.168.1.1/24")

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "192.168.1.1/24", status.IP)
	assert.True(t, status.Warmup)
}

// NewRequest Error Tests for IP Address methods
func TestGetIPAddresses_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetIPAddresses(context.TODO())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestGetAssignedIPAddresses_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetAssignedIPAddresses(context.TODO())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestGetRemainingIPCount_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetRemainingIPCount(context.TODO())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestGetIPAddress_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetIPAddress(context.TODO(), "192.168.1.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestGetIPPools_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetIPPools(context.TODO())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestCreateIPPool_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.CreateIPPool(context.TODO(), "test-pool")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestGetIPPool_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetIPPool(context.TODO(), "test-pool")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestUpdateIPPool_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.UpdateIPPool(context.TODO(), "old-pool", "new-pool")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestDeleteIPPool_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	err := client.DeleteIPPool(context.TODO(), "test-pool")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestAddIPToPool_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	err := client.AddIPToPool(context.TODO(), "test-pool", "192.168.1.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestRemoveIPFromPool_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	err := client.RemoveIPFromPool(context.TODO(), "test-pool", "192.168.1.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestStartIPWarmup_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.StartIPWarmup(context.TODO(), "192.168.1.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestStopIPWarmup_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.StopIPWarmup(context.TODO(), "192.168.1.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestGetIPWarmupStatus_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetIPWarmupStatus(context.TODO(), "192.168.1.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

func TestGetAllIPWarmupStatus_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetAllIPWarmupStatus(context.TODO())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trailing slash")

	client.baseURL = originalBaseURL
}

