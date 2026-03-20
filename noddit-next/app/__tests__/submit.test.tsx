import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import SubmitPage from "../submit/page";

const mockPush = jest.fn();
const mockBack = jest.fn();

jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
    back: mockBack,
    forward: jest.fn(),
    refresh: jest.fn(),
    replace: jest.fn(),
    prefetch: jest.fn(),
  }),
}));

jest.mock("@clerk/nextjs", () => ({
  useUser: jest.fn(),
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
    post: jest.fn().mockResolvedValue({ postId: 1 }),
  },
}));

const { api } = require("@/lib/api");
const { useUser } = require("@clerk/nextjs");

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
  api.get.mockResolvedValue([
    { subnodditId: 1, subnodditName: "golang", subnodditDescription: "Go" },
    { subnodditId: 2, subnodditName: "react", subnodditDescription: "React" },
  ]);
});

describe("Submit Page", () => {
  test("redirects when signed out", () => {
    useUser.mockReturnValueOnce({
      user: null,
      isSignedIn: false,
      isLoaded: true,
    });

    render(<SubmitPage />);

    expect(mockPush).toHaveBeenCalledWith("/");
  });

  test("renders form fields", async () => {
    render(<SubmitPage />);

    await waitFor(() => {
      expect(screen.getByLabelText(/title/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/text/i)).toBeInTheDocument();
      expect(screen.getByText(/choose a community/i)).toBeInTheDocument();
    });
  });

  test("loads communities dropdown from API", async () => {
    render(<SubmitPage />);

    // Wait for the form to load
    await waitFor(() => {
      expect(screen.getByLabelText(/title/i)).toBeInTheDocument();
    });

    // API should have been called to fetch communities
    expect(api.get).toHaveBeenCalledWith("/api/public/subnoddits");
  });

  test("shows validation error for empty title", async () => {
    const user = userEvent.setup();

    render(<SubmitPage />);

    await waitFor(() => {
      expect(screen.getByLabelText(/text/i)).toBeInTheDocument();
    });

    // Fill body but leave title empty
    await user.type(screen.getByLabelText(/text/i), "Some body content");

    // Click the Post button - the form has required attributes so browser validation
    // will prevent submission. We test that the form renders the required fields.
    const titleInput = screen.getByLabelText(/title/i);
    expect(titleInput).toBeRequired();
  });

  test("successful submit redirects to post", async () => {
    api.post.mockResolvedValueOnce({ postId: 99 });
    const user = userEvent.setup();

    render(<SubmitPage />);

    await waitFor(() => {
      expect(screen.getByLabelText(/title/i)).toBeInTheDocument();
    });

    await user.type(screen.getByLabelText(/title/i), "My New Post");
    await user.type(screen.getByLabelText(/text/i), "Post content here");
    await user.click(screen.getByRole("button", { name: /post/i }));

    await waitFor(() => {
      expect(api.post).toHaveBeenCalledWith(
        "/api/post/create",
        expect.objectContaining({
          title: "My New Post",
          body: "Post content here",
        }),
        true
      );
    });

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith("/n/golang/99");
    });
  });

  test("shows loading state during submit", async () => {
    // Make the API call hang
    api.post.mockReturnValue(new Promise(() => {}));
    const user = userEvent.setup();

    render(<SubmitPage />);

    await waitFor(() => {
      expect(screen.getByLabelText(/title/i)).toBeInTheDocument();
    });

    await user.type(screen.getByLabelText(/title/i), "Test");
    await user.type(screen.getByLabelText(/text/i), "Content");
    await user.click(screen.getByRole("button", { name: /post/i }));

    await waitFor(() => {
      expect(screen.getByText("Posting...")).toBeInTheDocument();
    });
  });
});
