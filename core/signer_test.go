package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuthorizationHMACSHA256Golden(t *testing.T) {
	t.Parallel()
	raw, err := os.ReadFile(filepath.Join("testdata", "authorization_hmac_sha256.json"))
	if err != nil {
		t.Fatal(err)
	}
	var vec struct {
		URI                     string `json:"uri"`
		URIViaKong              string `json:"uri_via_kong"`
		Body                    string `json:"body"`
		Secret                  string `json:"secret"`
		Signature               string `json:"signature"`
		SignatureViaKong        string `json:"signature_via_kong"`
		AuthorizationHeader     string `json:"authorization_header"`
		AuthorizationHeaderKong string `json:"authorization_header_via_kong"`
		Nonce                   string `json:"nonce"`
		Timestamp               int64  `json:"timestamp"`
		BodyIncludesNewline     bool   `json:"body_includes_trailing_newline"`
	}
	if err := json.Unmarshal(raw, &vec); err != nil {
		t.Fatal(err)
	}
	if !vec.BodyIncludesNewline || !strings.HasSuffix(vec.Body, "\n") {
		t.Fatal("黄金向量 body 必须含 json.Encoder 尾随换行")
	}
	msg := fmt.Sprintf(SignatureFormat, vec.URI, vec.Timestamp, vec.Nonce, vec.Body)
	got := SignMessage(vec.Secret, msg)
	if got != vec.Signature {
		t.Fatalf("signature got %q want %q", got, vec.Signature)
	}
	hdr := SignAuthorization("10000000000", "shim_test", vec.Secret, vec.URI, vec.Body, vec.Timestamp, vec.Nonce)
	if hdr != vec.AuthorizationHeader {
		t.Fatalf("header got %q want %q", hdr, vec.AuthorizationHeader)
	}
	if vec.URIViaKong != "" {
		kongMsg := fmt.Sprintf(SignatureFormat, vec.URIViaKong, vec.Timestamp, vec.Nonce, vec.Body)
		kongSig := SignMessage(vec.Secret, kongMsg)
		if kongSig != vec.SignatureViaKong {
			t.Fatalf("kong signature got %q want %q", kongSig, vec.SignatureViaKong)
		}
		kongHdr := SignAuthorization("10000000000", "shim_test", vec.Secret, vec.URIViaKong, vec.Body, vec.Timestamp, vec.Nonce)
		if kongHdr != vec.AuthorizationHeaderKong {
			t.Fatalf("kong header got %q want %q", kongHdr, vec.AuthorizationHeaderKong)
		}
	}
	if strings.Contains(strings.ToLower(hdr), "rsa") {
		t.Fatal("Authorization 头不得出现 RSA")
	}
}

func TestGETEmptyBodySignature(t *testing.T) {
	t.Parallel()
	msg := fmt.Sprintf(SignatureFormat, "/v1/chat/messages", int64(1735689600), "n1", "")
	got := SignMessage("test-auth-secret", msg)
	want := "XmLNL2QLDtJMWoe/HnacSHs1t5oUKh7Q56yMCo2T0ZQ="
	if got != want {
		t.Fatalf("GET empty body signature got %q want %q", got, want)
	}
}

func TestUserVerifySignGolden(t *testing.T) {
	t.Parallel()
	got := UserVerifySign("shim_test", "tok", "test-auth-secret", 1735689600)
	want := "bI5G%2BNdANAgwcvqvBGDTUHZe990%3D"
	if got != want {
		t.Fatalf("UserVerify sign got %q want %q", got, want)
	}
}

func TestSignerTypeNameHasNoRSA(t *testing.T) {
	t.Parallel()
	var s HMACSHA256Signer
	_ = s
	name := fmt.Sprintf("%T", s)
	if strings.Contains(strings.ToLower(name), "rsa") {
		t.Fatalf("类型名含 RSA: %s", name)
	}
}
