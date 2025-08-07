package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"

	"dev_profiler/internal/config"
	"dev_profiler/internal/dto"
	"dev_profiler/internal/services"
	"dev_profiler/internal/ui"
	"dev_profiler/internal/utils"
)

// MainController handles the main window business logic
type MainController struct {
	app           fyne.App
	window        fyne.Window
	ui            *ui.MainWindowUI
	configUI      *ui.ConfigWindowUI
	githubService *services.GitHubService
	openaiService *services.OpenAIService
	config        *config.Config
}

// NewMainController creates a new MainController instance
func NewMainController() *MainController {
	return &MainController{}
}

// Run starts the application
func (ctrl *MainController) Run() {
	// Create application
	ctrl.app = app.New()

	// Load configuration
	var err error
	ctrl.config, err = config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		ctrl.config = config.DefaultConfig()
	}

	// Initialize services
	ctrl.githubService = services.NewGitHubService(ctrl.config.GitHub)
	ctrl.openaiService = services.NewOpenAIService(ctrl.config.OpenAI)

	// Create main window
	ctrl.window = ctrl.app.NewWindow(utils.AppName)
	ctrl.window.SetMaster()
	ctrl.window.Resize(fyne.NewSize(800, 600))
	ctrl.window.CenterOnScreen()

	// Create UI
	ctrl.ui = ui.NewMainWindowUI()
	ctrl.window.SetContent(ctrl.ui.CreateMainLayout())

	// Set up callbacks
	ctrl.setupCallbacks()

	// Show window and run
	ctrl.window.ShowAndRun()
}

// setupCallbacks sets up all UI callbacks
func (ctrl *MainController) setupCallbacks() {
	// Analyze button callback
	ctrl.ui.SetAnalyzeButtonCallback(func() {
		ctrl.analyzeGitHubProfile()
	})

	// Config button callback
	ctrl.ui.SetConfigButtonCallback(func() {
		ctrl.showConfigWindow()
	})

	// Save button callback
	ctrl.ui.SetSaveButtonCallback(func() {
		ctrl.saveResults()
	})

	// Clear button callback
	ctrl.ui.SetClearButtonCallback(func() {
		ctrl.ui.ClearResults()
	})
}

// analyzeGitHubProfile performs GitHub profile analysis
func (ctrl *MainController) analyzeGitHubProfile() {
	username := strings.TrimSpace(ctrl.ui.GetUsername())
	if username == "" {
		dialog.ShowError(fmt.Errorf("please enter a GitHub username"), ctrl.window)
		return
	}

	// Disable analyze button and show progress
	ctrl.ui.SetAnalyzeButtonEnabled(false)
	ctrl.ui.ShowProgress()
	ctrl.ui.SetStatus("Starting analysis...")
	ctrl.ui.SetProgress(0.1)

	// Perform analysis in goroutine
	go func() {
		defer func() {
			ctrl.ui.HideProgress()
			ctrl.ui.SetAnalyzeButtonEnabled(true)
		}()

		ctx := context.Background()

		// Update progress
		ctrl.ui.SetStatus("Fetching user information...")
		ctrl.ui.SetProgress(0.2)

		// Perform the audit
		result, err := ctrl.githubService.PerformFullAudit(ctx, username)
		if err != nil {
			ctrl.ui.SetStatus("Analysis failed")
			dialog.ShowError(fmt.Errorf("analysis failed: %v", err), ctrl.window)
			return
		}

		ctrl.ui.SetProgress(0.7)
		ctrl.ui.SetStatus("Generating JSON report...")

		// Convert result to JSON for display
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			ctrl.ui.SetStatus("Failed to generate report")
			dialog.ShowError(fmt.Errorf("failed to generate report: %v", err), ctrl.window)
			return
		}

		// Save debug JSON if enabled
		if ctrl.config.GitHub.SaveDebugJSON {
			ctrl.ui.SetStatus("Saving debug JSON file...")
			err = ctrl.saveDebugJSON(result, username)
			if err != nil {
				// Log error but don't stop the process
				fmt.Printf("Warning: Failed to save debug JSON: %v\n", err)
			}
		}

		// Check if OpenAI is configured for LLM analysis
		if ctrl.config.OpenAI.APIKey != "" {
			ctrl.ui.SetProgress(0.8)
			ctrl.ui.SetStatus("Performing LLM analysis...")

			// Perform OpenAI analysis
			markdownAnalysis, err := ctrl.openaiService.AnalyzeGitHubData(result)
			if err != nil {
				// Show warning but continue with JSON report
				ctrl.ui.SetStatus("LLM analysis failed, showing JSON report")
				dialog.ShowInformation("LLM Analysis Failed", fmt.Sprintf("OpenAI analysis failed: %v\n\nShowing JSON report instead.", err), ctrl.window)
			} else {
				ctrl.ui.SetProgress(0.9)
				ctrl.ui.SetStatus("Converting markdown to HTML and generating report...")

				// Convert markdown to HTML
				htmlContent := ctrl.openaiService.ConvertMarkdownToHTML(markdownAnalysis, username)

				// Save report
				err = ctrl.saveHTMLReport(htmlContent, username)
				if err != nil {
					dialog.ShowError(fmt.Errorf("failed to save HTML report: %v", err), ctrl.window)
					return
				}

				// Display success message in the UI (HTML content is saved to file)
				ctrl.ui.SetProgress(1.0)
				ctrl.ui.SetStatus("LLM analysis completed successfully - HTML report saved")
				ctrl.ui.SetResults(fmt.Sprintf("**Analysis Complete!**\n\nHTML report generated and opened in browser.\n\n**Summary:** Professional technical assessment completed for user '%s'. The detailed report includes:\n\n- User profile overview\n- Repository analysis\n- Code quality assessment\n- Experience level mapping\n- Hiring recommendations\n\nThe full report has been saved and automatically opened in your default browser.", username))
				return
			}
		}

		// Update UI with JSON results (fallback or no OpenAI)
		ctrl.ui.SetProgress(1.0)
		if ctrl.config.OpenAI.APIKey == "" {
			ctrl.ui.SetStatus("Analysis completed - Configure OpenAI for LLM analysis")
		} else {
			ctrl.ui.SetStatus("Analysis completed with JSON report")
		}
		ctrl.ui.SetResults(string(jsonData))
	}()
}

// showConfigWindow shows the configuration dialog
func (ctrl *MainController) showConfigWindow() {
	// Create config UI
	ctrl.configUI = ui.NewConfigWindowUI()
	configContent := ctrl.configUI.CreateConfigLayout()
	ctrl.configUI.LoadConfig(ctrl.config.GitHub, ctrl.config.OpenAI)

	// Create a custom dialog with the config content
	configDialog := dialog.NewCustomWithoutButtons("Configuration", configContent, ctrl.window)
	configDialog.Resize(fyne.NewSize(850, 600)) // Optimized size for scrollable content

	// Set up config callbacks
	ctrl.configUI.SetSaveButtonCallback(func() {
		ctrl.saveConfigFromDialog()
		configDialog.Hide()
	})

	ctrl.configUI.SetCancelButtonCallback(func() {
		configDialog.Hide()
	})

	configDialog.Show()
}

// saveConfigFromDialog saves the configuration from modal dialog
func (ctrl *MainController) saveConfigFromDialog() {
	newGitHubConfig, newOpenAIConfig, err := ctrl.configUI.GetConfig()
	if err != nil {
		dialog.ShowError(fmt.Errorf("invalid configuration: %v", err), ctrl.window)
		return
	}

	// Update config
	ctrl.config.GitHub = newGitHubConfig
	ctrl.config.OpenAI = newOpenAIConfig

	// Save to file
	if err := config.SaveConfig(ctrl.config); err != nil {
		dialog.ShowError(fmt.Errorf("failed to save configuration: %v", err), ctrl.window)
		return
	}

	// Update services
	ctrl.githubService = services.NewGitHubService(ctrl.config.GitHub)
	ctrl.openaiService = services.NewOpenAIService(ctrl.config.OpenAI)

	// Show success
	dialog.ShowInformation("Configuration Saved", "Configuration has been saved successfully.", ctrl.window)
}

// saveHTMLReport saves the HTML report to a file and opens it in browser
func (ctrl *MainController) saveHTMLReport(htmlContent, username string) error {
	// Create reports directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	
	reportsDir := filepath.Join(homeDir, "github_reports")
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}
	
	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_github_assessment_%s.html", username, timestamp)
	filePath := filepath.Join(reportsDir, filename)
	
	// Write HTML file
	if err := os.WriteFile(filePath, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("failed to write HTML report: %w", err)
	}
	
	// Try to open in system browser
	go func() {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "linux":
			cmd = exec.Command("xdg-open", filePath)
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath)
		case "darwin":
			cmd = exec.Command("open", filePath)
		default:
			return // Unsupported OS
		}
		
		if err := cmd.Run(); err != nil {
			// Silently fail - user can manually open the file
			return
		}
	}()
	
	return nil
}

// saveDebugJSON saves the raw GitHub audit data to a JSON file for debugging
func (ctrl *MainController) saveDebugJSON(result *dto.AuditResult, username string) error {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	
	reportsDir := filepath.Join(homeDir, "github_reports")
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}
	
	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_github_audit_debug_%s.json", username, timestamp)
	filePath := filepath.Join(reportsDir, filename)
	
	// Convert result to JSON with pretty formatting
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal audit data to JSON: %w", err)
	}
	
	// Write JSON file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write debug JSON file: %w", err)
	}
	
	fmt.Printf("Debug JSON saved to: %s\n", filePath)
	return nil
}

// saveResults saves the analysis results to a file
func (ctrl *MainController) saveResults() {
	// Create file dialog
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		// Get results text from RichText
		results := ctrl.ui.ResultsRichText.String()
		if results == "" || results == "Analysis results will appear here..." {
			dialog.ShowError(fmt.Errorf("no results to save"), ctrl.window)
			return
		}

		// Write to file
		_, err = writer.Write([]byte(results))
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to save file: %v", err), ctrl.window)
			return
		}

		dialog.ShowInformation("File Saved", "Results have been saved successfully.", ctrl.window)
	}, ctrl.window)
}
