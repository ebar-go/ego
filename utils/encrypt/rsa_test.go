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

var privateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDO3bPODcV97hyu
xNuAT7lg/CL/AitYzkZapXQJfndIZhFj578YfNgSNr/SHZPAHauCXQzhBqCu86Nq
cimyaoUe50mUMia9vhDFbwgQVULiaFkS/fGxvFF+rwGtvUPF8LyIEiLNPeKTQR6/
v3kpYiM0lUMWbT/hook74GOHr4HKCf4/cW1yAbS2Ap05Uc90CUvF3XjhuLbjrI3n
5bRVnT5GKxBzc3TEvB/IQRgHxsFHCgr/+2i1Ne290YlBVyMIzZKRqxufoGQj0Mfk
pXNR76YirUoZdrjabSwl6FyS/af+AUQhJeLold56tOInGld4EePvtDeGrscyKY7L
GA9dAYpbAgMBAAECggEAdSZEfzp5TzmbcLX3AJElkHD5eKTH24DlgswUDJRcBNoY
mxVQuRNqYdYzh1BMrg2fZTJA6uaP9MPxFYUVa/666KdemdhU7DtI0CZy0J0YRZOG
biT1zQuysyV0s+ltChmtCvoKT2TufSnxufE7Ml5rRYoJ9hdkh+k+AnSuqNaRj9JX
yujxWoIQGxYDXIkXlq1PythtkNpSvn9a5Fh88zMwgy1DglCcAFpgTf/5JmvVeT/n
48Tw+45U1sAblqJ84Vak/HhYiqe8UFMHLGcy/lQBrpgIfW6INENRK0Ff13Ic2bz/
IHwaDmQ8hhxCkBzl9FRxZpUpsTuxc+z6t4qFOK1wYQKBgQD1A0Ja4e5G5Xgb8534
Y7F2aY4IFH0UmI0FZzi7+csAREzNjZyF0mb4NWnXEtDOpLu2YAt1MwSbtVvi+Vxb
2glkFWqsJynLUgiwDwtrXz1TbZ2/rMXR/DnstKYvTLIQb6o8RedlqPt/iXMt+paz
zA4J7+3rTMNjx8cVeKKjup+0uQKBgQDYJIUSIAjCzLuSZvzNgVknkN53wm6VOT5l
yCdP4OBUSMOTj+FKn66j75lG9IMy7TSC+rRyXfUBKEdkNc+E7MXSLE/2RYrzwgSh
Z8ZDVxeUrCMYU7psUYvYajgZRAT1Nqvy5B6saUpdi/O94hsclF36c9SyvJROth45
VfRDZhsVswKBgDr079CisQ2KRh6jvo14n8lYmP7Ev1xnYPe94N8Kuphz1u9XdiSV
foWXhMJmGqy+4xR9hARNpHw7ZcL2Mg1AKCZXKPYH7nyoXsLOu/a4Ui9zHxRyZJ+k
y+NzjNGw6OAfnp0mTQofYXVNA9Q6imz1WyN1ApEuRY4LEpLOFoTDcY6xAoGAcYE8
IHiSITpChm9u8ryqhQyex2VjjRmymuCxRFFjfN95VVSJixawL4bzhz+AZo4KtX3S
pPySXTk5xHY9tCBjAiwjEcETZ07L/7bvdw2VZI3BIFVX4Oox9kRRkXMW527+fV8u
fHwOgXGtXloOwsNnVs7dM7+0YmFhHdr1my9TqeECgYEAxnQEzLhQxtX9dbRtCOpY
LRqOLYoh6XTZggXUiWHAmpEBOA17/a4a6cJl6j0tPDUyrbHKkvLyHWB0FsUOh8W6
nxJIGKYM8pnhqPfcBUvIdS02/a6g9FQgUlceRDuYNVwkSa0MyHboBRSJ9fvtxh16
LP+kNrupW4Eue8riBhHpfvE=
-----END PRIVATE KEY-----
`)

func TestRsaEncrypt(t *testing.T) {
	str := "Hello hero"
	encrypt, err := RsaEncrypt(publicKey, []byte(str))
	base64Encode := Base64Encode(encrypt)
	fmt.Println(err, base64Encode)

	decrypt, err := RsaDecryptPkcs8(privateKey, encrypt)
	fmt.Println(string(decrypt), err)
}

func TestRsaDecrypt(t *testing.T) {
	encrypt := "duMC1gD8VgX7AuGktqVH28F4vwBej2f2S5AVOxkdxuZyUnVjuoysDqvjj7MH602uIsvPaN3JpWVbKWBAh8/aBaAXLcy1K9mgN/q/ih5YBnwE304n/lFdvndo6V0qAZAeACxLot7xHq1/QcaoKD0K7RZE2kpHT2uYNYqPvvzz/ILC+iebbO3MPM6XfF1X4CgcgAf0N6nnASZATu5ik3HcCQ1gagllk5THWzi4BJRsKob0ncU5SLtJktK6r7qsOqP7RtuE1/OzT9sU6oNQkUgpndVkscXXnXeHzaMSZzorlxd9FLjYkZICCPRVCehXp3FSthbcnXl2tv1TxnRJ8AGp8g=="
	decrypt, err := RsaDecryptPkcs8(privateKey, Base64Decode(encrypt))
	fmt.Println(string(decrypt), err)
}