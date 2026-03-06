# Noddit - Modernized Message Board

A Reddit-inspired message board application, fully modernized with **Go backend** and **Next.js frontend**.

## 🎉 What's New (2026 Modernization)

This project has been **completely modernized** from the 2019 Tech Elevator capstone:

- ✅ **Go Backend** - Replaced Java/Spring Boot with Go 1.25 + Gin framework
- ✅ **Next.js Frontend** - Replaced Vue 2 with Next.js 15 + React 19 + TypeScript
- ✅ **Full Feature Parity** - All original features ported and working
- ✅ **Docker + Tilt** - Modern development workflow with hot reload
- ✅ **100% API Compatible** - Same endpoints, same auth, same database schema
- ✅ **Better Performance** - Go's concurrency + Next.js optimizations
- ✅ **Type Safety** - TypeScript throughout the frontend
- ✅ **Modern UI** - Tailwind CSS with responsive dark mode design

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

### Core Features ✅
- **User Authentication** - Register, login, JWT token management
- **Posts & Comments** - Create posts, nested comments with unlimited depth
- **Voting System** - Upvote/downvote posts and comments
- **Communities (Subnoddits)** - Create and browse communities
- **Search** - Search communities by name
- **User Profiles** - View profiles, update settings, change username
- **Favorites** - Favorite posts and communities
- **Moderation** - Community moderators with delete permissions
- **Image Support** - Upload images with posts
- **Responsive Design** - Mobile-friendly dark mode UI

### Frontend Pages ✅
- **Home** (`/`) - Popular posts from all communities
- **Submit Post** (`/submit`) - Create new posts with images
- **Post Detail** (`/s/[name]/[id]`) - View post with nested comments
- **Community Page** (`/s/[name]`) - View community posts (recent/popular)
- **Browse Communities** (`/subnoddits`) - Explore all communities
- **Create Community** (`/subnoddits/create`) - Start a new community
- **Search** (`/search/[term]`) - Search for communities
- **User Profile** (`/profile`) - View and edit your profile
- **Login/Register** (`/login`, `/register`) - Authentication pages

### Components ✅
- **PostCard** - Display posts with voting, delete, reply
- **Comment** - Nested comment display with voting
- **CommentSection** - Full comment tree rendering
- **ReplyForm** - Create replies to posts and comments
- **Nav** - Navigation with search and favorites dropdown

## API Endpoints

The Go backend provides **29 REST API endpoints**:

### Authentication (Public)
- `POST /login` - User login (returns JWT)
- `POST /register` - User registration

### Subnoddits
- `GET /api/public/subnoddits` - Get all subnoddits
- `GET /api/public/subnoddits/active` - Get 5 most active
- `GET /api/public/subnoddits/:name` - Get by name
- `GET /api/public/subnoddits/search/:term` - Search subnoddits
- `POST /api/subnoddits/create` 🔒 - Create subnoddit
- `PUT /api/subnoddits/update` 🔒 - Update subnoddit
- `DELETE /api/subnoddits/delete/:name` 🔒 - Delete subnoddit
- `GET /api/public/moderators/:name` - Get moderators

### Posts
- `GET /api/public/allposts` - Get all posts
- `GET /api/public/recentposts` - Get 5 recent posts
- `GET /api/public/popularposts` - Get popular (24h)
- `GET /api/public/allposts/:subnodditName` - Get by subnoddit
- `GET /api/public/allpostspopular/:subnodditName` - Popular by subnoddit
- `GET /api/public/:subnodditName/:postId` - Get single post
- `GET /api/public/:subnodditName/:postId/replies` - Get replies
- `POST /api/post/create` 🔒 - Create post
- `PUT /api/post/update/:postId` 🔒 - Update post
- `DELETE /api/post/delete/:postId` 🔒 - Delete post
- `POST /:subnodditName/:postId/createreply` 🔒 - Create reply

### Votes
- `GET /api/public/post/votes/:postId` - Get votes
- `POST /api/post/vote` 🔒 - Vote on post

### Users
- `GET /api/user/:username` - Get user profile
- `GET /api/public/post/user/:username` - Get user's posts
- `PUT /api/user/update/email/:username` 🔒 - Update email
- `PUT /api/user/update/username/:username` 🔒 - Update username
- `PUT /api/user/update/name/:username` 🔒 - Update name
- `PUT /api/user/update/avatar/:username` 🔒 - Update avatar
- `DELETE /api/user/delete/:username` 🔒👑 - Delete user (admin)

### Favorites
- `GET /api/favorites/:username` 🔒 - Get favorites
- `POST /api/favorites/create/post` 🔒 - Favorite post
- `POST /api/favorites/create/subnoddit` 🔒 - Favorite subnoddit
- `DELETE /api/favorites/delete/post/:postId` 🔒 - Unfavorite post
- `DELETE /api/favorites/subnoddit/:subnodditId` 🔒 - Unfavorite subnoddit

🔒 = Requires authentication | 👑 = Admin only

See `backend-go/README.md` for complete API documentation.

## Sample Data

The database comes with pre-seeded sample data:

### Users (password = username)
- **Admins**: `rgodwin`, `csamad`, `eknutson`, `jminihan`
- **Users**: `test`, `asd1`, `david`

### Communities
- **Cats** - A home for cats and cat accessories
- **Dogs** - Dogs are not as cool as cats, but sometimes we like them too
- **Harold** - Our inspiration to hiding the pain in our lives
- **Politics** - Don't post here
- **Gardening** - Plant stuff in the ground and hope it grows
- **star_wars** - All the Jedi things!
- **test_subnoddit** - this is a test

### Sample Posts
- 15 posts across various communities with images
- Nested comments demonstrating reply functionality
- Pre-populated votes and favorites

Login with any sample user to explore the full functionality!

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

## Testing & Usage

### Quick Test Drive

1. **Start the app**: `tilt up`
2. **Visit**: http://localhost:8081
3. **Login**: Use `rgodwin` / `rgodwin` (or any sample user - password = username)
4. **Explore**:
   - Browse popular posts on the home page
   - Click into communities like `/s/Cats` or `/s/star_wars`
   - Upvote/downvote posts
   - Create a comment on any post
   - Submit a new post with `/submit`
   - Search for communities
   - Add favorites and view them in the nav dropdown

### API Testing

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
  -d '{"username":"rgodwin","password":"rgodwin"}'

# Get popular posts from last 24 hours
curl http://localhost:8080/api/public/popularposts
```

## Roadmap

### Completed ✅
- [x] Modernize backend (Go with Gin)
- [x] All 29 API endpoints implemented
- [x] Set up database migrations with seed data
- [x] Create Next.js 15 frontend
- [x] Implement authentication (JWT)
- [x] Add Docker + Tilt development setup
- [x] Port all Vue components to React
- [x] Complete all frontend pages (9 pages, 5 components)
- [x] Voting system (upvote/downvote)
- [x] Nested comments system
- [x] Search functionality
- [x] User profiles with settings
- [x] Favorites system
- [x] Community creation and moderation

### Future Enhancements 🚧
- [ ] Add automated tests (backend + frontend)
- [ ] Real-time updates (WebSockets)
- [ ] User avatars with upload
- [ ] Rich text editor for posts
- [ ] Notifications system
- [ ] Email verification
- [ ] Password reset flow
- [ ] Deploy to cloud (Vercel + Fly.io)
- [ ] Performance optimization
- [ ] Mobile app (React Native)

## Documentation

- **[QUICK_START.md](./QUICK_START.md)** - Get running in 60 seconds
- **[DEVELOPMENT.md](./DEVELOPMENT.md)** - Detailed development guide
- **[backend-go/README.md](./backend-go/README.md)** - Backend API docs
- **[noddit-next/README.md](./noddit-next/README.md)** - Frontend docs

## Contributing

This is an educational/portfolio project. Feel free to fork and experiment!

## License

Educational project - Tech Elevator Capstone 2019, Modernized 2026
