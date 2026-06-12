# TravelSphere

A full-stack destination discovery and trip planner built with the Beego framework (Go).

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.26+, Beego v2 |
| Templating | Beego SSR templates (`.tpl`) |
| Storage | In-memory wishlist store (no database) |
| External APIs | REST Countries, OpenTripMap, WeatherAPI |
| Frontend | Vanilla JS (Fetch API), CSS custom properties |

---

## Prerequisites

- Go 1.26+
- bee CLI

```bash
go install github.com/beego/bee/v2@latest
```

---

## Setup

```bash
# 1. Clone the repository
git clone https://github.com/Shahriar-Hasan123/TravelSphere-ShahriarHasan
cd TravelSphere-ShahriarHasan

# 2. Install dependencies
go mod tidy

# 3. Configure the application
cp conf/app.conf.example conf/app.conf
cp .env.example .env
# Edit .env and fill in your API keys

# 4. Run the development server
bee run
```
The application will be available at `http://localhost:8080`.

---

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `RESTCOUNTRIES_BASE_URL` | No | Defaults to `https://api.restcountries.com/countries/v5` |
| `RESTCOUNTRIES_API_KEY` | **Yes** | API key for REST Countries from [restcountries.com](https://restcountries.com/) |
| `OPENTRIPMAP_BASE_URL` | No | Defaults to `https://api.opentripmap.com/0.1/en` |
| `OPENTRIPMAP_API_KEY` | **Yes** | Get a free key at [opentripmap.org](https://dev.opentripmap.org/product) |
| `WEATHERAPI_BASE_URL` | No | Defaults to `http://api.weatherapi.com/v1` |
| `WEATHERAPI_KEY` | **Yes** | Get a free key at [weatherapi.com](https://www.weatherapi.com/) |

---

## Authentication

Session-based authentication with no user database. Enter any non-empty username on the login page to create a session.

---

## Wishlist Storage

Managed entirely in memory via a thread-safe service layer (`sync.RWMutex`, `sync.Once`). Data persists for the lifetime of the server process and resets on restart. No database, SQLite, or ORM is used anywhere in this project.

---

## URL Slug Format

| Country | Slug | URL |
| Albania | `albania` | `/countries/albania` |
| United States | `united-states` | `/countries/united-states` |
| Bangladesh | `bangladesh` | `/countries/bangladesh` |

---

## Project Structure

```
TravelSphere/
├── controllers/        # SSR page controllers and JSON API controllers
│   └── api/            # /api/* JSON endpoints
├── filters/            # Logging and authentication middleware
├── models/             # Domain entities and DTOs
├── routers/            # Route registration — SSR and API separated
├── services/           # All business logic
├── utils/              # Formatters, validators, response helpers
│   └── clients/        # HTTP clients for external APIs
├── views/              # Beego .tpl templates
├── static/             # CSS and JavaScript
│   ├── css/
│   └── js/
├── conf/               # Beego configuration
└── tests/              # Unit test files
```

---

## Pages

| Route | Auth Required | Description |
|---|---|---|
| `GET /` | No | Home page — featured destinations and popular attractions |
| `GET /countries` | No | Country Explorer — search and region filter |
| `GET /countries/:slug` | No | Destination detail — attractions and weather |
| `GET /wishlist` | **Yes** | Travel wishlist — add, edit, delete entries |
| `GET /dashboard` | **Yes** | Dashboard — saved trip stats and destination list |
| `GET /login` | No | Login page |
| `GET /logout` | No | Clears session and redirects to home |

---

## API Endpoints

### Countries

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/countries` | List all countries. Supports `search` and `region` query params |
| `GET` | `/api/countries/:slug` | Single country detail by slug |
| `GET` | `/api/countries/suggestions?q=` | Autocomplete suggestions for home page search |

### Wishlist

> All wishlist endpoints require an authenticated session.

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/wishlist` | Get all wishlist entries for the authenticated user |
| `POST` | `/api/wishlist` | Create a new wishlist entry |
| `PUT` | `/api/wishlist/:id` | Update note and status — returns the updated item |
| `DELETE` | `/api/wishlist/:id` | Delete a wishlist entry — returns `204 No Content` |

### Dashboard

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/dashboard/summary` | Returns `total`, `planned`, and `visited` counts |

### Attractions

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/attractions?lat=&lon=` | Attractions near given coordinates |

---

## Request & Response Format

All `/api/*` endpoints return JSON.

**Success**
```json
{
  "status": "ok",
  "data": {}
}
```

**Created**
```json
{
  "status": "created",
  "data": {}
}
```

**Error**
```json
{
  "status": "error",
  "message": "description of the error",
  "code": 400
}
```

---

## Running Tests

```bash
go test ./...
```

With coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
```

Coverage targets: `utils` 100% · `utils/clients` 95%+ · `services` 95%+ · `filters` 95%+

---

## Docker (Optional)

**Prerequisites:** Install [Docker Desktop](https://www.docker.com/products/docker-desktop/) — this includes both Docker and docker-compose.

---

### First Time Setup

**Step 1 — Create your environment file**

Copy the example environment file and fill in your API keys:

```bash
cp .env.example .env
```

Open `.env` and set your API keys:

```
OPENTRIPMAP_API_KEY=your_opentripmap_key_here
WEATHERAPI_KEY=your_weatherapi_key_here
```

**Step 2 — Build the Docker image**

```bash
docker build -t travelsphere:latest .
```

This may take a minute the first time while it downloads dependencies.

**Step 3 — Start the application**

```bash
docker-compose up -d
```

**Step 4 — Open in your browser**

```
http://localhost:8080
```

---

### Starting the App (After First Setup)

Every time you want to run the app after the first setup:

```bash
docker-compose up -d
```

Then open `http://localhost:8080` in your browser.

---

### Stopping the App

```bash
docker-compose down
```

---

### Restarting the App

```bash
docker-compose restart
```

---

### Checking Logs (If Something Goes Wrong)

```bash
docker-compose logs -f
```

Press `Ctrl + C` to stop watching logs.

---

### Rebuilding After Code Changes

If you change the source code and want to rebuild the image:

```bash
docker-compose down
docker build -t travelsphere:latest .
docker-compose up -d
```

---

### Standalone Run (Without docker-compose)

```bash
docker run -p 8080:8080 \
  -e OPENTRIPMAP_API_KEY=your_key \
  -e WEATHERAPI_KEY=your_key \
  travelsphere:latest
```

Then open `http://localhost:8080` in your browser.

---

### Notes

- `OPENTRIPMAP_API_KEY` is required for attractions to load on destination pages.
- `WEATHERAPI_KEY` is optional — weather panel shows a placeholder if not set.
- Never commit your `.env` file — it is listed in `.gitignore` and `.dockerignore`.
---

## Git Branching Strategy

```
main
 └── dev
      ├── feature-1/base-mvc-structure
      ├── feature-2/country-explorer
      ├── feature-3/destination-detail
      ├── feature-4/auth-filters
      ├── feature-5/wishlist
      ├── feature-6/dashboard-home
      ├── feature-7/username-only-auth
      ├── feature-8/rest-compliant-wishlist-api
      └── feature-9/unit-tests
```

---

## External API References

- [REST Countries](https://restcountries.com/) — country data, flags, languages, currencies
- [OpenTripMap](https://dev.opentripmap.org/product) — tourist attractions and landmarks
- [WeatherAPI](https://www.weatherapi.com/) — current weather and forecast (optional)