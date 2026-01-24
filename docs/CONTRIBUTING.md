# Contributing to TuneSlap

Thank you for your interest in contributing to TuneSlap! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.24+ installed
- Node.js 22+ and Yarn installed
- Docker and Docker Compose installed
- MongoDB and Redis (or use Docker Compose)

### Development Setup

1. Fork and clone the repository:

```bash
git clone https://github.com/yourusername/tuneslap.git
cd tuneslap
```

2. Set up the backend:

```bash
cd server
go mod download
cp example.env .env
# Edit .env with your configuration
```

3. Set up the frontend:

```bash
cd frontend
yarn install
```

4. Start dependencies with Docker Compose:

```bash
# From project root
docker-compose up -d mongodb redis
```

5. Run the backend:

```bash
cd server
go run main.go
```

6. Run the frontend (in a separate terminal):

```bash
cd frontend
yarn dev
```

## Development Workflow

### Making Changes

1. Create a new branch from `main`:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

2. Make your changes following the code standards below

3. Test your changes locally

4. Commit your changes (see Commit Guidelines)

5. Push to your fork and create a Pull Request

### Code Standards

#### Go (Backend)

- Follow Go formatting: run `go fmt ./...` before committing
- Run `go vet ./...` to check for issues
- Write tests for new functionality
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small

#### TypeScript/React (Frontend)

- Follow the existing code style
- Use TypeScript for type safety
- Run `yarn lint` before committing
- Use functional components with hooks
- Follow React best practices

### Commit Guidelines

We use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(media): add audio trimming support
fix(auth): resolve token expiration issue
docs(readme): update installation instructions
```

### Pull Request Process

1. Update documentation if needed
2. Add tests for new features
3. Ensure all tests pass
4. Update CHANGELOG.md if applicable (or use changesets)
5. Request review from maintainers
6. Address any feedback
7. Once approved, maintainers will merge

### Using Changesets

We use [Changesets](https://github.com/changesets/changesets) for version management and release automation. **Changesets are required for releases to happen.**

#### Why Changesets?

Changesets let us:

- Track what changed between releases
- Automatically determine version bumps (patch/minor/major)
- Generate changelogs
- Trigger automated builds and deployments

#### How the Release Flow Works

1. **You create a PR** with your code changes
2. **You add a changeset** describing what changed
3. **PR gets merged** to main
4. **GitHub Action creates a "Version Packages" PR** that bumps versions and updates the changelog
5. **When the Version Packages PR is merged**, Docker images are built and pushed, and deployment is triggered

**Without a changeset, steps 4-5 won't happen** and your changes won't be released.

#### Creating a Changeset

Before submitting your PR, run:

```bash
yarn changeset
```

You'll be prompted to:

1. Select which packages changed (use spacebar to select, enter to confirm)
2. Choose the version bump type:
   - `patch` - Bug fixes, small changes (0.0.X)
   - `minor` - New features, non-breaking changes (0.X.0)
   - `major` - Breaking changes (X.0.0)
3. Write a summary of your changes

This creates a markdown file in `.changeset/` with a random name. **Commit this file with your PR.**

#### Example

```bash
$ yarn changeset

What kind of change is this? (patch/minor/major)
> patch

Please enter a summary for this change:
> Fix audio player not pausing when switching tracks

# Creates .changeset/friendly-lions-dance.md
```

#### When to Skip Changesets

Not every PR needs a changeset. Skip them for:

- Documentation-only changes
- CI/tooling changes that don't affect the app
- Refactors with no user-facing changes

Add `[skip changeset]` to your PR title if intentionally skipping.

## Project Structure

```
tuneslap/
├── frontend/          # Next.js frontend application
├── server/            # Go backend API
│   ├── handlers/      # HTTP handlers
│   ├── models/        # Data models
│   ├── repositories/  # Database access layer
│   ├── services/      # Business logic
│   ├── mempeg/        # Media processing (audio/image)
│   ├── tasks/         # Background tasks (Asynq)
│   └── router/        # Route definitions
├── docs/              # Documentation
│   ├── media/         # Media system docs (uploads, storage, processing)
│   ├── CONTRIBUTING.md
│   └── DEPLOYMENT.md
├── scripts/           # Utility scripts
└── openapi/           # API specification (browse at openapi.tuneslap.com)
```

### API Reference

The full API documentation is available at [openapi.tuneslap.com](https://openapi.tuneslap.com). Use the explorer to browse endpoints, view request/response schemas, and test API calls.

## Testing

### Backend Tests

```bash
cd server
go test ./...
```

### Frontend Tests

```bash
cd frontend
yarn test
```

## Code Review Guidelines

When reviewing code, consider:

- Does the code solve the problem?
- Is it readable and maintainable?
- Are there edge cases not handled?
- Is error handling appropriate?
- Are tests adequate?
- Does it follow project conventions?

## Questions?

If you have questions:

- Open an issue for discussion
- Check existing issues and PRs
- Reach out to maintainers

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
