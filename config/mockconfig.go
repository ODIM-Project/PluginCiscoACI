//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

// Package config ...
package config

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	lutilconf "github.com/ODIM-Project/ODIM/lib-utilities/config"
)

var (
	rsaPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
	redisPassword = "F2n0YuRgavd/tYeanHI94cJR5r/C+FUaGQJBetQOxed1pLXxnWKAMmVLjs+jCBGq" +
		"C66YfEZ+DK5ZIg9QmuQGEoahwSVWC+Pa8hNrqIDgBYXP4cyyEAE0XE0j8amyf049aqhxxTYXfzov5Km4t/" +
		"Tzqru3CJ2CcUGzRmq1WfbfuMqx+tAZGw4UY1SW9IDoHwXaqsKld9uwiYq6lBqJpYzNGcCgrVyHwQg" +
		"hTrYlypQocsDdVY7/bFzg3amIHdStmzF+mvpNolhtkgrXeq1ov7stepdgpzOF39Fe5DDO+OG53wyR" +
		"4OMBAZ2NjX5LLQkhNEUpAA7GM8ajtOuJGecO506St2ASatcojJqRDHbIzNhzAxY3wtB0bx1S5TS1jl" +
		"kW1VTFXqFNjnKd3j7Q/YZXJY5a+zX/PhIZLnCp+yWY2/qU7s4BZjex8jNRikFTRzhqDGfKP1hFar8" +
		"qLr2D0FSRrDK4NtViUMUv5PaWygHtRk8e0fSnNhTSGv5kzr/fwEE4S/ayo5+9pqjgjr4iu+d6oRSq" +
		"2dQVQIfdm25Lqfw8RnmeveVVKuQk7xT/T0pcKmwfYuN0R4UjJ44BiBXgI9e/pgVTHzOrzHfLT5ekk1" +
		"eSx/fIuTMe38lxCD+L6mfhi4zsI/IkaQsjgvR70n5RlpsT8ndNLBtmfNS4NB3Ls2Cbp0AFyTY="

	localhost = "127.0.0.1"
	hostCA    = []byte(`-----BEGIN CERTIFICATE-----
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
9AjGBBg=
-----END CERTIFICATE-----`)
	hostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIG1jCCBL6gAwIBAgIUYttkxU3D5hOJyBWaZHYhHUgRcfwwDQYJKoZIhvcNAQEL
BQAwcDELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRMwEQYDVQQHDApDYWxpZm9y
bmlhMQwwCgYDVQQKDANIUEUxGDAWBgNVBAsMD1RlbGNvIFNvbHV0aW9uczEXMBUG
A1UEAwwOT0RJTVJBX1JPT1RfQ0EwHhcNMjIwNTExMDcxMDEzWhcNMzIwNTA4MDcx
MDEzWjBwMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExEzARBgNVBAcMCkNhbGlm
b3JuaWExDDAKBgNVBAoMA0hQRTEYMBYGA1UECwwPVGVsY28gU29sdXRpb25zMRcw
FQYDVQQDDA5PRElNUkFfU1ZDX0NSVDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCC
AgoCggIBALvx70x1niMxDMRHxUSJZDfQYtNQuCySjYZ/CnCwuJTAYQ+2+ogC0elU
ZBGzbB3m+4ca4MttEilinzsVdHATvOx8zqjqAJNRZWd+JsBpp6Y1sOMrQHuAHCFx
eLK1EkYLjXifq7ScJwQrm6MNxDL/Wa9pa5f3j2sxKnqf7SrPBqzrXKsaXMhOtoa+
vOTZkCRCeNn+XjK7MJRi0OCEjhKCppzSt1lp6xURm/K8ELqowoKSqLkEZIxcgGxD
u7U9lrf6CL4+TODaZ9qlQ7xUmggffOONnvlxRYNUPEINSLgWuFU+F94x8DAtNvtm
6iOOXD+AzJ1IwhChQnbx5iKCtyDcJoMWh0xEiapMBJPeDAyD6D2p7IdbqphIqa/9
gFDN7jw9GBmWaOlT8z8W3mojcEdekRpwCua9hKKwFbUYs7FItyjn3hPmBFA5VqLc
Cov8cN7PJJNPrcOhxTL9UWN2AmMlpaiQE320wXwjduu74OsuZomDAgDiecKiZUrm
pCqv8fJ20W2RkSvzWM754m+1mgqOAfPzggOOHKiIA4oFc2o4R+ZF4pmLtuAEcNvx
iIMPETIP0qGot2oZp0tlNe13xITxj9KkrBQdlNs8+aAVdxi/n75rfrDEArcaDue7
+hXbHJ0ZBowRCYfp7YVwwnvSnva4db729QOE4FDDx9xNyyl8PxdnAgMBAAGjggFm
MIIBYjAdBgNVHQ4EFgQUrOVT19Um6aIJPXDFIKocnyrhALkwga0GA1UdIwSBpTCB
ooAUO6rgMJaE7Ug1L12DTP7Zh9TX8pOhdKRyMHAxCzAJBgNVBAYTAlVTMQswCQYD
VQQIDAJDQTETMBEGA1UEBwwKQ2FsaWZvcm5pYTEMMAoGA1UECgwDSFBFMRgwFgYD
VQQLDA9UZWxjbyBTb2x1dGlvbnMxFzAVBgNVBAMMDk9ESU1SQV9ST09UX0NBghRw
rQR8cLxuujIPUdZRqxJKyZGPKzAOBgNVHQ8BAf8EBAMCBeAwHQYDVR0lBBYwFAYI
KwYBBQUHAwIGCCsGAQUFBwMBMGIGA1UdEQRbMFmCC2NsdXN0ZXJub2Rlgg5yZWRp
cy1pbm1lbW9yeYIMcmVkaXMtb25kaXNrgglpbG9wbHVnaW6CEGlsb3BsdWdpbi1l
dmVudHOCCWxvY2FsaG9zdIcEfwAAATANBgkqhkiG9w0BAQsFAAOCAgEAIKkIlcZg
tC2RiNJnvOnwDpLin0Ygy5BZbHVizo82RFAhHI2UPSpSNaRi6/c9gVGKy9RX7sDR
w3a7SAIcD0NDgLddvemfFJ/yLmQ4OJ8J9+1R4+PszwmzYXFBEKWr5WzTVNOSBvi3
INItVaWeI81m/dXVNQ7PHiVkFhpEqW/HsXuG/VSKff4e1jnU/6a7Zc2qnrZvFRha
Q/HtIu42eIMFTtNgFEPQkeD3OIsFLcRSIP+uPu3V/GmOWPAPJhzgrJBk2y82g9j9
gmofhYWiL1DWwV0P/2LeAIQct9txqfsxNX3LqtVCHSZmfeTjN02KHFpLiTsOI28r
LJBi+6auCn5oLIEhEhuD3o7Wg/UAxsbekmXVCJmCwl3ez1HIXdbuuytbiSK2pylp
HiBkXiXywsqUOqQRx1l4uRamnpwZRY3Ox5PFKGa5UnXgaeQMu0KtWxd7H+LsC0f3
LmkskCGEfQ2TDA3zoJMy2+adyOua300JzG74AoAAnWnGZ8CS2ONdkA6E0AlyYuS5
PBYAdaXyoT21GdSgSOl/oTa4I4zM/St3mp1AJSVlCqLg0mpPmIT1/g7my0La1C/a
xsO9Vj1rJ+m8HoLuTecLuyR/z/zfagyiyODBfb2mSMFM2b5XQE7WFd5x3BTYUP5Y
zHq9dIL4UN8D74DLasF9SlPYC99+DnPz8Mk=
-----END CERTIFICATE-----`)
	hostPrivKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)
	hostPubKey = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA+Q01x6+Jgpi8Tdo3ul3P
iuacjeJsFr32gmdo1l0oJp4FweMtiTm9k/pggQMrUGcPqV0ANI4h0DVHx+RdR27R
DV1bqMJqOghTAKbCOw4Wh0X9vbITiuxhLPsYMhSsOxY0Au/YVJd8/ZZQ7QHqKf5V
SvF0ickFrKP4Rkd9LecAXTlXwVxBYBMw8Y1OZVuv272soFMVt8BblCKxNEM8pX8H
KcFoXMrWZr2tgrxDJi7Ry1zNJ4S/gBkPYJfNj8+lPOwwc1nUKIIbzAGkN67h3Q9I
ZRlyyM8D7ayZEKk3tfhNvD9lAip24yORWQDocQ8+wsjerXtTJU/bdqDpLPeAvTdv
QViOzzrMvikIpw9YzbRN6i17jA26BEI0yOgtLLcHOA2ah+K/0kDpHINR3YX0TNfL
SeMEWylQq2Sv6cRO9u0iaRih4GHfWOkc0R/4VaRftS2TJmGEVKT2h7XEnlspCnDw
OdPafLtPKL+aNAdoDnS9fALkAb/lGskmsM/tmSrS4tjSgYhdsYMQp2FseyhtHfsC
4hLW0AnQn6ckdlr2kwXGOc+kpDcoWtc9V16rkaCoNjTu76P8nWIjEasNBZvm3unV
o79i25P4izfyzAG/tdOA+NVbArJEjBaHge0ekJKajHFPLMaaJ1kptItWS1PGVy5d
ZlgyKGJ8O0RB8M1vofMdLfMCAwEAAQ==
-----END PUBLIC KEY-----`)
)

// SetUpMockConfig set ups a mock ration for unit testing
func SetUpMockConfig(t *testing.T) error {
	Data.RootServiceUUID = "3bd1f589-117a-4cf9-89f2-da44ee8e2325"
	Data.FirmwareVersion = "1.0"
	Data.SessionTimeoutInMinutes = 30
	Data.PluginConf = &PluginConf{
		ID:       "GRF",
		Host:     localhost,
		Port:     "45001",
		UserName: "admin",
		Password: "O01bKrP7Tzs7YoO3YvQt4pRa2J_R6HI34ZfP4MxbqNIYAVQVt2ewGXmhjvBfzMifM7bHFccXKGmdHvj3hY44Hw==",
	}
	Data.LoadBalancerConf = &LoadBalancerConf{
		Host: localhost,
		Port: "45002",
	}
	Data.EventConf = &EventConf{
		DestURI:      "/redfishEventListener",
		ListenerHost: localhost,
		ListenerPort: "45002",
	}
	Data.MessageBusConf = &MessageBusConf{
		EmbType:  "Kafka",
		EmbQueue: []string{"REDFISH-EVENTS-TOPIC"},
	}
	Data.KeyCertConf = &KeyCertConf{
		RootCACertificate: hostCA,
		PrivateKey:        hostPrivKey,
		Certificate:       hostCert,
		RSAPrivateKey:     []byte(rsaPrivateKey),
	}
	Data.URLTranslation = &URLTranslation{
		NorthBoundURL: map[string]string{
			"ODIM": "redfish",
		},
		SouthBoundURL: map[string]string{
			"redfish": "ODIM",
		},
	}
	Data.TLSConf = &TLSConf{
		VerifyPeer: true,
		MinVersion: "TLS_1.2",
		MaxVersion: "TLS_1.2",
		PreferredCipherSuites: []string{
			"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
			"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
			"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
			"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
			"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
		},
	}
	lutilconf.SetVerifyPeer(Data.TLSConf.VerifyPeer)
	lutilconf.SetTLSMinVersion(Data.TLSConf.MinVersion)
	lutilconf.SetTLSMaxVersion(Data.TLSConf.MaxVersion)
	lutilconf.SetPreferredCipherSuites(Data.TLSConf.PreferredCipherSuites)

	Data.DBConf = &DBConf{
		Protocol:                     "tcp",
		Host:                         "ValidHost",
		Port:                         "ValidPort",
		MinIdleConns:                 2,
		PoolSize:                     4,
		RedisHAEnabled:               true,
		SentinelPort:                 "5678",
		MasterSet:                    "ValidMasterSet",
		RedisOnDiskEncryptedPassword: redisPassword,
	}

	return nil
}

// GetPublicKey provides the public key configured in MockConfig
func GetPublicKey() []byte {
	return hostPubKey
}

// GetRandomPort provides a random port between a range
func GetRandomPort() string {
	minPort := 45100
	maxPort := 65535
	rand.Seed(time.Now().UnixNano())
	port := (rand.Intn(maxPort-minPort+1) + minPort)
	return fmt.Sprintf("%d", port)
}
