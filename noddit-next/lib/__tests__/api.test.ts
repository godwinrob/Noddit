import { api, setClerkTokenGetter } from "../api";

// Mock fetch globally
const mockFetch = jest.fn();
global.fetch = mockFetch;

beforeEach(() => {
  mockFetch.mockReset();
  // Reset token getter
  setClerkTokenGetter(async () => null);
});

describe("ApiClient", () => {
  test("GET request sends correct URL and method", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ data: "test" }),
    });

    await api.get("/api/public/posts");

    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/public/posts",
      expect.objectContaining({ method: "GET" })
    );
  });

  test("POST request sends body and correct method", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ id: 1 }),
    });

    await api.post("/api/post/create", { title: "Test" });

    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/post/create",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ title: "Test" }),
      })
    );
  });

  test("PUT request sends body and correct method", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ message: "updated" }),
    });

    await api.put("/api/user/update", { name: "New" });

    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/user/update",
      expect.objectContaining({
        method: "PUT",
        body: JSON.stringify({ name: "New" }),
      })
    );
  });

  test("DELETE request sends correct method", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
    });

    await api.delete("/api/post/delete/1");

    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/post/delete/1",
      expect.objectContaining({ method: "DELETE" })
    );
  });

  test("auth header included when token getter is set and auth=true", async () => {
    setClerkTokenGetter(async () => "test-jwt-token");

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    });

    await api.get("/api/favorites/user1", true);

    const headers = mockFetch.mock.calls[0][1].headers;
    expect(headers["Authorization"]).toBe("Bearer test-jwt-token");
  });

  test("no auth header when auth=false", async () => {
    setClerkTokenGetter(async () => "test-jwt-token");

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    });

    await api.get("/api/public/posts", false);

    const headers = mockFetch.mock.calls[0][1].headers;
    expect(headers["Authorization"]).toBeUndefined();
  });

  test("throws error on non-ok response for GET", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      statusText: "Not Found",
    });

    await expect(api.get("/api/post/999")).rejects.toThrow("API Error: Not Found");
  });

  test("POST error parsing extracts error message from JSON", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      statusText: "Bad Request",
      json: () => Promise.resolve({ error: "Invalid input" }),
    });

    await expect(api.post("/api/post/create", {})).rejects.toThrow("Invalid input");
  });
});
