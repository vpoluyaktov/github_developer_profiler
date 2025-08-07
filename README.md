# GitHub Developer Profiler

A comprehensive GitHub user technical assessment tool built with Go and Fyne GUI framework. This application performs automated audits of GitHub users' public repositories to provide detailed technical skill assessments and experience level evaluations.

## Features

### Comprehensive GitHub Analysis
- **User Profile Analysis**: Extracts complete GitHub profile information including account creation date, follower counts, and subscription plans
- **Repository Assessment**: Analyzes both original and forked repositories with detailed metrics
- **Code Quality Evaluation**: Samples and analyzes code files to assess programming skills and best practices
- **Commit History Analysis**: Reviews recent commits to understand development patterns and activity levels
- **Technology Stack Detection**: Identifies programming languages and technologies used across repositories

### Smart Sampling & Analysis
- **Configurable Sampling**: Customizable number of repositories, commits, and files to analyze
- **Recent Activity Focus**: Prioritizes recently active repositories for more relevant assessments
- **Reproducible Results**: Uses seeded random sampling for consistent evaluation results
- **Fork Contribution Analysis**: Distinguishes between original work and fork contributions
- **Robust Error Handling**: Gracefully handles private repository access issues, API rate limits, and network errors

### Modern GUI Interface
- **Intuitive Design**: Clean, modern interface built with Fyne v2.6.1
- **Real-time Progress**: Live progress indicators during analysis
- **Configuration Management**: Easy-to-use configuration window for all parameters
- **Results Export**: Save analysis results to JSON files
- **Cross-platform**: Runs on Windows, macOS, and Linux

## Installation

### Prerequisites
- Go 1.21 or later
- Git

### Building from Source

```bash
git clone <repository-url>
cd dev_profiler
go mod tidy
go build -o dev_profiler .
```

### Running the Application

```bash
./dev_profiler
```

Or run directly with Go:

```bash
go run .
```

## Configuration

The application uses a configuration file stored in your home directory (`dev_profiler_config.json`). You can configure:

| Parameter | Default | Description |
|-----------|---------|-------------|
| `token` | "" | GitHub API token for authentication (supports private repos) |
| `sampled_repo_count` | 5 | Number of repositories to analyze in detail |
| `commits_per_repo` | 10 | Number of recent commits to examine per repository |
| `sample_file_count` | 10 | Number of code files to sample and analyze |
| `analysis_years` | 5 | Years of repository activity to consider |
| `include_private_repos` | false | Include private repositories if token allows |
| `random_seed` | 42 | Seed for reproducible sampling |

### GitHub Token Setup

1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Generate a new token with `repo` scope (for private repos) or `public_repo` scope
3. Enter the token in the application's configuration window

## Usage

1. **Launch the Application**: Run the executable or use `go run .`
2. **Enter Username**: Type the GitHub username you want to analyze
3. **Configure Settings** (Optional): Click "Configuration" to adjust analysis parameters
4. **Start Analysis**: Click "Analyze GitHub Profile" to begin the assessment
5. **View Results**: The analysis results will appear in JSON format in the results area
6. **Save Results**: Use the "Save Results" button to export the analysis to a file

## Architecture

The application follows a clean, modular architecture:

```
internal/
├── app/           # Application entry point
├── controllers/   # Business logic and UI event handling
├── dto/           # Data transfer objects and structures
├── config/        # Configuration management
├── services/      # GitHub API integration and analysis logic
├── ui/            # GUI components and layouts
└── utils/         # Utility functions and version info
```

## Dependencies

- **Fyne v2.6.1**: Cross-platform GUI framework
- **go-github v62**: GitHub API client library
- **oauth2**: OAuth2 authentication for GitHub API

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions, please open an issue on the project repository.
