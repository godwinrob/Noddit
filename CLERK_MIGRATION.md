# Clerk Migration Guide

## ✅ Completed

### Frontend (Next.js)
1. ✅ Installed `@clerk/nextjs` SDK
2. ✅ Created `middleware.ts` with `clerkMiddleware()`
3. ✅ Updated `app/layout.tsx` with `<ClerkProvider>`
4. ✅ Updated Nav component to use Clerk components:
   - `useUser()` hook for auth state
   - `<SignInButton>` and `<SignUpButton>` for login/signup
   - `<UserButton>` for user profile dropdown
5. ✅ Updated API client (`lib/api.ts`) to use Clerk tokens
6. ✅ Created `ClerkTokenProvider` to sync Clerk tokens with API client
7. ✅ Removed old custom auth files (login/register pages, auth-context)

### Backend (Go)
1. ✅ Installed Clerk Go SDK (`github.com/clerk/clerk-sdk-go/v2`)
2. ✅ Created `pkg/auth/clerk.go` for Clerk JWT validation
3. ✅ Updated middleware to validate Clerk tokens
4. ✅ Updated `.env` with Clerk configuration

## 🚀 How to Use (Keyless Mode)

### Step 1: Start the App
```bash
# Terminal 1: Start Next.js frontend
cd noddit-next
npm run dev

# Terminal 2: Start Go backend
cd backend-go
go run cmd/api/main.go
```

### Step 2: Test Auth Flow
1. Open http://localhost:8081
2. Click "Sign Up" button in the nav
3. Create a test account (Clerk handles this automatically)
4. You'll see the `<UserButton>` appear with your profile icon

### Step 3: Claim Your Application (Optional)
When ready to deploy or access the Clerk Dashboard:
1. You'll see a "Configure your application" prompt in bottom-right
2. Click it to sign up/login to Clerk
3. This will link your app to your Clerk account
4. You'll get a proper `CLERK_SECRET_KEY` to add to `.env`

## 🔧 Next Steps (TODO)

### 1. User Synchronization
Currently, Clerk users are not synced with your local PostgreSQL `users` table. You need to:

**Option A: Webhook-Based Sync (Recommended)**
- Create a webhook endpoint in Go to listen for Clerk user events
- When a user signs up in Clerk, create a corresponding record in your `users` table
- Map Clerk `user_id` to your local user records

**Option B: On-Demand Sync**
- When a user makes their first API request, check if they exist in your DB
- If not, create them automatically
- Store Clerk's `user_id` as the primary identifier

### 2. Database Schema Update
Add a `clerk_id` column to link Clerk users with your local users:
```sql
ALTER TABLE users ADD COLUMN clerk_id VARCHAR(255) UNIQUE;
ALTER TABLE users ADD COLUMN clerk_username VARCHAR(255);
```

### 3. Update Backend Handlers
Currently, handlers expect `username` from context. Update them to:
1. Get Clerk user ID from middleware context
2. Query your database using `clerk_id` instead of `username`
3. Handle new user creation when they don't exist yet

### 4. Clerk Metadata for Roles
Clerk doesn't have built-in roles. For your `admin`/`super_admin` system:
- Use Clerk's **Public Metadata** to store user roles
- Update middleware to read role from Clerk metadata
- Set metadata via Clerk Dashboard or API

### 5. Environment Variables
When you claim your app, add to `.env`:
```bash
CLERK_SECRET_KEY=sk_live_your_actual_key_here
```

And to `noddit-next/.env.local`:
```bash
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_live_your_key_here
CLERK_SECRET_KEY=sk_live_your_key_here
```

## 📚 Key Differences from Old Auth

| Old (Self-Managed JWT) | New (Clerk) |
|------------------------|-------------|
| Username/password in PostgreSQL | Managed by Clerk |
| PBKDF2 password hashing | Handled by Clerk |
| Custom JWT generation | Clerk-issued JWTs |
| HS512 signing | RS256 signing (RSA) |
| localStorage token | Clerk session management |
| Custom login/register pages | Clerk components |

## 🔐 Security Improvements

1. ✅ **No password storage** - Clerk handles it securely
2. ✅ **Session management** - Automatic token refresh
3. ✅ **MFA support** - Available out of the box
4. ✅ **Social login** - Easy to add (Google, GitHub, etc.)
5. ✅ **Email verification** - Built-in
6. ✅ **Password reset** - Managed by Clerk

## 🐛 Debugging Tips

### Frontend Issues
- Check browser console for Clerk errors
- Verify Clerk components are rendering
- Check Network tab for API calls with Authorization headers

### Backend Issues
- Check Go server logs for auth failures
- Verify `CLERK_SECRET_KEY` is set (can be empty for keyless mode)
- Test token validation separately

### Common Errors
- **"Missing or invalid Authorization header"**: API client not getting Clerk token
- **"Invalid Clerk token"**: Backend can't validate token (check CLERK_SECRET_KEY)
- **User not found**: Database sync issue - need to implement user sync

## 📖 Resources

- [Clerk Next.js Docs](https://clerk.com/docs/nextjs)
- [Clerk Go SDK](https://github.com/clerk/clerk-sdk-go)
- [Clerk Dashboard](https://dashboard.clerk.com) (after claiming app)
- [Clerk Webhooks](https://clerk.com/docs/webhooks/overview) (for user sync)
