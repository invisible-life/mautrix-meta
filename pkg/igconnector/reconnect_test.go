package igconnector

import (
	"testing"
	"time"

	"go.mau.fi/mautrix-meta/pkg/messagix/cookies"
	"go.mau.fi/mautrix-meta/pkg/metaid"
)

func TestUnauthorizedReconnectLimit(t *testing.T) {
	client := &IGClient{
		Main:      &IGConnector{Config: Config{MinFullReconnectIntervalSeconds: 3600}},
		LoginMeta: &metaid.UserLoginMetadata{Cookies: &cookies.Cookies{}},
	}
	if !client.canReconnect() {
		t.Fatal("first unauthorized close must permit a session reload")
	}
	client.lastFullReconnect = time.Now()
	if client.canReconnect() {
		t.Fatal("a repeated unauthorized close must not cause a reconnect loop")
	}
	client.lastFullReconnect = time.Now().Add(-2 * time.Hour)
	if !client.canReconnect() {
		t.Fatal("reconnect must be allowed after the configured interval")
	}
	client.LoginMeta.Cookies = nil
	if client.canReconnect() {
		t.Fatal("a logged-out account without cookies must not reconnect")
	}
}
