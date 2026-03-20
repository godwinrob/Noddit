import { render, screen, waitFor } from "@testing-library/react";
import { ClerkTokenProvider, useNodditUser } from "../ClerkTokenProvider";

jest.mock("@/lib/api", () => ({
  setClerkTokenGetter: jest.fn(),
  api: {
    post: jest.fn(),
  },
}));

const { setClerkTokenGetter, api } = require("@/lib/api");

function UsernameDisplay() {
  const { username, isReady } = useNodditUser();
  return (
    <div>
      <span data-testid="username">{username ?? "none"}</span>
      <span data-testid="ready">{isReady ? "yes" : "no"}</span>
    </div>
  );
}

beforeEach(() => {
  jest.clearAllMocks();
});

describe("ClerkTokenProvider", () => {
  test("calls setClerkTokenGetter on mount", () => {
    api.post.mockResolvedValue({ isNew: false, username: "testuser", userId: 1 });

    render(
      <ClerkTokenProvider>
        <div>Child Content</div>
      </ClerkTokenProvider>
    );

    expect(setClerkTokenGetter).toHaveBeenCalledWith(expect.any(Function));
  });

  test("renders children", () => {
    api.post.mockResolvedValue({ isNew: false, username: "testuser", userId: 1 });

    render(
      <ClerkTokenProvider>
        <div>Child Content</div>
      </ClerkTokenProvider>
    );

    expect(screen.getByText("Child Content")).toBeInTheDocument();
  });

  test("sets username in context for existing user", async () => {
    api.post.mockResolvedValue({ isNew: false, username: "testuser", userId: 1 });

    render(
      <ClerkTokenProvider>
        <UsernameDisplay />
      </ClerkTokenProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId("username")).toHaveTextContent("testuser");
      expect(screen.getByTestId("ready")).toHaveTextContent("yes");
    });
  });

  test("shows username modal for new user", async () => {
    api.post.mockResolvedValue({ isNew: true });

    render(
      <ClerkTokenProvider>
        <div>App Content</div>
      </ClerkTokenProvider>
    );

    await waitFor(() => {
      expect(screen.getByText("Choose your username")).toBeInTheDocument();
    });
  });
});
