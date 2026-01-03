import { HttpClient } from "@angular/common/http";
import { Injectable, inject, signal, computed } from "@angular/core";
import type { Observable } from "rxjs";
import { tap } from "rxjs/operators";

/**
 * Authentication response from the API.
 */
export interface AuthResponse {
  access_token: string;
  refresh_token?: string;
  expires_in: number;
  token_type: string;
}

/**
 * Login request payload.
 */
export interface LoginRequest {
  username: string;
  password: string;
}

/**
 * Registration request payload.
 */
export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

/**
 * User information decoded from JWT.
 */
export interface UserInfo {
  user_id: number;
  username: string;
}

const ACCESS_TOKEN_KEY = "sabakan_access_token";
const REFRESH_TOKEN_KEY = "sabakan_refresh_token";

/**
 * Service for managing authentication.
 */
@Injectable({ providedIn: "root" })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = "/auth";

  /** Current access token */
  private readonly _accessToken = signal<string | undefined>(this.getStoredToken());

  /** Whether the user is authenticated */
  readonly isAuthenticated = computed(() => this._accessToken() !== undefined);

  /** Current user information */
  readonly currentUser = computed<UserInfo | undefined>(() => {
    const token = this._accessToken();
    if (token === undefined || token === "") {
      return undefined;
    }
    return this.decodeToken(token);
  });

  /**
   * Logs in with username and password.
   * @param credentials - Login credentials
   * @returns Observable of AuthResponse
   */
  login(credentials: LoginRequest): Observable<AuthResponse> {
    return this.http
      .post<AuthResponse>(`${this.baseUrl}/login`, credentials)
      .pipe(tap((response) => this.storeTokens(response)));
  }

  /**
   * Registers a new user.
   * @param data - Registration data
   * @returns Observable of the response
   */
  register(data: RegisterRequest): Observable<{ message: string; user_id: number }> {
    return this.http.post<{ message: string; user_id: number }>(`${this.baseUrl}/register`, data);
  }

  /**
   * Refreshes the access token using the refresh token.
   * @returns Observable of AuthResponse
   */
  refresh(): Observable<AuthResponse> {
    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
    return this.http
      .post<AuthResponse>(`${this.baseUrl}/refresh`, {
        refresh_token: refreshToken,
      })
      .pipe(tap((response) => this.storeTokens(response)));
  }

  /**
   * Logs out the current user.
   * @returns Observable that completes on logout
   */
  logout(): Observable<{ message: string }> {
    return this.http
      .post<{ message: string }>(`${this.baseUrl}/logout`, null)
      .pipe(tap(() => this.clearTokens()));
  }

  /**
   * Clears tokens without calling the API (for local logout).
   */
  clearTokens(): void {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    this._accessToken.set(undefined);
  }

  /**
   * Gets the current access token for HTTP interceptor.
   * @returns Access token or undefined
   */
  getAccessToken(): string | undefined {
    return this._accessToken();
  }

  /**
   * Stores tokens from OAuth callback.
   * @param accessToken - Access token
   * @param refreshToken - Refresh token
   */
  storeOAuthTokens(accessToken: string, refreshToken: string): void {
    this.storeTokens({
      access_token: accessToken,
      refresh_token: refreshToken,
      expires_in: 900,
      token_type: "Bearer",
    });
  }

  /**
   * Returns the OAuth authorization URL for a provider.
   * @param provider - OAuth provider name (google, discord)
   * @returns OAuth authorization URL
   */
  getOAuthUrl(provider: "google" | "discord"): string {
    return `${this.baseUrl}/oauth/${provider}`;
  }

  private storeTokens(response: AuthResponse): void {
    localStorage.setItem(ACCESS_TOKEN_KEY, response.access_token);
    if (response.refresh_token) {
      localStorage.setItem(REFRESH_TOKEN_KEY, response.refresh_token);
    }
    this._accessToken.set(response.access_token);
  }

  private getStoredToken(): string | undefined {
    return localStorage.getItem(ACCESS_TOKEN_KEY) ?? undefined;
  }

  private decodeToken(token: string): UserInfo | undefined {
    try {
      const parts = token.split(".");
      const payload = parts[1];
      if (payload === undefined) {
        return undefined;
      }
      const decoded = JSON.parse(atob(payload)) as { user_id: number; username: string };
      return {
        user_id: decoded.user_id,
        username: decoded.username,
      };
    } catch {
      return undefined;
    }
  }
}
