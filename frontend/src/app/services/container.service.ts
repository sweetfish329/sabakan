import { HttpClient } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import type { Observable } from "rxjs";
import type { Container, ContainerLogEntry } from "../models/container.model";

/**
 * Service for managing containers via the backend API.
 */
@Injectable({ providedIn: "root" })
export class ContainerService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = "/api/containers";

  /**
   * Retrieves a list of all containers.
   * @returns {Observable<Container[]>} Observable of Container array
   */
  list(): Observable<Container[]> {
    return this.http.get<Container[]>(this.baseUrl);
  }

  /**
   * Retrieves a specific container by ID.
   * @param {string} id - Container ID or name
   * @returns {Observable<Container>} Observable of Container
   */
  get(id: string): Observable<Container> {
    return this.http.get<Container>(`${this.baseUrl}/${encodeURIComponent(id)}`);
  }

  /**
   * Starts a container.
   * @param {string} id - Container ID or name
   * @returns {Observable<void>} Observable that completes when the container starts
   */
  start(id: string): Observable<void> {
    return this.http.post<void>(`${this.baseUrl}/${encodeURIComponent(id)}/start`, {});
  }

  /**
   * Stops a container.
   * @param {string} id - Container ID or name
   * @param {number} timeout - Timeout in seconds (default: 10)
   * @returns {Observable<void>} Observable that completes when the container stops
   */
  stop(id: string, timeout = 10): Observable<void> {
    return this.http.post<void>(
      `${this.baseUrl}/${encodeURIComponent(id)}/stop`,
      {},
      {
        params: { timeout: timeout.toString() },
      },
    );
  }

  /**
   * Retrieves container logs.
   * @param {string} id - Container ID or name
   * @param {number} lines - Number of log lines to retrieve (default: 100)
   * @returns {Observable<ContainerLogEntry[]>} Observable of log entries
   */
  logs(id: string, lines = 100): Observable<ContainerLogEntry[]> {
    return this.http.get<ContainerLogEntry[]>(`${this.baseUrl}/${encodeURIComponent(id)}/logs`, {
      params: { lines: lines.toString() },
    });
  }
}
