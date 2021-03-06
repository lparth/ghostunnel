/*-
 * Copyright 2015 Square Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"crypto/tls"
	"encoding/base64"
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCertificate = `
-----BEGIN CERTIFICATE-----
MIIDKDCCAhCgAwIBAgIJAPjKcAKZMSkUMA0GCSqGSIb3DQEBCwUAMCMxEjAQBgNV
BAMTCWxvY2FsaG9zdDENMAsGA1UECxMEdGVzdDAeFw0xNTEwMDcxODExNTlaFw0x
NjEwMDYxODExNTlaMCMxEjAQBgNVBAMTCWxvY2FsaG9zdDENMAsGA1UECxMEdGVz
dDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAK4EbZf3EMb/ciW5nGlN
yrf5Pcfz3ZnjWRy1kvBriuPD6NQSZaTWTPmJnbdS/Q5FH0p/6ZjdZKXf6f7WNnAz
JwW0XK7NT3N2DrWfgQqrrVvLAYlfqgHnC7Fxqq7FCpgWjf7L8wcQXfdIYkhdsE4n
osLmCRvx7qS+wuasb6nLzBtg7b99ZvO8K/sezrDIjwzemBWA1Vovztw/vGD4J4/h
D0hiOOqFGWstwFxB9oG4d/QJ45VttLMGuiZCY+A4IyBgPCxphrEec6zf8H4u/ceQ
bB8i1IMmD1VTsq9afeVhMKuoSn2Bs3VRB6c9FpL41/ftN5mYpZCteZH+qQ/DhK/y
Dz0CAwEAAaNfMF0wDAYDVR0TBAUwAwEB/zALBgNVHQ8EBAMCAqwwHQYDVR0lBBYw
FAYIKwYBBQUHAwIGCCsGAQUFBwMBMCEGA1UdEQQaMBiHBH8AAAGHEAAAAAAAAAAA
AAAAAAAAAAEwDQYJKoZIhvcNAQELBQADggEBABuBe5cuyZy6StCYebI3FLN3CEla
/3Hreul6i5giqkF90X6M+9eERZCqSqm2whBMSF4vG+1B6GX1K6S29PUOmTDWyasW
B0WlBgRiZld3JfFBuJu6xk1a8+XwwlGOgEsggepjkrAXbjbqnUMAKOJkjFIyIPvk
5p97SYDJYiOh7MmjyXUIzyNdqpL5WiUgKPTxXL+1tNzxH1jjxfVdjaNaNcOJuu20
9tsMqDZyTm2yZWOBUXbtqlaMQHrs5Ksz5EKk5/U5KfJehKss8oba2npg/6echTJU
nkOOZ6U4eEju7H1S46qlN9ZmUmSrrjwec3H7CnvxQ0ncEyZXlEiTlbO2JQI=
-----END CERTIFICATE-----`

var testKeystore, _ = base64.StdEncoding.DecodeString(`
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
CLdyMSoQneGHAgIIAA==`)

var testKeystoreNoPrivKey, _ = base64.StdEncoding.DecodeString(`
MIIEdAIBAzCCBDoGCSqGSIb3DQEHAaCCBCsEggQnMIIEIzCCBB8GCSqGSIb3DQEH
BqCCBBAwggQMAgEAMIIEBQYJKoZIhvcNAQcBMBwGCiqGSIb3DQEMAQYwDgQIMajQ
SoGOIUACAggAgIID2BD884Z9xHNWs//HQneMjFjUhH0J6UwuNgY/tYF3Ol3OoCaq
A9wRPSryJ8/YfAFs3EHpJb+0RYDgT3sJxrHEXwJuVR0hmUzP4rfUhq6z2fvjMYbL
eHWtiEi9aL+lv8Eoczszsqp8KMlbdthskjOIwZiiOBMxTu4Zv7eku0Cwz53U1NCm
wyAFe6n7D/09pvPTURaX0FjmhYdIN+Yb1CnFDC2rHHv3LMRofmYXu6cu0IvY9uvd
Z372R15uIfDr8oyvpQKPhJjEUQ9EfrhLVekWO48LWi4/XoJH6hpnZ8VFHw7M3QIO
CKkITThRhZ5gTTFJi+4/n0q+DzIan9SnniaDeXXgS3zvL7uIm5QWZV7SWIHvnXCA
Mkho6/iDrlu4l0zeLaAiLWg+fuXmnONmX/dGA+AXtK+1wY1dMmrV+kDqBJfuICIb
oXQYxaJ3TgzJnsRXNfXXtX2WSMBfdj+668NdaUcKf/goTrOcznVpGx8Pkm6oHUQ0
r0eK6iV4ApNS9ph7cHS67RQVqbw9PidWYCqYjfasZmcZvLyyVqrSfbUJnjzQARuA
Nblsj0AWGRQIvJHcnrw5Qc3zMtiJh8GhAXCJOKLKlbsjo/aJnn+3KiVwl8BpQhmz
NLbsx8DPdcWUIAxJves8S3UyDaJA6fj4nf1KMNqLu6vpnFptiIiF9pCQcvXTHYc6
tW5nZC5KExxME+Ldkh1Hsp/1DkOsfuFhVAKZrm4F/7Pz5W6BMHteQKeTX3fl1uz0
E/IT5/8yYOan/vYSQNWCcVpc4Z1jcbVkgRtBWqCZ88kq87jvYaFi28znp1qpxGka
q8DDrDQ2ZovXp1KBvLfBzQwRigupi1wQCeKu7pX+TbuTEMkbGPZh0U2dpGm/fxrZ
Nr//yF7N5NLzPWq7qApfZ0Z6DFi+NS5kU6S405ZNHgmwQV+VC324IjuWrEX+AprH
cUSL7wJ34HMsejTaaU4AiqYrN9MdIkn+qsGrQNurzEFJ6NRyAanUwUXyuDN4Yq5m
zgyC07LU5vRRmfdjyBsiJ+QyKLFU6zkQCyCdmENQJr90U596wE9nYDEWABGMSppe
wxQ5fj7+z1alRQZu6jHIal4JH2dJlMAP+MT6Ixokou2GJjuB7qznmpdYTGh2veyj
W7Wvo/eciyujzQ72eO2sRqhzX+SeP+i669ucbYlMBA6DCO101iINxi8LzgOEguWd
KYMB/SV5VsjIOckZuBIn8mMQIAqFGIvqeCS2qovntjHZMyuAbenOFLfi+WRg1KZZ
YAnq2h6R3bmXYwpZzI/S+E/0PQDXHArbsM4XgimleOle+O2bqjAxMCEwCQYFKw4D
AhoFAAQUMdr6fwPsXl5nAlbi51zv2YJHelkECLJyvuiCk4LpAgIIAA==`)

var testKeystorePassword = "password"

func TestBuildConfig(t *testing.T) {
	tmpKeystore, err := ioutil.TempFile("", "ghostunnel-test")
	panicOnError(err)

	tmpKeystoreNoPrivKey, err := ioutil.TempFile("", "ghostunnel-test")
	panicOnError(err)

	tmpCaBundle, err := ioutil.TempFile("", "ghostunnel-test")
	panicOnError(err)

	tmpKeystore.Write(testKeystore)
	tmpKeystoreNoPrivKey.Write(testKeystoreNoPrivKey)
	tmpCaBundle.WriteString(testCertificate)
	tmpCaBundle.WriteString("\n")

	tmpKeystore.Sync()
	tmpCaBundle.Sync()

	defer os.Remove(tmpKeystore.Name())
	defer os.Remove(tmpCaBundle.Name())
	defer os.Remove(tmpKeystoreNoPrivKey.Name())

	_, err = buildConfig("", tmpCaBundle.Name())
	assert.NotNil(t, err, "should fail to build config with no cipher suites")

	conf, err := buildConfig("AES,CHACHA", tmpCaBundle.Name())
	assert.Nil(t, err, "should be able to build TLS config")
	assert.NotNil(t, conf.RootCAs, "config must have CA certs")
	assert.NotNil(t, conf.ClientCAs, "config must have CA certs")
	assert.True(t, conf.MinVersion == tls.VersionTLS12, "must have correct TLS min version")

	conf, err = buildConfig("AES", "does-not-exist")
	assert.Nil(t, conf, "conf with invalid params should be nil")
	assert.NotNil(t, err, "should reject invalid CA cert bundle")

	cert, err := buildCertificate("", "")
	assert.NotNil(t, cert, "empty keystorePath should lead to empty certificate")
	assert.Nil(t, err, "empty keystorePath should not raise an error")

	cert, err = buildCertificate(tmpKeystore.Name(), "totes invalid")
	assert.Nil(t, cert, "cert with invalid params should be nil")
	assert.NotNil(t, err, "should reject invalid keystore pass")

	cert, err = buildCertificate("does-not-exist", testKeystorePassword)
	assert.Nil(t, cert, "cert with invalid params should be nil")
	assert.NotNil(t, err, "should reject missing keystore (not found)")

	cert, err = buildCertificate(tmpKeystoreNoPrivKey.Name(), "")
	assert.Nil(t, cert, "cert with invalid params should be nil")
	assert.NotNil(t, err, "should reject invalid keystore (no private key)")

	cert, err = buildCertificate("/dev/null", "")
	assert.Nil(t, cert, "cert with invalid params should be nil")
	assert.NotNil(t, err, "should reject invalid keystore (empty)")
}

func TestCipherSuitePreference(t *testing.T) {
	tmpCaBundle, err := ioutil.TempFile("", "ghostunnel-test")
	panicOnError(err)

	tmpCaBundle.WriteString(testCertificate)
	tmpCaBundle.WriteString("\n")

	tmpCaBundle.Sync()
	defer os.Remove(tmpCaBundle.Name())

	conf, err := buildConfig("XYZ", tmpCaBundle.Name())
	assert.NotNil(t, err, "should not be able to build TLS config with invalid cipher suite option")

	conf, err = buildConfig("", tmpCaBundle.Name())
	assert.NotNil(t, err, "should not be able to build TLS config wihout cipher suite selection")

	conf, err = buildConfig("CHACHA,AES", tmpCaBundle.Name())
	assert.Nil(t, err, "should be able to build TLS config")
	assert.True(t, conf.CipherSuites[0] == tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, "expecting ChaCha20")

	conf, err = buildConfig("AES,CHACHA", tmpCaBundle.Name())
	assert.Nil(t, err, "should be able to build TLS config")
	assert.True(t, conf.CipherSuites[0] == tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, "expecting AES")
}

func TestReload(t *testing.T) {
	tmpKeystore, err := ioutil.TempFile("", "ghostunnel-test")
	panicOnError(err)

	tmpKeystore.Write(testKeystore)
	tmpKeystore.Sync()

	defer os.Remove(tmpKeystore.Name())

	c, err := buildCertificate(tmpKeystore.Name(), testKeystorePassword)
	assert.Nil(t, err, "should be able to build certificate")

	c.reload()
}

func TestBuildConfigSystemRoots(t *testing.T) {
	if runtime.GOOS == "windows" {
		// System roots are not supported on Windows
		t.SkipNow()
		return
	}
	conf, err := buildConfig("AES", "")
	assert.Nil(t, err, "should be able to build TLS config")
	assert.NotNil(t, conf.RootCAs, "config must have CA certs")
	assert.NotNil(t, conf.ClientCAs, "config must have CA certs")
	assert.True(t, conf.MinVersion == tls.VersionTLS12, "must have correct TLS min version")
}

func TestTimeoutError(t *testing.T) {
	err := timeoutError{}
	assert.False(t, err.Error() == "", "Timeout error should have message")
	assert.True(t, err.Timeout(), "Timeout error should have Timeout() == true")
	assert.True(t, err.Temporary(), "Timeout error should have Temporary() == true")
}
