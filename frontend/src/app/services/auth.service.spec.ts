import { TestBed } from "@angular/core/testing";
import { HttpTestingController, provideHttpClientTesting } from "@angular/common/http/testing";
import { provideHttpClient } from "@angular/common/http";
import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";

import {
  AuthService,
  type AuthResponse,
  type LoginRequest,
  type RegisterRequest,
} from "./auth.service";

describe("AuthService", () => {
  let service: AuthService;
  let httpMock: HttpTestingController;

  // Mock localStorage
  const localStorageMock = (() => {
    let store: Record<string, string> = {};
    return {
      getItem: vi.fn((key: string) => store[key] ?? null),
      setItem: vi.fn((key: string, value: string) => {
        store[key] = value;
      }),
      removeItem: vi.fn((key: string) => {
        delete store[key];
      }),
      clear: vi.fn(() => {
        store = {};
      }),
    };
  })();

  beforeEach(() => {
    // Setup localStorage mock
    Object.defineProperty(globalThis, "localStorage", {
      value: localStorageMock,
      writable: true,
    });
    localStorageMock.clear();

    TestBed.configureTestingModule({
      providers: [AuthService, provideHttpClient(), provideHttpClientTesting()],
    });

    service = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  describe("login", () => {
    it("should send login request and store tokens", () => {
      const credentials: LoginRequest = {
        username: "testuser",
        password: "password123",
      };
      const mockResponse: AuthResponse = {
        access_token: "mock.access.token",
        refresh_token: "mock-refresh-token",
        expires_in: 900,
        token_type: "Bearer",
      };

      service.login(credentials).subscribe((response) => {
        expect(response).toEqual(mockResponse);
      });

      const req = httpMock.expectOne("/auth/login");
      expect(req.request.method).toBe("POST");
      expect(req.request.body).toEqual(credentials);
      req.flush(mockResponse);

      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        "sabakan_access_token",
        "mock.access.token",
      );
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        "sabakan_refresh_token",
        "mock-refresh-token",
      );
      expect(service.isAuthenticated()).toBe(true);
    });

    it("should handle login error", () => {
      const credentials: LoginRequest = {
        username: "wronguser",
        password: "wrongpassword",
      };

      service.login(credentials).subscribe({
        error: (error) => {
          expect(error.status).toBe(401);
        },
      });

      const req = httpMock.expectOne("/auth/login");
      req.flush({ message: "Invalid credentials" }, { status: 401, statusText: "Unauthorized" });

      expect(service.isAuthenticated()).toBe(false);
    });
  });

  describe("register", () => {
    it("should send registration request", () => {
      const registerData: RegisterRequest = {
        username: "newuser",
        email: "newuser@example.com",
        password: "password123",
      };
      const mockResponse = { message: "User created", user_id: 1 };

      service.register(registerData).subscribe((response) => {
        expect(response).toEqual(mockResponse);
      });

      const req = httpMock.expectOne("/auth/register");
      expect(req.request.method).toBe("POST");
      expect(req.request.body).toEqual(registerData);
      req.flush(mockResponse);
    });
  });

  describe("refresh", () => {
    it("should send refresh request with stored refresh token", () => {
      localStorageMock.getItem.mockReturnValueOnce("stored-refresh-token");

      const mockResponse: AuthResponse = {
        access_token: "new.access.token",
        refresh_token: "new-refresh-token",
        expires_in: 900,
        token_type: "Bearer",
      };

      service.refresh().subscribe((response) => {
        expect(response).toEqual(mockResponse);
      });

      const req = httpMock.expectOne("/auth/refresh");
      expect(req.request.method).toBe("POST");
      expect(req.request.body).toEqual({ refresh_token: "stored-refresh-token" });
      req.flush(mockResponse);
    });
  });

  describe("logout", () => {
    it("should send logout request and clear tokens", () => {
      service.logout().subscribe();

      const req = httpMock.expectOne("/auth/logout");
      expect(req.request.method).toBe("POST");
      req.flush({ message: "Logged out" });

      expect(localStorageMock.removeItem).toHaveBeenCalledWith("sabakan_access_token");
      expect(localStorageMock.removeItem).toHaveBeenCalledWith("sabakan_refresh_token");
      expect(service.isAuthenticated()).toBe(false);
    });
  });

  describe("clearTokens", () => {
    it("should clear tokens without API call", () => {
      service.clearTokens();

      expect(localStorageMock.removeItem).toHaveBeenCalledWith("sabakan_access_token");
      expect(localStorageMock.removeItem).toHaveBeenCalledWith("sabakan_refresh_token");
      expect(service.isAuthenticated()).toBe(false);
    });
  });

  describe("getAccessToken", () => {
    it("should return current access token", () => {
      // Initially undefined
      expect(service.getAccessToken()).toBeUndefined();

      // After login
      const mockResponse: AuthResponse = {
        access_token: "test.token",
        expires_in: 900,
        token_type: "Bearer",
      };

      service.login({ username: "test", password: "test" }).subscribe();
      httpMock.expectOne("/auth/login").flush(mockResponse);

      expect(service.getAccessToken()).toBe("test.token");
    });
  });

  describe("currentUser", () => {
    it("should decode JWT and return user info", () => {
      // Create a mock JWT with payload
      const payload = { user_id: 123, username: "testuser" };
      const encodedPayload = btoa(JSON.stringify(payload));
      const mockToken = `header.${encodedPayload}.signature`;

      const mockResponse: AuthResponse = {
        access_token: mockToken,
        expires_in: 900,
        token_type: "Bearer",
      };

      service.login({ username: "test", password: "test" }).subscribe();
      httpMock.expectOne("/auth/login").flush(mockResponse);

      const user = service.currentUser();
      expect(user).toEqual({ user_id: 123, username: "testuser" });
    });

    it("should return undefined for invalid token", () => {
      expect(service.currentUser()).toBeUndefined();
    });
  });

  describe("getOAuthUrl", () => {
    it("should return Google OAuth URL", () => {
      const url = service.getOAuthUrl("google");
      expect(url).toContain("/auth/google");
      expect(url).toContain("redirect_url=");
    });

    it("should return Discord OAuth URL", () => {
      const url = service.getOAuthUrl("discord");
      expect(url).toContain("/auth/discord");
      expect(url).toContain("redirect_url=");
    });
  });

  describe("storeOAuthTokens", () => {
    it("should store OAuth tokens", () => {
      service.storeOAuthTokens("oauth-access-token", "oauth-refresh-token");

      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        "sabakan_access_token",
        "oauth-access-token",
      );
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        "sabakan_refresh_token",
        "oauth-refresh-token",
      );
      expect(service.isAuthenticated()).toBe(true);
    });
  });
});
