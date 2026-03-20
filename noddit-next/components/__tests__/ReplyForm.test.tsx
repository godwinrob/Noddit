import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import ReplyForm from "../ReplyForm";

jest.mock("@/components/ClerkTokenProvider", () => ({
  useNodditUser: jest.fn(() => ({
    username: "testuser",
    userId: 1,
    isReady: true,
  })),
}));

jest.mock("@/lib/api", () => ({
  api: {
    post: jest.fn().mockResolvedValue({}),
  },
}));

const { api } = require("@/lib/api");

const defaultProps = {
  subnodditName: "golang",
  subnodditId: 1,
  parentPostId: 42,
  topLevelId: 42,
  onSuccess: jest.fn(),
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe("ReplyForm", () => {
  test("renders textarea and submit button", () => {
    render(<ReplyForm {...defaultProps} />);

    expect(screen.getByPlaceholderText("What are your thoughts?")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /reply/i })).toBeInTheDocument();
  });

  test("shows character count", () => {
    render(<ReplyForm {...defaultProps} />);

    expect(screen.getByText("0/2000")).toBeInTheDocument();
  });

  test("submit button disabled when textarea is empty", () => {
    render(<ReplyForm {...defaultProps} />);

    expect(screen.getByRole("button", { name: /reply/i })).toBeDisabled();
  });

  test("submit button enabled when textarea has content", async () => {
    const user = userEvent.setup();
    render(<ReplyForm {...defaultProps} />);

    await user.type(screen.getByPlaceholderText("What are your thoughts?"), "Hello!");

    expect(screen.getByRole("button", { name: /reply/i })).not.toBeDisabled();
  });

  test("calls API and onSuccess on submit", async () => {
    const onSuccess = jest.fn();
    const user = userEvent.setup();

    render(<ReplyForm {...defaultProps} onSuccess={onSuccess} />);

    await user.type(screen.getByPlaceholderText("What are your thoughts?"), "Great post!");
    await user.click(screen.getByRole("button", { name: /reply/i }));

    await waitFor(() => {
      expect(api.post).toHaveBeenCalledWith(
        "/golang/42/createreply",
        expect.objectContaining({
          body: "Great post!",
          username: "testuser",
          parentPostId: 42,
          topLevelId: 42,
          subnodditId: 1,
        }),
        true
      );
    });

    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalled();
    });
  });

  test("clears textarea after successful submit", async () => {
    const user = userEvent.setup();

    render(<ReplyForm {...defaultProps} />);

    const textarea = screen.getByPlaceholderText("What are your thoughts?");
    await user.type(textarea, "Reply text");
    await user.click(screen.getByRole("button", { name: /reply/i }));

    await waitFor(() => {
      expect(textarea).toHaveValue("");
    });
  });
});
