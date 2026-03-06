import { render, screen } from "@testing-library/react";
import { ClerkTokenProvider } from "../ClerkTokenProvider";

jest.mock("@/lib/api", () => ({
  setClerkTokenGetter: jest.fn(),
}));

const { setClerkTokenGetter } = require("@/lib/api");

describe("ClerkTokenProvider", () => {
  test("calls setClerkTokenGetter on mount", () => {
    render(
      <ClerkTokenProvider>
        <div>Child Content</div>
      </ClerkTokenProvider>
    );

    expect(setClerkTokenGetter).toHaveBeenCalledWith(expect.any(Function));
  });

  test("renders children", () => {
    render(
      <ClerkTokenProvider>
        <div>Child Content</div>
      </ClerkTokenProvider>
    );

    expect(screen.getByText("Child Content")).toBeInTheDocument();
  });
});
