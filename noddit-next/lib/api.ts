const API_URL = process.env.NEXT_PUBLIC_API_URL;

// Validate API URL is configured
if (!API_URL) {
  console.error(
    "NEXT_PUBLIC_API_URL is not set. Please set it in your .env.local file."
  );
  if (typeof window !== "undefined") {
    // Only throw in browser, not during build
    throw new Error(
      "API URL not configured. Check your environment variables."
    );
  }
}

// Global token getter - will be set by Clerk
let clerkTokenGetter: (() => Promise<string | null>) | null = null;

export function setClerkTokenGetter(getter: () => Promise<string | null>) {
  clerkTokenGetter = getter;
}

class ApiClient {
  private async getToken(): Promise<string | null> {
    if (clerkTokenGetter) {
      return await clerkTokenGetter();
    }
    return null;
  }

  private async getHeaders(includeAuth: boolean = false): Promise<HeadersInit> {
    const headers: HeadersInit = {
      "Content-Type": "application/json",
    };

    if (includeAuth) {
      const token = await this.getToken();
      if (token) {
        headers["Authorization"] = `Bearer ${token}`;
      }
    }

    return headers;
  }

  async get<T>(path: string, auth: boolean = false): Promise<T> {
    const response = await fetch(`${API_URL}${path}`, {
      method: "GET",
      headers: await this.getHeaders(auth),
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }

    return response.json();
  }

  async post<T>(path: string, data: unknown, auth: boolean = false): Promise<T> {
    const response = await fetch(`${API_URL}${path}`, {
      method: "POST",
      headers: await this.getHeaders(auth),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `API Error: ${response.statusText}`);
    }

    return response.json();
  }

  async put<T>(path: string, data: unknown, auth: boolean = true): Promise<T> {
    const response = await fetch(`${API_URL}${path}`, {
      method: "PUT",
      headers: await this.getHeaders(auth),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }

    return response.json();
  }

  async delete(path: string, auth: boolean = true): Promise<void> {
    const response = await fetch(`${API_URL}${path}`, {
      method: "DELETE",
      headers: await this.getHeaders(auth),
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }
  }
}

export const api = new ApiClient();
