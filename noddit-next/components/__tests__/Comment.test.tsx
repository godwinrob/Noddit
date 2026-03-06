import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Comment from "../Comment";

jest.mock("@/lib/api", () => ({
  api: {
    get: jest.fn().mockResolvedValue([]),
    post: jest.fn().mockResolvedValue({}),
  },
}));

jest.mock("../ReplyForm", () => {
  return function MockReplyForm({ onSuccess }: { onSuccess: () => void }) {
    return (
      <div data-testid="reply-form">
        <button onClick={onSuccess}>Submit Reply</button>
      </div>
    );
  };
});

const mockComment = {
  postId: 10,
  body: "This is a comment",
  username: "commentuser",
  postScore: 7,
  createdDate: "2024-06-15T12:00:00Z",
  subnodditId: 1,
};

const defaultProps = {
  comment: mockComment,
  subnodditName: "golang",
  topLevelId: 1,
  onRefresh: jest.fn(),
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe("Comment", () => {
  test("renders comment body and username", () => {
    render(<Comment {...defaultProps} />);

    expect(screen.getByText("This is a comment")).toBeInTheDocument();
    expect(screen.getByText(/u\/commentuser/)).toBeInTheDocument();
  });

  test("renders vote buttons", () => {
    render(<Comment {...defaultProps} />);

    const buttons = screen.getAllByRole("button");
    // Should have upvote, downvote, and reply buttons
    expect(buttons.length).toBeGreaterThanOrEqual(2);
  });

  test("renders score", () => {
    render(<Comment {...defaultProps} />);

    expect(screen.getByText("7")).toBeInTheDocument();
  });

  test("shows reply button when signed in", () => {
    render(<Comment {...defaultProps} />);

    expect(screen.getByText("Reply")).toBeInTheDocument();
  });

  test("toggles reply form on click", async () => {
    const user = userEvent.setup();
    render(<Comment {...defaultProps} />);

    await user.click(screen.getByText("Reply"));

    expect(screen.getByTestId("reply-form")).toBeInTheDocument();
    expect(screen.getByText("Cancel")).toBeInTheDocument();
  });

  test("hides reply button when not signed in", () => {
    const { useUser } = require("@clerk/nextjs");
    useUser.mockReturnValueOnce({ user: null, isSignedIn: false, isLoaded: true });

    render(<Comment {...defaultProps} />);

    expect(screen.queryByText("Reply")).not.toBeInTheDocument();
  });
});
