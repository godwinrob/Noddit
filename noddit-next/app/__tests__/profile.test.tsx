import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import ProfilePage from "../profile/page";

const mockPush = jest.fn();
const mockSignOut = jest.fn();

jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
    back: jest.fn(),
    forward: jest.fn(),
    refresh: jest.fn(),
    replace: jest.fn(),
    prefetch: jest.fn(),
  }),
}));

jest.mock("@clerk/nextjs", () => ({
  useUser: jest.fn(),
  useClerk: jest.fn(),
}));

jest.mock("@/components/ClerkTokenProvider", () => ({
  useNodditUser: jest.fn(() => ({
    username: "testuser",
    userId: 1,
    isReady: true,
  })),
}));

jest.mock("@/lib/api", () => ({
  api: {
    get: jest.fn().mockResolvedValue([]),
    put: jest.fn().mockResolvedValue({ message: "Updated" }),
  },
}));

jest.mock("@/components/PostCard", () => {
  return function MockPostCard({ post }: { post: { title: string } }) {
    return <div data-testid="post-card">{post.title}</div>;
  };
});

const { api } = require("@/lib/api");
const { useUser, useClerk } = require("@clerk/nextjs");

beforeEach(() => {
  jest.clearAllMocks();
  useUser.mockReturnValue({
    user: {
      id: "user_123",
      username: "testuser",
      primaryEmailAddress: { emailAddress: "test@example.com" },
      publicMetadata: {},
    },
    isSignedIn: true,
    isLoaded: true,
  });
  useClerk.mockReturnValue({
    signOut: mockSignOut,
  });
  api.get.mockResolvedValue([]);
  // Mock window.alert
  jest.spyOn(window, "alert").mockImplementation(() => {});
});

describe("Profile Page", () => {
  test("redirects when signed out", () => {
    useUser.mockReturnValueOnce({
      user: null,
      isSignedIn: false,
      isLoaded: true,
    });

    render(<ProfilePage />);

    expect(mockPush).toHaveBeenCalledWith("/");
  });

  test("renders user posts", async () => {
    api.get.mockResolvedValueOnce([
      {
        postId: 1,
        title: "My Post",
        body: "content",
        username: "testuser",
        subnodditName: "golang",
        postScore: 10,
        createdDate: "2024-01-01",
        subnodditId: 1,
      },
    ]);

    render(<ProfilePage />);

    await waitFor(() => {
      expect(screen.getByText("My Post")).toBeInTheDocument();
    });
  });

  test("shows empty state when no posts", async () => {
    api.get.mockResolvedValueOnce([]);

    render(<ProfilePage />);

    await waitFor(() => {
      expect(screen.getByText(/you haven't posted anything yet/i)).toBeInTheDocument();
    });
  });

  test("displays username", async () => {
    render(<ProfilePage />);

    await waitFor(() => {
      expect(screen.getByText("testuser")).toBeInTheDocument();
    });
  });

  test("update username calls API and signs out", async () => {
    api.put.mockResolvedValueOnce({ message: "Updated" });
    const user = userEvent.setup();

    render(<ProfilePage />);

    // Open the change username section
    await user.click(screen.getByText("Change Username"));

    // Type new username
    const input = screen.getByPlaceholderText("New username");
    await user.type(input, "newname");

    // Click update
    await user.click(screen.getByText("Update Username"));

    await waitFor(() => {
      expect(api.put).toHaveBeenCalledWith(
        "/api/user/update/username/testuser",
        { newUsername: "newname" },
        true
      );
    });

    await waitFor(() => {
      expect(mockSignOut).toHaveBeenCalled();
    });
  });
});
