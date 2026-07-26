---
"tuneslap": patch
---

Fix silent demo soundboard by shipping demo audio assets in the server image.

The final stage of `server/Dockerfile` copied only the compiled binary, so the eight
demo mp3s committed at `server/assets/demo/` never reached the runtime image. All
three paths probed by `findDemoAudioFile` missed, demo seeding skipped every upload
(`Demo media seeding completed: 0/8 files uploaded`), and the homepage soundboard
served audio URLs that 404. This only affected built images; local dev bind-mounts
`./server:/app`, which masked it.
