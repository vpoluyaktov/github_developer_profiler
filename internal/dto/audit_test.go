package dto

import (
	"testing"
	"time"

	"dev_profiler/internal/config"
)

func TestUserInfo(t *testing.T) {
	userInfo := UserInfo{
		Username:         "testuser",
		Name:             "Test User",
		Company:          "Test Company",
		Location:         "Test Location",
		Email:            "test@example.com",
		PublicRepos:      10,
		PublicGists:      5,
		Followers:        100,
		Following:        50,
		CreatedAt:        time.Now().AddDate(-2, 0, 0),
		UpdatedAt:        time.Now(),
		SubscriptionPlan: "free",
	}
	
	// Test basic field assignments
	if userInfo.Username != "testuser" {
		t.Error("Username field not set correctly")
	}
	
	if userInfo.PublicRepos != 10 {
		t.Error("PublicRepos field not set correctly")
	}
	
	if userInfo.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero time")
	}
}

func TestRepository(t *testing.T) {
	repo := Repository{
		Name:            "test-repo",
		Description:     "Test repository",
		Stars:           25,
		Fork:            false,
		ForkSource:      "",
		CreatedAt:       time.Now().AddDate(-1, 0, 0),
		UpdatedAt:       time.Now(),
		LanguagesUsed:   []string{"Go", "JavaScript"},
		CommitCount:     50,
		FileCount:       20,
		UserCommits:     30,
		IsSignificant:   true,
		IncludeAnalysis: true,
	}
	
	// Test basic field assignments
	if repo.Name != "test-repo" {
		t.Error("Name field not set correctly")
	}
	
	if repo.Stars != 25 {
		t.Error("Stars field not set correctly")
	}
	
	if len(repo.LanguagesUsed) != 2 {
		t.Error("LanguagesUsed should have 2 elements")
	}
	
	if repo.LanguagesUsed[0] != "Go" || repo.LanguagesUsed[1] != "JavaScript" {
		t.Error("LanguagesUsed not set correctly")
	}
}

func TestFileAnalysis(t *testing.T) {
	fileAnalysis := FileAnalysis{
		Repo:        "test-repo",
		Path:        "src/main.go",
		Language:    "Go",
		Size:        1024,
		Lines:       100,
		Content:     "package main\n\nfunc main() {}",
		HasTests:    true,
	}
	
	// Test basic field assignments
	if fileAnalysis.Repo != "test-repo" {
		t.Error("Repo field not set correctly")
	}
	
	if fileAnalysis.Language != "Go" {
		t.Error("Language field not set correctly")
	}
	
	if fileAnalysis.Lines != 100 {
		t.Error("Lines field not set correctly")
	}
	
	if !fileAnalysis.HasTests {
		t.Error("HasTests field not set correctly")
	}
}

func TestCommitDetail(t *testing.T) {
	commitDetail := CommitDetail{
		Repo:      "test-repo",
		SHA:       "abc123def456",
		Message:   "Initial commit",
		Author:    "Test Author",
		Date:      time.Now(),
		FilesChanged: []string{"main.go", "README.md"},
	}
	
	// Test basic field assignments
	if commitDetail.Repo != "test-repo" {
		t.Error("Repo field not set correctly")
	}
	
	if commitDetail.SHA != "abc123def456" {
		t.Error("SHA field not set correctly")
	}
	
	if len(commitDetail.FilesChanged) != 2 {
		t.Error("FilesChanged field not set correctly")
	}
	
	if commitDetail.FilesChanged[0] != "main.go" {
		t.Error("First file in FilesChanged not set correctly")
	}
}

func TestRepoStatistics(t *testing.T) {
	stats := RepoStatistics{
		TotalRepos:       50,
		OriginalRepos:    30,
		ForkedRepos:      20,
		SignificantForks: 5,
		TotalStars:       1000,
		AnalysisRepos:    10,
	}
	
	// Test basic field assignments
	if stats.TotalRepos != 50 {
		t.Error("TotalRepos field not set correctly")
	}
	
	if stats.OriginalRepos != 30 {
		t.Error("OriginalRepos field not set correctly")
	}
	
	// Test logical consistency
	if stats.OriginalRepos + stats.ForkedRepos != stats.TotalRepos {
		t.Error("OriginalRepos + ForkedRepos should equal TotalRepos")
	}
	
	if stats.SignificantForks > stats.ForkedRepos {
		t.Error("SignificantForks cannot exceed ForkedRepos")
	}
}

func TestRepoStats(t *testing.T) {
	originalRepo := &Repository{
		Name:  "original-repo",
		Fork:  false,
		Stars: 10,
	}
	
	forkedRepo := &Repository{
		Name:  "forked-repo",
		Fork:  true,
		Stars: 5,
	}
	
	stats := RepoStats{
		Statistics: RepoStatistics{
			TotalRepos:    2,
			OriginalRepos: 1,
			ForkedRepos:   1,
			TotalStars:    15,
		},
		OriginalRepos: []*Repository{originalRepo},
		ForkedRepos:   []*Repository{forkedRepo},
	}
	
	// Test structure
	if len(stats.OriginalRepos) != 1 {
		t.Error("Should have 1 original repo")
	}
	
	if len(stats.ForkedRepos) != 1 {
		t.Error("Should have 1 forked repo")
	}
	
	if stats.OriginalRepos[0].Fork {
		t.Error("Original repo should not be marked as fork")
	}
	
	if !stats.ForkedRepos[0].Fork {
		t.Error("Forked repo should be marked as fork")
	}
}

func TestAnalysisSummary(t *testing.T) {
	summary := AnalysisSummary{
		ReposAnalyzedForCode:     []string{"repo1", "repo2", "repo3"},
		ReposWithErrors:          []string{"error-repo"},
		TotalReposAttempted:      4,
		SuccessfulAnalysisCount:  3,
	}
	
	// Test basic field assignments
	if len(summary.ReposAnalyzedForCode) != 3 {
		t.Error("Should have 3 successfully analyzed repos")
	}
	
	if len(summary.ReposWithErrors) != 1 {
		t.Error("Should have 1 repo with errors")
	}
	
	if summary.TotalReposAttempted != 4 {
		t.Error("TotalReposAttempted should be 4")
	}
	
	if summary.SuccessfulAnalysisCount != 3 {
		t.Error("SuccessfulAnalysisCount should be 3")
	}
	
	// Test logical consistency
	if summary.SuccessfulAnalysisCount != len(summary.ReposAnalyzedForCode) {
		t.Error("SuccessfulAnalysisCount should match length of ReposAnalyzedForCode")
	}
	
	expectedTotal := len(summary.ReposAnalyzedForCode) + len(summary.ReposWithErrors)
	if summary.TotalReposAttempted != expectedTotal {
		t.Error("TotalReposAttempted should equal successful + error repos")
	}
}

func TestAuditResult(t *testing.T) {
	userInfo := UserInfo{
		Username:    "testuser",
		PublicRepos: 10,
	}
	
	repoStats := RepoStats{
		Statistics: RepoStatistics{
			TotalRepos: 10,
		},
	}
	
	auditResult := AuditResult{
		UserInfo:        userInfo,
		RepoStats:       repoStats,
		FileAnalysis:    []*FileAnalysis{},
		CommitDetails:   []*CommitDetail{},
		AuditParameters: config.GitHubConfig{},
		AnalysisSummary: AnalysisSummary{},
	}
	
	// Test structure initialization
	if auditResult.UserInfo.Username != "testuser" {
		t.Error("UserInfo not set correctly")
	}
	
	if auditResult.RepoStats.Statistics.TotalRepos != 10 {
		t.Error("RepoStats not set correctly")
	}
	
	if auditResult.FileAnalysis == nil {
		t.Error("FileAnalysis should be initialized")
	}
	
	if auditResult.CommitDetails == nil {
		t.Error("CommitDetails should be initialized")
	}
}
