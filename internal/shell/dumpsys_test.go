package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWifiInfoSsidUnknown(t *testing.T) {
	res := `
mLinkProperties {LinkAddresses: [ ] DnsAddresses: [ ] Domains: null MTU: 0 Routes: [ ]}
mWifiInfo SSID: <unknown ssid>, BSSID: <none>, MAC: 02:00:00:00:00:00, Security type: -1, Supplicant state: DISCONNECTED, Wi-Fi standard: 0, RSSI: -127, Link speed: -1Mbps, Tx Link speed: -1Mbps, Max Supported Tx Link speed: -1Mbps, Rx Link speed: -1Mbps, Max Supported Rx Link speed: -1Mbps, Frequency: -1MHz, Net ID: -1, Metered hint: false, score: 0, CarrierMerged: false, SubscriptionId: -1, IsPrimary: 0
mDhcpResultsParcelable baseConfiguration nullleaseDuration 0mtu 0serverAddress nullserverHostName nullvendorInfo null
`
	ssid, err := _parseWifiInfoSsid(res)
	assert.Equal(t, ssid, "")
	assert.Error(t, err)
}

func TestParseWifiInfoSsid(t *testing.T) {
	res := `
mLinkProperties {LinkAddresses: [ ] DnsAddresses: [ ] Domains: null MTU: 0 Routes: [ ]}
mWifiInfo SSID: "wifi ssid", BSSID: <none>, MAC: 02:00:00:00:00:00, Security type: -1, Supplicant state: DISCONNECTED, Wi-Fi standard: 0, RSSI: -127, Link speed: -1Mbps, Tx Link speed: -1Mbps, Max Supported Tx Link speed: -1Mbps, Rx Link speed: -1Mbps, Max Supported Rx Link speed: -1Mbps, Frequency: -1MHz, Net ID: -1, Metered hint: false, score: 0, CarrierMerged: false, SubscriptionId: -1, IsPrimary: 0
mDhcpResultsParcelable baseConfiguration nullleaseDuration 0mtu 0serverAddress nullserverHostName nullvendorInfo null
`
	ssid, err := _parseWifiInfoSsid(res)
	assert.Equal(t, ssid, "wifi ssid")
	assert.NoError(t, err)
}

func TestParseBatteryLevel(t *testing.T) {
	res := `
level: 80
`
	level, err := _parseBatteryLevel(res)
	assert.Equal(t, level, 80)
	assert.NoError(t, err)
}

func TestParseBatteryLevelZero(t *testing.T) {
	res := `
a: b
b: c
`
	level, err := _parseBatteryLevel(res)
	assert.Equal(t, level, 0)
	assert.Error(t, err)
}

func TestParseBatteryPoweredType(t *testing.T) {
	res := `
  AC powered: false
  USB powered: true
  Wireless powered: false
`
	powereds := _parseBatteryPoweredType(res)
	assert.Equal(t, len(powereds), 1)
	assert.Equal(t, powereds[0], "USB")
}

func TestParseBatteryPoweredTypeAll(t *testing.T) {
	res := `
  AC powered: true
  USB powered: true
  Wireless powered: true
`
	powereds := _parseBatteryPoweredType(res)
	assert.Equal(t, len(powereds), 3)
	assert.Equal(t, powereds[0], "AC")
	assert.Equal(t, powereds[1], "USB")
	assert.Equal(t, powereds[2], "Wireless")
}
