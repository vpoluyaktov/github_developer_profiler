package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"dev_profiler/internal/dto"
)

// ConfigWindowUI represents the UI components for the configuration window
type ConfigWindowUI struct {
	// GitHub configuration
	TokenEntry            *widget.Entry
	SampledRepoCountEntry *widget.Entry
	CommitsPerRepoEntry   *widget.Entry
	SampleFileCountEntry  *widget.Entry
	AnalysisYearsEntry    *widget.Entry
	IncludePrivateCheck   *widget.Check
	RandomSeedEntry       *widget.Entry
	SaveDebugJSONCheck    *widget.Check
	// OpenAI configuration
	OpenAIKeyEntry        *widget.Entry
	OpenAIModelEntry      *widget.Entry
	SystemPromptEntry     *widget.Entry
	HTMLTemplateEntry     *widget.Entry
	CSSStylesEntry        *widget.Entry
	// Buttons
	SaveButton            *widget.Button
	CancelButton          *widget.Button
}

// NewConfigWindowUI creates a new ConfigWindowUI instance
func NewConfigWindowUI() *ConfigWindowUI {
	return &ConfigWindowUI{}
}

// CreateConfigLayout creates the main configuration layout with tabs
func (ui *ConfigWindowUI) CreateConfigLayout() *fyne.Container {
	// Initialize components
	ui.initializeComponents()
	
	// Create tab contents
	credentialsTab := ui.createCredentialsTab()
	githubTab := ui.createGitHubTab()
	openaiTab := ui.createOpenAITab()
	systemPromptTab := ui.createSystemPromptTab()
	htmlTemplateTab := ui.createHTMLTemplateTab()
	cssStylesTab := ui.createCSSStylesTab()
	
	// Create tabs
	tabs := container.NewAppTabs()
	tabs.Append(container.NewTabItem("Credentials", credentialsTab))
	tabs.Append(container.NewTabItem("GitHub", githubTab))
	tabs.Append(container.NewTabItem("OpenAI", openaiTab))
	tabs.Append(container.NewTabItem("System Prompt", systemPromptTab))
	tabs.Append(container.NewTabItem("HTML Template", htmlTemplateTab))
	tabs.Append(container.NewTabItem("CSS Styles", cssStylesTab))
	
	// Create buttons section
	buttonsSection := ui.createButtonsSection()
	
	// Create main layout with tabs and buttons
	mainContent := container.NewBorder(
		nil, buttonsSection, nil, nil, tabs,
	)
	
	return container.NewPadded(mainContent)
}

// initializeComponents initializes all UI components
func (ui *ConfigWindowUI) initializeComponents() {
	// GitHub token
	ui.TokenEntry = widget.NewPasswordEntry()
	ui.TokenEntry.SetPlaceHolder("Enter GitHub API token (required)")
	
	// Analysis parameters
	ui.SampledRepoCountEntry = widget.NewEntry()
	ui.SampledRepoCountEntry.SetPlaceHolder("10")
	
	ui.CommitsPerRepoEntry = widget.NewEntry()
	ui.CommitsPerRepoEntry.SetPlaceHolder("50")
	
	ui.SampleFileCountEntry = widget.NewEntry()
	ui.SampleFileCountEntry.SetPlaceHolder("5")
	
	ui.AnalysisYearsEntry = widget.NewEntry()
	ui.AnalysisYearsEntry.SetPlaceHolder("2")
	
	ui.IncludePrivateCheck = widget.NewCheck("Include private repositories", nil)
	
	ui.RandomSeedEntry = widget.NewEntry()
	ui.RandomSeedEntry.SetPlaceHolder("42")
	
	ui.SaveDebugJSONCheck = widget.NewCheck("Save debug JSON file (raw GitHub data)", nil)
	
	// OpenAI configuration
	ui.OpenAIKeyEntry = widget.NewPasswordEntry()
	ui.OpenAIKeyEntry.SetPlaceHolder("Enter OpenAI API key (required for LLM analysis)")
	
	ui.OpenAIModelEntry = widget.NewEntry()
	ui.OpenAIModelEntry.SetPlaceHolder("gpt-4o")
	
	ui.SystemPromptEntry = widget.NewMultiLineEntry()
	ui.SystemPromptEntry.SetPlaceHolder("Enter custom system prompt for GitHub user evaluation...")
	ui.SystemPromptEntry.Wrapping = fyne.TextWrapOff // Disable line wrapping
	ui.SystemPromptEntry.TextStyle = fyne.TextStyle{Monospace: true} // Use monospace font
	// Remove fixed resize to allow proper expansion in tab layout
	
	ui.HTMLTemplateEntry = widget.NewMultiLineEntry()
	ui.HTMLTemplateEntry.SetPlaceHolder("Enter custom HTML template for reports...")
	ui.HTMLTemplateEntry.Wrapping = fyne.TextWrapOff // Disable line wrapping
	ui.HTMLTemplateEntry.TextStyle = fyne.TextStyle{Monospace: true} // Use monospace font
	
	ui.CSSStylesEntry = widget.NewMultiLineEntry()
	ui.CSSStylesEntry.SetPlaceHolder("Enter custom CSS styles for reports...")
	ui.CSSStylesEntry.Wrapping = fyne.TextWrapOff // Disable line wrapping
	ui.CSSStylesEntry.TextStyle = fyne.TextStyle{Monospace: true} // Use monospace font
	
	// Buttons
	ui.SaveButton = widget.NewButton("Save", nil)
	ui.SaveButton.Importance = widget.HighImportance
	
	ui.CancelButton = widget.NewButton("Cancel", nil)
}

// createParametersSection creates the analysis parameters section
func (ui *ConfigWindowUI) createParametersSection() *fyne.Container {
	title := widget.NewLabel("Analysis Parameters")
	title.TextStyle = fyne.TextStyle{Bold: true}
	
	// Create form-like layout for parameters
	sampledRepoLabel := widget.NewLabel("Repositories to analyze:")
	sampledRepoContainer := container.NewBorder(nil, nil, sampledRepoLabel, nil, ui.SampledRepoCountEntry)
	
	commitsLabel := widget.NewLabel("Commits per repository:")
	commitsContainer := container.NewBorder(nil, nil, commitsLabel, nil, ui.CommitsPerRepoEntry)
	
	filesLabel := widget.NewLabel("Files to sample per repo:")
	filesContainer := container.NewBorder(nil, nil, filesLabel, nil, ui.SampleFileCountEntry)
	
	yearsLabel := widget.NewLabel("Years of activity to consider:")
	yearsContainer := container.NewBorder(nil, nil, yearsLabel, nil, ui.AnalysisYearsEntry)
	
	seedLabel := widget.NewLabel("Random seed:")
	seedContainer := container.NewBorder(nil, nil, seedLabel, nil, ui.RandomSeedEntry)
	
	parametersSection := container.NewVBox(
		title,
		sampledRepoContainer,
		commitsContainer,
		filesContainer,
		yearsContainer,
		ui.IncludePrivateCheck,
		seedContainer,
		ui.SaveDebugJSONCheck,
	)
	
	return parametersSection
}

// createButtonsSection creates the buttons section
func (ui *ConfigWindowUI) createButtonsSection() *fyne.Container {
	buttonsContainer := container.NewHBox(
		ui.SaveButton,
		ui.CancelButton,
	)
	
	return container.NewCenter(buttonsContainer)
}

// createCredentialsTab creates the credentials configuration tab
func (ui *ConfigWindowUI) createCredentialsTab() *fyne.Container {
	// GitHub Token Section
	githubTitle := widget.NewLabel("GitHub API Configuration")
	githubTitle.TextStyle = fyne.TextStyle{Bold: true}
	
	tokenLabel := widget.NewLabel("GitHub Token:")
	tokenHelp := widget.NewLabel("Required: Provides API access and avoids rate limiting (60 requests/hour without token)")
	tokenHelp.TextStyle = fyne.TextStyle{Italic: true}
	
	githubSection := container.NewVBox(
		githubTitle,
		tokenLabel,
		ui.TokenEntry,
		tokenHelp,
		widget.NewSeparator(),
	)
	
	// OpenAI API Key Section
	openaiTitle := widget.NewLabel("OpenAI API Configuration")
	openaiTitle.TextStyle = fyne.TextStyle{Bold: true}
	
	openaiKeyLabel := widget.NewLabel("OpenAI API Key:")
	openaiKeyHelp := widget.NewLabel("Required: Enables LLM-powered analysis and HTML report generation")
	openaiKeyHelp.TextStyle = fyne.TextStyle{Italic: true}
	
	openaiSection := container.NewVBox(
		openaiTitle,
		openaiKeyLabel,
		ui.OpenAIKeyEntry,
		openaiKeyHelp,
	)
	
	return container.NewVBox(githubSection, openaiSection)
}

// createGitHubTab creates the GitHub configuration tab with analysis parameters
func (ui *ConfigWindowUI) createGitHubTab() *fyne.Container {
	parametersSection := ui.createParametersSection()
	return container.NewVBox(parametersSection)
}

// createOpenAITab creates the OpenAI configuration tab with model selector only
func (ui *ConfigWindowUI) createOpenAITab() *fyne.Container {
	title := widget.NewLabel("OpenAI Model Configuration")
	title.TextStyle = fyne.TextStyle{Bold: true}
	
	openaiModelLabel := widget.NewLabel("OpenAI Model:")
	openaiModelHelp := widget.NewLabel("Model to use for analysis (e.g., gpt-4o, gpt-4, gpt-3.5-turbo)")
	openaiModelHelp.TextStyle = fyne.TextStyle{Italic: true}
	
	modelSection := container.NewVBox(
		title,
		openaiModelLabel,
		ui.OpenAIModelEntry,
		openaiModelHelp,
	)
	
	return container.NewVBox(modelSection)
}

// createSystemPromptTab creates the system prompt editor tab
func (ui *ConfigWindowUI) createSystemPromptTab() *fyne.Container {
	title := widget.NewLabel("System Prompt Configuration")
	title.TextStyle = fyne.TextStyle{Bold: true}
	
	description := widget.NewLabel("Customize the system prompt used for GitHub user evaluation. Leave empty to use the default prompt with industry-standard developer level matrix.")
	description.Wrapping = fyne.TextWrapWord
	description.TextStyle = fyne.TextStyle{Italic: true}
	
	// Create header section
	headerSection := container.NewVBox(
		title,
		description,
		widget.NewSeparator(),
	)
	
	// Create the text area with proper multiline support
	textAreaContainer := container.NewScroll(ui.SystemPromptEntry)
	
	// Use border layout to ensure text area expands to fill available space
	return container.NewBorder(
		headerSection, nil, nil, nil, textAreaContainer,
	)
}

// createHTMLTemplateTab creates the HTML template editor tab
func (ui *ConfigWindowUI) createHTMLTemplateTab() *fyne.Container {
	title := widget.NewLabel("HTML Template Configuration")
	title.TextStyle = fyne.TextStyle{Bold: true}
	
	description := widget.NewLabel("Customize the HTML template used for report generation. Use {{.Username}}, {{.Content}}, and {{.CSSStyles}} as placeholders. Leave empty to use the default template.")
	description.Wrapping = fyne.TextWrapWord
	description.TextStyle = fyne.TextStyle{Italic: true}
	
	// Create header section
	headerSection := container.NewVBox(
		title,
		description,
		widget.NewSeparator(),
	)
	
	// Create the text area with proper multiline support
	textAreaContainer := container.NewScroll(ui.HTMLTemplateEntry)
	
	// Use border layout to ensure text area expands to fill available space
	return container.NewBorder(
		headerSection, nil, nil, nil, textAreaContainer,
	)
}

// createCSSStylesTab creates the CSS styles editor tab
func (ui *ConfigWindowUI) createCSSStylesTab() *fyne.Container {
	title := widget.NewLabel("CSS Styles Configuration")
	title.TextStyle = fyne.TextStyle{Bold: true}
	
	description := widget.NewLabel("Customize the CSS styles used for report styling. This CSS will be embedded in the HTML template. Leave empty to use the default styles.")
	description.Wrapping = fyne.TextWrapWord
	description.TextStyle = fyne.TextStyle{Italic: true}
	
	// Create header section
	headerSection := container.NewVBox(
		title,
		description,
		widget.NewSeparator(),
	)
	
	// Create the text area with proper multiline support
	textAreaContainer := container.NewScroll(ui.CSSStylesEntry)
	
	// Use border layout to ensure text area expands to fill available space
	return container.NewBorder(
		headerSection, nil, nil, nil, textAreaContainer,
	)
}

// LoadConfig loads configuration into the UI
func (ui *ConfigWindowUI) LoadConfig(githubConfig *dto.GitHubConfig, openaiConfig *dto.OpenAIConfig) {
	// Load GitHub configuration
	ui.TokenEntry.SetText(githubConfig.Token)
	ui.SampledRepoCountEntry.SetText(strconv.Itoa(githubConfig.SampledRepoCount))
	ui.CommitsPerRepoEntry.SetText(strconv.Itoa(githubConfig.CommitsPerRepo))
	ui.SampleFileCountEntry.SetText(strconv.Itoa(githubConfig.SampleFileCount))
	ui.AnalysisYearsEntry.SetText(strconv.Itoa(githubConfig.AnalysisYears))
	ui.IncludePrivateCheck.SetChecked(githubConfig.IncludePrivateRepo)
	ui.RandomSeedEntry.SetText(strconv.Itoa(githubConfig.RandomSeed))
	ui.SaveDebugJSONCheck.SetChecked(githubConfig.SaveDebugJSON)
	
	// Load OpenAI configuration
	ui.OpenAIKeyEntry.SetText(openaiConfig.APIKey)
	ui.OpenAIModelEntry.SetText(openaiConfig.Model)
	
	// Load system prompt - show default if no custom prompt is configured
	if openaiConfig.SystemPrompt != "" {
		ui.SystemPromptEntry.SetText(openaiConfig.SystemPrompt)
	} else {
		// Display default system prompt for discoverability and easier editing
		ui.SystemPromptEntry.SetText(dto.DefaultSystemPrompt())
	}
	
	// Load HTML template - show default if no custom template is configured
	if openaiConfig.HTMLTemplate != "" {
		ui.HTMLTemplateEntry.SetText(openaiConfig.HTMLTemplate)
	} else {
		// Display default HTML template for discoverability and easier editing
		ui.HTMLTemplateEntry.SetText(dto.DefaultHTMLTemplate())
	}
	
	// Load CSS styles - show default if no custom styles are configured
	if openaiConfig.CSSStyles != "" {
		ui.CSSStylesEntry.SetText(openaiConfig.CSSStyles)
	} else {
		// Display default CSS styles for discoverability and easier editing
		ui.CSSStylesEntry.SetText(dto.DefaultCSSStyles())
	}
}

// GetConfig returns the configuration from the UI
func (ui *ConfigWindowUI) GetConfig() (*dto.GitHubConfig, *dto.OpenAIConfig, error) {
	githubConfig := &dto.GitHubConfig{}
	openaiConfig := &dto.OpenAIConfig{}
	
	// Get GitHub configuration
	githubConfig.Token = ui.TokenEntry.Text
	
	var err error
	githubConfig.SampledRepoCount, err = strconv.Atoi(ui.SampledRepoCountEntry.Text)
	if err != nil {
		return nil, nil, err
	}
	
	githubConfig.CommitsPerRepo, err = strconv.Atoi(ui.CommitsPerRepoEntry.Text)
	if err != nil {
		return nil, nil, err
	}
	
	githubConfig.SampleFileCount, err = strconv.Atoi(ui.SampleFileCountEntry.Text)
	if err != nil {
		return nil, nil, err
	}
	
	githubConfig.AnalysisYears, err = strconv.Atoi(ui.AnalysisYearsEntry.Text)
	if err != nil {
		return nil, nil, err
	}
	
	githubConfig.RandomSeed, err = strconv.Atoi(ui.RandomSeedEntry.Text)
	if err != nil {
		return nil, nil, err
	}
	
	githubConfig.IncludePrivateRepo = ui.IncludePrivateCheck.Checked
	githubConfig.SaveDebugJSON = ui.SaveDebugJSONCheck.Checked
	
	// Get OpenAI configuration
	openaiConfig.APIKey = ui.OpenAIKeyEntry.Text
	openaiConfig.Model = ui.OpenAIModelEntry.Text
	
	// Always save the current system prompt text
	// The LoadConfig logic will handle showing default when appropriate
	openaiConfig.SystemPrompt = ui.SystemPromptEntry.Text
	
	// Always save the current HTML template and CSS styles text
	// The LoadConfig logic will handle showing default when appropriate
	openaiConfig.HTMLTemplate = ui.HTMLTemplateEntry.Text
	openaiConfig.CSSStyles = ui.CSSStylesEntry.Text
	
	return githubConfig, openaiConfig, nil
}

// ResetToDefaults resets all fields to default values
func (ui *ConfigWindowUI) ResetToDefaults() {
	githubDefaults := dto.DefaultGitHubConfig()
	openaiDefaults := dto.DefaultOpenAIConfig()
	ui.LoadConfig(githubDefaults, openaiDefaults)
}

// SetSaveButtonCallback sets the callback for the save button
func (ui *ConfigWindowUI) SetSaveButtonCallback(callback func()) {
	ui.SaveButton.OnTapped = callback
}

// SetCancelButtonCallback sets the callback for the cancel button
func (ui *ConfigWindowUI) SetCancelButtonCallback(callback func()) {
	ui.CancelButton.OnTapped = callback
}


