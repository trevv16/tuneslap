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
cp .env.example .env
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

For version management, we use Changesets:

1. Create a changeset:
```bash
yarn changeset
```

2. Select the type of change (patch, minor, major)
3. Describe your changes
4. Commit the changeset file

The changeset will be included in the next release.

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
└── openapi/           # API specification
```

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
