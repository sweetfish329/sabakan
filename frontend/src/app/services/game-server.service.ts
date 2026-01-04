import { HttpClient } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import type { Observable } from "rxjs";
import type {
  GameServer,
  CreateGameServerRequest,
  UpdateGameServerRequest,
} from "../models/game-server.model";

/**
 * Service for managing game servers via the backend API.
 */
@Injectable({ providedIn: "root" })
export class GameServerService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = "/api/game-servers";

  /**
   * Retrieves a list of all game servers for the current user.
   * @returns {Observable<GameServer[]>} Observable of GameServer array
   */
  list(): Observable<GameServer[]> {
    return this.http.get<GameServer[]>(this.baseUrl);
  }

  /**
   * Retrieves a specific game server by slug.
   * @param {string} slug - Game server slug
   * @returns {Observable<GameServer>} Observable of GameServer
   */
  get(slug: string): Observable<GameServer> {
    return this.http.get<GameServer>(`${this.baseUrl}/${encodeURIComponent(slug)}`);
  }

  /**
   * Creates a new game server.
   * @param {CreateGameServerRequest} data - Server creation data
   * @returns {Observable<GameServer>} Observable of created GameServer
   */
  create(data: CreateGameServerRequest): Observable<GameServer> {
    return this.http.post<GameServer>(this.baseUrl, data);
  }

  /**
   * Updates an existing game server.
   * @param {string} slug - Game server slug
   * @param {UpdateGameServerRequest} data - Server update data
   * @returns {Observable<GameServer>} Observable of updated GameServer
   */
  update(slug: string, data: UpdateGameServerRequest): Observable<GameServer> {
    return this.http.put<GameServer>(`${this.baseUrl}/${encodeURIComponent(slug)}`, data);
  }

  /**
   * Deletes a game server.
   * @param {string} slug - Game server slug
   * @returns {Observable<void>} Observable that completes when deleted
   */
  delete(slug: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${encodeURIComponent(slug)}`);
  }
}
