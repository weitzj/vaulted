// Copyright 2018 SumUp Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sumup-oss/go-pkgs/os/ostest"
	theseusTestUtils "github.com/sumup-oss/go-pkgs/testutils"
	"github.com/sumup-oss/vaulted/pkg/aes"
	"github.com/sumup-oss/vaulted/pkg/base64"
	"github.com/sumup-oss/vaulted/pkg/hcl"
	"github.com/sumup-oss/vaulted/pkg/ini"
	"github.com/sumup-oss/vaulted/pkg/pkcs7"
	"github.com/sumup-oss/vaulted/pkg/rsa"
	"github.com/sumup-oss/vaulted/pkg/terraform"
	"github.com/sumup-oss/vaulted/pkg/terraform_encryption_migration"
	"github.com/sumup-oss/vaulted/pkg/vaulted/content"
	"github.com/sumup-oss/vaulted/pkg/vaulted/header"
	"github.com/sumup-oss/vaulted/pkg/vaulted/passphrase"
	"github.com/sumup-oss/vaulted/pkg/vaulted/payload"
)

func TestNewTerraformCmd(t *testing.T) {
	t.Parallel()

	osExecutor := ostest.NewFakeOsExecutor(t)
	b64Svc := base64.NewBase64Service()
	rsaSvc := rsa.NewRsaService(osExecutor)
	aesSvc := aes.NewAesService(pkcs7.NewPkcs7Service())
	encPassphraseSvc := passphrase.NewEncryptedPassphraseService(b64Svc, rsaSvc)
	legacyEncContentSvc := content.NewLegacyEncryptedContentService(b64Svc, aesSvc)
	encContentSvc := content.NewV1EncryptedContentService(b64Svc, aesSvc)
	encPayloadSvc := payload.NewEncryptedPayloadService(
		header.NewHeaderService(),
		encPassphraseSvc,
		encContentSvc,
	)
	hclSvc := hcl.NewHclService()
	tfSvc := terraform.NewTerraformService()
	tfEncMigrationSvc := terraform_encryption_migration.NewTerraformEncryptionMigrationService(tfSvc)

	actual := NewTerraformCmd(
		osExecutor,
		rsaSvc,
		ini.NewIniService(),
		encPassphraseSvc,
		legacyEncContentSvc,
		encPayloadSvc,
		hclSvc,
		tfSvc,
		tfEncMigrationSvc,
	)

	assert.Equal(t, "terraform", actual.Use)
	assert.Equal(t, "Terraform resources related commands", actual.Short)
	assert.Equal(t, "Terraform resources related commands", actual.Long)
}

func TestTerraformCmd_Execute(t *testing.T) {
	t.Parallel()

	outputBuff := &bytes.Buffer{}

	osExecutor := ostest.NewFakeOsExecutor(t)
	osExecutor.On("Stdout").Return(outputBuff)

	b64Svc := base64.NewBase64Service()
	rsaSvc := rsa.NewRsaService(osExecutor)
	aesSvc := aes.NewAesService(pkcs7.NewPkcs7Service())
	encPassphraseSvc := passphrase.NewEncryptedPassphraseService(b64Svc, rsaSvc)
	encContentSvc := content.NewV1EncryptedContentService(b64Svc, aesSvc)
	encPayloadSvc := payload.NewEncryptedPayloadService(
		header.NewHeaderService(),
		encPassphraseSvc,
		encContentSvc,
	)
	hclSvc := hcl.NewHclService()
	tfSvc := terraform.NewTerraformService()
	legacyEncContentSvc := content.NewLegacyEncryptedContentService(b64Svc, aesSvc)
	tfEncMigrationSvc := terraform_encryption_migration.NewTerraformEncryptionMigrationService(tfSvc)

	cmdInstance := NewTerraformCmd(
		osExecutor,
		rsaSvc,
		ini.NewIniService(),
		encPassphraseSvc,
		legacyEncContentSvc,
		encPayloadSvc,
		hclSvc,
		tfSvc,
		tfEncMigrationSvc,
	)

	_, err := theseusTestUtils.RunCommandInSameProcess(
		cmdInstance,
		[]string{},
		outputBuff,
	)

	assert.Equal(t, "Use `--help` to see available commands", outputBuff.String())
	assert.Nil(t, err)

	osExecutor.AssertExpectations(t)
}
