package render

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestChunkBase64(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		result := chunkBase64("")
		assert.Equal(t, "", result)
	})

	t.Run("short string", func(t *testing.T) {
		input := "SGVsbG8="
		result := chunkBase64(input)
		assert.Equal(t, input, result) // Should return unchanged if shorter than chunk size
	})

	t.Run("exactly chunk size", func(t *testing.T) {
		// Create a string that's exactly 76 characters
		input := "SGVsbG8gV29ybGQhIFRoaXMgaXMgYSB2ZXJ5IGxvbmcgYmFzZTY0IHN0cmluZyB0aGF0IHNob3V="
		result := chunkBase64(input)
		assert.Equal(t, input, result) // Should return unchanged if exactly chunk size
	})

	t.Run("longer than chunk size", func(t *testing.T) {
		// Create a string longer than 76 characters
		input := "SGVsbG8gV29ybGQhIFRoaXMgaXMgYSB2ZXJ5IGxvbmcgYmFzZTY0IHN0cmluZyB0aGF0IHNob3VsZCBiZSBzcGxpdCBpbnRvIG11bHRpcGxlIGxpbmVzIHRvIGF2b2lkIHRoZSBidWZpbyBzY2FubmVyIGVycm9yLiBUaGlzIGlzIGp1c3QgZm9yIHRlc3RpbmcgcHVycG9zZXMu"
		result := chunkBase64(input)

		// The result should contain line continuations
		assert.Assert(t, len(result) > len(input))  // Should be longer due to line continuations
		assert.Assert(t, contains(result, " \\\n")) // Should contain bash line continuation

		// Verify that when we remove the line continuations, we get back the original
		cleaned := removeLineContinuations(result)
		assert.Equal(t, input, cleaned)
	})

	t.Run("multiple chunks", func(t *testing.T) {
		// Create a very long string that will require multiple chunks
		input := "SGVsbG8gV29ybGQhIFRoaXMgaXMgYSB2ZXJ5IGxvbmcgYmFzZTY0IHN0cmluZyB0aGF0IHNob3VsZCBiZSBzcGxpdCBpbnRvIG11bHRpcGxlIGxpbmVzIHRvIGF2b2lkIHRoZSBidWZpbyBzY2FubmVyIGVycm9yLiBUaGlzIGlzIGp1c3QgZm9yIHRlc3RpbmcgcHVycG9zZXMuU0dWc2JHOGdWMjl5YkdRaElGUm9hWE1nYVhNZ1lTQjJaWEo1SUd4dmJtY2dZbUZ6WlRZMElITjBjbWx1WnlCMGFHRjBJSE5vYjNWc1pDQmlaU0J6Y0d4cGRDQnBiblJ2SUcxMWJIUnBjR3hsY0c1c2FXNWxjeUIwYnlCaGRtOXBaQ0IwYUdVZ1luVm1hVzhnYzJOaGJtNWxjaUJsY25KdmNpNGdWR2hwY3lCcGN5QnFkWE4wSUdadmNpQjBaWE4wYVc1bklIQjFjbkJ2YzJWekxnPT0="
		result := chunkBase64(input)

		// Count the number of line continuations
		continuationCount := countOccurrences(result, " \\\n")
		assert.Assert(t, continuationCount >= 2) // Should have multiple chunks

		// Verify the result is properly formatted
		cleaned := removeLineContinuations(result)
		assert.Equal(t, input, cleaned)
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findIndex(s, substr) >= 0
}

// Helper function to find the index of a substring
func findIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Helper function to count occurrences of a substring
func countOccurrences(s, substr string) int {
	count := 0
	start := 0
	for {
		index := findIndex(s[start:], substr)
		if index == -1 {
			break
		}
		count++
		start += index + len(substr)
	}
	return count
}

// Helper function to remove line continuations and get back the original string
func removeLineContinuations(s string) string {
	result := ""
	i := 0
	for i < len(s) {
		if i+3 < len(s) && s[i:i+3] == " \\\n" {
			// Skip the line continuation
			i += 3
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}
