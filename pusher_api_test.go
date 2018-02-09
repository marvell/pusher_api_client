package pusher_api_client

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	client := &Client{
		Debug:         true,
		AppID:         "ee987526a24ba107824c",
		Cluster:       "eu",
		ClientName:    "go-wex_client",
		ClientVersion: "0.1",
	}

	ch := client.Subscribe(Event("depth"))
	go func() {
		for msg := range ch {
			fmt.Printf("TEST: %+v", msg)
		}
	}()

	err := client.Connect()
	if err != nil {
		t.Error(err)
	}

	defer client.Close()

	time.Sleep(10 * time.Second)
}
