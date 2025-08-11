# GitHub Developer Profiler
A GitHub developer assessment tool that analyzes users' public repositories to generate detailed insights about their coding skills, 
technical expertise, and development patterns. Built with Go and Fyne, 
it combines intelligent repository analysis with AI-powered code evaluation to create professional reports showcasing developer capabilities, 
language proficiency, code quality metrics, and project complexity assessments.

## Features

### GitHub Data Analysis
- **Profile Info**: Gets basic GitHub profile data like when account was created and follower count
- **Repository Analysis**: Looks at both original and forked repositories
- **Code Review**: Samples code files to check skills and practices
- **Commit History**: Checks recent commits to see how the user codes
- **Language Detection**: Finds which programming languages and tools the user knows

### Sampling & Analysis
- **Adjustable Sampling**: Choose how many repositories, commits, and files to analyze
- **Recent Activity**: Focuses on recently updated repositories for better results
- **Consistent Results**: Uses fixed random sampling for repeatable results
- **Fork Analysis**: Tells the difference between original work and forked projects

### AI Features
- **OpenAI Integration**: Uses OpenAI's GPT models to analyze GitHub profiles
- **Editable System Prompt**: Change the instructions given to the AI
- **HTML Reports**: Create HTML reports using your own templates
- **CSS Styling**: Change how reports look with your own CSS
- **Syntax Highlighting**: Code in reports is colored for easier reading


### User Interface
- **Simple Design**: Easy-to-use interface built with Fyne v2.6.1
- **Progress Bars**: Shows analysis progress in real-time
- **Settings Tabs**: Organized settings window with tabs
- **First-Run Setup**: Setup wizard for new users
- **Save Results**: Save results as HTML reports and JSON files
- **Works Everywhere**: Runs on Windows, macOS, and Linux

## Screenshots

### Main Window
<img width="806" height="634" alt="Screenshot 2025-08-11 at 11 17 49" src="https://github.com/user-attachments/assets/24fc80ef-fd97-40a8-aac1-4a966eb57263" />

### Credentials Settings Window
<img width="810" height="638" alt="Screenshot 2025-08-11 at 11 26 54" src="https://github.com/user-attachments/assets/b5e18484-2778-4463-98da-478d48c4c26f" />


### Github Analysis Settings Window
<img width="807" height="638" alt="Screenshot 2025-08-11 at 11 27 05" src="https://github.com/user-attachments/assets/cc0915f4-9e7c-4309-a093-11d3e2545c85" />

### OpenAI Settings Page
<img width="804" height="636" alt="Screenshot 2025-08-11 at 11 27 21" src="https://github.com/user-attachments/assets/0f5abc2f-e9f8-413c-84c8-25956803d6d4" />


### System Prompt Editor
<img width="808" height="634" alt="Screenshot 2025-08-11 at 11 27 31" src="https://github.com/user-attachments/assets/c07a7c65-f03e-4cc1-a059-e079c75897f6" />


### HTML Template Editor
<img width="809" height="636" alt="Screenshot 2025-08-11 at 11 27 52" src="https://github.com/user-attachments/assets/9406fa50-32a8-4b71-987a-f34645535cfa" />


### Example Evaluation Report
<img width="1578" height="934" alt="Screenshot 2025-08-11 at 12 01 50" src="https://github.com/user-attachments/assets/16b65937-89f2-43e4-a08d-a1f9d383be72" />

*Note: A sample evaluation report is also available in PDF format:*
[report sample GitHub Developer Assessment](https://vpoluyaktov.github.io/github_developer_profiler/report_sample.html)


## Installation

### Download Ready-to-Use App

1. Go to the [GitHub Releases page] (https://github.com/vpoluyaktov/github_developer_profiler/releases)
2. Download the latest version for your system (Windows, macOS, or Linux)
3. Extract the files to any folder
4. Run the program (`github_developer_profiler` or `github_developer_profiler.exe`)

#### Note for macOS Users

When opening the app on macOS, you'll see a security warning because the app uses a self-signed certificate. To run it:

1. Right-click (or Control-click) on the app icon
2. Select "Open" from the menu
3. Click "Open" in the popup window
4. After the first time, you can open it normally

### What You Need to Build from Source
- Go 1.21 or later
- Git
- Fyne requirements (check [Fyne docs](https://developer.fyne.io/started/) for your system)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/vpoluyaktov/github_developer_profiler.git
cd github_developer_profiler

# Install dependencies
go mod tidy

# Build the application
go build -o github_developer_profiler .
```

### Running the Application

```bash
# Run the built executable
./github_developer_profiler
```

Or run directly with Go:

```bash
go run .
```

### Command Line Options

```bash
# Show version information
./github_developer_profiler --version
```

## Settings

The app saves settings in your home folder (`~/.dev_profiler/config.json`). When you first run the app, a setup wizard helps you configure it. You can set:

| Parameter | Default | Description |
|-----------|---------|-------------|
| `token` | "" | GitHub API token for authentication (supports private repos) |
| `sampled_repo_count` | 5 | Number of repositories to analyze in detail |
| `commits_per_repo` | 10 | Number of recent commits to examine per repository |
| `sample_file_count` | 10 | Number of code files to sample and analyze |
| `analysis_years` | 5 | Years of repository activity to consider |
| `include_private_repos` | false | Include private repositories if token allows |
| `random_seed` | 42 | Seed for reproducible sampling |
| `save_debug_json` | false | Save raw analysis data as JSON for debugging |
| `openai_api_key` | "" | OpenAI API key for AI-powered analysis |
| `openai_model` | "gpt-4" | OpenAI model to use for analysis |
| `system_prompt` | *template* | Customizable system prompt for AI analysis |
| `html_template` | *template* | Customizable HTML template for reports |
| `css_styles` | *template* | Customizable CSS styles for reports |

### GitHub Token Setup

1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Generate a new token with `repo` scope (for private repos) or `public_repo` scope
3. Enter the token in the application's configuration window

### OpenAI API Key Setup

1. Go to [OpenAI's website](https://platform.openai.com/) and sign up or log in
2. Navigate to the [API keys section](https://platform.openai.com/api-keys)
3. Click "Create new secret key" and give it a name
4. Copy the key immediately (you won't be able to see it again)
5. Enter the key in the OpenAI tab of the configuration window
6. Choose a model (like "gpt-4" or "gpt-3.5-turbo")

**Note**: OpenAI API usage incurs costs based on your usage. Check [OpenAI's pricing page](https://openai.com/pricing) for current rates.

## How to Use

1. **Start the App**: Run the program or use `go run .`
   - First-time users will see a setup wizard
2. **Enter Username**: Type in the GitHub username to analyze
3. **Change Settings** (Optional): Click "Configuration" to adjust options
4. **Run Analysis**: Click "Analyze GitHub Profile" to start
5. **See Results**: Results will show in the main window
   - With OpenAI API key: you get AI analysis
   - Without API key: you get raw JSON data
6. **Save Results**: Click "Save Results" to export
   - HTML reports save to `~/github_reports/` folder
   - Reports open in your web browser
   - Debug files save to the same folder if enabled


## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions, please open an issue on the project repository.
