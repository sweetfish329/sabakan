import { TestBed } from "@angular/core/testing";
import { HttpClient, provideHttpClient, withInterceptors } from "@angular/common/http";
import { HttpTestingController, provideHttpClientTesting } from "@angular/common/http/testing";
import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import type { Mock } from "vitest";

import { AuthService } from "../../services/auth.service";
import { authInterceptor } from "./auth.interceptor";

describe("authInterceptor", () => {
  let httpClient: HttpClient;
  let httpMock: HttpTestingController;
  let authServiceMock: {
    getAccessToken: Mock;
    refresh: Mock;
    clearTokens: Mock;
  };

  beforeEach(() => {
    authServiceMock = {
      getAccessToken: vi.fn(),
      refresh: vi.fn(),
      clearTokens: vi.fn(),
    };

    TestBed.configureTestingModule({
      providers: [
        provideHttpClient(withInterceptors([authInterceptor])),
        provideHttpClientTesting(),
        { provide: AuthService, useValue: authServiceMock },
      ],
    });

    httpClient = TestBed.inject(HttpClient);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  describe("Authorization Header", () => {
    it("should add Authorization header when access token exists", () => {
      authServiceMock.getAccessToken.mockReturnValue("test-token");

      httpClient.get("/api/containers").subscribe();

      const req = httpMock.expectOne("/api/containers");
      expect(req.request.headers.get("Authorization")).toBe("Bearer test-token");
      req.flush({});
    });

    it("should not add Authorization header when no token exists", () => {
      authServiceMock.getAccessToken.mockReturnValue(undefined);

      httpClient.get("/api/containers").subscribe();

      const req = httpMock.expectOne("/api/containers");
      expect(req.request.headers.has("Authorization")).toBe(false);
      req.flush({});
    });

    it("should not add Authorization header for auth endpoints", () => {
      authServiceMock.getAccessToken.mockReturnValue("test-token");

      httpClient.post("/auth/login", {}).subscribe();

      const req = httpMock.expectOne("/auth/login");
      expect(req.request.headers.has("Authorization")).toBe(false);
      req.flush({});
    });

    it("should not add Authorization header for refresh endpoint", () => {
      authServiceMock.getAccessToken.mockReturnValue("test-token");

      httpClient.post("/auth/refresh", {}).subscribe();

      const req = httpMock.expectOne("/auth/refresh");
      expect(req.request.headers.has("Authorization")).toBe(false);
      req.flush({});
    });
  });

  describe("401 Error Handling", () => {
    it("should pass through non-401 errors", () => {
      authServiceMock.getAccessToken.mockReturnValue("test-token");

      let errorReceived = false;
      httpClient.get("/api/containers").subscribe({
        error: (error) => {
          errorReceived = true;
          expect(error.status).toBe(500);
        },
      });

      const req = httpMock.expectOne("/api/containers");
      req.flush({ message: "Server Error" }, { status: 500, statusText: "Internal Server Error" });

      expect(errorReceived).toBe(true);
    });

    it("should clear tokens on 401 for auth endpoints", () => {
      authServiceMock.getAccessToken.mockReturnValue("test-token");

      httpClient.post("/auth/refresh", {}).subscribe({
        error: () => {
          // Expected error
        },
      });

      const req = httpMock.expectOne("/auth/refresh");
      req.flush({ message: "Unauthorized" }, { status: 401, statusText: "Unauthorized" });

      expect(authServiceMock.clearTokens).toHaveBeenCalled();
    });
  });
});
