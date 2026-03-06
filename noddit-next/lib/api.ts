const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

class ApiClient {
  private getToken(): string | null {
    if (typeof window !== "undefined") {
      return localStorage.getItem("token");
    }
    return null;
  }

  private getHeaders(includeAuth: boolean = false): HeadersInit {
    const headers: HeadersInit = {
      "Content-Type": "application/json",
    };

    if (includeAuth) {
      const token = this.getToken();
      if (token) {
        headers["Authorization"] = `Bearer ${token}`;
      }
    }

    return headers;
  }

  async get<T>(path: string, auth: boolean = false): Promise<T> {
    const response = await fetch(`${API_URL}${path}`, {
      method: "GET",
      headers: this.getHeaders(auth),
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }

    return response.json();
  }

  async post<T>(path: string, data: unknown, auth: boolean = false): Promise<T> {
    const response = await fetch(`${API_URL}${path}`, {
      method: "POST",
      headers: this.getHeaders(auth),
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
      headers: this.getHeaders(auth),
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
      headers: this.getHeaders(auth),
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }
  }
}

export const api = new ApiClient();
