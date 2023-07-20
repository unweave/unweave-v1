//nolint:paralleltest,testpackage
package endpointsrv

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unweave/unweave/api/types"
)

func TestAbsoluteURL(t *testing.T) {
	type testCase struct {
		name      string
		baseURL   string
		reference string
		expected  string
	}

	testCases := []testCase{
		{
			name:      "absolute reference",
			baseURL:   "https://example.com",
			reference: "http://somewhere.com/ok",
			expected:  "http://somewhere.com/ok",
		},
		{
			name:      "relative reference",
			baseURL:   "https://example.com/path",
			reference: "/relative",
			expected:  "https://example.com/relative",
		},
		{
			name:      "empty reference",
			baseURL:   "https://example.com/path",
			reference: "",
			expected:  "https://example.com/path",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			u, err := absoluteURL(test.reference, test.baseURL)

			require.NoError(t, err)
			require.Equal(t, test.expected, u)
		})
	}
}

func TestResolveManifestURLs(t *testing.T) {
	type testCase struct {
		name        string
		input       types.EvalManifest
		endpointURL string
		evalURL     string
		expected    types.EvalManifest
	}

	testCases := []testCase{
		{
			name:        "defaults",
			input:       types.EvalManifest{},
			endpointURL: "endpoint.example",
			evalURL:     "eval.example",
			expected: types.EvalManifest{
				DatasetURL: "https://eval.example/dataset",
				AssertURL:  "https://eval.example/assert",
				RunURL:     "https://endpoint.example/",
			},
		},
		{
			name: "paths",
			input: types.EvalManifest{
				DatasetURL: "/updated-dataset",
				AssertURL:  "/updated-assert",
				RunURL:     "/run",
			},
			endpointURL: "endpoint.example",
			evalURL:     "eval.example",
			expected: types.EvalManifest{
				DatasetURL: "https://eval.example/updated-dataset",
				AssertURL:  "https://eval.example/updated-assert",
				RunURL:     "https://eval.example/run",
			},
		},
		{
			name: "absolute urls",
			input: types.EvalManifest{
				DatasetURL: "http://somewhere/updated-dataset",
				AssertURL:  "http://somewhere/updated-assert",
				RunURL:     "http://somewhere/run",
			},
			endpointURL: "endpoint.example",
			evalURL:     "eval.example",
			expected: types.EvalManifest{
				DatasetURL: "http://somewhere/updated-dataset",
				AssertURL:  "http://somewhere/updated-assert",
				RunURL:     "http://somewhere/run",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			manifest, err := resolveManifestURLs(
				test.input,
				types.Endpoint{HTTPAddress: test.endpointURL},
				types.Eval{HTTPEndpoint: test.evalURL},
			)

			require.NoError(t, err)
			require.Equal(t, test.expected, manifest)
		})
	}
}
