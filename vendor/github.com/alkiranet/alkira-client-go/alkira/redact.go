// Copyright (C) 2026 Alkira Inc. All Rights Reserved.

package alkira

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

const redactedValue = "[REDACTED]"

// Headers whose values are session or API secrets.
var sensitiveHeaders = []string{"Authorization", "Proxy-Authorization", "Cookie", "Set-Cookie", "X-Api-Key"}

// JSON keys (lower-cased, without '_' or '-') whose values are always secret.
var sensitiveKeys = map[string]bool{
	"key":                true,
	"credentials":        true,
	"licensekey":         true,
	"registrationkey":    true,
	"fmcregistrationkey": true,
	"masterkey":          true,
	"sickey":             true,
	"activationkey":      true,
	"accountkey":         true,
	"privatekeyid":       true,
}

// Substrings of a normalised JSON key that mark its value as secret.
var sensitiveKeyParts = []string{
	"password", "secret", "token", "privatekey", "presharedkey", "apikey",
	"authkey", "authenticationkey", "authvalue", "authcode", "sharedkey",
}

// Suffixes that mark a key as a reference or a flag rather than the secret itself.
var nonSecretKeySuffixes = []string{"credentialid", "enabled", "type", "uri"}

func isSensitiveKey(key string) bool {
	k := strings.ToLower(strings.NewReplacer("_", "", "-", "").Replace(key))
	for _, suffix := range nonSecretKeySuffixes {
		if strings.HasSuffix(k, suffix) {
			return false
		}
	}
	if sensitiveKeys[k] {
		return true
	}
	for _, part := range sensitiveKeyParts {
		if strings.Contains(k, part) {
			return true
		}
	}
	return false
}

// redactHeaders returns a copy of h with session and API secrets masked.
func redactHeaders(h http.Header) http.Header {
	out := h.Clone()
	for _, name := range sensitiveHeaders {
		if _, ok := out[http.CanonicalHeaderKey(name)]; ok {
			out.Set(name, redactedValue)
		}
	}
	return out
}

// redactBody returns a JSON body with the values of secret fields masked, for logs and
// error messages. A body that is not JSON, or has nothing to mask, is returned unchanged.
func redactBody(data []byte) string {
	var v interface{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&v); err != nil {
		return string(data)
	}

	v, changed := redactValue(v)
	if !changed {
		return string(data)
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return redactedValue
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

func redactValue(v interface{}) (interface{}, bool) {
	changed := false
	switch val := v.(type) {
	case map[string]interface{}:
		for k, child := range val {
			if isSensitiveKey(k) {
				if child != nil && child != "" {
					val[k] = redactedValue
					changed = true
				}
				continue
			}
			if redacted, c := redactValue(child); c {
				val[k] = redacted
				changed = true
			}
		}
	case []interface{}:
		for i, child := range val {
			if redacted, c := redactValue(child); c {
				val[i] = redacted
				changed = true
			}
		}
	}
	return v, changed
}
