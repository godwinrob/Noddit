# Noddit Next.js Frontend

Modern Next.js 15 frontend for the Noddit message board application.

## Tech Stack

- **Next.js 15** - React framework with App Router
- **React 19** - UI library
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **Server Components** - Performance optimization

## Features

- ✅ User authentication (login/register)
- ✅ Browse recent posts
- ✅ View active communities
- ✅ Responsive dark mode design
- ✅ JWT token management
- ✅ Type-safe API client
- 🚧 Post creation (coming soon)
- 🚧 Commenting system (coming soon)
- 🚧 Voting system (coming soon)
- 🚧 User profiles (coming soon)

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn
- Running Go backend on port 8080

### Installation

```bash
# Install dependencies
npm install

# Copy environment file
cp .env.local.example .env.local

# Update .env.local if needed
# NEXT_PUBLIC_API_URL=http://localhost:8080
```

### Development

```bash
# Run development server (port 8081)
npm run dev
```

Open [http://localhost:8081](http://localhost:8081) in your browser.

### Production Build

```bash
# Build for production
npm run build

# Start production server
npm start
```

## Project Structure

```
noddit-next/
├── app/
│   ├── layout.tsx           # Root layout with auth provider
│   ├── page.tsx             # Home page (recent posts)
│   ├── login/page.tsx       # Login page
│   ├── register/page.tsx    # Registration page
│   └── globals.css          # Global styles
├── components/
│   └── Nav.tsx              # Navigation bar
├── lib/
│   ├── api.ts               # API client
│   └── auth-context.tsx     # Authentication context
└── public/                  # Static assets
```

## API Integration

The frontend communicates with the Go backend via REST API:

- Authentication endpoints: `/login`, `/register`
- Public endpoints: `/api/public/*`
- Protected endpoints: `/api/*` (require JWT token)

## Environment Variables

- `NEXT_PUBLIC_API_URL` - Backend API URL (default: http://localhost:8080)

## Development Roadmap

### Phase 1: Core Features (Current)
- [x] Project setup
- [x] Authentication
- [x] Home page with recent posts
- [x] Navigation

### Phase 2: Content Creation
- [ ] Create post page
- [ ] Create subnoddit page
- [ ] Image upload support

### Phase 3: Interaction
- [ ] Post detail page
- [ ] Comment system
- [ ] Voting (upvote/downvote)
- [ ] Reply to comments

### Phase 4: User Features
- [ ] User profiles
- [ ] Edit profile
- [ ] Favorites
- [ ] User posts history

### Phase 5: Discovery
- [ ] Browse all subnoddits
- [ ] Search functionality
- [ ] Popular posts page
- [ ] Sorting options

### Phase 6: Polish
- [ ] Loading states
- [ ] Error handling
- [ ] Toast notifications
- [ ] Pagination
- [ ] Image previews
- [ ] Markdown support

## Migration from Vue

This Next.js app replaces the old Vue 2 frontend. Key improvements:

- Modern React 19 with Server Components
- Better TypeScript support
- Improved performance with Next.js 15
- Cleaner, more maintainable code
- Better SEO capabilities
- Simpler state management

## Contributing

When adding new features:
1. Create components in `components/`
2. Create pages in `app/`
3. Add API calls to `lib/api.ts`
4. Update types as needed

## License

Part of the Noddit project - Tech Elevator capstone (2019)
Modernized with Next.js (2026)
