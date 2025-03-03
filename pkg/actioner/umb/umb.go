package umb

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/devtools-qe-incubator/eventmanager/pkg/configuration/providers"
	"github.com/devtools-qe-incubator/eventmanager/pkg/services/messaging/umb"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util/logging"
)

func Send(providersFilePath, destination, eventFilePath string) error {
	if err := setupUMBClient(providersFilePath); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	message, err := ioutil.ReadFile(eventFilePath)
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	messageUnFormatted := strings.ReplaceAll(string(message), "\n", "")
	return umb.SendBytes(destination, []byte(messageUnFormatted))
}

func setupUMBClient(providersFilePath string) (err error) {
	providersInfo, err := providers.LoadFile(providersFilePath)
	if err != nil {
		return err
	}
	if util.IsEmpty(providersInfo.UMB) {
		return fmt.Errorf("umb provider configuration is required")
	}
	userCertificate, userKey, certificateAuthority, err := providers.ParseUMBFiles(providersInfo.UMB)
	if err != nil {
		return err
	}
	err = umb.InitClient(
		providersInfo.UMB.ConsumerID,
		providersInfo.UMB.Driver,
		strings.Split(providersInfo.UMB.Brokers, ","),
		userCertificate,
		userKey,
		certificateAuthority)
	return
}
