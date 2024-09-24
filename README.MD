### Basic authentication system and user role-based access control in Go, using Fiber

This repository provides a basic setup for user authentication and management routes using the Fiber framework in Go. The app includes signup, login, email verification, user management, and various middleware to protect routes based on user roles and ownership.

### Features
- User Authentication
- POST /api/signup - Signup new users.
- POST /api/signin - Login existing users.
- GET /api/verify - Verify user email addresses after signup.
- User Management
- POST /api/seed/admin - Seed an admin user.
- GET /api/users - Get a list of all users (Admin only).
- GET /api/users/:id - Get user details by ID (Admin or the user themselves).
- PATCH /api/users/:id/approve - Approve user accounts (Admin only).
- PATCH /api/users/:id - Edit user details (Admin or the user themselves).
- DELETE /api/users/:id - Delete a user account (Admin or the user themselves).
- Password Management
- POST /api/forgot-password - Forgot password route to initiate password reset.

### Middleware
- AuthMiddleware: Protects routes to ensure only authenticated users can access them.
- AdminOnlyMiddleware: Restricts access to certain routes to admin users only.
- OwnDataOrAdminMiddleware: Ensures a user can access only their own data or allows an admin to access any user’s data.
