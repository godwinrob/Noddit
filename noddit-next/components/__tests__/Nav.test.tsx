import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Nav from "../Nav";

const mockPush = jest.fn();

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

jest.mock("next/link", () => {
  return function MockLink({ children, href }: { children: React.ReactNode; href: string }) {
    return <a href={href}>{children}</a>;
  };
});

jest.mock("@/lib/api", () => ({
  api: {
    get: jest.fn().mockResolvedValue([]),
  },
}));

const { api } = require("@/lib/api");

beforeEach(() => {
  jest.clearAllMocks();
  api.get.mockResolvedValue([]);
});

describe("Nav", () => {
  test("renders logo link", () => {
    render(<Nav />);

    expect(screen.getByText("Noddit")).toBeInTheDocument();
    expect(screen.getByText("Noddit").closest("a")).toHaveAttribute("href", "/");
  });

  test("renders communities link", () => {
    render(<Nav />);

    expect(screen.getByText("Communities")).toBeInTheDocument();
  });

  test("shows signed-in UI with user info", () => {
    render(<Nav />);

    expect(screen.getByText("Create Post")).toBeInTheDocument();
    expect(screen.getByText(/u\/testuser/)).toBeInTheDocument();
  });

  test("shows signed-out UI", () => {
    const { useUser } = require("@clerk/nextjs");
    useUser.mockReturnValueOnce({ user: null, isSignedIn: false, isLoaded: true });

    render(<Nav />);

    expect(screen.getByText("Log In")).toBeInTheDocument();
    expect(screen.getByText("Sign Up")).toBeInTheDocument();
    expect(screen.queryByText("Create Post")).not.toBeInTheDocument();
  });

  test("renders favorites dropdown when user has favorites", async () => {
    api.get.mockResolvedValueOnce([
      { subnodditName: "golang" },
      { subnodditName: "react" },
    ]);

    render(<Nav />);

    await waitFor(() => {
      expect(screen.getByText(/Favorites/)).toBeInTheDocument();
    });
  });

  test("search submits and navigates", async () => {
    const user = userEvent.setup();

    render(<Nav />);

    const searchInput = screen.getByPlaceholderText("Search communities...");
    await user.type(searchInput, "golang{enter}");

    expect(mockPush).toHaveBeenCalledWith("/search/golang");
  });

  test("search clears after submit", async () => {
    const user = userEvent.setup();

    render(<Nav />);

    const searchInput = screen.getByPlaceholderText("Search communities...");
    await user.type(searchInput, "golang{enter}");

    expect(searchInput).toHaveValue("");
  });
});
