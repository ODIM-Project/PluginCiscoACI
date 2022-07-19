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
MIIJKwIBAAKCAgEA4RV1LmHtH/X23G+Qz45w8wmfDwkwnsCsrVU45lU67+fUdoy1
90mmxU5i7bT8Mj/312K4SiE+FtgvF8T0i+UStXG/l9FokSoeLfLE2pFGpz3+CIk9
4wQjpMgi8SH4A8wrbb4rR7Z5jiFfIrOi+zwC1zjnhK9yiWe9e308GGtXuXVmtqfQ
LOVvupIr1YJ5W1dnF2SS5r4OPf+i0r7v0D12WYHmDlxkc0Mr2mHnaAujDzj1OZsQ
q9MeNwdGCOfDYx820vvQNyM+uYkPX+aGrVJDO3GT4X0jr/dDsVxtTxHRdY3E/H7Y
U2PDviW6sFbzUtd8sYw3msoYpkY/Wp22OvH6sM0iwg+cTLy+npoAbgOhuHCpcgO2
Juq6h7rmijnWx7HqAW2oJBXpex0qtcyKAp69NMLRGw6CC678g2sVa0vxvpQKlQxz
29SBPBmK1u+45bnhYuXunrjhPxkNyjHXRRrLO/1m7qI3hLL/fFe9dPYJecourbRI
tLA7CX+J84YCgt74b4RCamTcY9pcGtSiKIch1eF2eLuh7TScIVtsofM5T1cxTisf
JPtC7K9oFhDcbAIPh6XFsOMpr4cPfL0MycBT2ZVmRHLPn09jXtddz6R8ozepdkqJ
HbgAmD3Cr48pvel49oM3osBylmE5Xk+9eSDerBSffnU0FaAws4t2BvaiHZsCAwEA
AQKCAgEAyozkxriZCwns/LHpPt6QBiXCXWWHu1ToD5OBgMVyJDIboBNALSi6SxQf
MoqL6SxnfAv6i7sehLBGsL0s1Ddwfpe+MoDf+MJOJksxmv7g9d9zm3rllkVDTiZM
S3KmHcS90CQyDnbHLIAbfL7rC+sVI1ix/1VjXQNeIKKyUcdHSj28EOMzEzPlN6AS
kjC3xNsCiqqXB85AQsqpW703Uc39ks6ymHnMa20nKX6xH5BZTHmVNCG2/ukdZ6fD
/n+R9MFCNNsmpHezGoOcslBhIdfFaNjsmx5h3xhEcncaZu1B8OeDPTVotqIwpAyP
0+BrV0FTlPL5lvIG/Jp6qLEELEdVr9TZsQBE+BETXYlNPRon2dhNGsjscCDTppdF
oDYWiCSxv2rJ7aYf1eYR3cjo5eFbCJHzZVUlUQP/LhUn4rL/Et+0lzrzMlNWNg/F
Ev7/H4PNrTDa3OsdgDVouC4hILUtHud4cVfracng4nSqLCxLgKlljzs8TAHKFt+l
JA30LxIPo1xsW71ijTGA6FdZbTxUlA8boVzPE2A2c7AdN/I5CC/g3OJdC2VodfaY
0RmPxqh4dcgnO8pm925I0YGdVfwf9BKXngyhc0pbqOhV5aHgm3qWrfq+0SUVuDA+
JSkmh4IEj03KvakDuOA0HBTWvzenUKlFdkLOP4p915SQq4zXuAECggEBAPL0r8sy
EaedYrLBtNpz/VUJCNWMcERHXl0xH7rygaD+by+iID59/v/a0JDQn/bhd+Vojh8Q
jjsLHbuJMhVVtgF+Db1AKym92EhyWwu05vBtrEjwLFIqlpD+IjXZYQdyY9WXAsYd
NHwuTv1JrnbBAw6PxjpAoivi10AxaDMIhnN0/BZz3ObywLO+wMLwl1cKJr7tr8x0
uwziXZMXnp05K79GqVYeuVp0RWOU5tN0DFri5pX8+0bMsY19WE+FoBo9rSot6JLu
lNiS5nLKnlfy3yAawuzT7nasDJZ5UBPLoJ4m2mjaqB4e3bgijp4g1RramU4S1z6+
b7IIifUMtTNS1DsCggEBAO0rIQRsfIytL44pKlXiW64j3ryCFiiWIhh7QJC2Kk5X
nuXHZvMA0udiNsrntxCWhT823Fnvuh0OxDDpL4MjTOsjI8lkY0To2mowXQlNV8/I
yaZD+ly8lQC63aQYN+Byu3Ow+hiKQzQhsBt5U/Fb/jZig1LPgcmPyiqHD53hIpdp
qSHlpRAvmVcrejCFyuChrvrg6twTSh6D7CzTPeLqE9vNJA3W96H4n2H8REYEFYMQ
KVLjOROH58wCYKq4XwQrE+QnnkCO8vAmyPCYaQkA9QHPQPjyY75z8xxss1klghn1
G5tKu5/1rNYtMONB/P+ZFMmCYFX4n9mcRvSNbxUcJiECggEBAL+xeibL+YwTrPU3
yzd1vxNiDntX1Ji66uSCxvNdNhRNzHJ77A8CoLlE77zjLuO/IDd8mG5ARMinS61V
YZPdzb49tB93StcjeEwpFlcVRAW9surVvVKTUbtTGLD+NAWJJuY2wTSJhIjajO5i
PWprfbr2i8QYjRwtXgLDOODTQCpGykP45Pm/3XW08yicZfyCAPIyXbvm+lL/JC/T
ug15N2AzI5bUpRCOntUkfj+m17y6PI9pTOWeyhTGKnCMETfDJCcck92iqwR6W6OE
5Qylj5EoLFZqHUO7Gi97xkfoKXG/XCLRK0agufX4Jijz5NDMW5tzWCukXELPY/Ja
NXoqR1MCggEBAMags1NIJHuQ494Uvd8V56CNbBLGhBZTvpRwTR+lYQMhwPNCMAde
bkPY7ni63Yen+Ep8AMnVyzJg1pD8Co2yt83KLUOSrszck2gRvyl2PA/KYo+8KOcY
DVaCKfQvUETK8hEvbBW3XhdAC4TG9TWTzPDxSnjFTzZnFXLOkJayIc1bcYnxEW/f
3XWy9O/Ebaf54Vk9m5TbFt09sUPNWuw7DIyuXv60RcrCNYHTy74z12xf0awYnwmr
bcdfSmRQa0tLZKpVP+VjkzTr1qghjP48bfWpBQo5vq2X4EizBPWpQy/IJunFCiQq
lij9yg7aii/qng0yAsqdogqXJpnUBe9RFuECggEBAIdXDuRf33r9nYXHh3HM4iKK
3FDIAXJ2/aN4R5rFphRoatOFpKx97EkUIbJSxfRQxEDujmU0tbUL3YNglQCOhi26
OM5wOqORIeTwS4+L4vv9M7MabGZiG8l7TXwkxFBDYwEqqwjAeU5qY55f/pZaSN9E
QIU+TwYUqYN7xaKMklUzubA95XSJ6J0WWdy+zeJN+X5txcqVFpBq0wIM0Au0UPOp
dLfmAFsFh1pNv8p7MQyOaQo0kwZZuUu92YXU8tC9dNCKTd8sWP+CkRtjDZRykXo5
/vohwYCB6eglzR4vo7W2Ukms3oEwfiCywInGpfYYE3peuHDN83GVsXdLjuBFrNk=
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
6ATq
-----END CERTIFICATE-----`)
	hostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIFsjCCA5qgAwIBAgIUNN0LnDt13/3nM1GipqN9Pku+ZdEwDQYJKoZIhvcNAQEN
BQAwUjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRMwEQYDVQQHDApDYWxpZm9y
bmlhMQ8wDQYDVQQKDAZPRElNUkExEDAOBgNVBAMMB1Jvb3QgQ0EwHhcNMjExMDI2
MDkzNDE5WhcNMjMwMzEwMDkzNDE5WjBBMQswCQYDVQQGEwJVUzELMAkGA1UECAwC
Q0ExDzANBgNVBAoMBk9ESU1SQTEUMBIGA1UEAwwLU2VydmVyIENlcnQwggIiMA0G
CSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC1vI9gNVOQnpoSsm2Bonk41xbcJ1CR
5tj6wswEOBaDNXHtYOq85vMT29ou4r5DYanN05Fhd+uJC22Gf4nbGoLkkt1AuIIb
d2U3PY9YJvNfQ2R2GAvTSBTuB5wZ6Co7JEX9kRlVbCcQrF1TJjTI0MgSsHf5U4rg
pr70YN2WecpUieQ3GG+x94NpupwwGJWeh8JIgDxAdh8m/cebOT1VO5xxtPpcd8n4
TxlroPLbgz3UFyMiElmih3+w+L89pOuuPgivBdIVKNsYLYICFw9s09cgwlsi6qhy
Esa0yMjLARuaxz2VFVPwDH9+JmGHWWIUfjXnQ3WixGTrpAI1iboN24DrFtjdIRfP
RgVpuDP7HSN2/CNIXxQ2RQF0fB/z4g1KJS56UgGFCvgDl2BG3aOj+tlQofRvgxT9
jdBBkeiCLgixkVs4nu73rB4Bq4RHRosrpt4kT4nmioMuUcLJRsGB1zt0+xQ+DR6L
vsnzbD4bS3fYADuTzDzXjBem9Z2RYUX9wUtXeiOqqNXepb/AeHmZLzUlVlSu5h6c
PxR9adSVw8NBSR1JXn0hqLFjOVEDk3wMicQnzRXYUrpJgh0tUqu6gHs7DHUsYQPj
qzXxMZbpJyMZSJirqu43VkyVhyIswnzGErwDs6gEM9qUlrQ7L4Zcpy+Af9SqpZPG
aB6p6/m1a4T7XwIDAQABo4GQMIGNMAkGA1UdEwQCMAAwCwYDVR0PBAQDAgXgMB0G
A1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4EFgQUdWXGZkyYvxcX
A8WCUJ9Df/KA4w0wHwYDVR0jBBgwFoAU2x5lPlMmcsAOk5p96W+fpo2kPeYwFAYD
VR0RBA0wC4IJbG9jYWxob3N0MA0GCSqGSIb3DQEBDQUAA4ICAQBpMuSlw82RHEFJ
YOFy8iPM8bxTThrZRzCpubePKdRPOKEor8WXuHyiKNhbaaPdvgHgPvTBYjFGOsJn
gjWbUc8WwhyCj/1U+m/xSiOTzgbs2sE5Qa4/v4s7Xai3GRAeTC3Y+QDTpaZQWkRT
mSoZQbELzjPYtyxcblUUzLjjupjjld5vlNKKoPyl6WdtRrnkkvgABEFAysgWyW1l
yKO/xfe6sBS7H/UWqZOZZQQIy8svmCzNEPXsr1UwjQiJB74KgZrRuomm/FaA+fVW
jpAQ+LsN6Dhe9dv8J4e4MgHA9uaiFRwI1z+xeQnUHXJUbWHQNWpbI62ZTH3LftFT
+f54Obh3DdVyt+zugBZoIHUaAB+eKcTcruVEwDcDRFS0ObyqBkIej2ED+1iTFw3l
VoASse8XMdn734a3Uso65yUoTCj/Z54/N09m24ZRXyoyahgWwl0Kc+2lINcEEj7l
zAAjO94sWIWvFnA1c2rgx9zFQd7PU2zIHOP/A3tse3yyoszvS38eiTZkvP2Tj0mu
AzjfunVr6C7YjfYWnAa3tGYTpSlKnKVnCl7vIfNtthC2i3EcXACQcCecQOBpp6A2
7pCccBHHUvPMuETjEnCC5wSrSYA0MSHOreS0to7w+NGAc5/YjINN+7Zo6JooALNk
KoQRsFQ0udOGCu+Z5uZfMr0nre618w==
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
