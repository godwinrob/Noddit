import { render, screen } from "@testing-library/react";
import CommentSection from "../CommentSection";

// Mock the Comment component to simplify testing
jest.mock("../Comment", () => {
  return function MockComment({ comment }: { comment: { body: string; postId: number } }) {
    return <div data-testid={`comment-${comment.postId}`}>{comment.body}</div>;
  };
});

const mockComments = [
  {
    postId: 10,
    body: "First comment",
    username: "user1",
    postScore: 5,
    createdDate: "2024-01-01T00:00:00Z",
    subnodditId: 1,
  },
  {
    postId: 11,
    body: "Second comment",
    username: "user2",
    postScore: 3,
    createdDate: "2024-01-02T00:00:00Z",
    subnodditId: 1,
  },
];

const defaultProps = {
  subnodditName: "golang",
  topLevelId: 1,
  onRefresh: jest.fn(),
};

describe("CommentSection", () => {
  test("renders all comments", () => {
    render(<CommentSection comments={mockComments} {...defaultProps} />);

    expect(screen.getByTestId("comment-10")).toBeInTheDocument();
    expect(screen.getByTestId("comment-11")).toBeInTheDocument();
    expect(screen.getByText("First comment")).toBeInTheDocument();
    expect(screen.getByText("Second comment")).toBeInTheDocument();
  });

  test("shows empty state when no comments", () => {
    render(<CommentSection comments={[]} {...defaultProps} />);

    expect(screen.getByText(/no comments yet/i)).toBeInTheDocument();
  });

  test("passes correct props to Comment components", () => {
    render(<CommentSection comments={mockComments} {...defaultProps} />);

    // Both comments should be rendered
    expect(screen.getAllByTestId(/comment-/)).toHaveLength(2);
  });
});
