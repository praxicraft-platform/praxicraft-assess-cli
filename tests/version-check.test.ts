import { describe, expect, test } from "bun:test";
import { compareSemver } from "../src/utils/version-check";

describe("compareSemver", () => {
  test("orders patch versions", () => {
    expect(compareSemver("2.0.7", "2.0.6")).toBe(1);
    expect(compareSemver("2.0.6", "2.0.7")).toBe(-1);
    expect(compareSemver("2.0.6", "2.0.6")).toBe(0);
  });

  test("orders minor and major versions", () => {
    expect(compareSemver("2.1.0", "2.0.9")).toBe(1);
    expect(compareSemver("3.0.0", "2.9.9")).toBe(1);
  });

  test("strips leading v prefix", () => {
    expect(compareSemver("v2.0.7", "2.0.6")).toBe(1);
  });
});

describe("checkForUpdate", () => {
  test("returns null when registry version is not newer", async () => {
    const { checkForUpdate } = await import("../src/utils/version-check");
    const fetchImpl = async () =>
      ({
        ok: true,
        json: async () => ({ version: "0.0.1" }),
      }) as Response;

    await expect(checkForUpdate(fetchImpl)).resolves.toBeNull();
  });

  test("returns latest when registry version is newer", async () => {
    const { checkForUpdate } = await import("../src/utils/version-check");
    const fetchImpl = async () =>
      ({
        ok: true,
        json: async () => ({ version: "99.0.0" }),
      }) as Response;

    await expect(checkForUpdate(fetchImpl)).resolves.toEqual({
      current: expect.any(String),
      latest: "99.0.0",
    });
  });

  test("returns null on fetch failure", async () => {
    const { checkForUpdate } = await import("../src/utils/version-check");
    const fetchImpl = async () => {
      throw new Error("offline");
    };

    await expect(checkForUpdate(fetchImpl)).resolves.toBeNull();
  });
});
