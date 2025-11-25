package chroma

import (
	"encoding/json"
	"testing"

	chromav2 "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/stretchr/testify/require"
)

func TestDocumentMetadataToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    map[string]any
		expected map[string]any
	}{
		{
			name: "string value",
			input: map[string]any{
				"location": "kitchen",
			},
			expected: map[string]any{
				"location": "kitchen",
			},
		},
		{
			name: "int value",
			input: map[string]any{
				"count": 42,
			},
			expected: map[string]any{
				"count": int64(42), // GetInt returns int64
			},
		},
		{
			name: "float value",
			input: map[string]any{
				"score": 3.14,
			},
			expected: map[string]any{
				"score": 3.14,
			},
		},
		{
			name: "bool value",
			input: map[string]any{
				"active": true,
			},
			expected: map[string]any{
				"active": true,
			},
		},
		{
			name: "mixed types",
			input: map[string]any{
				"name":   "test",
				"count":  100,
				"score":  9.5,
				"active": false,
			},
			expected: map[string]any{
				"name":   "test",
				"count":  int64(100),
				"score":  9.5,
				"active": false,
			},
		},
		{
			name:     "empty metadata",
			input:    map[string]any{},
			expected: map[string]any{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create DocumentMetadata from the input map
			dm, err := chromav2.NewDocumentMetadataFromMap(tc.input)
			require.NoError(t, err)

			// Convert back using our function
			result := documentMetadataToMap(dm)

			// Verify all expected keys and values are present
			require.Len(t, result, len(tc.expected))
			for key, expectedVal := range tc.expected {
				actualVal, ok := result[key]
				require.True(t, ok, "key %q should exist in result", key)
				require.Equal(t, expectedVal, actualVal, "value for key %q should match", key)
			}
		})
	}
}

func TestDocumentMetadataToMap_JSONMarshal(t *testing.T) {
	t.Parallel()

	// This test verifies that the converted metadata can be properly JSON marshaled,
	// which was the original bug - MetadataValue structs would serialize as {}
	input := map[string]any{
		"location":    "office",
		"square_feet": 500,
		"rating":      4.5,
		"is_open":     true,
	}

	dm, err := chromav2.NewDocumentMetadataFromMap(input)
	require.NoError(t, err)

	result := documentMetadataToMap(dm)

	// Marshal to JSON
	jsonBytes, err := json.Marshal(result)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled map[string]any
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	require.NoError(t, err)

	// Verify values are preserved (note: JSON numbers become float64)
	require.Equal(t, "office", unmarshaled["location"])
	require.Equal(t, float64(500), unmarshaled["square_feet"]) // JSON numbers are float64
	require.Equal(t, 4.5, unmarshaled["rating"])
	require.Equal(t, true, unmarshaled["is_open"])
}

func TestDocumentMetadataToMap_NilInput(t *testing.T) {
	t.Parallel()

	// Test with nil DocumentMetadata
	result := documentMetadataToMap(nil)
	require.NotNil(t, result)
	require.Empty(t, result)
}
