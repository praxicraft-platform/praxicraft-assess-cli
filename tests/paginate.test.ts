import { describe, expect, test } from "bun:test";
import {
  advanceQuery,
  fetchAllPages,
  nextCursor,
  nextPageNumber,
  resultsRows,
  stripAllFlag,
} from "../src/utils/paginate";

describe("paginate", () => {
  test("nextCursor extracts cursor from next URL", () => {
    expect(
      nextCursor({
        next: "https://assess.praxicraft.com/api/v1/public/assessments/?cursor=abc123&page_size=20",
        results: [],
      }),
    ).toBe("abc123");
    expect(nextCursor({ next: null })).toBe("");
  });

  test("nextPageNumber extracts page from next URL", () => {
    expect(nextPageNumber({ next: "https://example.com/results/?page=2" })).toBe("2");
  });

  test("advanceQuery sets cursor query param", () => {
    const query: Record<string, string> = {};
    expect(advanceQuery(query, { next: "https://x/?cursor=zz" })).toBe(true);
    expect(query.cursor).toBe("zz");
  });

  test("fetchAllPages merges results and clears next", async () => {
    let calls = 0;
    const out = await fetchAllPages(async (query) => {
      calls++;
      if (calls === 1) {
        expect(query.page_size).toBe("100");
        return {
          next: "https://x/?cursor=c2",
          results: [{ slug: "a" }, { slug: "b" }],
        };
      }
      expect(query.cursor).toBe("c2");
      return {
        next: null,
        results: [{ slug: "c" }],
      };
    });

    expect(resultsRows(out)).toHaveLength(3);
    expect((out as Record<string, unknown>).next).toBeNull();
    expect(calls).toBe(2);
  });

  test("stripAllFlag removes --all", () => {
    expect(stripAllFlag(["foo", "--all"])).toEqual({ args: ["foo"], all: true });
    expect(stripAllFlag(["foo"])).toEqual({ args: ["foo"], all: false });
  });
});
