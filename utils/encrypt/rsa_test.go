package encrypt

import (
	"fmt"
	"testing"
)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzt2zzg3Ffe4crsTbgE+5
YPwi/wIrWM5GWqV0CX53SGYRY+e/GHzYEja/0h2TwB2rgl0M4QagrvOjanIpsmqF
HudJlDImvb4QxW8IEFVC4mhZEv3xsbxRfq8Brb1DxfC8iBIizT3ik0Eev795KWIj
NJVDFm0/4aKJO+Bjh6+Bygn+P3FtcgG0tgKdOVHPdAlLxd144bi246yN5+W0VZ0+
RisQc3N0xLwfyEEYB8bBRwoK//totTXtvdGJQVcjCM2Skasbn6BkI9DH5KVzUe+m
Iq1KGXa42m0sJehckv2n/gFEISXi6JXeerTiJxpXeBHj77Q3hq7HMimOyxgPXQGK
WwIDAQAB
-----END PUBLIC KEY-----
`)

func TestRsaEncrypt(t *testing.T) {
	str := "Hello xx"
	encrypt, err := RsaEncrypt(publicKey, []byte(str))
	base64Encode := Base64Encode(encrypt)
	fmt.Println(err, base64Encode)
}