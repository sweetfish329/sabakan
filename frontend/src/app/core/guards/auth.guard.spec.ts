import { TestBed } from "@angular/core/testing";
import { Router } from "@angular/router";
import { describe, it, expect, beforeEach, vi, type Mock } from "vitest";

import { AuthService } from "../../services/auth.service";
import { authGuard } from "./auth.guard";
import { guestGuard } from "./guest.guard";

describe("Auth Guards", () => {
  let authServiceMock: { isAuthenticated: Mock };
  let routerMock: { navigate: Mock; createUrlTree: Mock };

  beforeEach(() => {
    authServiceMock = {
      isAuthenticated: vi.fn(),
    };

    routerMock = {
      navigate: vi.fn(),
      createUrlTree: vi.fn().mockImplementation((commands: string[]) => ({
        toString: () => commands.join("/"),
      })),
    };

    TestBed.configureTestingModule({
      providers: [
        { provide: AuthService, useValue: authServiceMock },
        { provide: Router, useValue: routerMock },
      ],
    });
  });

  describe("authGuard", () => {
    it("should allow access when user is authenticated", () => {
      authServiceMock.isAuthenticated.mockReturnValue(true);

      const result = TestBed.runInInjectionContext(() => authGuard());

      expect(result).toBe(true);
    });

    it("should redirect to /login when user is not authenticated", () => {
      authServiceMock.isAuthenticated.mockReturnValue(false);

      const result = TestBed.runInInjectionContext(() => authGuard());

      expect(routerMock.createUrlTree).toHaveBeenCalledWith(["/login"]);
      expect(result).toEqual({ toString: expect.any(Function) });
    });
  });

  describe("guestGuard", () => {
    it("should allow access when user is not authenticated", () => {
      authServiceMock.isAuthenticated.mockReturnValue(false);

      const result = TestBed.runInInjectionContext(() => guestGuard());

      expect(result).toBe(true);
    });

    it("should redirect to / when user is authenticated", () => {
      authServiceMock.isAuthenticated.mockReturnValue(true);

      const result = TestBed.runInInjectionContext(() => guestGuard());

      expect(routerMock.createUrlTree).toHaveBeenCalledWith(["/"]);
      expect(result).toEqual({ toString: expect.any(Function) });
    });
  });
});
