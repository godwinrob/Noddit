# Development Guide

This guide covers local development workflows for Noddit.

## Prerequisites

Choose one of these setups:

### Option 1: Docker + Tilt (Recommended)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- [Tilt](https://docs.tilt.dev/install.html)

### Option 2: Native Development
- Go 1.25+
- Node.js 18+
- PostgreSQL 12+

---

## Quick Start with Tilt (Recommended)

Tilt provides the best local development experience with hot reloading and service orchestration.

### 1. Install Tilt

**macOS:**
```bash
brew install tilt-dev/tap/tilt
```

**Windows:**
```powershell
scoop bucket add tilt-dev https://github.com/tilt-dev/scoop-bucket
scoop install tilt
```

**Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
```

### 2. Start Everything

```bash
tilt up
```

This will:
- Start PostgreSQL database
- Build and run Go backend
- Build and run Next.js frontend
- Watch for file changes and auto-reload

**Press SPACE** to open the Tilt UI in your browser.

### 3. Access the Application

- **Frontend:** http://localhost:8081
- **Backend API:** http://localhost:8080
- **Tilt UI:** http://localhost:10350

### 4. Stop Everything

Press `q` in the Tilt terminal or:

```bash
tilt down
```

---

## Docker Compose (Without Tilt)

If you prefer Docker Compose without Tilt:

### Start All Services
```bash
docker-compose up --build
```

### Start Specific Services
```bash
# Backend + Database only
docker-compose up postgres backend

# Frontend only (requires backend)
docker-compose up frontend
```

### Stop Services
```bash
docker-compose down
```

### Clean Up (Remove volumes)
```bash
docker-compose down -v
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
```

---

## Native Development (No Docker)

For faster iteration during active development:

### 1. Start PostgreSQL

**Using Docker:**
```bash
docker run --name noddit-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres1 \
  -e POSTGRES_DB=userdb \
  -p 5432:5432 \
  -d postgres:16-alpine
```

**Or use local PostgreSQL:**
```bash
createdb -U postgres userdb
cd backend-go/scripts
./setup_db.bat  # Windows
# or
./setup_db.sh   # Linux/Mac
```

### 2. Start Go Backend

```bash
cd backend-go
go run cmd/api/main.go
```

Backend runs on http://localhost:8080

### 3. Start Next.js Frontend

```bash
cd noddit-next
npm run dev
```

Frontend runs on http://localhost:8081

---

## Development Workflows

### Backend Development

#### Hot Reload with Tilt
Tilt automatically rebuilds when you change `.go` files:

```bash
tilt up
# Edit files in backend-go/
# Tilt detects changes and rebuilds
```

#### Manual Testing
```bash
# Run backend locally
cd backend-go
go run cmd/api/main.go

# Test endpoints
curl http://localhost:8080/api/public/recentposts
```

#### Run Tests
```bash
cd backend-go
go test ./...
```

### Frontend Development

#### Hot Reload with Tilt
Tilt automatically rebuilds Next.js:

```bash
tilt up
# Edit files in noddit-next/
# Tilt detects changes and rebuilds
```

#### Fast Refresh (Native)
For fastest iteration, run Next.js natively:

```bash
cd noddit-next
npm run dev
# Fast Refresh enabled - instant updates
```

#### Build for Production
```bash
cd noddit-next
npm run build
npm start
```

### Database Management

#### Apply Migrations
```bash
cd backend-go/scripts
./setup_db.bat  # Windows
./setup_db.sh   # Linux/Mac
```

#### Connect to Database
```bash
# If using Docker Compose
docker-compose exec postgres psql -U postgres -d userdb

# If using local PostgreSQL
psql -U postgres -d userdb
```

#### Load Sample Data
```bash
psql -U postgres -d userdb -f backend/database/dbexport.pgsql
```

#### Reset Database
```bash
# Drop and recreate
docker-compose down -v
docker-compose up postgres
cd backend-go/scripts
./setup_db.bat
```

---

## Tilt Configuration

### Run Specific Services

Create `tilt_config.json`:

```json
{
  "to-run": ["backend"]
}
```

Or use command line:

```bash
# Backend only
tilt up -- --to-run=backend

# Frontend only
tilt up -- --to-run=frontend

# Both (default)
tilt up
```

### Enable Local Development Mode

Edit `Tiltfile` and uncomment the `local_resource` sections for faster iteration:

```python
local_resource(
    'backend-dev',
    serve_cmd='cd backend-go && go run cmd/api/main.go',
    deps=['./backend-go'],
    labels=['backend'],
    resource_deps=['postgres'],
)
```

This runs Go/Node natively instead of in containers (faster rebuilds).

---

## Environment Variables

### Backend (.env)
```env
PORT=8080
FRONTEND_URL=http://localhost:8081
DB_HOST=localhost  # or 'postgres' in Docker
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres1
DB_NAME=userdb
JWT_SECRET=your_base64_secret
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## Troubleshooting

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080  # Mac/Linux
netstat -ano | findstr :8080  # Windows

# Kill the process or change ports in .env
```

### Database Connection Failed
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Restart database
docker-compose restart postgres
```

### Tilt Issues
```bash
# Clean everything and restart
tilt down
docker-compose down -v
tilt up
```

### Go Build Errors
```bash
cd backend-go
go mod tidy
go clean -cache
```

### Next.js Build Errors
```bash
cd noddit-next
rm -rf .next node_modules
npm install
npm run build
```

---

## Performance Tips

1. **Use Native Mode for Active Development**
   - Run Postgres in Docker
   - Run Go/Next.js natively for fastest hot reload

2. **Use Tilt for Full-Stack Testing**
   - Test all services together
   - Verify Docker builds work

3. **Use Docker Compose for Production-Like Testing**
   - Test with exact production configurations
   - Verify environment variables

---

## Next Steps

- Read [CONTRIBUTING.md](./CONTRIBUTING.md) for contribution guidelines
- Check [backend-go/README.md](./backend-go/README.md) for backend specifics
- Check [noddit-next/README.md](./noddit-next/README.md) for frontend specifics
