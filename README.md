# Campaign Website

Campaign website built with htmx, Go/Fiber, and templ.

## Quick Start

**Unix/Linux/macOS:**
```bash
# Install dependencies
./build.sh install

# Generate static site and run server
./build.sh run
```

**Windows:**
```cmd
REM Install dependencies
build.bat install

REM Generate static site and run server
build.bat run
```

Visit http://localhost:3000

## Commands

**Unix/Linux/macOS:** Use `./build.sh <command>`  
**Windows:** Use `build.bat <command>`

- `install` - Install Go dependencies and templ
- `generate` - Generate templ code and static HTML
- `run` - Run the server
- `build` - Build binary
- `clean` - Remove generated files
- `test` - Show test commands

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
