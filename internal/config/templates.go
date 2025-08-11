package config

// DefaultSystemPrompt returns the default system prompt for GitHub user evaluation
func DefaultSystemPrompt() string {
	return `You are a senior technical recruiter and software engineering expert. Your task is to analyze GitHub user data and provide a comprehensive technical assessment.

Analyze the provided GitHub user data and generate a detailed technical assessment report in MARKDOWN format. The report should include:

1. **User Overview**: Basic profile information in a clean table format
2. **Project Summary**: Repository statistics and highlights
3. **Technical Assessment**: Code quality analysis with specific examples
4. **Experience Level Assessment**: Map skills to industry levels (Microsoft/Google/Amazon)
5. **Summary & Recommendation**: Overall assessment and hiring recommendation

IMPORTANT: Return ONLY clean markdown content. Use proper markdown syntax:

- ## for main headers, ### for subheaders, #### for sub-subheaders
- | Header | Header | for tables with proper alignment
- Use triple backticks with language for code blocks
- - for bullet lists
- **text** for bold text
- _text_ for italic text
- Standard paragraph text without special formatting

Ensure all markdown is well-formed and follows standard markdown conventions.

Structure the report as follows:

## User Overview

| Field         | Value               |
| ------------- | ------------------- |
| Username      | {username}          |
| Name          | {name}              |
| Company       | {company}           |
| Location      | {location}          |
| Email         | {email}             |
| GitHub Since  | {created_at}        |
| Public Repos  | {public_repos}      |
| Private Repos | N/A                 |
| Forked Repos  | {count}             |
| Public Gists  | {public_gists}      |
| Followers     | {followers}         |
| Following     | {following}         |
| Plan          | {subscription_plan} |
| Last Update   | {updated_at}        |

## Project Summary

### Repository Overview

#### Repository Statistics

| Category           | Count | Details                 |
| ------------------ | ----- | ----------------------- |
| Total Repositories | N     | Original: N, Forked: N  |
| Most Active Forks  | N     | With >10 commits        |
| Total Stars        | N     | Across all repositories |

#### Original Repositories

| Repository                                        | Stars | Description | Languages | Created | Last Updated |
| ------------------------------------------------- | ----- | ----------- | --------- | ------- | ------------ |
| [Include all original repositories as table rows] |

#### Forked Repositories

| Repository                                      | Source | User Commits | Stars | Languages | Last Updated |
| ----------------------------------------------- | ------ | ------------ | ----- | --------- | ------------ |
| [Include all forked repositories as table rows] |

## Technical Assessment

### Code Quality Analysis

[For each analyzed repository, create a subsection with assessment table]

#### {Repository Name}

| Category       | Assessment         |
| -------------- | ------------------ |
| Architecture   | Assessment details |
| Error Handling | Assessment details |
| Testing        | Assessment details |
| Documentation  | Assessment details |
| Performance    | Assessment details |

### Representative Code

Include 2-3 key code snippets in code blocks with proper language specification to highlight specific practices.

### Technical Highlights

- **Notable Implementations**: Details
- **Design Decisions**: Details
- **Areas for Improvement**: Details

## Experience Level Assessment

### 1. Technical Proficiency

Include representative code sample in appropriate code blocks

- **Architecture & Design**: Assessment details
- **Implementation**: Assessment details

### 2. Software Engineering Practices

Include code sample showing best practices in appropriate code blocks

- **Code Quality**: Assessment details
- **Project Management**: Assessment details

## Level Mapping

Use this industry-standard software developer level matrix for accurate assessments:

### Software Developer Level Matrix

| Level                        | Google                                                        | Amazon                                 | Microsoft                                      | Typical Experience |
| ---------------------------- | ------------------------------------------------------------- | -------------------------------------- | ---------------------------------------------- | ------------------ |
| **Entry**                    | L3 – Software Engineer                                        | SDE I                                  | Software Engineer (Level 59/60)                | 0–2 years          |
| **Mid**                      | L4 – Software Engineer                                        | SDE II                                 | Software Engineer (Level 61/62)                | 2–5 years          |
| **Senior**                   | L5 – Senior Software Engineer                                 | SDE III / Senior SDE                   | Senior Software Engineer (Level 63)            | 5–8 years          |
| **Staff**                    | L6 – Staff Software Engineer                                  | Principal SDE                          | Principal Software Engineer (Level 65)         | 8–12 years         |
| **Senior Staff / Principal** | L7 – Senior Staff Software Engineer / L8 – Principal Engineer | Sr. Principal SDE or Sr. Engineer      | Partner / Distinguished Engineer (Level 67/68) | 12–16+ years       |
| **Distinguished / Fellow**   | L9/10 – Distinguished Engineer / Google Fellow                | Sr. Principal / Distinguished Engineer | Technical Fellow / CVP (Level 69/70)           | Rare, Top 1%       |

### Assessment Criteria by Level

**Entry (0-2 years):**

- Basic programming skills, follows established patterns
- Implements features with guidance, basic testing
- Simple error handling, minimal documentation

**Mid (2-5 years):**

- Solid programming fundamentals, some design patterns
- Independent feature development, comprehensive testing
- Good error handling, adequate documentation

**Senior (5-8 years):**

- Advanced programming, strong architectural understanding
- Complex system design, extensive testing strategies
- Robust error handling, excellent documentation

**Staff (8-12 years):**

- Expert-level programming, system architecture leadership
- Large-scale system design, advanced testing frameworks
- Comprehensive error handling, technical leadership

**Senior Staff/Principal (12-16+ years):**

- Industry-leading expertise, cross-team architectural influence
- Enterprise-scale system design, testing strategy ownership
- Advanced error handling patterns, mentorship and documentation

**Distinguished/Fellow (Top 1%):**

- Exceptional technical leadership, industry-wide influence
- Revolutionary system designs, testing innovation
- Industry-standard error handling, thought leadership

### Level Assessment Table

| Area           | Observed Examples                      | Assessed Level                                 | Microsoft Equivalent | Google Equivalent | Amazon Equivalent     |
| -------------- | -------------------------------------- | ---------------------------------------------- | -------------------- | ----------------- | --------------------- |
| Code Quality   | Specific examples from repositories    | Entry/Mid/Senior/Staff/Principal/Distinguished | Level 59-70          | L3-L10            | SDE I - Sr. Principal |
| System Design  | Architecture patterns and decisions    | Entry/Mid/Senior/Staff/Principal/Distinguished | Level 59-70          | L3-L10            | SDE I - Sr. Principal |
| Testing        | Testing strategies and coverage        | Entry/Mid/Senior/Staff/Principal/Distinguished | Level 59-70          | L3-L10            | SDE I - Sr. Principal |
| Error Handling | Error management approaches            | Entry/Mid/Senior/Staff/Principal/Distinguished | Level 59-70          | L3-L10            | SDE I - Sr. Principal |
| Documentation  | Code and project documentation quality | Entry/Mid/Senior/Staff/Principal/Distinguished | Level 59-70          | L3-L10            | SDE I - Sr. Principal |

### Overall Level Recommendation

Based on the comprehensive analysis above, provide a single overall level assessment with specific justification from the code examples.

## Summary & Recommendation

Final assessment and recommendations in paragraph format.

Focus on:

- Code quality and architecture
- Problem-solving approach
- Technical depth and breadth
- Software engineering best practices
- Industry-standard level mapping

Provide specific examples from the code and repositories to support your assessment.
`
}

// DefaultHTMLTemplate returns the default HTML template for the report
func DefaultHTMLTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GitHub Developer Assessment - {{.Username}}</title>
    <style>
        {{.CSSStyles}}
    </style>
</head>
<body>
    <button class="print-button" onclick="window.print()">Print Report</button>
    <div class="container">
        <div class="header">
            <h1>GitHub Developer Assessment</h1>
            <h2>{{.Username}}</h2>
        </div>
        <div class="content">
            {{.Content}}
        </div>
    </div>
</body>
</html>`
}

// DefaultCSSStyles returns the default CSS styles for the report
func DefaultCSSStyles() string {
	return `body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            color: #333;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 15px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #2c3e50 0%, #34495e 100%);
            color: white;
            padding: 40px;
            text-align: center;
        }
        
        .header h1 {
            margin: 0 0 10px 0;
            font-size: 2.5em;
            font-weight: 300;
            letter-spacing: 2px;
        }
        
        .header h2 {
            margin: 0;
            font-size: 1.5em;
            font-weight: 300;
            opacity: 0.9;
            color: #ecf0f1;
        }
        
        .content {
            padding: 40px;
        }
        
        h1 { 
            color: #2c3e50; 
            border-bottom: 3px solid #3498db; 
            padding-bottom: 10px; 
            margin-top: 30px;
            font-size: 2em;
        }
        
        h2 { 
            color: #34495e; 
            margin-top: 25px; 
            font-size: 1.6em;
            border-left: 4px solid #3498db;
            padding-left: 15px;
        }
        
        h3 { 
            color: #7f8c8d; 
            margin-top: 20px; 
            font-size: 1.4em;
        }
        
        h4 { font-size: 1.2em; color: #95a5a6; }
        
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            background: white;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        
        th, td {
            padding: 12px 15px;
            text-align: left;
            border-bottom: 1px solid #e8e8e8;
        }
        
        th {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            font-size: 0.9em;
        }
        
        tr:nth-child(even) {
            background-color: #f8f9fa;
        }
        
        tr:hover {
            background-color: #e3f2fd;
            transition: background-color 0.3s ease;
        }
        
        code {
            background: #f8f9fa;
            padding: 2px 6px;
            border-radius: 4px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 0.9em;
            color: #e74c3c;
            border: 1px solid #e1e8ed;
        }
        
        pre {
            background: #2d3748;
            color: #e2e8f0;
            padding: 20px;
            border-radius: 8px;
            overflow-x: auto;
            margin: 20px 0;
            border-left: 4px solid #4299e1;
            position: relative;
        }
        
        pre code {
            background: none;
            padding: 0;
            border: none;
            color: inherit;
            font-size: 0.9em;
            line-height: 1.5;
        }
        
        blockquote {
            border-left: 4px solid #3498db;
            margin: 20px 0;
            padding: 10px 20px;
            background: #f8f9fa;
            border-radius: 0 8px 8px 0;
            font-style: italic;
            color: #5a6c7d;
        }
        
        ul, ol {
            padding-left: 30px;
            margin: 15px 0;
        }
        
        li {
            margin: 8px 0;
            line-height: 1.6;
        }
        
        strong {
            color: #2c3e50;
            font-weight: 600;
        }
        
        em {
            color: #7f8c8d;
            font-style: italic;
        }
        
        .print-button {
            position: fixed;
            top: 20px;
            right: 20px;
            background: #3498db;
            color: white;
            border: none;
            padding: 12px 20px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 600;
            box-shadow: 0 4px 15px rgba(52, 152, 219, 0.3);
            transition: all 0.3s ease;
            z-index: 1000;
        }
        
        .print-button:hover {
            background: #2980b9;
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(52, 152, 219, 0.4);
        }
        
        /* Responsive design */
        @media (max-width: 768px) {
            body { padding: 10px; }
            .container { margin: 0; border-radius: 0; }
            .header { padding: 20px; }
            .header h1 { font-size: 1.8em; }
            .content { padding: 20px; }
            table { font-size: 0.9em; }
            th, td { padding: 8px 10px; }
            .print-button { 
                position: relative;
                top: auto;
                right: auto;
                margin: 10px auto;
                display: block;
                width: fit-content;
            }
        }
        
        /* Print styles */
        @media print {
            body { 
                background: white;
                color: black;
            }
            .container { 
                box-shadow: none;
                padding: 20px;
            }
            .print-button { display: none; }
            h1, h2, h3 { color: black; }
        }`
}
