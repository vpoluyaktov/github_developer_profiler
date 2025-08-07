package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MainWindowUI represents the UI components for the main window
type MainWindowUI struct {
	// GitHub user input
	UsernameEntry  *widget.Entry
	AnalyzeButton  *widget.Button
	
	// Configuration
	ConfigButton *widget.Button
	
	// Progress and status
	ProgressBar    *widget.ProgressBar
	StatusLabel    *widget.Label
	
	// Results display
	ResultsRichText *widget.RichText
	ResultsScroll   *container.Scroll
	SaveButton      *widget.Button
	ClearButton     *widget.Button
}

// NewMainWindowUI creates a new MainWindowUI instance
func NewMainWindowUI() *MainWindowUI {
	return &MainWindowUI{}
}

// CreateMainLayout creates the main window layout
func (ui *MainWindowUI) CreateMainLayout() *fyne.Container {
	// Initialize components
	ui.initializeComponents()
	
	configPanel := ui.createConfigPanel()
	buttonPanel := ui.createButtonPanel()
	outputPanel := ui.createOutputPanel()

	// Use Border layout with output panel as center to allow it to expand
	mainLayout := container.NewBorder(
		configPanel, // top
		buttonPanel, // bottom
		nil,         // left
		nil,         // right
		outputPanel, // center - this will expand to fill available space
	)

	// Reduced padding for more compact layout
	paddedLayout := container.NewPadded(mainLayout)

	mainBg := canvas.NewRectangle(theme.Color(theme.ColorNameBackground))
	return container.NewStack(mainBg, paddedLayout)
}

// initializeComponents initializes all UI components
func (ui *MainWindowUI) initializeComponents() {
	// Username input
	ui.UsernameEntry = widget.NewEntry()
	ui.UsernameEntry.SetPlaceHolder("Enter GitHub username...")
	
	// Buttons
	ui.AnalyzeButton = widget.NewButton("Analyze GitHub Profile", nil)
	ui.AnalyzeButton.Importance = widget.HighImportance
	
	ui.ConfigButton = widget.NewButton("Configuration", nil)
	ui.ConfigButton.SetIcon(theme.SettingsIcon())
	
	ui.SaveButton = widget.NewButton("Save Results", nil)
	ui.SaveButton.SetIcon(theme.DocumentSaveIcon())
	ui.SaveButton.Disable()
	
	ui.ClearButton = widget.NewButton("Clear", nil)
	ui.ClearButton.SetIcon(theme.DeleteIcon())
	
	// Progress and status
	ui.ProgressBar = widget.NewProgressBar()
	ui.ProgressBar.Hide()
	
	ui.StatusLabel = widget.NewLabel("Ready to analyze GitHub profiles")
	
	// Results display - Use RichText for better dark theme visibility
	ui.ResultsRichText = widget.NewRichTextFromMarkdown("Analysis results will appear here...")
	ui.ResultsRichText.Wrapping = fyne.TextWrapWord
	ui.ResultsScroll = container.NewScroll(ui.ResultsRichText)
	ui.ResultsScroll.SetMinSize(fyne.NewSize(600, 400))
}

// createConfigPanel creates the configuration panel
func (ui *MainWindowUI) createConfigPanel() *fyne.Container {
	// GitHub Analysis Section
	githubHeader := CreateSectionHeader("GitHub Developer Profiler")
	
	// Subtitle
	subtitle := widget.NewLabel("Comprehensive GitHub user technical assessment tool")
	subtitle.Alignment = fyne.TextAlignCenter
	
	// Username input
	usernameLabel := widget.NewLabel("GitHub Username:")
	usernameContainer := container.NewBorder(nil, nil, usernameLabel, nil, ui.UsernameEntry)
	
	githubSection := container.NewVBox(githubHeader, subtitle, usernameContainer)
	
	return githubSection
}

// createButtonPanel creates the action buttons panel
func (ui *MainWindowUI) createButtonPanel() *fyne.Container {
	// Create action header
	actionHeader := CreateSectionHeader("Actions")

	// Create a layout for buttons with spacers between them
	buttonPanel := container.NewHBox(
		layout.NewSpacer(),
		ui.AnalyzeButton,
		FixedSpacer(5, 5, 0, 0),
		ui.ConfigButton,
	)

	return container.NewVBox(actionHeader, buttonPanel)
}

// createOutputPanel creates the output console panel
func (ui *MainWindowUI) createOutputPanel() *fyne.Container {
	// Status and progress section
	statusHeader := CreateSectionHeader("Analysis Status")
	statusContainer := container.NewVBox(
		ui.StatusLabel,
		ui.ProgressBar,
	)
	statusSection := container.NewVBox(statusHeader, statusContainer)
	
	// Results section
	resultsHeader := CreateSectionHeader("Analysis Results")
	
	// Configure RichText for better display
	ui.ResultsRichText.Wrapping = fyne.TextWrapWord
	ui.ResultsScroll.SetMinSize(fyne.NewSize(0, 200))
	
	// Results footer with save/clear buttons
	resultsFooterBar := container.NewBorder(nil, nil, nil, 
		container.NewHBox(ui.SaveButton, ui.ClearButton))
	
	// Use Border layout so ResultsScroll expands
	resultsSection := container.NewBorder(
		resultsHeader,    // top
		resultsFooterBar, // bottom
		nil,              // left
		nil,              // right
		ui.ResultsScroll, // center (expands)
	)
	
	// Combine status and results
	return container.NewBorder(
		statusSection, // top
		nil,           // bottom
		nil,           // left
		nil,           // right
		resultsSection, // center (expands)
	)
}

// createHeaderSection creates the header section with title and logo
func (ui *MainWindowUI) createHeaderSection() *fyne.Container {
	// Title
	title := widget.NewLabel("GitHub Developer Profiler")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	
	// Subtitle
	subtitle := widget.NewLabel("Comprehensive GitHub user technical assessment tool")
	subtitle.Alignment = fyne.TextAlignCenter
	
	// Create header container
	header := container.NewVBox(
		title,
		subtitle,
	)
	
	return container.NewCenter(header)
}

// createInputSection creates the input section
func (ui *MainWindowUI) createInputSection() *fyne.Container {
	// Username input with label
	usernameLabel := widget.NewLabel("GitHub Username:")
	usernameContainer := container.NewBorder(nil, nil, usernameLabel, nil, ui.UsernameEntry)
	
	// Buttons container
	buttonsContainer := container.NewHBox(
		ui.AnalyzeButton,
		layout.NewSpacer(),
		ui.ConfigButton,
	)
	
	// Input section
	inputSection := container.NewVBox(
		usernameContainer,
		buttonsContainer,
	)
	
	return inputSection
}

// createProgressSection creates the progress section
func (ui *MainWindowUI) createProgressSection() *fyne.Container {
	// Progress container
	progressContainer := container.NewVBox(
		ui.StatusLabel,
		ui.ProgressBar,
	)
	
	return progressContainer
}

// createResultsSection creates the results section
func (ui *MainWindowUI) createResultsSection() *fyne.Container {
	// Results header
	resultsLabel := widget.NewLabel("Analysis Results:")
	resultsLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	// Results buttons
	resultsButtons := container.NewHBox(
		ui.SaveButton,
		ui.ClearButton,
		layout.NewSpacer(),
	)
	
	// Results header container
	resultsHeader := container.NewBorder(nil, nil, resultsLabel, resultsButtons)
	
	// Results section
	resultsSection := container.NewVBox(
		resultsHeader,
		ui.ResultsScroll,
	)
	
	return resultsSection
}

// SetAnalyzeButtonCallback sets the callback for the analyze button
func (ui *MainWindowUI) SetAnalyzeButtonCallback(callback func()) {
	ui.AnalyzeButton.OnTapped = callback
}

// SetConfigButtonCallback sets the callback for the config button
func (ui *MainWindowUI) SetConfigButtonCallback(callback func()) {
	ui.ConfigButton.OnTapped = callback
}

// SetSaveButtonCallback sets the callback for the save button
func (ui *MainWindowUI) SetSaveButtonCallback(callback func()) {
	ui.SaveButton.OnTapped = callback
}

// SetClearButtonCallback sets the callback for the clear button
func (ui *MainWindowUI) SetClearButtonCallback(callback func()) {
	ui.ClearButton.OnTapped = callback
}

// GetUsername returns the entered username
func (ui *MainWindowUI) GetUsername() string {
	return ui.UsernameEntry.Text
}

// SetStatus updates the status label
func (ui *MainWindowUI) SetStatus(status string) {
	fyne.Do(func() {
		ui.StatusLabel.SetText(status)
	})
}

// ShowProgress shows the progress bar
func (ui *MainWindowUI) ShowProgress() {
	fyne.Do(func() {
		ui.ProgressBar.Show()
	})
}

// HideProgress hides the progress bar
func (ui *MainWindowUI) HideProgress() {
	fyne.Do(func() {
		ui.ProgressBar.Hide()
	})
}

// SetProgress sets the progress bar value (0.0 to 1.0)
func (ui *MainWindowUI) SetProgress(value float64) {
	fyne.Do(func() {
		ui.ProgressBar.SetValue(value)
	})
}

// SetResults sets the results text
func (ui *MainWindowUI) SetResults(results string) {
	fyne.Do(func() {
		ui.ResultsRichText.ParseMarkdown(results)
		ui.SaveButton.Enable()
	})
}

// ClearResults clears the results text
func (ui *MainWindowUI) ClearResults() {
	fyne.Do(func() {
		ui.ResultsRichText.ParseMarkdown("Analysis results will appear here...")
		ui.SaveButton.Disable()
	})
}

// SetAnalyzeButtonEnabled enables/disables the analyze button
func (ui *MainWindowUI) SetAnalyzeButtonEnabled(enabled bool) {
	fyne.Do(func() {
		if enabled {
			ui.AnalyzeButton.Enable()
		} else {
			ui.AnalyzeButton.Disable()
		}
	})
}
