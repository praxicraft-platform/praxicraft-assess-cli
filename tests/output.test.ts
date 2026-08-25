import { describe, expect, test } from "bun:test";
import {
  extractOrgName,
  formatOutput,
} from "../src/utils/output";

describe("extractOrgName", () => {
  test("prefers name over slug", () => {
    expect(extractOrgName({ name: "Acme Corp", slug: "acme" })).toBe("Acme Corp");
  });

  test("falls back through legacy keys", () => {
    expect(extractOrgName({ organisation_name: "Beta" })).toBe("Beta");
    expect(extractOrgName({ slug: "gamma" })).toBe("gamma");
  });

  test("handles missing fields", () => {
    expect(extractOrgName({})).toBe("(unknown organisation)");
  });
});

describe("formatOutput table", () => {
  test("formats object as key/value rows", () => {
    const text = formatOutput({ name: "Acme", slug: "acme", plan: "starter" }, "table");
    expect(text).toContain("name");
    expect(text).toContain("Acme");
    expect(text).toContain("slug");
    expect(text).toContain("acme");
  });

  test("formats list envelope as table", () => {
    const text = formatOutput(
      {
        results: [
          { slug: "backend-screen", status: "active", title: "Backend" },
          { slug: "data-pipeline", status: "draft", title: "Data" },
        ],
      },
      "table",
    );
    expect(text).toContain("SLUG");
    expect(text).toContain("backend-screen");
    expect(text).toContain("data-pipeline");
  });

  test("empty list shows (empty)", () => {
    expect(formatOutput({ results: [] }, "table")).toBe("(empty)");
  });

  test("null shows (null)", () => {
    expect(formatOutput(null, "table")).toBe("(null)");
  });

  test("nested non-scalars render as indented JSON", () => {
    const text = formatOutput({ name: "Acme", settings: { a: 1 } }, "table");
    expect(text).toContain("settings");
    expect(text).toContain('"a": 1');
  });

  test("unwraps success envelope", () => {
    const text = formatOutput(
      { status: "success", data: { name: "Wrapped Org", slug: "wrapped" } },
      "table",
    );
    expect(text).toContain("Wrapped Org");
  });

  test("json format unchanged", () => {
    const data = { slug: "x" };
    expect(formatOutput(data, "json")).toBe(JSON.stringify(data, null, 2));
  });

  test("aligns columns with variable-width cells", () => {
    const text = formatOutput(
      {
        results: [
          { slug: "short", status: "active", email: "a@example.com" },
          { slug: "much-longer-slug", status: "draft", email: "ada@example.com" },
        ],
      },
      "table",
    );
    const lines = text.split("\n");
    expect(lines[0]?.indexOf("STATUS")).toBeGreaterThan(lines[0]?.indexOf("SLUG") ?? -1);
    // Column starts should line up row-to-row.
    const slugCol = lines[1]?.search(/\S/) ?? 0;
    const slugCol2 = lines[2]?.search(/\S/) ?? 0;
    expect(slugCol2).toBe(slugCol);
  });

  test("object keys align to widest key", () => {
    const text = formatOutput(
      { invite_limit: 100, slug: "backend-screen", name: "Acme" },
      "table",
    );
    const lines = text.split("\n").filter(Boolean);
    const valueStarts = lines.map((line) => line.search(/\S+\s+\S/));
    expect(new Set(valueStarts).size).toBeLessThanOrEqual(2);
  });

  test("shortens uuid ids in list rows", () => {
    const text = formatOutput(
      {
        results: [{ id: "550e8400-e29b-41d4-a716-446655440000", slug: "test" }],
      },
      "table",
    );
    expect(text).toContain("550e8400");
    expect(text).not.toContain("550e8400-e29b-41d4-a716-446655440000");
  });
});
