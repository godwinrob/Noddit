# Quick Start Guide

Get Noddit running in 60 seconds with Tilt, or 5 minutes manually.

## ⚡ Fastest: Using Tilt

### 1. Install Tilt

**macOS:**
```bash
brew install tilt-dev/tap/tilt
```

**Windows (with Scoop):**
```powershell
scoop bucket add tilt-dev https://github.com/tilt-dev/scoop-bucket
scoop install tilt
```

### 2. Start Everything

```bash
tilt up
```

**That's it!**

- Frontend: http://localhost:8081
- Backend: http://localhost:8080
- Press SPACE for Tilt UI

### 3. Stop

Press `q` or run:
```bash
tilt down
```

---

## 🐳 Alternative: Docker Compose

### 1. Start All Services

```bash
docker-compose up --build
```

### 2. Access

- Frontend: http://localhost:8081
- Backend: http://localhost:8080

### 3. Stop

```bash
docker-compose down
```

---

## 🔧 Manual Setup (No Docker)

### 1. Database

```bash
cd backend-go/scripts
./setup_db.bat  # Windows
# or
./setup_db.sh   # Linux/Mac

# Optional: Load sample data (posts, users, communities)
./seed_db.bat   # Windows
# or
./seed_db.sh    # Linux/Mac
```

### 2. Backend (Terminal 1)

```bash
cd backend-go
go run cmd/api/main.go
```

### 3. Frontend (Terminal 2)

```bash
cd noddit-next
npm run dev
```

### 4. Access

- Frontend: http://localhost:8081
- Backend: http://localhost:8080

---

## 📝 First Steps

1. Visit http://localhost:8081
2. Click "Sign Up" to create an account
3. Browse posts and communities
4. Try creating a post (backend ready, UI coming soon!)

---

## 🆘 Troubleshooting

**Port 8080 in use:**
```bash
lsof -i :8080  # Find process
kill -9 <PID>  # Kill it
```

**Database issues:**
```bash
docker-compose down -v  # Remove volumes
docker-compose up postgres  # Restart DB
```

**Tilt issues:**
```bash
tilt down
docker-compose down -v
tilt up
```

---

## 📚 Next Steps

- [DEVELOPMENT.md](./DEVELOPMENT.md) - Detailed development guide
- [backend-go/README.md](./backend-go/README.md) - Backend documentation
- [noddit-next/README.md](./noddit-next/README.md) - Frontend documentation
