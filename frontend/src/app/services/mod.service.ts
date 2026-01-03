import { HttpClient } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import type { Observable } from "rxjs";

/**
 * Mod data model.
 */
export interface Mod {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  name: string;
  slug: string;
  description?: string;
  sourceUrl?: string;
  version?: string;
}

/**
 * Request to create a new mod.
 */
export interface CreateModRequest {
  name: string;
  slug: string;
  description?: string;
  sourceUrl?: string;
  version?: string;
}

/**
 * Request to update an existing mod.
 */
export interface UpdateModRequest {
  name?: string;
  slug?: string;
  description?: string;
  sourceUrl?: string;
  version?: string;
}

/**
 * Service for managing mods.
 */
@Injectable({ providedIn: "root" })
export class ModService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = "/api/mods";

  /**
   * Fetches all mods.
   * @returns {Observable<Mod[]>} Observable of mod array
   */
  list(): Observable<Mod[]> {
    return this.http.get<Mod[]>(this.baseUrl);
  }

  /**
   * Fetches a mod by ID.
   * @param {number} id - Mod ID
   * @returns {Observable<Mod>} Observable of mod
   */
  get(id: number): Observable<Mod> {
    return this.http.get<Mod>(`${this.baseUrl}/${id}`);
  }

  /**
   * Creates a new mod.
   * @param {CreateModRequest} data - Mod creation data
   * @returns {Observable<Mod>} Observable of created mod
   */
  create(data: CreateModRequest): Observable<Mod> {
    return this.http.post<Mod>(this.baseUrl, data);
  }

  /**
   * Updates an existing mod.
   * @param {number} id - Mod ID
   * @param {UpdateModRequest} data - Mod update data
   * @returns {Observable<Mod>} Observable of updated mod
   */
  update(id: number, data: UpdateModRequest): Observable<Mod> {
    return this.http.put<Mod>(`${this.baseUrl}/${id}`, data);
  }

  /**
   * Deletes a mod by ID.
   * @param {number} id - Mod ID
   * @returns {Observable<void>} Observable that completes on deletion
   */
  delete(id: number): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${id}`);
  }
}
