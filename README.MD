### Basic authentication system and user role-based access control in Go, using Fiber

This repository provides a basic setup for user authentication and management routes using the Fiber framework in Go. The app includes signup, login, email verification, user management, and various middleware to protect routes based on user roles and ownership.

### Features
- User Authentication
- Signup new users.
- Login existing users.
- Verify user email addresses after signup.
- User Management
- Seed an admin user.
- Get a list of all users (Admin only).
- Get user details by ID (Admin or the user themselves).
- Approve user accounts (Admin only).
- Edit user details (Admin or the user themselves).
- Delete a user account (Admin or the user themselves).
- Password Management
- Forgot password route to initiate password reset.

### Middleware
- Protects routes to ensure only authenticated users can access them.
- Restricts access to certain routes to admin users only.
- Ensures a user can access only their own data or allows an admin to access any user’s data.
