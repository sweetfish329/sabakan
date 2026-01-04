import { TestBed } from "@angular/core/testing";
import { HttpTestingController, provideHttpClientTesting } from "@angular/common/http/testing";
import { provideHttpClient } from "@angular/common/http";
import { describe, it, expect, beforeEach, afterEach } from "vitest";

import { ModService, type Mod, type CreateModRequest, type UpdateModRequest } from "./mod.service";

describe("ModService", () => {
  let service: ModService;
  let httpMock: HttpTestingController;

  const mockMod: Mod = {
    ID: 1,
    CreatedAt: "2024-01-01T00:00:00Z",
    UpdatedAt: "2024-01-01T00:00:00Z",
    name: "Test Mod",
    slug: "test-mod",
    description: "A test mod",
    version: "1.0.0",
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [ModService, provideHttpClient(), provideHttpClientTesting()],
    });

    service = TestBed.inject(ModService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  describe("list", () => {
    it("should fetch all mods", () => {
      const mockMods: Mod[] = [mockMod, { ...mockMod, ID: 2, name: "Mod 2", slug: "mod-2" }];

      service.list().subscribe((mods) => {
        expect(mods).toEqual(mockMods);
        expect(mods.length).toBe(2);
      });

      const req = httpMock.expectOne("/api/mods");
      expect(req.request.method).toBe("GET");
      req.flush(mockMods);
    });
  });

  describe("get", () => {
    it("should fetch a mod by ID", () => {
      service.get(1).subscribe((mod) => {
        expect(mod).toEqual(mockMod);
      });

      const req = httpMock.expectOne("/api/mods/1");
      expect(req.request.method).toBe("GET");
      req.flush(mockMod);
    });
  });

  describe("create", () => {
    it("should create a new mod", () => {
      const createData: CreateModRequest = {
        name: "New Mod",
        slug: "new-mod",
        description: "A new mod",
        version: "1.0.0",
      };

      const createdMod: Mod = { ...mockMod, ...createData, ID: 3 };

      service.create(createData).subscribe((mod) => {
        expect(mod.name).toBe("New Mod");
      });

      const req = httpMock.expectOne("/api/mods");
      expect(req.request.method).toBe("POST");
      expect(req.request.body).toEqual(createData);
      req.flush(createdMod);
    });
  });

  describe("update", () => {
    it("should update an existing mod", () => {
      const updateData: UpdateModRequest = {
        name: "Updated Mod",
        description: "Updated description",
      };

      const updatedMod: Mod = { ...mockMod, ...updateData };

      service.update(1, updateData).subscribe((mod) => {
        expect(mod.name).toBe("Updated Mod");
      });

      const req = httpMock.expectOne("/api/mods/1");
      expect(req.request.method).toBe("PUT");
      expect(req.request.body).toEqual(updateData);
      req.flush(updatedMod);
    });
  });

  describe("delete", () => {
    it("should delete a mod", () => {
      service.delete(1).subscribe();

      const req = httpMock.expectOne("/api/mods/1");
      expect(req.request.method).toBe("DELETE");
      req.flush({});
    });
  });
});
