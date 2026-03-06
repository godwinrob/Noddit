# Noddit Go Backend

Modern Go implementation of the Noddit message board backend API.

## Tech Stack

- **Go 1.25+** - Backend language
- **Gin** - HTTP web framework
- **PostgreSQL** - Database
- **JWT** - Authentication
- **PBKDF2** - Password hashing

## Features

- RESTful API with 29 endpoints
- JWT-based authentication (6-hour expiration)
- PBKDF2 password hashing with salts
- Nested comment system
- Vote/scoring system
- User favorites
- Subnoddit (subreddit-like) communities
- Role-based access control (user, admin, super_admin)

## Project Structure

```
backend-go/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── database/
│   │   └── db.go                # Database connection
│   ├── handlers/
│   │   ├── auth.go              # Login/Register
│   │   ├── post.go              # Post CRUD + replies
│   │   ├── subnoddit.go         # Subnoddit CRUD
│   │   ├── vote.go              # Voting system
│   │   ├── user.go              # User management
│   │   └── favorites.go         # Favorites system
│   ├── middleware/
│   │   └── auth.go              # JWT middleware
│   └── models/
│       └── models.go            # Data models
├── pkg/
│   └── auth/
│       ├── jwt.go               # JWT token handling
│       └── password.go          # Password hashing
└── migrations/                  # Database migrations

```

## Setup

### Prerequisites

- Go 1.25 or higher
- PostgreSQL 12+
- Make (optional)

### Database Setup

1. Create database and apply schema:
```bash
cd scripts
./setup_db.bat  # Windows
# or
./setup_db.sh   # Linux/Mac
```

2. (Optional) Load sample data with example posts:
```bash
./seed_db.bat  # Windows
# or
./seed_db.sh   # Linux/Mac
```

This creates 7 sample users and posts about cats, dogs, Harold, and Star Wars!

### Environment Configuration

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Update `.env` with your configuration:
```env
PORT=8080
FRONTEND_URL=http://localhost:8081

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres1
DB_NAME=userdb
DB_SSLMODE=disable

JWT_SECRET=<your-base64-secret>
JWT_EXPIRATION_HOURS=6
```

### Running the Application

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

### Build for Production

```bash
go build -o noddit-api cmd/api/main.go
./noddit-api
```

## API Endpoints

### Authentication (Public)
- `POST /login` - User login (returns JWT token)
- `POST /register` - User registration

### Subnoddits
- `GET /api/public/subnoddits` - Get all subnoddits
- `GET /api/public/subnoddits/active` - Get 5 most active
- `GET /api/public/subnoddits/:name` - Get by name
- `GET /api/public/subnoddits/search/:term` - Search subnoddits
- `POST /api/subnoddits/create` - Create (auth required)
- `PUT /api/subnoddits/update` - Update (auth required)
- `DELETE /api/subnoddits/delete/:name` - Delete (auth required)
- `GET /api/public/moderators/:name` - Get moderators

### Posts
- `GET /api/public/allposts` - Get all posts
- `GET /api/public/recentposts` - Get 5 recent posts
- `GET /api/public/popularposts` - Get popular (24h)
- `GET /api/public/allposts/:subnodditName` - Get by subnoddit
- `GET /api/public/allpostspopular/:subnodditName` - Popular by subnoddit
- `GET /api/public/:subnodditName/:postId` - Get single post
- `GET /api/public/:subnodditName/:postId/replies` - Get replies
- `POST /api/post/create` - Create post (auth required)
- `PUT /api/post/update/:postId` - Update post (auth required)
- `DELETE /api/post/delete/:postId` - Delete post (auth required)
- `POST /:subnodditName/:postId/createreply` - Create reply (auth required)

### Votes
- `GET /api/public/post/votes/:postId` - Get votes
- `POST /api/post/vote` - Vote on post (auth required)

### Users
- `GET /api/user/:username` - Get user profile
- `GET /api/public/post/user/:username` - Get user's posts
- `PUT /api/user/update/email/:username` - Update email (auth required)
- `PUT /api/user/update/username/:username` - Update username (auth required)
- `PUT /api/user/update/name/:username` - Update name (auth required)
- `PUT /api/user/update/avatar/:username` - Update avatar (auth required)
- `DELETE /api/user/delete/:username` - Delete user (admin only)

### Favorites
- `GET /api/favorites/:username` - Get favorites (auth required)
- `POST /api/favorites/create/post` - Favorite post (auth required)
- `POST /api/favorites/create/subnoddit` - Favorite subnoddit (auth required)
- `DELETE /api/favorites/delete/post/:postId` - Unfavorite post (auth required)
- `DELETE /api/favorites/subnoddit/:subnodditId` - Unfavorite subnoddit (auth required)

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

Tokens are obtained via `/login` and are valid for 6 hours.

## Development

### Code Organization

- **cmd/api** - Application entry point
- **internal** - Private application code
- **pkg** - Public, reusable packages

### Adding New Endpoints

1. Define the handler in `internal/handlers/`
2. Add the route in `cmd/api/main.go`
3. Add authentication middleware if needed

## Migration from Java Backend

This Go backend is a complete drop-in replacement for the Java/Spring backend. All API endpoints match exactly, including:

- Same URL paths
- Same request/response formats
- Same authentication mechanism (JWT with HS512)
- Same password hashing (PBKDF2 with 100k iterations)
- Compatible with existing Vue.js frontend
- Uses same PostgreSQL database schema

## License

Created as Tech Elevator capstone project (2019)
Modernized with Go backend (2026)
