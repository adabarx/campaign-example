# Campaign Website

A modern campaign website built with htmx, Go/Fiber, and templ.

## Features

- 🚀 Static site generation for fast page loads
- ⚡ Dynamic donation system with htmx
- 📝 Blog with individual posts
- 💰 Real-time donation stats
- 🎯 No frontend framework needed

## Quick Start

```bash
# Install dependencies
make install

# Generate static site and run server
make run
```

Visit http://localhost:3000

## Commands

- `make install` - Install Go dependencies and templ
- `make generate` - Generate templ code and static HTML
- `make run` - Run the server
- `make build` - Build binary
- `make clean` - Remove generated files
- `make test` - Show test commands

## Project Structure

```
campaign/
├── main.go              # Fiber server + API routes
├── build.go             # Static site generator
├── templates/           # Templ components
├── data/                # Blog posts
└── public/              # Generated static files
```

## Testing

### Via Browser
1. Visit http://localhost:3000
2. Fill out the donation form
3. Watch stats update in real-time

### Via cURL
```bash
# Get stats
curl http://localhost:3000/api/stats

# Create donation
curl -X POST http://localhost:3000/api/donations \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=Alice&email=alice@example.com&amount=100"

# Get recent donors
curl http://localhost:3000/api/recent-donors
```

## Tech Stack

- **Backend**: Go + Fiber
- **Templates**: templ (type-safe Go templates)
- **Frontend**: htmx (14KB, no framework)
- **Architecture**: Static site generation + dynamic API

## License

MIT
