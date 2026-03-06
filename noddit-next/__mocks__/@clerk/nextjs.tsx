import React from "react";

const mockUser = {
  id: "user_123",
  username: "testuser",
  primaryEmailAddress: { emailAddress: "test@example.com" },
  publicMetadata: { role: "user" },
};

const mockGetToken = jest.fn().mockResolvedValue("mock-token");

export const useUser = jest.fn(() => ({
  user: mockUser,
  isSignedIn: true,
  isLoaded: true,
}));

export const useAuth = jest.fn(() => ({
  getToken: mockGetToken,
  isSignedIn: true,
  isLoaded: true,
  userId: "user_123",
}));

export const useClerk = jest.fn(() => ({
  signOut: jest.fn(),
}));

export const SignInButton = ({ children }: { children: React.ReactNode }) => (
  <div data-testid="sign-in-button">{children}</div>
);

export const SignUpButton = ({ children }: { children: React.ReactNode }) => (
  <div data-testid="sign-up-button">{children}</div>
);

export const UserButton = () => <div data-testid="user-button" />;

export const ClerkProvider = ({ children }: { children: React.ReactNode }) => (
  <>{children}</>
);
