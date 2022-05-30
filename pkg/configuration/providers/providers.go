package providers

import (
	"encoding/base64"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/file"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

func LoadFile(providersFilePath string) (*Providers, error) {
	var structuredProviders Providers
	if len(providersFilePath) > 0 {
		if err := file.LoadFileAsStruct(providersFilePath, &structuredProviders); err != nil {
			logging.Errorf("Can not load providers file: %v", err)
			return nil, err
		}
	}
	return &structuredProviders, nil
}

func ParseUMBFiles(source UMB) (userCertificate, userKey, certificateAuthority []byte, err error) {
	userCertificate, err =
		base64.StdEncoding.DecodeString(source.UserCertificate)
	if err != nil {
		return
	}
	userKey, err =
		base64.StdEncoding.DecodeString(source.UserKey)
	if err != nil {
		return
	}
	certificateAuthority, err =
		base64.StdEncoding.DecodeString(source.CertificateAuthority)
	if err != nil {
		return
	}
	return
}

func ParseGithubFiles(source Github) (appKey []byte, err error) {
	return base64.StdEncoding.DecodeString(source.AppKey)
}
