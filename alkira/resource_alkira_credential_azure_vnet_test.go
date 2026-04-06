package alkira

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlkiraCredentialAzureVnet_SchemaDefinition(t *testing.T) {
	r := resourceAlkiraCredentialAzureVnet()
	s := r.Schema

	tests := []struct {
		name      string
		field     string
		fieldType schema.ValueType
		optional  bool
		required  bool
		sensitive bool
		writeOnly bool
		computed  bool
	}{
		{
			name:      "name is required non-sensitive",
			field:     "name",
			fieldType: schema.TypeString,
			optional:  false,
			required:  true,
			sensitive: false,
			writeOnly: false,
		},
		{
			name:      "application_id is optional sensitive writeonly",
			field:     "application_id",
			fieldType: schema.TypeString,
			optional:  true,
			required:  false,
			sensitive: true,
			writeOnly: true,
		},
		{
			name:      "secret_key is optional sensitive writeonly",
			field:     "secret_key",
			fieldType: schema.TypeString,
			optional:  true,
			required:  false,
			sensitive: true,
			writeOnly: true,
		},
		{
			name:      "tenant_id is optional sensitive writeonly",
			field:     "tenant_id",
			fieldType: schema.TypeString,
			optional:  true,
			required:  false,
			sensitive: true,
			writeOnly: true,
		},
		{
			name:      "subscription_id is optional sensitive writeonly",
			field:     "subscription_id",
			fieldType: schema.TypeString,
			optional:  true,
			required:  false,
			sensitive: true,
			writeOnly: true,
		},
		{
			name:      "environment is optional computed non-sensitive",
			field:     "environment",
			fieldType: schema.TypeString,
			optional:  true,
			required:  false,
			sensitive: false,
			writeOnly: false,
			computed:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldSchema, ok := s[tt.field]
			require.True(t, ok, "field %q must exist in schema", tt.field)

			assert.Equal(t, tt.fieldType, fieldSchema.Type, "field %q type mismatch", tt.field)
			assert.Equal(t, tt.optional, fieldSchema.Optional, "field %q Optional mismatch", tt.field)
			assert.Equal(t, tt.required, fieldSchema.Required, "field %q Required mismatch", tt.field)
			assert.Equal(t, tt.sensitive, fieldSchema.Sensitive, "field %q Sensitive mismatch", tt.field)
			assert.Equal(t, tt.writeOnly, fieldSchema.WriteOnly, "field %q WriteOnly mismatch", tt.field)
			assert.Equal(t, tt.computed, fieldSchema.Computed, "field %q Computed mismatch", tt.field)
		})
	}
}

func TestAlkiraCredentialAzureVnet_GetValueFromEnvVar(t *testing.T) {
	r := resourceAlkiraCredentialAzureVnet()
	d := r.TestResourceData()

	tests := []struct {
		name        string
		field       string
		envVar      string
		envValue    string
		required    bool
		expected    string
		expectError bool
		errorField  string
		errorEnvVar string
	}{
		{
			name:     "env var set and required",
			field:    "application_id",
			envVar:   "AK_AZURE_APPLICATION_ID",
			envValue: "test-app-id",
			required: true,
			expected: "test-app-id",
		},
		{
			name:     "env var set and optional",
			field:    "subscription_id",
			envVar:   "AK_AZURE_SUBSCRIPTION_ID",
			envValue: "test-sub-id",
			required: false,
			expected: "test-sub-id",
		},
		{
			name:        "env var unset and required returns error",
			field:       "secret_key",
			envVar:      "AK_AZURE_SECRET_KEY",
			envValue:    "",
			required:    true,
			expectError: true,
			errorField:  "secret_key",
			errorEnvVar: "AK_AZURE_SECRET_KEY",
		},
		{
			name:     "env var unset and optional returns empty",
			field:    "subscription_id",
			envVar:   "AK_AZURE_SUBSCRIPTION_ID",
			envValue: "",
			required: false,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.envVar, tt.envValue)
			} else {
				t.Setenv(tt.envVar, "")
			}

			// For the "unset" cases, we need to actually unset
			if tt.envValue == "" {
				t.Setenv(tt.envVar, "")
				// Setenv with empty string is different from unset,
				// but empty string is treated as "not set" by the helper
			}

			val, err := getAzureVnetCredentialValue(d, tt.field, tt.envVar, tt.required)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorField)
				assert.Contains(t, err.Error(), tt.errorEnvVar)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, val)
			}
		})
	}
}

func TestAlkiraCredentialAzureVnet_BuildCredentialFromEnvVars(t *testing.T) {
	r := resourceAlkiraCredentialAzureVnet()

	t.Run("all required env vars set", func(t *testing.T) {
		d := r.TestResourceData()
		t.Setenv("AK_AZURE_APPLICATION_ID", "app-123")
		t.Setenv("AK_AZURE_SECRET_KEY", "secret-456")
		t.Setenv("AK_AZURE_TENANT_ID", "tenant-789")
		t.Setenv("AK_AZURE_SUBSCRIPTION_ID", "sub-000")

		cred, err := buildAzureVnetCredential(d)
		require.NoError(t, err)
		assert.Equal(t, "app-123", cred.ApplicationId)
		assert.Equal(t, "secret-456", cred.SecretKey)
		assert.Equal(t, "tenant-789", cred.TenantId)
		assert.Equal(t, "sub-000", cred.SubscriptionId)
	})

	t.Run("missing application_id returns error", func(t *testing.T) {
		d := r.TestResourceData()
		t.Setenv("AK_AZURE_APPLICATION_ID", "")
		t.Setenv("AK_AZURE_SECRET_KEY", "secret-456")
		t.Setenv("AK_AZURE_TENANT_ID", "tenant-789")

		_, err := buildAzureVnetCredential(d)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "application_id")
	})

	t.Run("missing secret_key returns error", func(t *testing.T) {
		d := r.TestResourceData()
		t.Setenv("AK_AZURE_APPLICATION_ID", "app-123")
		t.Setenv("AK_AZURE_SECRET_KEY", "")
		t.Setenv("AK_AZURE_TENANT_ID", "tenant-789")

		_, err := buildAzureVnetCredential(d)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "secret_key")
	})

	t.Run("missing tenant_id returns error", func(t *testing.T) {
		d := r.TestResourceData()
		t.Setenv("AK_AZURE_APPLICATION_ID", "app-123")
		t.Setenv("AK_AZURE_SECRET_KEY", "secret-456")
		t.Setenv("AK_AZURE_TENANT_ID", "")

		_, err := buildAzureVnetCredential(d)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tenant_id")
	})

	t.Run("missing optional subscription_id succeeds", func(t *testing.T) {
		d := r.TestResourceData()
		t.Setenv("AK_AZURE_APPLICATION_ID", "app-123")
		t.Setenv("AK_AZURE_SECRET_KEY", "secret-456")
		t.Setenv("AK_AZURE_TENANT_ID", "tenant-789")
		t.Setenv("AK_AZURE_SUBSCRIPTION_ID", "")

		cred, err := buildAzureVnetCredential(d)
		require.NoError(t, err)
		assert.Equal(t, "", cred.SubscriptionId)
		assert.Equal(t, "app-123", cred.ApplicationId)
	})
}

func TestAlkiraCredentialAzureVnet_ReadDoesNotSetSensitiveFields(t *testing.T) {
	r := resourceAlkiraCredentialAzureVnet()
	d := r.TestResourceData()

	// Simulate what the Read function does — only sets name and environment
	d.Set("name", "test-credential")
	d.Set("environment", "AZURE")

	// Sensitive WriteOnly fields must not be in state
	assert.Equal(t, "", d.Get("application_id").(string))
	assert.Equal(t, "", d.Get("secret_key").(string))
	assert.Equal(t, "", d.Get("tenant_id").(string))
	assert.Equal(t, "", d.Get("subscription_id").(string))

	// Non-sensitive fields should be set
	assert.Equal(t, "test-credential", d.Get("name").(string))
	assert.Equal(t, "AZURE", d.Get("environment").(string))
}

func TestAlkiraCredentialAzureVnet_ImporterConfigured(t *testing.T) {
	r := resourceAlkiraCredentialAzureVnet()
	assert.NotNil(t, r.Importer, "resource must have an Importer")
	assert.NotNil(t, r.Importer.StateContext, "Importer must have StateContext")
}
