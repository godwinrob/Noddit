import { render, screen, waitFor } from "@testing-library/react";
import Home from "../page";

jest.mock("next/link", () => {
  return function MockLink({ children, href }: { children: React.ReactNode; href: string }) {
    return <a href={href}>{children}</a>;
  };
});

jest.mock("@/lib/api", () => ({
  api: {
    get: jest.fn(),
  },
}));

jest.mock("@/components/PostCard", () => {
  return function MockPostCard({ post }: { post: { title: string } }) {
    return <div data-testid="post-card">{post.title}</div>;
  };
});

const { api } = require("@/lib/api");

beforeEach(() => {
  jest.clearAllMocks();
});

describe("Home Page", () => {
  test("shows loading state initially", () => {
    api.get.mockReturnValue(new Promise(() => {})); // Never resolves
    render(<Home />);

    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  test("renders popular posts", async () => {
    api.get.mockImplementation((path: string) => {
      if (path.includes("popularposts")) {
        return Promise.resolve([
          {
            postId: 1,
            title: "Hot Post",
            body: "content",
            username: "user1",
            subnodditName: "golang",
            postScore: 100,
            createdDate: "2024-01-01",
            subnodditId: 1,
          },
        ]);
      }
      return Promise.resolve([]);
    });

    render(<Home />);

    await waitFor(() => {
      expect(screen.getByText("Hot Post")).toBeInTheDocument();
    });
  });

  test("renders active communities sidebar", async () => {
    api.get.mockImplementation((path: string) => {
      if (path.includes("active")) {
        return Promise.resolve([
          {
            subnodditId: 1,
            subnodditName: "golang",
            subnodditDescription: "Go programming",
          },
        ]);
      }
      return Promise.resolve([]);
    });

    render(<Home />);

    await waitFor(() => {
      expect(screen.getByText("n/golang")).toBeInTheDocument();
      expect(screen.getByText("Go programming")).toBeInTheDocument();
    });
  });

  test("shows empty state when no popular posts", async () => {
    api.get.mockResolvedValue([]);

    render(<Home />);

    await waitFor(() => {
      expect(screen.getByText(/no popular posts today/i)).toBeInTheDocument();
    });
  });

  test("shows action links", async () => {
    api.get.mockResolvedValue([]);

    render(<Home />);

    await waitFor(() => {
      expect(screen.getByText("Create Post")).toBeInTheDocument();
      expect(screen.getByText("Create Community")).toBeInTheDocument();
    });
  });
});
