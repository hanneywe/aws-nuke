package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// e2eTestCase defines a single end-to-end test case that parses a DSL,
// generates code, writes to a temp dir, and verifies expected files and patterns.
type e2eTestCase struct {
	name          string
	dslFile       string
	expectedFiles []string
	contentChecks map[string][]string // file suffix -> patterns that must appear
	absentChecks  map[string][]string // file suffix -> patterns that must NOT appear
}

func e2eTestCases() []e2eTestCase {
	return []e2eTestCase{
		{
			name:    "flat-simple",
			dslFile: "testdata/flat-simple.yaml",
			expectedFiles: []string{
				"resources/mailmanager-relay.go",
				"resources/mailmanager.go",
				"resources/mailmanager_mock_test.go",
				"resources/mailmanager-relay_mock_test.go",
			},
			contentChecks: map[string][]string{
				"mailmanager-relay.go": {
					"MailManagerRelayResource",
					"MailManagerRelayLister",
					"type MailManagerRelay struct",
					"func (r *MailManagerRelay) Remove(",
					"func (r *MailManagerRelay) Properties()",
					"func (r *MailManagerRelay) String()",
					"registry.Register",
					"ListRelays",
					"DeleteRelay",
					"NewListRelaysPaginator",
				},
				"mailmanager.go": {
					"MailmanagerClient",
					"ListRelays(",
					"DeleteRelay(",
				},
				"mailmanager_mock_test.go": {
					"mockMailmanagerClient",
					"ListRelays(",
					"DeleteRelay(",
				},
				"mailmanager-relay_mock_test.go": {
					"Test_Mock_MailManagerRelay_List",
					"Test_Mock_MailManagerRelay_List_Empty",
					"Test_Mock_MailManagerRelay_Remove",
					"Test_Mock_MailManagerRelay_Properties",
					"Test_Mock_MailManagerRelay_String",
				},
			},
		},
		{
			name:    "flat-with-describe",
			dslFile: "testdata/flat-with-describe.yaml",
			expectedFiles: []string{
				"resources/eks-clusters.go",
			},
			contentChecks: map[string][]string{
				"eks-clusters.go": {
					"EKSClusterResource",
					"EKSClusterLister",
					"type EKSCluster struct",
					"DescribeCluster",
					"*eks.Client",
					"func (r *EKSCluster) Remove(",
					"func (r *EKSCluster) String()",
				},
			},
			absentChecks: map[string][]string{
				"eks-clusters.go": {
					"EksAPI",
				},
			},
		},
		{
			name:    "flat-with-tags",
			dslFile: "testdata/flat-with-tags.yaml",
			expectedFiles: []string{
				"resources/vpclattice-target-group.go",
				"resources/vpclattice.go",
				"resources/vpclattice_mock_test.go",
				"resources/vpclattice-target-group_mock_test.go",
			},
			contentChecks: map[string][]string{
				"vpclattice-target-group.go": {
					"VPCLatticeTargetGroupResource",
					"VPCLatticeTargetGroupLister",
					"type VPCLatticeTargetGroup struct",
					"ListTagsForResource",
					"ListTargetGroups",
					"DeleteTargetGroup",
				},
				"vpclattice.go": {
					"ListTagsForResource(",
				},
			},
		},
		{
			name:    "nested-two-level",
			dslFile: "testdata/nested-two-level.yaml",
			expectedFiles: []string{
				"resources/s3vectors-vector-index.go",
				"resources/s3vectors.go",
				"resources/s3vectors_mock_test.go",
				"resources/s3vectors-vector-index_mock_test.go",
			},
			contentChecks: map[string][]string{
				"s3vectors-vector-index.go": {
					"S3VectorIndexResource",
					"S3VectorIndexLister",
					"type S3VectorIndex struct",
					"ListVectorBuckets",
					"ListVectorIndexes",
					"DeleteVectorIndex",
				},
			},
		},
		{
			name:    "nested-three-level",
			dslFile: "testdata/nested-three-level.yaml",
			expectedFiles: []string{
				"resources/mediapackagev2-origin-endpoint.go",
				"resources/mediapackagev2.go",
				"resources/mediapackagev2_mock_test.go",
				"resources/mediapackagev2-origin-endpoint_mock_test.go",
			},
			contentChecks: map[string][]string{
				"mediapackagev2-origin-endpoint.go": {
					"MediaPackageV2OriginEndpointResource",
					"MediaPackageV2OriginEndpointLister",
					"type MediaPackageV2OriginEndpoint struct",
					"ListChannelGroups",
					"ListChannels",
					"ListOriginEndpoints",
					"DeleteOriginEndpoint",
				},
			},
		},
		{
			name:    "singleton",
			dslFile: "testdata/singleton.yaml",
			expectedFiles: []string{
				"resources/emr-block-public-access-configuration.go",
				"resources/emr.go",
				"resources/emr_mock_test.go",
				"resources/emr-block-public-access-configuration_mock_test.go",
			},
			contentChecks: map[string][]string{
				"emr-block-public-access-configuration.go": {
					"EMRBlockPublicAccessConfigurationResource",
					"EMRBlockPublicAccessConfigurationLister",
					"type EMRBlockPublicAccessConfiguration struct",
					"GetBlockPublicAccessConfiguration",
					"PutBlockPublicAccessConfiguration",
				},
			},
		},
		{
			name:    "with-filters",
			dslFile: "testdata/with-filters.yaml",
			expectedFiles: []string{
				"resources/chimesdkvoice-phone-number.go",
				"resources/chimesdkvoice.go",
				"resources/chimesdkvoice_mock_test.go",
				"resources/chimesdkvoice-phone-number_mock_test.go",
			},
			contentChecks: map[string][]string{
				"chimesdkvoice-phone-number.go": {
					"ChimeSDKVoicePhoneNumberResource",
					"func (r *ChimeSDKVoicePhoneNumber) Filter()",
					"DeleteInProgress",
					"Cancel" + "led",
					"ReleaseInProgress",
					"already being deleted or cancel" + "led",
				},
			},
		},
		{
			name:    "with-settings",
			dslFile: "testdata/with-settings.yaml",
			expectedFiles: []string{
				"resources/eks-clusters.go",
			},
			contentChecks: map[string][]string{
				"eks-clusters.go": {
					"EKSClusterResource",
					"DisableDeletionProtection",
					"func (r *EKSCluster) Settings(",
					"settings",
					"protection",
					"*eks.Client",
					"registry.Register",
				},
			},
		},
		{
			name:    "with-pre-deletion",
			dslFile: "testdata/with-pre-deletion.yaml",
			expectedFiles: []string{
				"resources/vpclattice-target-group.go",
				"resources/vpclattice.go",
				"resources/vpclattice_mock_test.go",
				"resources/vpclattice-target-group_mock_test.go",
			},
			contentChecks: map[string][]string{
				"vpclattice-target-group.go": {
					"VPCLatticeTargetGroupResource",
					"ListTargets",
					"DeregisterTargets",
					"DeleteTargetGroup",
				},
			},
		},
		{
			name:    "with-overrides",
			dslFile: "testdata/with-overrides.yaml",
			expectedFiles: []string{
				"resources/apigatewayv2-portal.go",
				"resources/apigatewayv2.go",
				"resources/apigatewayv2_mock_test.go",
				"resources/apigatewayv2-portal_mock_test.go",
			},
			contentChecks: map[string][]string{
				"apigatewayv2-portal.go": {
					"APIGatewayV2PortalResource",
					"DisablePortal",
					"DeletePortal",
					"apigatewayv2types.PublishStatusPublished",
					`apigatewayv2types "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"`,
				},
			},
		},
		{
			name:    "with-integration-test",
			dslFile: "testdata/with-integration-test.yaml",
			expectedFiles: []string{
				"resources/mailmanager-relay.go",
				"resources/mailmanager.go",
				"resources/mailmanager_mock_test.go",
				"resources/mailmanager-relay_mock_test.go",
				"resources/mailmanager-relay_test.go",
			},
			contentChecks: map[string][]string{
				"mailmanager-relay_test.go": {
					"//go:build integration",
					"TestMailManagerRelay",
					"SetupSuite",
					"TearDownSuite",
					"TestList",
					"TestRemove",
					"CreateRelay",
				},
			},
		},
	}
}

func TestE2E_GenerateFromSampleDSLs(t *testing.T) {
	for _, tc := range e2eTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)

			// 1. Parse the YAML DSL
			def, err := Parse(tc.dslFile)
			r.NoError(err, "failed to parse DSL file %s", tc.dslFile)

			// 2. Validate the parsed definition
			validationErrs := Validate(def)
			r.Empty(validationErrs, "validation errors for %s: %v", tc.dslFile, validationErrs)

			// 3. Generate all files
			tmpDir := t.TempDir()
			files, err := GenerateAll(def, tmpDir, GenerateOpts{})
			r.NoError(err, "failed to generate files for %s", tc.dslFile)

			// 4. Write files to temp directory
			manifest, err := WriteFiles(files, WriteOpts{
				OutputDir: tmpDir,
				Force:     true,
			})
			r.NoError(err, "failed to write files for %s", tc.dslFile)
			a.NotNil(manifest)

			// 5. Verify expected files exist
			for _, expectedFile := range tc.expectedFiles {
				fullPath := filepath.Join(tmpDir, expectedFile)
				_, err := os.Stat(fullPath)
				a.NoError(err, "expected file %s does not exist", expectedFile)
			}

			// 6. Verify generated code contains expected patterns
			for fileSuffix, patterns := range tc.contentChecks {
				fullPath := findFileWithSuffix(tmpDir, fileSuffix)
				r.NotEmpty(fullPath, "could not find file ending with %s", fileSuffix)

				content, err := os.ReadFile(fullPath)
				r.NoError(err, "failed to read %s", fullPath)

				for _, pattern := range patterns {
					a.Contains(string(content), pattern,
						"file %s should contain %q", fileSuffix, pattern)
				}
			}

			// 7. Verify absent patterns (things that should NOT appear)
			for fileSuffix, patterns := range tc.absentChecks {
				fullPath := findFileWithSuffix(tmpDir, fileSuffix)
				if fullPath == "" {
					continue
				}

				content, err := os.ReadFile(fullPath)
				r.NoError(err, "failed to read %s", fullPath)

				for _, pattern := range patterns {
					a.NotContains(string(content), pattern,
						"file %s should NOT contain %q", fileSuffix, pattern)
				}
			}
		})
	}
}

// findFileWithSuffix walks the temp directory and returns the first file
// whose path ends with the given suffix.
func findFileWithSuffix(root, suffix string) string {
	var found string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && filepath.Base(path) == filepath.Base(suffix) {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	return found
}
