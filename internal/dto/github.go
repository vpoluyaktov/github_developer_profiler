package dto

import "time"

// GitHubConfig holds GitHub API configuration
type GitHubConfig struct {
	Token              string `json:"token"`
	SampledRepoCount   int    `json:"sampled_repo_count"`
	CommitsPerRepo     int    `json:"commits_per_repo"`
	SampleFileCount    int    `json:"sample_file_count"`
	AnalysisYears      int    `json:"analysis_years"`
	IncludePrivateRepo bool   `json:"include_private_repos"`
	RandomSeed         int    `json:"random_seed"`
	SaveDebugJSON      bool   `json:"save_debug_json"`
}

// OpenAIConfig holds OpenAI API configuration
type OpenAIConfig struct {
	APIKey string `json:"api_key"`
	Model  string `json:"model"`
}

// DefaultGitHubConfig returns default configuration
func DefaultGitHubConfig() *GitHubConfig {
	return &GitHubConfig{
		Token:              "",
		SampledRepoCount:   10,
		CommitsPerRepo:     50,
		SampleFileCount:    3,
		AnalysisYears:      5,
		IncludePrivateRepo: false,
		RandomSeed:         42,
		SaveDebugJSON:      false,
	}
}

// DefaultOpenAIConfig returns default OpenAI configuration
func DefaultOpenAIConfig() *OpenAIConfig {
	return &OpenAIConfig{
		APIKey: "",
		Model:  "gpt-4o",
	}
}

// UserInfo represents GitHub user information
type UserInfo struct {
	Username         string    `json:"username"`
	Name             string    `json:"name"`
	Company          string    `json:"company"`
	Location         string    `json:"location"`
	Email            string    `json:"email"`
	CreatedAt        time.Time `json:"created_at"`
	PublicRepos      int       `json:"public_repos"`
	PublicGists      int       `json:"public_gists"`
	Followers        int       `json:"followers"`
	Following        int       `json:"following"`
	SubscriptionPlan string    `json:"subscription_plan"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Repository represents a GitHub repository
type Repository struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Stars           int               `json:"stars"`
	Fork            bool              `json:"fork"`
	ForkSource      string            `json:"fork_source,omitempty"`
	UserCommits     int               `json:"user_commits,omitempty"`
	LanguagesUsed   []string          `json:"languages_used"`
	FileCount       int               `json:"file_count"`
	CommitCount     int               `json:"commit_count"`
	IncludeAnalysis bool              `json:"include_in_analysis"`
	IsSignificant   bool              `json:"is_significant,omitempty"`
}

// RepoStatistics holds repository statistics
type RepoStatistics struct {
	TotalRepos      int `json:"total_repos"`
	OriginalRepos   int `json:"original_repos"`
	ForkedRepos     int `json:"forked_repos"`
	SignificantForks int `json:"significant_forks"`
	TotalStars      int `json:"total_stars"`
	AnalysisRepos   int `json:"analysis_repos"`
}

// CommitDetail represents commit information
type CommitDetail struct {
	Repo      string    `json:"repo"`
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	Date      time.Time `json:"date"`
	Author    string    `json:"author"`
	FilesChanged []string `json:"files_changed"`
}

// FileAnalysis represents code file analysis
type FileAnalysis struct {
	Repo         string `json:"repo"`
	Path         string `json:"path"`
	Language     string `json:"language"`
	Size         int    `json:"size"`
	Lines        int    `json:"lines"`
	HasTests     bool   `json:"has_tests"`
	Content      string `json:"content"`
}

// AuditResult represents the complete audit result
type AuditResult struct {
	UserInfo        UserInfo         `json:"user_info"`
	RepoStats       RepoStats        `json:"repo_stats"`
	FileAnalysis    []*FileAnalysis  `json:"file_analysis"`
	CommitDetails   []*CommitDetail  `json:"commit_details"`
	AuditParameters GitHubConfig     `json:"audit_parameters"`
	AnalysisSummary AnalysisSummary  `json:"analysis_summary"`
}

// RepoStats holds repository statistics and lists
type RepoStats struct {
	Statistics    RepoStatistics `json:"statistics"`
	OriginalRepos []*Repository  `json:"original_repos"`
	ForkedRepos   []*Repository  `json:"forked_repos"`
}

// AnalysisSummary holds analysis summary information
type AnalysisSummary struct {
	ReposAnalyzedForCode       []string `json:"repos_analyzed_for_code"`
	ReposWithErrors           []string `json:"repos_with_errors"`
	ReposWithoutUserCommits   []string `json:"repos_without_user_commits"`
	TotalReposAttempted       int      `json:"total_repos_attempted"`
	SuccessfulAnalysisCount   int      `json:"successful_analysis_count"`
}
