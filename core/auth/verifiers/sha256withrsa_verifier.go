// Package verifiers 微信支付 API v3 Go SDK 数字签名验证器
package verifiers

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/wechatpay-apiv3/wechatpay-go/core/cert"
)

// SHA256WithRSAVerifier SHA256WithRSA 数字签名验证器
type SHA256WithRSAVerifier struct {
	// Certificates 微信支付平台证书Map，key: 平台证书序列号， value: 微信支付平台证书
	certProvider cert.CertificateProvider
}

func (verifier *SHA256WithRSAVerifier) CertificateProvider() cert.CertificateProvider {
	return verifier.certProvider
}

// Verify 对数字签名信息进行验证
func (verifier *SHA256WithRSAVerifier) Verify(ctx context.Context, serialNumber, message, signature string) error {
	err := checkParameter(ctx, serialNumber, message, signature)
	if err != nil {
		return err
	}
	if verifier.certProvider == nil {
		return fmt.Errorf("verifier has no validator")
	}
	certificate, ok := verifier.certProvider.GetCertificate(serialNumber)
	if !ok {
		return fmt.Errorf("certificate[%s] not found in verifier", serialNumber)
	}
	hashed := sha256.Sum256([]byte(message))
	err = rsa.VerifyPKCS1v15(certificate.PublicKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], []byte(signature))
	if err != nil {
		return fmt.Errorf("verifty signature with public key err:%s", err.Error())
	}
	return nil
}

func checkParameter(ctx context.Context, serialNumber, message, signature string) error {
	if ctx == nil {
		return fmt.Errorf("context is nil, verifier need input context.Context")
	}
	if strings.TrimSpace(serialNumber) == "" {
		return fmt.Errorf("serialNumber is empty, verifier need input serialNumber")
	}
	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("message is empty, verifier need input message")
	}
	if strings.TrimSpace(signature) == "" {
		return fmt.Errorf("signature is empty, verifier need input signature")
	}
	return nil
}

func NewSHA256WithRSAVerifier(provider cert.CertificateProvider) *SHA256WithRSAVerifier {
	return &SHA256WithRSAVerifier{certProvider: provider}
}
