import { ApiError } from "./api";

/** Cap --all walks to avoid runaway loops (matches legacy Go CLI). */
export const MAX_LIST_PAGES = 250;

const DEFAULT_PAGE_SIZE = "100";

function unwrapApiPayload(value: unknown): unknown {
  if (!value || typeof value !== "object" || Array.isArray(value)) return value;
  const record = value as Record<string, unknown>;
  if (record.status === "success" && "data" in record) return record.data;
  return value;
}

function asRecord(value: unknown): Record<string, unknown> | null {
  if (!value || typeof value !== "object" || Array.isArray(value)) return null;
  return value as Record<string, unknown>;
}

/** Extract list rows from a Public API page payload. */
export function resultsRows(page: unknown): Record<string, unknown>[] {
  const payload = unwrapApiPayload(page);
  if (Array.isArray(payload)) {
    return payload.filter(
      (row): row is Record<string, unknown> =>
        !!row && typeof row === "object" && !Array.isArray(row),
    );
  }
  const record = asRecord(payload);
  if (!record) return [];
  const results = record.results;
  if (!Array.isArray(results)) return [];
  return results.filter(
    (row): row is Record<string, unknown> =>
      !!row && typeof row === "object" && !Array.isArray(row),
  );
}

export function nextCursor(page: unknown): string {
  const record = asRecord(unwrapApiPayload(page));
  const next = record?.next;
  if (typeof next !== "string" || !next.trim()) return "";
  try {
    return new URL(next).searchParams.get("cursor")?.trim() ?? "";
  } catch {
    return "";
  }
}

export function nextPageNumber(page: unknown): string {
  const record = asRecord(unwrapApiPayload(page));
  const next = record?.next;
  if (typeof next !== "string" || !next.trim()) return "";
  try {
    return new URL(next).searchParams.get("page")?.trim() ?? "";
  } catch {
    return "";
  }
}

export function advanceQuery(
  query: Record<string, string>,
  page: unknown,
): boolean {
  const cursor = nextCursor(page);
  if (cursor) {
    query.cursor = cursor;
    delete query.page;
    return true;
  }
  const pageNum = nextPageNumber(page);
  if (pageNum) {
    query.page = pageNum;
    delete query.cursor;
    return true;
  }
  return false;
}

function rowIdentity(row: Record<string, unknown>): string {
  for (const key of ["id", "invite_token", "slug", "token"]) {
    const value = row[key];
    if (typeof value === "string" && value.trim()) return `${key}:${value.trim()}`;
  }
  return "";
}

/** Follow cursor/page links and merge `results` into one payload. */
export async function fetchAllPages(
  fetchPage: (query: Record<string, string>) => Promise<unknown>,
  initialQuery: Record<string, string> = {},
): Promise<unknown> {
  const query = { ...initialQuery };
  if (!query.page_size) query.page_size = DEFAULT_PAGE_SIZE;

  const merged: Record<string, unknown>[] = [];
  const seen = new Set<string>();
  let extra: Record<string, unknown> | null = null;

  for (let pageIndex = 0; pageIndex < MAX_LIST_PAGES; pageIndex++) {
    const raw = await fetchPage({ ...query });
    const payload = unwrapApiPayload(raw);
    const record = asRecord(payload);

    if (record && !extra) {
      extra = {};
      for (const [key, value] of Object.entries(record)) {
        if (key === "results" || key === "next" || key === "previous") continue;
        extra[key] = value;
      }
    }

    const rows = resultsRows(raw);
    if (rows.length === 0 && Array.isArray(payload)) {
      for (const row of payload) {
        if (row && typeof row === "object" && !Array.isArray(row)) {
          merged.push(row as Record<string, unknown>);
        }
      }
      break;
    }

    for (const row of rows) {
      const id = rowIdentity(row);
      if (id) {
        if (seen.has(id)) continue;
        seen.add(id);
      }
      merged.push(row);
    }

    if (!advanceQuery(query, raw)) break;

    if (pageIndex === MAX_LIST_PAGES - 1) {
      throw new ApiError(
        `Stopped after ${MAX_LIST_PAGES} pages; narrow the list or increase filters.`,
        0,
        "PAGINATION_LIMIT",
      );
    }
  }

  return {
    next: null,
    previous: null,
    results: merged,
    ...(extra ?? {}),
  };
}

export function stripAllFlag(args: string[]): { args: string[]; all: boolean } {
  const all = args.includes("--all");
  return { args: args.filter((arg) => arg !== "--all"), all };
}

export async function fetchList(
  fetchPage: (query: Record<string, string>) => Promise<unknown>,
  opts?: { all?: boolean; query?: Record<string, string> },
): Promise<unknown> {
  const query = opts?.query ?? {};
  if (opts?.all) return fetchAllPages(fetchPage, query);
  return fetchPage(query);
}
