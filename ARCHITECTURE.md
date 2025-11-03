# Architecture: Campaign Example

For developers new to web development, reading *Hypermedia Systems*.

---

## Table of Contents

1. [What This Project Does](#what-this-project-does)
2. [Project Structure](#project-structure)
3. [How a Request Works](#how-a-request-works)
4. [The Three Core Parts](#the-three-core-parts)
5. [Static Generation Pipeline](#static-generation-pipeline)
6. [API Handlers](#api-handlers)
7. [Data Models](#data-models)
8. [Key Templ Components](#key-templ-components)
9. [htmx Attributes](#htmx-attributes-in-this-project)
10. [Making Changes](#making-changes)
11. [Testing](#testing)
12. [Hypermedia Principles](#how-this-implements-hypermedia-principles)
13. [Typical Flow](#typical-flow)
14. [Build & Run](#build--run)
15. [Dependencies](#dependencies)
16. [Limitations](#limitations-poc)
17. [Further Reading](#further-reading--resources)

---

## What This Project Does

Campaign website where people can:
- Read static information (landing page, about, blog)
- Donate money
- See live donation counter and recent supporters

Built using hypermedia principles: the server sends complete HTML, the browser displays it, no client-side framework.

---

## Project Structure

```
campaign/
├── main.go                   # Fiber server, API handlers, donation state
├── build.go                  # Generates HTML files from templates + data
├── go.mod / go.sum           # Dependencies
├── build.sh                  # Build automation (Unix/Linux/macOS)
├── build.bat                 # Build automation (Windows)
│
├── templates/                # Templ components (HTML with Go)
│   ├── layout.templ          # Base HTML wrapper
│   ├── home.templ            # Home page
│   ├── about.templ           # About page
│   ├── blog.templ            # Blog list and posts
│   ├── donation_stats.templ  # Stats/counter fragment
│   └── donation_success.templ # Success message fragment
│
├── data/
│   └── posts.go              # Blog posts (hardcoded for PoC)
│
├── static-vendor/            # Vendored static assets
│   └── htmx.min.js           # htmx library (downloaded during build)
│
└── public/                   # Generated output (created by build.go)
    ├── index.html            # Generated home page
    ├── about.html
    ├── blog.html
    ├── blog/*.html           # Individual posts
    ├── js/
    │   └── htmx.min.js       # Copied from static-vendor/
    ├── style.css             # Blank stylesheet
    └── [all generated files]
```

---

## How a Request Works

### User visits homepage

```
1. Browser: GET http://localhost:3000/
2. Fiber: Looks for public/index.html, finds it, sends it
3. Browser: Displays HTML
4. HTML contains: <script src="/js/htmx.min.js"></script>
5. htmx library loads (47KB, served locally)
6. htmx scans page, finds: <div hx-get="/api/stats" hx-trigger="load">
7. htmx: Makes GET request to /api/stats
8. Fiber: Runs getStats() function
9. getStats(): Calculates total donations, calls Templ component, generates HTML
10. Fiber: Sends HTML back to browser
11. htmx: Swaps HTML into the page
12. User sees: Live donation stats appear (no page reload)
```

### User submits donation form

```
1. User fills form and clicks "Donate"
2. Browser would normally reload page
3. But htmx intercepts it (hx-post="/api/donations" on form)
4. htmx: Sends form data via background request (AJAX)
5. Fiber: Runs createDonation() handler
6. createDonation(): Validates, saves donation, renders success HTML
7. Fiber response includes: HX-Trigger: donationComplete header
8. htmx: Swaps success message into page
9. htmx: Sees HX-Trigger header, knows to cascade updates
10. htmx: Auto-fetches /api/stats and /api/recent-donors
11. User sees: Success message + updated counters + updated donor list
    All without page reload. Done in ~100ms.
```

---

## The Three Core Parts

### Templ: Server-Side HTML Components

Templ lets you write HTML with Go logic. It compiles to Go code.

**Example:** `templates/donation_stats.templ`
```templ
templ DonationStats(total float64, count int) {
  <div class="stats">
    <div class="stat">
      Total Raised: ${ fmt.Sprintf("%.2f", total) }
    </div>
    <div class="stat">
      Supporters: { count }
    </div>
  </div>
}
```

Used by both:
- `build.go` → generates static HTML files for the home page
- `main.go` handlers → generates HTML responses for API requests (`/api/stats`)

Same template, two uses. This is the key to the architecture.

### htmx: Hypermedia Client Library

47KB JavaScript library (served locally) that adds interaction to HTML without page reloads.

**What it does:**
- Intercepts form submissions and link clicks
- Makes background HTTP requests (AJAX)
- Swaps HTML responses into the page
- All controlled by HTML attributes

**Example:** `public/index.html`
```html
<!-- Load stats when page loads -->
<div hx-get="/api/stats" hx-trigger="load" hx-swap="innerHTML">
  Loading...
</div>

<!-- Submit donation form without reload -->
<form hx-post="/api/donations" hx-target="#result" hx-swap="innerHTML">
  <input name="name" placeholder="Your Name" required/>
  <input name="amount" placeholder="Amount" required/>
  <button type="submit">Donate</button>
</form>
<div id="result"></div>
```

No JavaScript code needed. HTML attributes tell htmx what to do.

### Fiber: Hypermedia Web Server

Go-based web server. Its job: listen for requests, generate HTML, send it back.

**Key difference:** Fiber sends HTML, not JSON. The browser doesn't parse data and render it. The server sends ready-to-display HTML.

**Example:** `main.go`
```go
func main() {
  app := fiber.New()
  
  // Serve static files (home page, blog, etc.)
  app.Static("/", "./public")
  
  // API endpoints (return HTML fragments)
  app.Get("/api/stats", getStats)
  app.Post("/api/donations", createDonation)
  
  app.Listen(":3000")
}

func getStats(c *fiber.Ctx) error {
  // Calculate data
  state.mu.Lock()
  total := calculateTotal(state.donations)
  count := len(state.donations)
  state.mu.Unlock()
  
  // Generate HTML using Templ
  c.Set("Content-Type", "text/html; charset=utf-8")
  return templates.DonationStats(total, count).Render(
    c.Context(), 
    c.Response().BodyWriter(),
  )
}
```

Handler pattern:
1. Get/calculate data
2. Call Templ component with data
3. Render component to HTML
4. Set Content-Type header
5. Send HTML back

---

## Static Generation Pipeline

When you run `./build.sh generate` (or `build.bat generate` on Windows):

```
1. Build script ensures static-vendor/htmx.min.js exists (downloads if missing)
2. build.go reads data/posts.go (blog posts)
3. build.go copies static-vendor/htmx.min.js to public/js/htmx.min.js
4. For each blog post, calls: templates.BlogPost(post)
5. Templ renders component to HTML string
6. build.go writes HTML to: public/blog/[slug].html
7. For home page, calls: templates.Home()
8. Templ renders to HTML
9. build.go writes to: public/index.html
10. Result: public/ folder fills with static HTML files
11. Fiber serves these files when /blog/[slug].html is requested
```

All HTML files are created once at build time. They never change until you rebuild.

The htmx library is downloaded once (if not present) and copied into the public directory during each build, eliminating the need for external CDN dependencies at runtime.

---

## API Handlers

### `getStats()`
Returns HTML snippet with donation stats.

When called:
- Calculate total from all donations
- Calculate supporter count
- Render `templates.DonationStats(total, count)`
- Send HTML fragment back

Called by:
- Browser on page load (via `hx-trigger="load"`)
- After donation completes (cascade trigger)

### `createDonation()`
Handles form submission.

When called:
- Parse form data into Donation struct
- Validate (name, email, amount present)
- Lock state, append to donations, unlock
- Render `templates.DonationSuccess()` with donor's name
- Send response with `HX-Trigger: donationComplete` header

The header tells htmx: "Something important changed, refresh related content."

---

## Data Models

### Donations (Runtime State)

In-memory storage for this PoC. Lost when server restarts.

```go
type Donation struct {
  ID        int
  Name      string
  Email     string
  Amount    float64
  Message   string
  CreatedAt time.Time
}

type AppState struct {
  donations []Donation
  nextID    int
  mu        sync.Mutex  // Protects concurrent access
}
```

Always use lock/unlock when accessing:
```go
state.mu.Lock()
defer state.mu.Unlock()
// Safe to read/write state.donations here
```

### Blog Posts (Static Data)

Loaded at build time from `data/posts.go`:

```go
type BlogPost struct {
  Slug    string
  Title   string
  Content string
  Date    string
}
```

Used by `build.go` to generate HTML. Never changes at runtime.

---

## Key Templ Components

| Component | File | Purpose |
|-----------|------|---------|
| `Layout(title, content)` | layout.templ | Base HTML wrapper, includes all pages |
| `Home()` | home.templ | Home page (includes form, stats, donors div) |
| `BlogList(posts)` | blog.templ | Blog listing page |
| `BlogPost(post)` | blog.templ | Individual blog post page |
| `DonationStats(total, count)` | donation_stats.templ | Counter snippet (used by static gen and API) |
| `RecentDonors(donations)` | donation_stats.templ | Donor list snippet |
| `DonationSuccess(name, amount)` | donation_success.templ | Success message after donation |

---

## htmx Attributes in This Project

| Attribute | Element | Purpose |
|-----------|---------|---------|
| `hx-get="/api/stats"` | stats div | Fetch stats from server |
| `hx-trigger="load"` | stats div | Trigger on page load |
| `hx-trigger="donationComplete from:body"` | stats/donors div | Listen for cascade trigger |
| `hx-post="/api/donations"` | form | Submit via background request |
| `hx-target="#result"` | form | Put response in this div |
| `hx-swap="innerHTML"` | multiple | Replace element's content |

Server response header:
| Header | Purpose |
|--------|---------|
| `HX-Trigger: donationComplete` | Tell htmx other elements should refresh |

---

## Making Changes

### Add a blog post

Edit `data/posts.go`:
```go
{
  Slug:    "new-post",
  Title:   "New Blog Post",
  Content: "Content goes here",
  Date:    "2025-11-03",
}
```

Run `./build.sh generate` (or `build.bat generate` on Windows). New post appears at `/blog/new-post.html` and in blog list.

### Add a form field

Edit `templates/home.templ`, add input to form:
```html
<input type="tel" name="phone" placeholder="Phone" required/>
```

Edit `main.go`, add field to Donation struct:
```go
type Donation struct {
  ...
  Phone   string
  ...
}
```

The form parser automatically picks up the new field.

### Add a new page

Create `templates/contact.templ`:
```templ
package templates

templ Contact() {
  @Layout("Contact", contactContent())
}

templ contactContent() {
  <section>
    <h1>Contact Us</h1>
    <p>Email: contact@example.com</p>
  </section>
}
```

Edit `build.go`, add to `generateStaticSite()`:
```go
renderToFile("public/contact.html", templates.Contact())
```

Run `./build.sh generate` (or `build.bat generate` on Windows). Page appears at `/contact.html`.

### Add an API endpoint

Edit `main.go`, add handler:
```go
func getNewData(c *fiber.Ctx) error {
  data := fetchData()
  c.Set("Content-Type", "text/html; charset=utf-8")
  return templates.NewComponent(data).Render(c.Context(), c.Response().BodyWriter())
}
```

Register it:
```go
app.Get("/api/new-data", getNewData)
```

Use in HTML:
```html
<div hx-get="/api/new-data" hx-trigger="load"></div>
```

---

## Testing

### Manual browser testing

```bash
./build.sh install    # or build.bat install on Windows
./build.sh generate   # or build.bat generate on Windows
./build.sh run        # or build.bat run on Windows
```

Visit `http://localhost:3000`:
- Open DevTools (F12) → Network tab
- Click "Donate" button
- Watch POST request to `/api/donations`
- See success message swap in
- See stats/donor list update automatically

### CLI testing

```bash
# Get home page
curl http://localhost:3000

# Get stats API
curl http://localhost:3000/api/stats

# Submit donation
curl -X POST http://localhost:3000/api/donations \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=Alice&email=alice@example.com&amount=100"

# Get updated donor list
curl http://localhost:3000/api/recent-donors
```

---

## How This Implements Hypermedia Principles

**From *Hypermedia Systems*:**

| Principle | How It Works Here |
|-----------|-------------------|
| HTML is the contract | Server sends HTML with `hx-*` attributes, browser follows them |
| Server decides UI | Server chooses what HTML to send, browser just displays |
| HATEOAS | Response includes `HX-Trigger` header telling client what to fetch next |
| Progressive enhancement | Page works without htmx (form reloads), htmx makes it smooth |
| No client logic | Browser doesn't decide what to fetch—HTML tells it |
| Hypermedia controls | Links, forms, htmx attributes are the interface |

---

## Typical Flow

1. **Request comes in** → Fiber receives it
2. **Call handler** → `getStats()`, `createDonation()`, etc.
3. **Prepare data** → Calculate totals, validate inputs
4. **Render Templ** → Call `templates.Component(data)`
5. **Set headers** → `Content-Type: text/html`, optional `HX-Trigger`
6. **Send HTML** → Browser receives it
7. **htmx swaps** → Updates page (no reload)
8. **Cascade** → If `HX-Trigger` header, other elements auto-refresh

All HTML. All the time.

---

## Build & Run

**Unix/Linux/macOS:**
```bash
# Install dependencies
./build.sh install

# Generate static HTML files
./build.sh generate

# Run server (includes auto-generate)
./build.sh run

# Visit http://localhost:3000
```

**Windows:**
```cmd
REM Install dependencies
build.bat install

REM Generate static HTML files
build.bat generate

REM Run server (includes auto-generate)
build.bat run

REM Visit http://localhost:3000
```

The build scripts automatically download htmx v1.9.10 if not present in `static-vendor/`.

---

## Dependencies

### Runtime Dependencies
- `github.com/gofiber/fiber/v2` → Web server
- `github.com/a-h/templ` → Template engine

### Build-Time Dependencies
- `htmx@1.9.10` → Downloaded from unpkg.com during build if not present
  - Stored in: `static-vendor/htmx.min.js`
  - Copied to: `public/js/htmx.min.js` during generation
  - Served locally at runtime (no CDN dependency)

---

## Limitations (PoC)

- Donations stored in memory (lost on restart)
- No database
- No authentication
- No persistence

For production, replace in-memory state with database queries. Everything else stays the same.

---

## Further Reading & Resources

### Books & Articles

- **[Hypermedia Systems](https://hypermedia.systems/)** — Carson Gross, Adam Argyle, Deniz Akşimşer. The foundational book for this approach.
- **[Building Web Apps That Work Everywhere](https://www.oreilly.com/library/view/building-web-apps/9781492053903/)** — Navigating progressive enhancement strategies.

### Templ Documentation

- **[Templ Official Docs](https://templ.guide/)** — Complete reference, syntax, and best practices
- **[Templ Getting Started](https://templ.guide/syntax-and-usage/getting-started/)** — Quick start guide
- **[Templ Components](https://templ.guide/syntax-and-usage/template-composition/)** — How to compose and reuse components
- **[Templ on GitHub](https://github.com/a-h/templ)** — Source code and examples

### Fiber Documentation

- **[Fiber Official Docs](https://docs.gofiber.io/)** — Complete API reference
- **[Fiber Getting Started](https://docs.gofiber.io/guide/getting-started/)** — Quick start
- **[Fiber Static Files](https://docs.gofiber.io/api/app/#static)** — Serving static content
- **[Fiber Context](https://docs.gofiber.io/api/ctx/)** — Request/response handling
- **[Fiber on GitHub](https://github.com/gofiber/fiber)** — Source code and community

### htmx Documentation

- **[htmx Official Docs](https://htmx.org/docs/)** — Complete reference with examples
- **[htmx Attributes](https://htmx.org/reference/#attributes)** — All available attributes
- **[htmx Examples](https://htmx.org/examples/)** — Common patterns and use cases
- **[htmx Response Headers](https://htmx.org/reference/#response_headers)** — Server-side headers like `HX-Trigger`
- **[htmx on GitHub](https://github.com/bigskysoftware/htmx)** — Source code

### Go Resources

- **[Go Official Website](https://go.dev/)** — Language documentation and download
- **[A Tour of Go](https://go.dev/tour/welcome/1)** — Interactive introduction to Go
- **[Effective Go](https://go.dev/doc/effective_go)** — Best practices and idioms
- **[Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)** — Understanding goroutines and channels

### Hypermedia & REST Concepts

- **[REST Constraints (Wikipedia)](https://en.wikipedia.org/wiki/Representational_state_transfer#Architectural_constraints)** — Understanding REST principles
- **[HATEOAS Explanation](https://en.wikipedia.org/wiki/HATEOAS)** — What HATEOAS means and why it matters
- **[Roy Fielding's REST Dissertation](https://www.ics.uci.edu/~fielding/pubs/dissertation/top.htm)** — Original REST specification
- **[HTML: The Standard](https://html.spec.whatwg.org/)** — Official HTML specification

### Related Projects & Examples

- **[Fiber Examples](https://github.com/gofiber/recipes)** — Community recipes and examples
- **[htmx Examples Repository](https://github.com/bigskysoftware/htmx-examples)** — Official htmx examples
- **[HTMX + Go Examples](https://github.com/search?q=htmx+go+fiber)** — Community projects using this stack

### Community & Support

- **[Hypermedia Discussion](https://discourse.hypermedia.systems/)** — Community forum for hypermedia approaches
- **[htmx Discord](https://discord.gg/rtV2YCw6ZQ)** — htmx community chat
- **[Fiber Discord](https://discord.gg/gofiber)** — Fiber framework community
- **[Go Community](https://go.dev/wiki/GoUserGroups)** — Local and online Go communities

---

**Last Updated:** 2025-11-03 03:46:57 UTC  
**Repo:** [adabarx/campaign-example](https://github.com/adabarx/campaign-example)
