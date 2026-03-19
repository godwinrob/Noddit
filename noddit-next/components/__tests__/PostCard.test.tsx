import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import PostCard from "../PostCard";

jest.mock("next/image", () => {
  return function MockImage({ alt, src }: { alt: string; src: string }) {
    return <img alt={alt} src={src} />;
  };
});

jest.mock("next/link", () => {
  return function MockLink({ children, href }: { children: React.ReactNode; href: string }) {
    return <a href={href}>{children}</a>;
  };
});

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
    post: jest.fn().mockResolvedValue({}),
    delete: jest.fn().mockResolvedValue(undefined),
  },
}));

jest.mock("../ReplyForm", () => {
  return function MockReplyForm() {
    return <div data-testid="reply-form">Reply Form</div>;
  };
});

const { api } = require("@/lib/api");

const mockPost = {
  postId: 1,
  title: "Test Post Title",
  body: "Test post body content",
  username: "testuser",
  subnodditName: "golang",
  postScore: 42,
  createdDate: "2024-06-15T12:00:00Z",
  subnodditId: 1,
};

beforeEach(() => {
  jest.clearAllMocks();
  api.get.mockResolvedValue([]);
});

describe("PostCard", () => {
  test("renders title and metadata", () => {
    render(<PostCard post={mockPost} />);

    expect(screen.getByText("Test Post Title")).toBeInTheDocument();
    expect(screen.getByText(/n\/golang/)).toBeInTheDocument();
    expect(screen.getByText(/u\/testuser/)).toBeInTheDocument();
  });

  test("renders score", () => {
    render(<PostCard post={mockPost} />);

    expect(screen.getByText("42")).toBeInTheDocument();
  });

  test("title is link when not showFullBody", () => {
    render(<PostCard post={mockPost} />);

    const link = screen.getByText("Test Post Title").closest("a");
    expect(link).toHaveAttribute("href", "/s/golang/1");
  });

  test("shows full body when showFullBody is true", () => {
    render(<PostCard post={mockPost} showFullBody={true} />);

    expect(screen.getByText("Test post body content")).toBeInTheDocument();
  });

  test("upvote button calls API and updates score", async () => {
    api.post.mockResolvedValueOnce({ vote: "upvote", score: 43 });
    const user = userEvent.setup();

    render(<PostCard post={mockPost} />);

    // Find upvote button (first button with triangle)
    const upvoteBtn = screen.getAllByRole("button")[0];
    await user.click(upvoteBtn);

    await waitFor(() => {
      expect(api.post).toHaveBeenCalledWith(
        "/api/post/vote",
        expect.objectContaining({
          postId: 1,
          vote: "upvote",
        }),
        true
      );
    });

    // Score should update to 43
    await waitFor(() => {
      expect(screen.getByText("43")).toBeInTheDocument();
    });
  });

  test("vote buttons disabled after voting", async () => {
    api.post.mockReturnValueOnce(new Promise(() => {})); // Never resolves — keeps isVoting=true
    const user = userEvent.setup();

    render(<PostCard post={mockPost} />);

    const upvoteBtn = screen.getAllByRole("button")[0];
    await user.click(upvoteBtn);

    await waitFor(() => {
      expect(upvoteBtn).toBeDisabled();
    });
  });

  test("shows delete button for post author", async () => {
    render(<PostCard post={mockPost} />);

    // The author is "testuser" which matches our Clerk mock user
    await waitFor(() => {
      expect(screen.getByText("Delete")).toBeInTheDocument();
    });
  });

  test("shows reply button in full body mode", () => {
    render(<PostCard post={mockPost} showFullBody={true} />);

    expect(screen.getByText("Reply")).toBeInTheDocument();
  });

  test("toggles reply form", async () => {
    const user = userEvent.setup();

    render(<PostCard post={mockPost} showFullBody={true} />);

    await user.click(screen.getByText("Reply"));

    expect(screen.getByTestId("reply-form")).toBeInTheDocument();
  });

  test("renders image when imageAddress is provided", () => {
    const postWithImage = {
      ...mockPost,
      imageAddress: "https://example.com/image.jpg",
    };

    render(<PostCard post={postWithImage} showFullBody={true} />);

    const img = screen.getByAltText("Test Post Title");
    expect(img).toBeInTheDocument();
  });
});
