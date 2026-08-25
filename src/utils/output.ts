/** Human-friendly CLI/TUI output (table default, JSON when piped). */

export type OutputFormat = "table" | "json";

const MAX_CELL = 56;
const MAX_SLUG = 64;

const LIST_ENVELOPE_KEYS = [
  "results",
  "integrations",
  "items",
  "data",
  "webhooks",
  "deliveries",
  "enrollments",
  "members",
  "cases",
  "invites",
  "assessments",
  "pipelines",
  "interviews",
];

const OBJECT_KEY_ORDER = [
  "email",
  "name",
  "slug",
  "title",
  "status",
  "invite_token",
  "token",
  "take_url",
  "url",
  "id",
  "assessment_slug",
  "assessment",
  "created_at",
  "updated_at",
  "expires_at",
  "sent_at",
  "provider",
  "role",
  "plan",
  "livemode",
  "ats_allowed",
  "invite_limit",
  "invites_used",
  "invites_remaining",
];

const ROW_COLUMN_ORDER = [
  "slug",
  "email",
  "name",
  "title",
  "status",
  "provider",
  "is_active",
  "has_api_key",
  "auth_mode",
  "invite_token",
  "case_count",
  "time_limit_minutes",
  "invitation_count",
  "id",
  "connect_url",
  "webhook_url",
  "created_at",
];

const FULL_VALUE_KEYS = new Set(["invite_token", "token", "take_url", "url", "id"]);

export function defaultCliFormat(): OutputFormat {
  return process.stdout.isTTY ? "table" : "json";
}

export function parseOutputFlag(argv: string[]): { args: string[]; format?: OutputFormat } {
  const args = [...argv];
  let format: OutputFormat | undefined;
  for (let i = 0; i < args.length; i++) {
    if (args[i] === "--output" || args[i] === "-o") {
      const next = args[i + 1]?.toLowerCase();
      if (next === "json" || next === "table") {
        format = next;
        args.splice(i, 2);
        i--;
      }
    }
  }
  return { args, format };
}

export function extractOrgName(data: unknown): string {
  const obj = asRecord(unwrapApiPayload(data));
  if (!obj) return "(unknown organisation)";
  for (const key of ["name", "organisation_name", "org_name", "slug"]) {
    const value = obj[key];
    if (typeof value === "string" && value.trim()) return value.trim();
  }
  return "(unknown organisation)";
}

export function formatOutput(data: unknown, format: OutputFormat = "table"): string {
  const normalized = normalize(data);
  if (format === "json") {
    return JSON.stringify(normalized, null, 2);
  }
  return formatTable(normalized);
}

function normalize(value: unknown): unknown {
  if (value === undefined) return null;
  return JSON.parse(JSON.stringify(value));
}

/** Public API may nest payloads under `data` when present. */
function unwrapApiPayload(value: unknown): unknown {
  const normalized = normalize(value);
  if (!normalized || typeof normalized !== "object" || Array.isArray(normalized)) {
    return normalized;
  }
  const record = normalized as Record<string, unknown>;
  if (record.status === "success" && "data" in record) {
    return record.data;
  }
  return normalized;
}

function formatTable(data: unknown): string {
  const payload = unwrapApiPayload(data);

  if (payload === null || payload === undefined) return "(null)";
  if (typeof payload === "string") return payload;
  if (typeof payload === "number" || typeof payload === "boolean") return String(payload);

  if (Array.isArray(payload)) {
    return formatRows(payload);
  }

  if (typeof payload === "object") {
    const record = payload as Record<string, unknown>;
    const envelope = extractListEnvelope(record);
    if (envelope) {
      const lines = [formatRows(envelope.rows)];
      if (Object.keys(envelope.rest).length > 0) {
        lines.push("", "—", formatObject(envelope.rest));
      }
      return lines.join("\n");
    }
    return formatObject(record);
  }

  return JSON.stringify(payload, null, 2);
}

function extractListEnvelope(
  record: Record<string, unknown>,
): { rows: unknown[]; rest: Record<string, unknown> } | null {
  for (const key of LIST_ENVELOPE_KEYS) {
    const value = record[key];
    if (Array.isArray(value)) {
      return splitEnvelope(record, key, value);
    }
  }

  let foundKey: string | null = null;
  let rows: unknown[] | null = null;
  for (const [key, value] of Object.entries(record)) {
    if (!Array.isArray(value) || value.length === 0) continue;
    if (!value.every((row) => row && typeof row === "object" && !Array.isArray(row))) continue;
    if (foundKey) return null;
    foundKey = key;
    rows = value;
  }
  if (!foundKey || !rows) return null;
  return splitEnvelope(record, foundKey, rows);
}

function splitEnvelope(
  record: Record<string, unknown>,
  listKey: string,
  rows: unknown[],
): { rows: unknown[]; rest: Record<string, unknown> } {
  const rest: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(record)) {
    if (key !== listKey) rest[key] = value;
  }

  const meta = rest.meta;
  if (meta && typeof meta === "object" && !Array.isArray(meta)) {
    delete rest.meta;
    for (const [key, value] of Object.entries(meta as Record<string, unknown>)) {
      if (!(key in rest) && isScalar(value)) rest[key] = value;
    }
  }

  return { rows, rest };
}

function formatObject(record: Record<string, unknown>): string {
  if (Object.keys(record).length === 0) return "(empty)";

  const keys = orderedKeys(record, OBJECT_KEY_ORDER);
  const lines: string[] = [];

  for (const key of keys) {
    const value = record[key];
    if (isScalar(value)) {
      lines.push(`${key.padEnd(18)} ${formatObjectValue(key, value)}`);
      continue;
    }
    lines.push(`${key}`);
    lines.push(JSON.stringify(value, null, 2).split("\n").map((l) => `  ${l}`).join("\n"));
  }

  return lines.join("\n");
}

function formatRows(rows: unknown[]): string {
  if (rows.length === 0) return "(empty)";

  const objects = rows.filter(
    (row): row is Record<string, unknown> => !!row && typeof row === "object" && !Array.isArray(row),
  );
  if (objects.length === 0) {
    return rows.map((row) => String(row)).join("\n");
  }

  const columns = pickColumns(objects[0]!, ROW_COLUMN_ORDER);
  if (columns.length === 0) {
    return JSON.stringify(rows, null, 2);
  }

  const header = columns.map((c) => c.toUpperCase());
  const body = objects.map((row) => columns.map((col) => formatCell(col, row[col])));
  return renderTable([header, ...body]);
}

function pickColumns(first: Record<string, unknown>, preferred: string[]): string[] {
  const cols: string[] = [];
  const seen = new Set<string>();

  for (const key of preferred) {
    if (key in first && isScalar(first[key])) {
      cols.push(key);
      seen.add(key);
    }
  }

  if (cols.length < 4) {
    const extras = Object.keys(first)
      .filter((key) => !seen.has(key) && isScalar(first[key]))
      .sort();
    for (const key of extras) {
      if (cols.length >= 8) break;
      cols.push(key);
    }
  }

  return cols;
}

function orderedKeys(record: Record<string, unknown>, preferred: string[]): string[] {
  const keys: string[] = [];
  const seen = new Set<string>();
  for (const key of preferred) {
    if (key in record) {
      keys.push(key);
      seen.add(key);
    }
  }
  const rest = Object.keys(record)
    .filter((key) => !seen.has(key))
    .sort();
  return [...keys, ...rest];
}

function renderTable(rows: string[][]): string {
  if (rows.length === 0) return "";
  const widths = rows[0]!.map((_, col) =>
    Math.max(...rows.map((row) => visibleLength(row[col] ?? ""))),
  );
  return rows
    .map((row, rowIndex) =>
      row
        .map((cell, col) => {
          const padded = padVisible(cell ?? "", widths[col]!);
          return rowIndex === 0 ? padded : padded;
        })
        .join("  "),
    )
    .join("\n");
}

function visibleLength(text: string): number {
  return [...text].length;
}

function padVisible(text: string, width: number): string {
  const len = visibleLength(text);
  if (len >= width) return text;
  return text + " ".repeat(width - len);
}

function isScalar(value: unknown): boolean {
  return (
    value === null ||
    typeof value === "string" ||
    typeof value === "number" ||
    typeof value === "boolean"
  );
}

function formatObjectValue(key: string, value: unknown): string {
  if (value === null || value === undefined) return "";
  if (typeof value === "boolean") return value ? "true" : "false";
  if (typeof value === "number") return String(value);
  if (typeof value === "string") {
    return FULL_VALUE_KEYS.has(key) ? value : truncate(value, MAX_CELL);
  }
  return String(value);
}

function formatCell(column: string, value: unknown): string {
  if (value === null || value === undefined) return "";
  const text = typeof value === "boolean" ? (value ? "true" : "false") : String(value);
  switch (column) {
    case "id":
    case "invite_token":
      return shortenId(text);
    case "slug":
    case "email":
      return truncate(text, MAX_SLUG);
    default:
      return truncate(text, MAX_CELL);
  }
}

function shortenId(text: string): string {
  if (/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(text)) {
    return text.slice(0, 8);
  }
  if ([...text].length > 12) return truncate(text, 12);
  return text;
}

function truncate(text: string, max: number): string {
  const runes = [...text];
  if (runes.length <= max) return text;
  if (max <= 1) return "…";
  return runes.slice(0, max - 1).join("") + "…";
}

function asRecord(value: unknown): Record<string, unknown> | null {
  if (!value || typeof value !== "object" || Array.isArray(value)) return null;
  return value as Record<string, unknown>;
}
