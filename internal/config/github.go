package config

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
