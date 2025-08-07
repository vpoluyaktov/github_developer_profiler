package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"

	"dev_profiler/internal/dto"
)

// GitHubService handles GitHub API interactions
type GitHubService struct {
	client *github.Client
	config *dto.GitHubConfig
}

// NewGitHubService creates a new GitHub service
func NewGitHubService(config *dto.GitHubConfig) *GitHubService {
	var client *github.Client
	
	if config.Token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.Token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	return &GitHubService{
		client: client,
		config: config,
	}
}

// GetUser retrieves GitHub user information
func (s *GitHubService) GetUser(ctx context.Context, username string) (*dto.UserInfo, error) {
	user, _, err := s.client.Users.Get(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", username, err)
	}

	userInfo := &dto.UserInfo{
		Username:    user.GetLogin(),
		Name:        user.GetName(),
		Company:     user.GetCompany(),
		Location:    user.GetLocation(),
		Email:       user.GetEmail(),
		PublicRepos: user.GetPublicRepos(),
		PublicGists: user.GetPublicGists(),
		Followers:   user.GetFollowers(),
		Following:   user.GetFollowing(),
	}

	if user.CreatedAt != nil {
		userInfo.CreatedAt = user.CreatedAt.Time
	}
	if user.UpdatedAt != nil {
		userInfo.UpdatedAt = user.UpdatedAt.Time
	}
	if user.Plan != nil {
		userInfo.SubscriptionPlan = user.Plan.GetName()
	}

	return userInfo, nil
}

// ListRepositories retrieves user repositories
func (s *GitHubService) ListRepositories(ctx context.Context, username string) ([]*dto.Repository, error) {
	opt := &github.RepositoryListOptions{
		Type:        "all",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := s.client.Repositories.List(ctx, username, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %w", err)
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var repositories []*dto.Repository
	for _, repo := range allRepos {
		repository := &dto.Repository{
			Name:        repo.GetName(),
			Description: repo.GetDescription(),
			Stars:       repo.GetStargazersCount(),
			Fork:        repo.GetFork(),
		}

		if repo.CreatedAt != nil {
			repository.CreatedAt = repo.CreatedAt.Time
		}
		if repo.UpdatedAt != nil {
			repository.UpdatedAt = repo.UpdatedAt.Time
		}

		if repo.GetFork() && repo.Parent != nil {
			repository.ForkSource = repo.Parent.GetFullName()
		}

		repositories = append(repositories, repository)
	}

	return repositories, nil
}

// GetRepositoryLanguages retrieves languages used in a repository
func (s *GitHubService) GetRepositoryLanguages(ctx context.Context, username, repoName string) ([]string, error) {
	languages, _, err := s.client.Repositories.ListLanguages(ctx, username, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to get languages for %s/%s: %w", username, repoName, err)
	}

	var langs []string
	for lang := range languages {
		langs = append(langs, lang)
	}

	// Sort by usage (GitHub API returns map with byte counts)
	sort.Slice(langs, func(i, j int) bool {
		return languages[langs[i]] > languages[langs[j]]
	})

	return langs, nil
}

// GetRepositoryCommits retrieves recent commits for a repository
func (s *GitHubService) GetRepositoryCommits(ctx context.Context, username, repoName string, count int) ([]*dto.CommitDetail, error) {
	opt := &github.CommitsListOptions{
		Author:      username,
		ListOptions: github.ListOptions{PerPage: count},
	}

	commits, _, err := s.client.Repositories.ListCommits(ctx, username, repoName, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits for %s/%s: %w", username, repoName, err)
	}

	var commitDetails []*dto.CommitDetail
	for _, commit := range commits {
		detail := &dto.CommitDetail{
			Repo:    repoName,
			SHA:     commit.GetSHA(),
			Message: commit.GetCommit().GetMessage(),
			Author:  commit.GetCommit().GetAuthor().GetName(),
		}

		if commit.GetCommit().GetAuthor().Date != nil {
			detail.Date = commit.GetCommit().GetAuthor().Date.Time
		}

		// Get files changed in this commit
		commitDetail, _, err := s.client.Repositories.GetCommit(ctx, username, repoName, commit.GetSHA(), nil)
		if err == nil && commitDetail.Files != nil {
			for _, file := range commitDetail.Files {
				detail.FilesChanged = append(detail.FilesChanged, file.GetFilename())
			}
		}

		commitDetails = append(commitDetails, detail)
	}

	return commitDetails, nil
}

// GetRepositoryContents retrieves repository file contents for analysis
func (s *GitHubService) GetRepositoryContents(ctx context.Context, username, repoName string) ([]*dto.FileAnalysis, error) {
	// Get repository tree
	tree, _, err := s.client.Git.GetTree(ctx, username, repoName, "HEAD", true)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository tree: %w", err)
	}

	// Filter and sample files
	var codeFiles []*github.TreeEntry
	for _, entry := range tree.Entries {
		if entry.GetType() == "blob" && s.isCodeFile(entry.GetPath()) {
			codeFiles = append(codeFiles, entry)
		}
	}

	// Sample files if there are too many
	if len(codeFiles) > s.config.SampleFileCount {
		rand.Seed(int64(s.config.RandomSeed))
		rand.Shuffle(len(codeFiles), func(i, j int) {
			codeFiles[i], codeFiles[j] = codeFiles[j], codeFiles[i]
		})
		codeFiles = codeFiles[:s.config.SampleFileCount]
	}

	var fileAnalyses []*dto.FileAnalysis
	for _, file := range codeFiles {
		analysis := s.analyzeFile(ctx, username, repoName, file)
		if analysis != nil {
			fileAnalyses = append(fileAnalyses, analysis)
		}
	}

	return fileAnalyses, nil
}

// isCodeFile determines if a file is a code file worth analyzing
func (s *GitHubService) isCodeFile(path string) bool {
	codeExtensions := map[string]bool{
		".go":   true, ".py":   true, ".js":   true, ".ts":   true,
		".java": true, ".cpp":  true, ".c":    true, ".cs":   true,
		".rb":   true, ".php":  true, ".rs":   true, ".kt":   true,
		".swift": true, ".scala": true, ".r":    true, ".m":    true,
		".h":    true, ".hpp":  true, ".cc":   true, ".cxx":  true,
	}

	// Get file extension
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return false
	}
	ext := "." + strings.ToLower(parts[len(parts)-1])

	return codeExtensions[ext]
}

// analyzeFile analyzes a single file
func (s *GitHubService) analyzeFile(ctx context.Context, username, repoName string, file *github.TreeEntry) *dto.FileAnalysis {
	// Get file content
	content, _, _, err := s.client.Repositories.GetContents(ctx, username, repoName, file.GetPath(), nil)
	if err != nil {
		return nil
	}

	if content.Content == nil {
		return nil
	}

	// Decode content
	decoded, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		return nil
	}

	fileContent := string(decoded)
	lines := strings.Split(fileContent, "\n")

	analysis := &dto.FileAnalysis{
		Repo:     repoName,
		Path:     file.GetPath(),
		Language: s.detectLanguage(file.GetPath()),
		Size:     len(decoded),
		Lines:    len(lines),
		HasTests: s.hasTestIndicators(fileContent),
		Content:  fileContent,
	}

	return analysis
}

// detectLanguage detects programming language from file extension
func (s *GitHubService) detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	languages := map[string]string{
		".go":    "Go",
		".py":    "Python",
		".js":    "JavaScript",
		".ts":    "TypeScript",
		".java":  "Java",
		".cpp":   "C++",
		".c":     "C",
		".cs":    "C#",
		".rb":    "Ruby",
		".php":   "PHP",
		".rs":    "Rust",
		".kt":    "Kotlin",
		".swift": "Swift",
		".scala": "Scala",
		".r":     "R",
		".m":     "Objective-C",
		".h":     "C/C++ Header",
		".hpp":   "C++ Header",
	}

	if lang, exists := languages[ext]; exists {
		return lang
	}
	return "Unknown"
}

// hasTestIndicators checks if file contains test-related code
func (s *GitHubService) hasTestIndicators(content string) bool {
	testIndicators := []string{
		"test", "Test", "TEST",
		"assert", "Assert",
		"expect", "Expect",
		"should", "Should",
		"describe", "it(",
		"@Test", "unittest",
		"testing.T", "t.Run",
	}

	contentLower := strings.ToLower(content)
	for _, indicator := range testIndicators {
		if strings.Contains(contentLower, strings.ToLower(indicator)) {
			return true
		}
	}
	return false
}

// assessCodeQuality provides a simple code quality assessment
func (s *GitHubService) assessCodeQuality(content string) string {
	lines := strings.Split(content, "\n")
	
	// Simple metrics
	commentLines := 0
	emptyLines := 0
	longLines := 0
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			emptyLines++
		} else if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "/*") {
			commentLines++
		}
		if len(line) > 120 {
			longLines++
		}
	}
	
	totalLines := len(lines)
	if totalLines == 0 {
		return "Unknown"
	}
	
	commentRatio := float64(commentLines) / float64(totalLines)
	longLineRatio := float64(longLines) / float64(totalLines)
	
	if commentRatio > 0.2 && longLineRatio < 0.1 {
		return "Good"
	} else if commentRatio > 0.1 && longLineRatio < 0.2 {
		return "Fair"
	}
	return "Needs Improvement"
}

// assessComplexity provides a simple complexity assessment
func (s *GitHubService) assessComplexity(content string) string {
	// Count control flow statements
	complexityKeywords := []string{
		"if", "else", "for", "while", "switch", "case",
		"try", "catch", "finally", "throw",
	}
	
	contentLower := strings.ToLower(content)
	complexityCount := 0
	
	for _, keyword := range complexityKeywords {
		complexityCount += strings.Count(contentLower, keyword)
	}
	
	lines := len(strings.Split(content, "\n"))
	if lines == 0 {
		return "Unknown"
	}
	
	complexityRatio := float64(complexityCount) / float64(lines)
	
	if complexityRatio > 0.1 {
		return "High"
	} else if complexityRatio > 0.05 {
		return "Medium"
	}
	return "Low"
}

// PerformFullAudit performs a comprehensive GitHub user audit
func (s *GitHubService) PerformFullAudit(ctx context.Context, username string) (*dto.AuditResult, error) {
	// Get user information
	userInfo, err := s.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Get repositories
	repositories, err := s.ListRepositories(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories: %w", err)
	}

	// Filter repositories for analysis (recent activity)
	cutoffDate := time.Now().AddDate(-s.config.AnalysisYears, 0, 0)
	var analysisRepos []*dto.Repository
	var originalRepos []*dto.Repository
	var forkedRepos []*dto.Repository
	
	for _, repo := range repositories {
		if repo.Fork {
			forkedRepos = append(forkedRepos, repo)
		} else {
			originalRepos = append(originalRepos, repo)
		}
		
		if repo.UpdatedAt.After(cutoffDate) {
			analysisRepos = append(analysisRepos, repo)
		}
	}

	// Sample repositories for detailed analysis - select N most recently updated
	if len(analysisRepos) > s.config.SampledRepoCount {
		// Sort repositories by UpdatedAt in descending order (most recent first)
		sort.Slice(analysisRepos, func(i, j int) bool {
			return analysisRepos[i].UpdatedAt.After(analysisRepos[j].UpdatedAt)
		})
		// Take the N most recently updated repositories
		analysisRepos = analysisRepos[:s.config.SampledRepoCount]
	}

	// Analyze repositories in detail
	var allCommits []*dto.CommitDetail
	var allFileAnalyses []*dto.FileAnalysis
	var analyzedRepos []string
	var reposWithErrors []string

	for _, repo := range analysisRepos {
		// Get repository languages
		languages, err := s.GetRepositoryLanguages(ctx, username, repo.Name)
		if err == nil {
			repo.LanguagesUsed = languages
		}

		// Get commits
		commits, err := s.GetRepositoryCommits(ctx, username, repo.Name, s.config.CommitsPerRepo)
		if err != nil {
			reposWithErrors = append(reposWithErrors, repo.Name)
			continue
		}
		
		if len(commits) == 0 {
			continue // Skip repos without user commits
		}

		allCommits = append(allCommits, commits...)
		repo.CommitCount = len(commits)

		// Get file analysis
		fileAnalyses, err := s.GetRepositoryContents(ctx, username, repo.Name)
		if err != nil {
			reposWithErrors = append(reposWithErrors, repo.Name)
			continue
		}

		allFileAnalyses = append(allFileAnalyses, fileAnalyses...)
		repo.FileCount = len(fileAnalyses)
		repo.IncludeAnalysis = true
		analyzedRepos = append(analyzedRepos, repo.Name)
	}

	// Calculate statistics
	totalStars := 0
	significantForks := 0
	for _, repo := range repositories {
		totalStars += repo.Stars
		if repo.Fork && repo.UserCommits > 10 {
			significantForks++
			repo.IsSignificant = true
		}
	}

	// Build result
	result := &dto.AuditResult{
		UserInfo: *userInfo,
		RepoStats: dto.RepoStats{
			Statistics: dto.RepoStatistics{
				TotalRepos:       len(repositories),
				OriginalRepos:    len(originalRepos),
				ForkedRepos:      len(forkedRepos),
				SignificantForks: significantForks,
				TotalStars:       totalStars,
				AnalysisRepos:    len(analysisRepos),
			},
			OriginalRepos: originalRepos,
			ForkedRepos:   forkedRepos,
		},
		FileAnalysis:  allFileAnalyses,
		CommitDetails: allCommits,
		AuditParameters: s.getSanitizedConfig(),
		AnalysisSummary: dto.AnalysisSummary{
			ReposAnalyzedForCode:     analyzedRepos,
			ReposWithErrors:          reposWithErrors,
			TotalReposAttempted:      len(analysisRepos),
			SuccessfulAnalysisCount:  len(analyzedRepos),
		},
	}

	return result, nil
}

// getSanitizedConfig returns a copy of the config with sensitive information removed
func (s *GitHubService) getSanitizedConfig() dto.GitHubConfig {
	// Create a copy of the config without the token
	sanitizedConfig := *s.config
	sanitizedConfig.Token = "[REDACTED]"
	return sanitizedConfig
}
