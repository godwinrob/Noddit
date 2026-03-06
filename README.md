# Noddit - Modernized Message Board

A Reddit-inspired message board application, fully modernized with **Go backend** and **Next.js frontend**.

## 🎉 What's New (2026 Modernization)

- ✅ **Go Backend** - Replaced Java/Spring with modern Go + Gin framework
- ✅ **Next.js Frontend** - Replaced Vue 2 with Next.js 15 + React 19
- ✅ **TypeScript** - Full type safety throughout
- ✅ **Modern Stack** - Latest tools and best practices
- ✅ **API Compatible** - Drop-in replacement for original Java backend

## Original Project (2019)

Tech Elevator Capstone Project by:
- Rob Godwin
- Craig Samad
- Jay Minihan
- Emma Knutson

## Tech Stack

### Backend (Go)
- **Go 1.25+** - Primary language
- **Gin** - HTTP web framework
- **PostgreSQL** - Database
- **JWT** - Authentication (HS512)
- **PBKDF2** - Password hashing

### Frontend (Next.js)
- **Next.js 15** - React framework
- **React 19** - UI library
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling

### Database
- **PostgreSQL 12+**

## Quick Start

### ⚡ Fastest: Using Tilt (Recommended)

```bash
tilt up
```

**That's it!** Frontend at http://localhost:8081, backend at http://localhost:8080.

Press SPACE for Tilt UI, `q` to quit.

### 🐳 Alternative: Docker Compose

```bash
docker-compose up --build
```

### 🔧 Manual Setup

See [QUICK_START.md](./QUICK_START.md) for detailed instructions.

## Project Structure

```
noddit/
├── backend-go/              # Modern Go backend
│   ├── cmd/api/             # Application entry point
│   ├── internal/            # Internal packages
│   ├── pkg/                 # Public packages
│   ├── migrations/          # Database migrations
│   └── README.md
│
├── noddit-next/             # Modern Next.js frontend
│   ├── app/                 # Next.js App Router pages
│   ├── components/          # React components
│   ├── lib/                 # Utilities and API client
│   └── README.md
│
├── backend/                 # Legacy Java backend (archived)
├── frontend/                # Legacy Vue frontend (archived)
└── README.md               # This file
```

## Features

### Implemented ✅
- User authentication (register/login)
- JWT token management
- Browse recent posts
- View active communities
- Responsive design
- Dark mode UI

### Coming Soon 🚧
- Create posts and comments
- Voting system (upvote/downvote)
- User profiles
- Subnoddit creation
- Search functionality
- Favorites system

## API Endpoints

The Go backend provides 29 REST API endpoints:

### Authentication
- `POST /login` - User login
- `POST /register` - User registration

### Posts
- `GET /api/public/allposts` - Get all posts
- `GET /api/public/recentposts` - Get recent posts
- `GET /api/public/popularposts` - Get popular posts
- `POST /api/post/create` - Create post (auth)
- And more...

### Communities (Subnoddits)
- `GET /api/public/subnoddits` - Get all communities
- `GET /api/public/subnoddits/active` - Get active communities
- `POST /api/subnoddits/create` - Create community (auth)
- And more...

See `backend-go/README.md` for complete API documentation.

## Environment Configuration

### Backend (.env)
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=userdb
JWT_SECRET=your_base64_secret
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Migration from Legacy Stack

### From Java/Spring to Go
- All 29 endpoints ported with exact compatibility
- Same JWT authentication mechanism
- Same password hashing (PBKDF2)
- Same database schema
- Zero breaking changes for clients

### From Vue 2 to Next.js
- Modern React 19 with Server Components
- Better TypeScript support
- Improved performance
- SEO-friendly
- Simpler state management

## Development

### With Tilt (Hot Reload)
```bash
tilt up
# Edit files - Tilt auto-rebuilds and reloads
```

### Native Development
See [DEVELOPMENT.md](./DEVELOPMENT.md) for:
- Running services natively
- Database management
- Testing workflows
- Troubleshooting

### Docker Compose
```bash
docker-compose up --build
```

## Testing the API

Use the Go backend with:
- **Frontend**: http://localhost:8081
- **Postman/curl**: http://localhost:8080
- **Browser**: Navigate to API endpoints directly

Example:
```bash
# Get recent posts
curl http://localhost:8080/api/public/recentposts

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

## Roadmap

- [x] Modernize backend (Go)
- [x] Set up database migrations
- [x] Create Next.js frontend foundation
- [x] Implement authentication
- [x] Add Docker + Tilt setup
- [ ] Complete all frontend pages
- [ ] Add tests
- [ ] Deploy to cloud (Vercel + Fly.io)
- [ ] Performance optimization

## Documentation

- **[QUICK_START.md](./QUICK_START.md)** - Get running in 60 seconds
- **[DEVELOPMENT.md](./DEVELOPMENT.md)** - Detailed development guide
- **[backend-go/README.md](./backend-go/README.md)** - Backend API docs
- **[noddit-next/README.md](./noddit-next/README.md)** - Frontend docs

## Contributing

This is an educational/portfolio project. Feel free to fork and experiment!

## License

Educational project - Tech Elevator Capstone 2019, Modernized 2026
