# Contributing to go-geolocation

Thank you for your interest in contributing to go-geolocation! We welcome contributions from the community and are pleased to have you join us.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code:

- **Be respectful**: Treat everyone with respect and kindness
- **Be inclusive**: Welcome newcomers and encourage diverse perspectives
- **Be collaborative**: Work together constructively and professionally
- **Be patient**: Help others learn and grow

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title**: Summarize the problem in the title
- **Detailed description**: Explain what you expected vs. what actually happened
- **Steps to reproduce**: List the steps to reproduce the behavior
- **Environment**: Include Go version, web framework (Gin/Echo/Fiber), and module version
- **Code samples**: Include relevant code snippets or configuration

### Suggesting Features

Feature requests are welcome! Please:

- Check existing issues for similar requests
- Explain the use case and why it would be beneficial
- Provide examples of how the feature would work
- Consider how it fits with the project's goals

### Development Setup

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/go-geolocation.git
   cd go-geolocation
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```
4. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

### Running Tests

Before submitting changes, ensure all tests pass:

```bash
# Run the test suite
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run linter (if golangci-lint is installed)
golangci-lint run

# Format code
go fmt ./...
```

### Code Standards

We follow Go best practices and coding standards:

- **gofmt**: Code must be formatted with `go fmt`
- **golint**: Follow Go linting recommendations
- **go vet**: Code must pass `go vet` checks
- **High Test Coverage**: All new code should include tests
- **Documentation**: Public functions and types must have GoDoc comments

### Adding New Framework Adapters

To add a new web framework adapter:

1. Create a new package in the project root (e.g., `myadapter/`).
2. Implement middleware following the framework's conventions.
3. Provide a `FromContext` function to extract geolocation data.
4. Add comprehensive tests in `myadapter/myadapter_test.go`.
5. Document usage in the README.
6. Add example usage in the `examples/` directory.

### Writing Tests

We use Go's built-in testing framework. When adding new features:

1. **Write tests first** (Test-Driven Development)
2. **Cover edge cases** and error conditions
3. **Use descriptive test names** that explain what is being tested
4. **Mock external dependencies** when necessary

Example test structure:

```go
func TestParseClientInfo_Mobile(t *testing.T) {
    r, _ := http.NewRequest("GET", "/", nil)
    r.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X)")

    info := ParseClientInfo(r)

    if info.Device != "Mobile" {
        t.Errorf("expected device to be 'Mobile', got '%s'", info.Device)
    }
}
```

### Submitting Changes

1. **Commit your changes** with clear, descriptive messages:

   ```bash
   git commit -m "Add feature: configurable rate limiting"
   ```

2. **Push to your fork**:

   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request** on GitHub with:
   - Clear title and description
   - Reference to any related issues
   - Screenshots/examples if applicable
   - Confirmation that tests pass

### Pull Request Guidelines

- **One feature per PR**: Keep changes focused and atomic
- **Update documentation**: Include relevant documentation updates
- **Backward compatibility**: Avoid breaking changes when possible
- **Performance**: Consider the performance impact of changes
- **Security**: Be mindful of security implications

### Commit Message Format

Use clear, descriptive commit messages:

```text
type(scope): description

Examples:
feat(adapters): add support for Fiber framework
fix(middleware): handle empty User-Agent headers
docs(readme): update configuration examples
test(core): add edge case tests for language parsing
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `style`, `chore`

### Documentation

- Update relevant documentation for any changes
- Include code examples for new features
- Update the CHANGELOG.md for notable changes
- Ensure README.md stays current

## Getting Help

- **Issues**: For bugs and feature requests
- **Discussions**: For questions and general discussion
- **Email**: [contact@rumenx.com](mailto:contact@rumenx.com) for private inquiries

## Recognition

Contributors are recognized in:

- CHANGELOG.md for significant contributions
- GitHub contributors page
- Release notes for major features

## Development Philosophy

go-geolocation aims to be:

- **Framework-agnostic**: Works with any Go web framework
- **Secure by default**: Built-in security features
- **Easy to extend**: Clean architecture and interfaces
- **Well-tested**: High test coverage and quality
- **Performance-focused**: Efficient and scalable

Thank you for contributing to go-geolocation! 🚀
