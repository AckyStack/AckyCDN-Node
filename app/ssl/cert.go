package ssl

import (
	"ackycdn-node/app/vhost"
	"crypto/tls"
	"errors"
)

func getCertificateByDomainName(domainName string) (*tls.Certificate, error) {
	vhc := vhost.GetConfigMem(domainName)
	if vhc == nil {
		return nil, errors.New("vhc config not found")
	}
	c, err := tls.X509KeyPair(vhc.TlsConfig.Certificate, vhc.TlsConfig.Key)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
