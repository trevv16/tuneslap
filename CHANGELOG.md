# tuneslap

## 0.1.1

### Patch Changes

- a27a79a: Fix silent demo soundboard by shipping demo audio assets in the server image.

  The final stage of `server/Dockerfile` copied only the compiled binary, so the eight
  demo mp3s committed at `server/assets/demo/` never reached the runtime image. All
  three paths probed by `findDemoAudioFile` missed, demo seeding skipped every upload
  (`Demo media seeding completed: 0/8 files uploaded`), and the homepage soundboard
  served audio URLs that 404. This only affected built images; local dev bind-mounts
  `./server:/app`, which masked it.

## 0.1.0

### Minor Changes

- e0e710d: Initial public release of TuneSlap - the collaborative soundboard for creators.

  **Core Technologies:**

  - Frontend: Next.js 16, React 19, TypeScript, TanStack Query
  - Backend: Go (Fiber), MongoDB, Redis, FFmpeg
  - Real-time collaboration and Web Audio API

  **What's Included:**

  - Self-hosted soundboard with instant playback
  - Team collaboration with real-time sync
  - Audio editing tools (trim, fade, speed, loop)
  - Customizable boards with drag-and-drop organization

  Check the [issues tab](https://github.com/trevv16/tuneslap/issues) for upcoming features and roadmap.
