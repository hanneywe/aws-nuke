package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
	ModelOpus     = "us.anthropic.claude-opus-4-6-v1"
	BedrockRegion = "us-west-2"
	MaxRetries    = 5

	SDKGitHubRawBase = "https://raw.githubusercontent.com/aws/aws-sdk-go-v2/main/service"

	ResultsTSVPath = "tools/resource-discovery/results.tsv"
)

type triageResult struct {
	Missing []string
	Covered []coveredEntry
	Skipped []skippedEntry
}

type coveredEntry struct {
	APIEntry         string
	ExistingResource string
}

type skippedEntry struct {
	APIEntry string
	Reason   string
}

type missingResource struct {
	APIEntry     string
	GoPackage    string
	DeleteOp     string
	ListOp       string
	NukeService  string
	NukeResource string
	ResourceName string
	Scope        string
}

type invalidEntry struct {
	Resource missingResource
	Reason   string
}

func main() {
	inputFile := flag.String("input", "", "Path to text file with API entries (runs discovery, writes results.tsv)")
	generate := flag.Bool("generate", false, "Read results.tsv and run codegen for each entry")
	flag.Parse()

	if !*generate && *inputFile == "" {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  Discover:  resource-discovery --input <file>")
		fmt.Fprintln(os.Stderr, "  Generate:  resource-discovery --generate")
		os.Exit(1)
	}

	if *generate {
		runGenerateFromTSV()
		return
	}

	runDiscovery(*inputFile)
}

func runGenerateFromTSV() {
	resources, err := readTSV(ResultsTSVPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", ResultsTSVPath, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Loaded %d resources from %s\n", len(resources), ResultsTSVPath)

	ctx := context.Background()
	client, err := newBedrockClient(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating bedrock client: %v\n", err)
		os.Exit(1)
	}
	dslReference := loadDSLReference()
	runGeneration(ctx, client, resources, dslReference)
}

func runDiscovery(inputFile string) {
	apiLines, err := readLines(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Loaded %d API entries\n", len(apiLines))

	existingResources := extractExistingResources()
	fmt.Fprintf(os.Stderr, "Found %d existing resources\n", len(existingResources))

	ctx := context.Background()
	client, err := newBedrockClient(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating bedrock client: %v\n", err)
		os.Exit(1)
	}

	// Pass 1: Triage
	fmt.Fprintf(os.Stderr, "Pass 1: Triaging %d entries...\n", len(apiLines))
	triage, err := triageEntries(ctx, client, apiLines, existingResources)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error triaging: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "  %d missing, %d covered, %d skipped\n",
		len(triage.Missing), len(triage.Covered), len(triage.Skipped))

	// Pass 2: Enrich missing entries with SDK details
	sdkPackages := discoverSDKPackages()
	fmt.Fprintf(os.Stderr, "Pass 2: Enriching %d missing entries...\n", len(triage.Missing))
	enriched, err := enrichMissing(ctx, client, triage.Missing, sdkPackages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error enriching: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "  %d enriched\n", len(enriched))

	// Pass 3: Validate against GitHub SDK source
	validated, invalid := validateAgainstGitHub(enriched)

	if len(invalid) > 0 {
		fmt.Fprintf(os.Stderr, "\n%d entries failed validation (excluded from results):\n", len(invalid))
		for i := range invalid {
			fmt.Fprintf(os.Stderr, "  %s: %s\n", invalid[i].Resource.APIEntry, invalid[i].Reason)
		}
	}

	if err := writeTSV(ResultsTSVPath, validated); err != nil {
		fmt.Fprintf(os.Stderr, "error writing TSV: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "\nWrote %d entries to %s\n", len(validated), ResultsTSVPath)
	fmt.Fprintf(os.Stderr, "Review/edit the file, then run with --generate to create resources.\n")
}

// --- Pass 1: Triage ---

func triageEntries(ctx context.Context, client *bedrockruntime.Client, apiLines, existingResources []string) (*triageResult, error) {
	result := &triageResult{}
	batchSize := 50
	resourceList := strings.Join(existingResources, "\n")

	for i := 0; i < len(apiLines); i += batchSize {
		end := i + batchSize
		if end > len(apiLines) {
			end = len(apiLines)
		}
		batch := apiLines[i:end]
		fmt.Fprintf(os.Stderr, "  batch %d-%d of %d...\n", i+1, end, len(apiLines))

		prompt := "For each API entry, respond with exactly one line:\n" +
			"MISSING\t<entry>\n" +
			"COVERED\t<entry>\t<existing resource name>\n" +
			"SKIP\t<entry>\t<reason>\n\n" +
			"MISSING = no existing resource covers this and it creates deletable state\n" +
			"COVERED = an existing resource already handles cleanup\n" +
			"SKIP = config-only (Put/Set with no delete), or not a deletable resource\n\n" +
			"## API entries\n" + strings.Join(batch, "\n") + "\n\n" +
			"## Existing resources\n" + resourceList + "\n\n" +
			"Output ONLY the TSV lines. No headers, no explanation."

		response, err := invokeBedrockText(ctx, client, prompt)
		if err != nil {
			return nil, fmt.Errorf("batch %d-%d: %w", i+1, end, err)
		}

		batchResult := parseTriageResponse(response)
		result.Missing = append(result.Missing, batchResult.Missing...)
		result.Covered = append(result.Covered, batchResult.Covered...)
		result.Skipped = append(result.Skipped, batchResult.Skipped...)
	}

	return result, nil
}

func parseTriageResponse(response string) *triageResult {
	result := &triageResult{}
	for _, line := range strings.Split(response, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, "\t", 3)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MISSING":
			result.Missing = append(result.Missing, fields[1])
		case "COVERED":
			entry := coveredEntry{APIEntry: fields[1]}
			if len(fields) >= 3 {
				entry.ExistingResource = fields[2]
			}
			result.Covered = append(result.Covered, entry)
		case "SKIP":
			entry := skippedEntry{APIEntry: fields[1]}
			if len(fields) >= 3 {
				entry.Reason = fields[2]
			}
			result.Skipped = append(result.Skipped, entry)
		}
	}
	return result
}

// --- Pass 2: Enrich ---

func enrichMissing(ctx context.Context, client *bedrockruntime.Client, entries, sdkPackages []string) ([]missingResource, error) {
	if len(entries) == 0 {
		return nil, nil
	}

	var allResources []missingResource
	batchSize := 30
	packageList := strings.Join(sdkPackages, ", ")

	for i := 0; i < len(entries); i += batchSize {
		end := i + batchSize
		if end > len(entries) {
			end = len(entries)
		}
		batch := entries[i:end]
		fmt.Fprintf(os.Stderr, "  batch %d-%d of %d...\n", i+1, end, len(entries))

		prompt := "For each API entry, provide SDK details as a TSV line:\n" +
			"<entry>\t<goPackage>\t<deleteOp>\t<listOp>\t<nukeService>\t<nukeResource>\t<resourceName>\t<scope>\n\n" +
			"Rules:\n" +
			"- goPackage from: " + packageList + "\n" +
			"- nukeService: kebab-case file prefix (e.g. access-analyzer)\n" +
			"- nukeResource: kebab-case file suffix (e.g. analyzer)\n" +
			"- resourceName: PascalCase Go struct (e.g. AccessAnalyzerAnalyzer)\n" +
			"- scope: account or region\n" +
			"- deleteOp: the SDK Delete/Remove/Deregister operation\n" +
			"- listOp: the SDK List/Describe/Get operation\n\n" +
			"## Entries\n" + strings.Join(batch, "\n") + "\n\n" +
			"Output ONLY TSV lines. No headers, no explanation."

		response, err := invokeBedrockText(ctx, client, prompt)
		if err != nil {
			return nil, fmt.Errorf("batch %d-%d: %w", i+1, end, err)
		}

		for _, line := range strings.Split(response, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			fields := strings.Split(line, "\t")
			if len(fields) >= 8 {
				allResources = append(allResources, missingResource{
					APIEntry:     fields[0],
					GoPackage:    fields[1],
					DeleteOp:     fields[2],
					ListOp:       fields[3],
					NukeService:  fields[4],
					NukeResource: fields[5],
					ResourceName: fields[6],
					Scope:        fields[7],
				})
			}
		}
	}

	return allResources, nil
}

// --- Pass 3: GitHub validation ---

func validateAgainstGitHub(enriched []missingResource) ([]missingResource, []invalidEntry) {
	fmt.Fprintf(os.Stderr, "Pass 3: Validating against GitHub...\n")
	var validated []missingResource
	var invalid []invalidEntry

	for i := range enriched {
		resource := &enriched[i]
		fmt.Fprintf(os.Stderr, "  [%d/%d] %s/%s...",
			i+1, len(enriched), resource.GoPackage, resource.DeleteOp)

		if fetchSDKOperationFile(resource.GoPackage, resource.DeleteOp) == "" {
			fmt.Fprintf(os.Stderr, " delete not found\n")
			invalid = append(invalid, invalidEntry{
				*resource,
				fmt.Sprintf("%s not in %s", resource.DeleteOp, resource.GoPackage),
			})
			continue
		}
		if fetchSDKOperationFile(resource.GoPackage, resource.ListOp) == "" {
			fmt.Fprintf(os.Stderr, " list not found\n")
			invalid = append(invalid, invalidEntry{
				*resource,
				fmt.Sprintf("%s not in %s", resource.ListOp, resource.GoPackage),
			})
			continue
		}

		fmt.Fprintf(os.Stderr, " OK\n")
		validated = append(validated, *resource)
	}
	fmt.Fprintf(os.Stderr, "  %d valid, %d invalid\n", len(validated), len(invalid))
	return validated, invalid
}

func fetchSDKOperationFile(goPackage, operation string) string {
	url := fmt.Sprintf("%s/%s/api_op_%s.go", SDKGitHubRawBase, goPackage, operation)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return ""
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

// --- Generation ---

func runGeneration(ctx context.Context, client *bedrockruntime.Client, resources []missingResource, dslReference string) {
	for i := range resources {
		resource := &resources[i]
		yamlPath := filepath.Join("resources", fmt.Sprintf("%s-%s.yaml", resource.NukeService, resource.NukeResource))
		if _, err := os.Stat(yamlPath); err == nil {
			fmt.Fprintf(os.Stderr, "  SKIP (exists): %s\n", yamlPath)
			continue
		}

		fmt.Fprintf(os.Stderr, "\nGenerating: %s -> %s\n", resource.APIEntry, yamlPath)
		if err := generateAndIterate(ctx, client, resource, yamlPath, dslReference); err != nil {
			fmt.Fprintf(os.Stderr, "  FAILED: %v\n", err)
			continue
		}
		fmt.Fprintf(os.Stderr, "  SUCCESS: %s\n", yamlPath)
	}
}

func generateAndIterate(
	ctx context.Context, client *bedrockruntime.Client,
	resource *missingResource, yamlPath, dslReference string,
) error {
	sdkSource := buildSDKSourceContext(resource.GoPackage, resource.ListOp, resource.DeleteOp)
	yamlContent, err := generateYAML(ctx, client, resource, dslReference, sdkSource, "")
	if err != nil {
		return fmt.Errorf("initial generation: %w", err)
	}

	for attempt := 0; attempt < MaxRetries; attempt++ {
		if err := os.WriteFile(yamlPath, []byte(yamlContent), 0o600); err != nil {
			return fmt.Errorf("writing YAML: %w", err)
		}
		output, err := runCodegen(yamlPath)
		if err == nil {
			return nil
		}
		fmt.Fprintf(os.Stderr, "  attempt %d failed\n", attempt+1)
		cleanupGeneratedFiles(yamlPath)
		yamlContent, err = generateYAML(ctx, client, resource, dslReference, sdkSource, output)
		if err != nil {
			return fmt.Errorf("retry %d: %w", attempt+1, err)
		}
	}

	cleanupGeneratedFiles(yamlPath)
	os.Remove(yamlPath)
	return fmt.Errorf("failed after %d attempts", MaxRetries)
}

func generateYAML(
	ctx context.Context, client *bedrockruntime.Client,
	resource *missingResource, dslReference, sdkSource, previousError string,
) (string, error) {
	errorCtx := ""
	if previousError != "" {
		errorCtx = "\nPREVIOUS ATTEMPT FAILED:\n" + previousError + "\n"
	}

	prompt := fmt.Sprintf("Generate a YAML DSL file for aws-nuke codegen.\n\n"+
		"SDK package: %s\nDelete: %s\nList: %s\nService: %s\nResource: %s\nStruct: %s\nScope: %s\n\n"+
		"## SDK Source\n%s\n\n## DSL Reference\n%s\n%s\n"+
		"Respond with ONLY YAML. No fences, no explanation.",
		resource.GoPackage, resource.DeleteOp, resource.ListOp,
		resource.NukeService, resource.NukeResource, resource.ResourceName, resource.Scope,
		sdkSource, dslReference, errorCtx)

	response, err := invokeBedrockText(ctx, client, prompt)
	if err != nil {
		return "", err
	}
	return stripMarkdownFences(response), nil
}

func buildSDKSourceContext(goPackage, listOp, deleteOp string) string {
	modCache := getModCache()
	serviceDir := findSDKServiceDir(modCache, goPackage)
	if serviceDir == "" {
		return buildSDKSourceFromGitHub(goPackage, listOp, deleteOp)
	}

	var builder strings.Builder
	if listOp != "" {
		if content, err := readStructBlock(filepath.Join(serviceDir, fmt.Sprintf("api_op_%s.go", listOp)), listOp+"Output"); err == nil {
			fmt.Fprintf(&builder, "### %sOutput:\n%s\n\n", listOp, content)
		}
	}
	if deleteOp != "" {
		if content, err := readStructBlock(filepath.Join(serviceDir, fmt.Sprintf("api_op_%s.go", deleteOp)), deleteOp+"Input"); err == nil {
			fmt.Fprintf(&builder, "### %sInput:\n%s\n\n", deleteOp, content)
		}
	}

	output := builder.String()
	if len(output) > 30000 {
		output = output[:30000] + "\n..."
	}
	return output
}

func buildSDKSourceFromGitHub(goPackage, listOp, deleteOp string) string {
	var builder strings.Builder
	if listOp != "" {
		if content := fetchSDKOperationFile(goPackage, listOp); content != "" {
			fmt.Fprintf(&builder, "### %s.go:\n%s\n\n", listOp, extractStructFromContent(content, listOp+"Output"))
		}
	}
	if deleteOp != "" {
		if content := fetchSDKOperationFile(goPackage, deleteOp); content != "" {
			fmt.Fprintf(&builder, "### %s.go:\n%s\n\n", deleteOp, extractStructFromContent(content, deleteOp+"Input"))
		}
	}
	if builder.Len() == 0 {
		return "SDK source not available"
	}
	return builder.String()
}

func extractStructFromContent(content, structName string) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder
	inStruct := false
	braceDepth := 0
	for _, line := range lines {
		if !inStruct && strings.Contains(line, "type "+structName+" struct") {
			inStruct = true
			braceDepth = 0
		}
		if inStruct {
			result.WriteString(line)
			result.WriteString("\n")
			braceDepth += strings.Count(line, "{") - strings.Count(line, "}")
			if braceDepth <= 0 && strings.Contains(line, "}") {
				break
			}
		}
	}
	return result.String()
}

// --- Utilities ---

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func extractExistingResources() []string {
	var resources []string
	entries, err := os.ReadDir("resources")
	if err != nil {
		return resources
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.Contains(entry.Name(), "_test") {
			continue
		}
		content, err := os.ReadFile(filepath.Join("resources", entry.Name()))
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(content), "\n") {
			if strings.Contains(line, "Resource = \"") && strings.HasPrefix(strings.TrimSpace(line), "const ") {
				parts := strings.SplitN(line, "\"", 3)
				if len(parts) >= 2 {
					resources = append(resources, parts[1])
				}
			}
		}
	}
	return resources
}

func discoverSDKPackages() []string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return nil
	}
	var packages []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "aws-sdk-go-v2/service/") {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				segments := strings.Split(parts[0], "/")
				packages = append(packages, segments[len(segments)-1])
			}
		}
	}
	return packages
}

func loadDSLReference() string {
	data, err := os.ReadFile("docs/codegen.md")
	if err != nil {
		return ""
	}
	return string(data)
}

func getModCache() string {
	if v := os.Getenv("GOMODCACHE"); v != "" {
		return v
	}
	output, err := exec.CommandContext(context.Background(), "go", "env", "GOMODCACHE").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func findSDKServiceDir(modCache, goPackage string) string {
	serviceBase := filepath.Join(modCache, "github.com", "aws", "aws-sdk-go-v2", "service")
	entries, err := os.ReadDir(serviceBase)
	if err != nil {
		return ""
	}
	var best string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), goPackage+"@") {
			best = filepath.Join(serviceBase, entry.Name())
		}
	}
	return best
}

func readStructBlock(filePath, structName string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(content), "\n")
	var result strings.Builder
	inStruct := false
	braceDepth := 0
	for _, line := range lines {
		if !inStruct && strings.Contains(line, "type "+structName+" struct") {
			inStruct = true
			braceDepth = 0
		}
		if inStruct {
			result.WriteString(line)
			result.WriteString("\n")
			braceDepth += strings.Count(line, "{") - strings.Count(line, "}")
			if braceDepth <= 0 && strings.Contains(line, "}") {
				break
			}
		}
	}
	if result.Len() == 0 {
		return "", fmt.Errorf("struct %s not found", structName)
	}
	return result.String(), nil
}

func runCodegen(yamlPath string) (string, error) {
	cmd := exec.CommandContext(context.Background(), "go", "run", "./tools/codegen/", "--no-commit", yamlPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return "", nil
}

func cleanupGeneratedFiles(yamlPath string) {
	base := strings.TrimSuffix(yamlPath, ".yaml")
	os.Remove(base + ".go")
	os.Remove(base + "_mock_test.go")
}

func newBedrockClient(ctx context.Context) (*bedrockruntime.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(BedrockRegion))
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}
	return bedrockruntime.NewFromConfig(cfg), nil
}

func invokeBedrockText(ctx context.Context, client *bedrockruntime.Client, prompt string) (string, error) {
	type message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	type request struct {
		AnthropicVersion string    `json:"anthropic_version"`
		MaxTokens        int       `json:"max_tokens"`
		Temperature      float64   `json:"temperature"`
		Messages         []message `json:"messages"`
	}

	body := request{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        65536,
		Temperature:      0,
		Messages:         []message{{Role: "user", Content: prompt}},
	}

	reqJSON, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	modelID := ModelOpus
	contentType := "application/json"

	resp, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     &modelID,
		ContentType: &contentType,
		Body:        reqJSON,
	})
	if err != nil {
		return "", fmt.Errorf("invoking bedrock: %w", err)
	}

	type contentBlock struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	type bedrockResponse struct {
		Content []contentBlock `json:"content"`
	}

	var bedrockResp bedrockResponse
	if err := json.Unmarshal(resp.Body, &bedrockResp); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}
	if len(bedrockResp.Content) == 0 {
		return "", fmt.Errorf("empty response")
	}
	return bedrockResp.Content[0].Text, nil
}

func stripMarkdownFences(response string) string {
	response = strings.TrimSpace(response)
	if !strings.HasPrefix(response, "```") {
		return response
	}
	lines := strings.Split(response, "\n")
	endIdx := len(lines)
	for i := len(lines) - 1; i >= 1; i-- {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "```") {
			endIdx = i
			break
		}
	}
	return strings.TrimSpace(strings.Join(lines[1:endIdx], "\n"))
}

// --- Report ---

func writeTSV(path string, resources []missingResource) error {
	var b strings.Builder
	b.WriteString("APIEntry\tGoPackage\tDeleteOp\tListOp\tNukeService\tNukeResource\tResourceName\tScope\n")
	for i := range resources {
		r := &resources[i]
		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			r.APIEntry, r.GoPackage, r.DeleteOp, r.ListOp,
			r.NukeService, r.NukeResource, r.ResourceName, r.Scope)
	}
	return os.WriteFile(path, []byte(b.String()), 0o600)
}

func readTSV(path string) ([]missingResource, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var resources []missingResource
	scanner := bufio.NewScanner(file)

	// Skip header
	if scanner.Scan() {
		header := scanner.Text()
		if !strings.HasPrefix(header, "APIEntry\t") {
			return nil, fmt.Errorf("unexpected header: %s", header)
		}
	}

	lineNum := 1
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 8 {
			return nil, fmt.Errorf("line %d: expected 8 fields, got %d", lineNum, len(fields))
		}
		resources = append(resources, missingResource{
			APIEntry:     fields[0],
			GoPackage:    fields[1],
			DeleteOp:     fields[2],
			ListOp:       fields[3],
			NukeService:  fields[4],
			NukeResource: fields[5],
			ResourceName: fields[6],
			Scope:        fields[7],
		})
	}
	return resources, scanner.Err()
}
