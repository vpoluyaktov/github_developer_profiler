package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"dev_profiler/internal/config"
	"dev_profiler/internal/utils"
)

func TestGitHubServiceRepositoryAnalysis(t *testing.T) {
	// Skip integration test in CI/CD environments
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping integration test in CI/CD environment")
	}

	// Load configuration from user's home directory
	cfg, err := loadConfigFromHome()
	if err != nil {
		t.Skipf("Skipping test - config not available: %v", err)
		return
	}

	if cfg.GitHub.Token == "" {
		t.Skip("GitHub token not configured - skipping integration test")
	}

	// Create GitHub service
	githubService := NewGitHubService(cfg.GitHub)

	// Test username for analysis
	username := "vpoluyaktov"

	t.Run("GetUser", func(t *testing.T) {
		ctx := context.Background()
		userInfo, err := githubService.GetUser(ctx, username)
		if err != nil {
			t.Fatalf("Failed to get user info: %v", err)
		}

		fmt.Printf("User Info:\n")
		fmt.Printf("  Username: %s\n", userInfo.Username)
		fmt.Printf("  Name: %s\n", userInfo.Name)
		fmt.Printf("  Public Repos: %d\n", userInfo.PublicRepos)
		fmt.Printf("  Created: %s\n", userInfo.CreatedAt.Format("2006-01-02"))

		if userInfo.Username != username {
			t.Errorf("Expected username %s, got %s", username, userInfo.Username)
		}
	})

	t.Run("ListRepositories", func(t *testing.T) {
		ctx := context.Background()
		repos, err := githubService.ListRepositories(ctx, username)
		if err != nil {
			t.Fatalf("Failed to list repositories: %v", err)
		}

		fmt.Printf("\nRepositories (%d total):\n", len(repos))
		
		originalCount := 0
		forkedCount := 0
		recentCount := 0
		cutoffDate := time.Now().AddDate(-cfg.GitHub.AnalysisYears, 0, 0)

		for _, repo := range repos {
			fmt.Printf("  %s (fork: %t, updated: %s, stars: %d)\n", 
				repo.Name, repo.Fork, repo.UpdatedAt.Format("2006-01-02"), repo.Stars)
			
			if repo.Fork {
				forkedCount++
			} else {
				originalCount++
			}
			
			if repo.UpdatedAt.After(cutoffDate) {
				recentCount++
			}
		}

		fmt.Printf("\nRepository Statistics:\n")
		fmt.Printf("  Original repos: %d\n", originalCount)
		fmt.Printf("  Forked repos: %d\n", forkedCount)
		fmt.Printf("  Recent repos (last %d years): %d\n", cfg.GitHub.AnalysisYears, recentCount)

		if len(repos) == 0 {
			t.Error("Expected at least one repository")
		}
	})

	t.Run("AnalyzeGitHubUser", func(t *testing.T) {
		ctx := context.Background()
		
		fmt.Printf("\nRunning full GitHub user analysis...\n")
		fmt.Printf("Configuration:\n")
		fmt.Printf("  Analysis Years: %d\n", cfg.GitHub.AnalysisYears)
		fmt.Printf("  Sampled Repo Count: %d\n", cfg.GitHub.SampledRepoCount)
		fmt.Printf("  Sample File Count: %d\n", cfg.GitHub.SampleFileCount)
		fmt.Printf("  Commits Per Repo: %d\n", cfg.GitHub.CommitsPerRepo)

		auditResult, err := githubService.PerformFullAudit(ctx, username)
		if err != nil {
			t.Fatalf("Failed to analyze GitHub user: %v", err)
		}

		fmt.Printf("\nAnalysis Results:\n")
		fmt.Printf("  Total repos: %d\n", auditResult.RepoStats.Statistics.TotalRepos)
		fmt.Printf("  Original repos: %d\n", auditResult.RepoStats.Statistics.OriginalRepos)
		fmt.Printf("  Forked repos: %d\n", auditResult.RepoStats.Statistics.ForkedRepos)
		fmt.Printf("  Analysis repos: %d\n", auditResult.RepoStats.Statistics.AnalysisRepos)
		fmt.Printf("  Total stars: %d\n", auditResult.RepoStats.Statistics.TotalStars)

		fmt.Printf("\nFile Analysis:\n")
		fmt.Printf("  Total files analyzed: %d\n", len(auditResult.FileAnalysis))
		
		repoFileCount := make(map[string]int)
		for _, file := range auditResult.FileAnalysis {
			repoFileCount[file.Repo]++
			fmt.Printf("    %s: %s (%s, %d lines)\n", 
				file.Repo, file.Path, file.Language, file.Lines)
		}

		fmt.Printf("\nFiles per repository:\n")
		for repo, count := range repoFileCount {
			fmt.Printf("  %s: %d files\n", repo, count)
		}

		fmt.Printf("\nCommit Analysis:\n")
		fmt.Printf("  Total commits analyzed: %d\n", len(auditResult.CommitDetails))

		// Verify we have some analysis data
		if auditResult.RepoStats.Statistics.TotalRepos == 0 {
			t.Error("Expected at least one repository in analysis")
		}

		if len(auditResult.FileAnalysis) == 0 {
			t.Error("Expected at least one file analysis")
		}

		// Check that file analysis includes content
		hasContent := false
		for _, file := range auditResult.FileAnalysis {
			if file.Content != "" {
				hasContent = true
				fmt.Printf("    Content sample from %s: %d characters\n", 
					file.Path, len(file.Content))
				break
			}
		}

		if !hasContent {
			t.Error("Expected at least one file to have content")
		}

		// Save audit result for inspection
		auditJSON, err := json.MarshalIndent(auditResult, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal audit result: %v", err)
		}

		testOutputFile := "/tmp/github_audit_test_result.json"
		err = os.WriteFile(testOutputFile, auditJSON, 0644)
		if err != nil {
			t.Fatalf("Failed to write test output: %v", err)
		}

		fmt.Printf("\nTest audit result saved to: %s\n", testOutputFile)
		fmt.Printf("Audit JSON size: %d bytes\n", len(auditJSON))
	})
}

// loadConfigFromHome loads the configuration from the user's home directory
func loadConfigFromHome() (*config.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".config", "dev_profiler", "config.json")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at %s", configPath)
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg config.Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Decrypt the GitHub token
	if cfg.GitHub.Token != "" {
		decryptedToken, err := decryptToken(cfg.GitHub.Token)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt GitHub token: %w", err)
		}
		cfg.GitHub.Token = decryptedToken
	}

	return &cfg, nil
}

// decryptToken decrypts the stored GitHub token
func decryptToken(encryptedToken string) (string, error) {
	// Decode from base64 first
	ciphertext, err := utils.DecodeBase64(encryptedToken)
	if err != nil {
		// If base64 decode fails, assume it's plain text
		return encryptedToken, nil
	}
	
	// Decrypt using the utils package
	decryptedToken, err := utils.DecryptString(ciphertext)
	if err != nil {
		// If decryption fails, assume it's plain text
		return encryptedToken, nil
	}
	
	return decryptedToken, nil
}
