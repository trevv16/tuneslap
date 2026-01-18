# TuneSlap

A self-hosted soundboard application for podcasters and content creators. Create and manage interactive soundboards with customizable audio buttons, manipulate audio clips, and collaborate with your team.

## Features

### Soundboard Management
- Upload and play audio files (MP3, WAV, etc.) with instant playback
- Create multiple customizable soundboards
- Drag-and-drop interface for organizing sounds
- Adjustable button labels, colors, and sizes
- Assign custom images to buttons

### Audio Controls
- Trim audio clips
- Fade in/out effects
- Modify playback speed and pitch
- Looping support

### Collaboration
- Invite team members to boards
- Role-based access control (creator, editor, viewer)
- Real-time collaboration

### User Management
- Email/password authentication
- Social login (Google)
- Password reset and email verification
- User profile management

## Quick Start

The easiest way to get started is using Docker Compose:

```bash
# Clone the repository
git clone https://github.com/yourusername/tuneslap.git
cd tuneslap

# Copy environment file template
cp server/.env.example server/.env

# Edit server/.env with your configuration
# See docs/DEPLOYMENT.md for detailed setup instructions

# Start all services
docker-compose up -d
```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

For detailed deployment instructions, see [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md).

## Tech Stack

### Frontend
- **Next.js 15** (React framework)
- **TypeScript**
- **TailwindCSS**
- **TanStack Query** (data fetching and caching)
- **Web Audio API** for audio previews

### Backend
- **Go (Fiber)** for API server
- **MongoDB** for data storage
- **Redis** for caching and job queues
- **FFmpeg** for audio processing
- **Sharp** for image processing

## Development

See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for development setup and contribution guidelines.

## API Documentation

All API routes are prefixed with `/api/v1`.

### Authentication
- `POST /auth/signup` - Register a new user
- `POST /auth/login` - Authenticate user and return tokens
- `POST /auth/forgot` - Request password reset
- `POST /auth/reset` - Complete password reset

### Boards
- `GET /boards` - List all boards
- `POST /boards` - Create a new board
- `GET /boards/:boardId` - Get board details
- `PATCH /boards/:boardId` - Update board
- `DELETE /boards/:boardId` - Delete board

### Media
- `GET /media` - List all media files
- `POST /media` - Upload a new media file
- `GET /media/:mediaId` - Get media metadata
- `PATCH /media/:mediaId` - Update media
- `DELETE /media/:mediaId` - Delete media
- `POST /media/:mediaId/process` - Process media (trim, convert, etc.)

### Collaborators
- `GET /boards/:boardId/collaborators` - List collaborators
- `POST /boards/:boardId/collaborators` - Add collaborator
- `PATCH /boards/:boardId/collaborators/:collaboratorId` - Update role
- `DELETE /boards/:boardId/collaborators/:collaboratorId` - Remove collaborator

For complete API documentation, see the [OpenAPI specification](openapi/openapi.yaml).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please read [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## Roadmap

- [ ] Support for additional audio formats
- [ ] MIDI device support
- [ ] Stream deck integration
- [ ] Audio library with preloaded sounds
- [ ] Advanced audio effects and filters
- [ ] Mobile app support
