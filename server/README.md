# TuneSlap Server

A template for a GoLang backend using Fiber, MongoDB, a Makefile, and more

## Getting Started

### VIDEOS

- [GoLang Download & Setup](https://www.youtube.com/watch?v=Q7uh85_i0-M)
- [Video Breakdown](https://youtu.be/6C-2R92L01Q)

### Prerequisites

- [GoLang](https://golang.org/doc/install)
- [MongoDB](https://docs.mongodb.com/manual/installation/)

### Installing

0. Install extra packages: 
    ```go install github.com/cosmtrek/air@latest```
1. Clone the repo
2. Create your own .env file
3. ```make dev```
4. Server running on http://localhost:8082
4. Build OAS docs with ```make docs-build```
4. Run OAS docs with ```make docs```
4. view docs at http://localhost:8081

### Scripts

- ```make dev``` - runs the server in development mode