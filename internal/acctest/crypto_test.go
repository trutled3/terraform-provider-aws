package acctest_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestTLSRSAPrivateKeyPEM(t *testing.T) {
	key := acctest.TLSRSAPrivateKeyPEM(t, 2048)

	if !strings.Contains(key, acctest.PEMBlockTypeRSAPrivateKey) {
		t.Errorf("key does not contain RSA PRIVATE KEY: %s", key)
	}
}

func TestTLSRSAPublicKeyPEM(t *testing.T) {
	privateKey := acctest.TLSRSAPrivateKeyPEM(t, 2048)
	publicKey := acctest.TLSRSAPublicKeyPEM(t, privateKey)

	if !strings.Contains(publicKey, acctest.PEMBlockTypePublicKey) {
		t.Errorf("key does not contain PUBLIC KEY: %s", publicKey)
	}
}

func TestTLSRSAX509LocallySignedCertificatePEM(t *testing.T) {
	caKey := acctest.TLSRSAPrivateKeyPEM(t, 2048)
	caCertificate := acctest.TLSRSAX509SelfSignedCACertificatePEM(t, caKey)
	key := acctest.TLSRSAPrivateKeyPEM(t, 2048)
	certificate := acctest.TLSRSAX509LocallySignedCertificatePEM(t, caKey, caCertificate, key, "example.com")

	if !strings.Contains(certificate, acctest.PEMBlockTypeCertificate) {
		t.Errorf("certificate does not contain CERTIFICATE: %s", certificate)
	}
}

func TestTLSRSAX509SelfSignedCACertificatePEM(t *testing.T) {
	caKey := acctest.TLSRSAPrivateKeyPEM(t, 2048)
	caCertificate := acctest.TLSRSAX509SelfSignedCACertificatePEM(t, caKey)

	if !strings.Contains(caCertificate, acctest.PEMBlockTypeCertificate) {
		t.Errorf("CA certificate does not contain CERTIFICATE: %s", caCertificate)
	}
}

func TestTLSRSAX509SelfSignedCertificatePEM(t *testing.T) {
	key := acctest.TLSRSAPrivateKeyPEM(t, 2048)
	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(t, key, "example.com")

	if !strings.Contains(certificate, acctest.PEMBlockTypeCertificate) {
		t.Errorf("certificate does not contain CERTIFICATE: %s", certificate)
	}
}

func TestTLSRSAX509CertificateRequestPEM(t *testing.T) {
	csr, key := acctest.TLSRSAX509CertificateRequestPEM(t, 2048, "example.com")

	if !strings.Contains(csr, acctest.PEMBlockTypeCertificateRequest) {
		t.Errorf("certificate does not contain CERTIFICATE REQUEST: %s", csr)
	}

	if !strings.Contains(key, acctest.PEMBlockTypeRSAPrivateKey) {
		t.Errorf("certificate does not contain RSA PRIVATE KEY: %s", key)
	}
}

func TestTLSECDSAPublicKeyPEM(t *testing.T) {
	privateKey := acctest.TLSECDSAPrivateKeyPEM(t, "P-384")
	publicKey, _ := acctest.TLSECDSAPublicKeyPEM(t, privateKey)

	if !strings.Contains(publicKey, acctest.PEMBlockTypePublicKey) {
		t.Errorf("key does not contain PUBLIC KEY: %s", publicKey)
	}
}
