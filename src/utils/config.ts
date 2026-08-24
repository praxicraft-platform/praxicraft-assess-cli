import { homedir } from "node:os";
import { dirname, join } from "node:path";
import {
  mkdirSync,
  readFileSync,
  writeFileSync,
  existsSync,
  chmodSync,
} from "node:fs";

export const EnvAPIKey = "PRAXICRAFT_API_KEY";
export const EnvBaseURL = "PRAXICRAFT_API_BASE_URL";
export const EnvProfile = "PRAXICRAFT_PROFILE";
export const DefaultBaseURL = "https://assess.praxicraft.com";

export type Profile = { api_key?: string; base_url?: string };
export type ConfigFile = {
  default_profile?: string;
  profiles?: Record<string, Profile>;
};

export type Resolved = {
  profile: string;
  apiKey: string;
  baseURL: string;
};

export function configPath(): string {
  if (process.env.PRAXICRAFT_CONFIG) return process.env.PRAXICRAFT_CONFIG;
  return join(homedir(), ".config", "praxicraft", "config.toml");
}

function parseTomlSimple(raw: string): ConfigFile {
  const out: ConfigFile = { profiles: {} };
  let section = "";
  let sawKv = false;
  for (const line of raw.split(/\r?\n/)) {
    const t = line.trim();
    if (!t || t.startsWith("#")) continue;
    const sec = t.match(/^\[([^\]]+)\]$/);
    if (sec) {
      section = sec[1]!;
      continue;
    }
    // Accept double- or single-quoted values
    const kv = t.match(/^([A-Za-z0-9_]+)\s*=\s*(?:"([^"]*)"|'([^']*)')$/);
    if (!kv) continue;
    sawKv = true;
    const key = kv[1]!;
    const value = (kv[2] ?? kv[3] ?? "") as string;
    if (!section && key === "default_profile") {
      out.default_profile = value;
      continue;
    }
    const m = section.match(/^profiles\.(.+)$/);
    if (m) {
      const name = m[1]!;
      out.profiles ??= {};
      out.profiles[name] ??= {};
      if (key === "api_key") out.profiles[name].api_key = value;
      if (key === "base_url") out.profiles[name].base_url = value;
    }
  }
  const nonempty = raw.replace(/#[^\n]*/g, "").trim().length > 0;
  if (nonempty && !sawKv && Object.keys(out.profiles ?? {}).length === 0) {
    throw new Error(
      `Could not parse config at ${configPath()}. Expected TOML like:\n` +
        `  default_profile = "default"\n` +
        `  [profiles.default]\n` +
        `  api_key = "ct_live_…"\n` +
        `  base_url = "https://assess.praxicraft.com"`,
    );
  }
  return out;
}

function encodeToml(f: ConfigFile): string {
  const lines: string[] = [];
  if (f.default_profile) lines.push(`default_profile = "${f.default_profile}"`);
  lines.push("");
  for (const [name, p] of Object.entries(f.profiles ?? {})) {
    lines.push(`[profiles.${name}]`);
    if (p.api_key) lines.push(`api_key = "${p.api_key}"`);
    if (p.base_url) lines.push(`base_url = "${p.base_url}"`);
    lines.push("");
  }
  return lines.join("\n");
}

export function loadConfig(): ConfigFile {
  const path = configPath();
  if (!existsSync(path)) return { profiles: {} };
  return parseTomlSimple(readFileSync(path, "utf8"));
}

export function saveConfig(f: ConfigFile): void {
  const path = configPath();
  const dir = dirname(path);
  mkdirSync(dir, { recursive: true, mode: 0o700 });
  try {
    chmodSync(dir, 0o700);
  } catch {
    /* best-effort on platforms that ignore mode */
  }
  writeFileSync(path, encodeToml(f), { mode: 0o600 });
  try {
    chmodSync(path, 0o600);
  } catch {
    /* ignore */
  }
}

function first(...vals: Array<string | undefined | null>): string {
  for (const v of vals) {
    if (v && v.trim()) return v.trim();
  }
  return "";
}

/** Loose shape check for Assess public API keys. */
export function isValidApiKeyShape(key: string): boolean {
  const k = key.trim();
  return /^(ct_live_|ct_test_)[A-Za-z0-9_-]{8,}$/.test(k);
}

export function resolveCredentials(opts?: {
  profile?: string;
  apiKey?: string;
  baseURL?: string;
  requireKey?: boolean;
}): Resolved {
  const f = loadConfig();
  const explicitProfile = first(opts?.profile, process.env[EnvProfile]);
  const profile = first(explicitProfile, f.default_profile, "default");
  const p = f.profiles?.[profile] ?? {};
  if (explicitProfile && !f.profiles?.[profile]) {
    throw new Error(`Unknown profile "${profile}". Run /login or configure.`);
  }
  const apiKey = first(opts?.apiKey, process.env[EnvAPIKey], p.api_key);
  const baseURL = first(opts?.baseURL, process.env[EnvBaseURL], p.base_url, DefaultBaseURL);
  if (opts?.requireKey !== false && !apiKey) {
    throw new Error("Sign in first. Run /login to get started.");
  }
  let normalizedBase = baseURL.replace(/\/$/, "");
  if (normalizedBase && !/^https?:\/\//i.test(normalizedBase)) {
    throw new Error(`Invalid base URL "${normalizedBase}". Include https://`);
  }
  return { profile, apiKey, baseURL: normalizedBase };
}

export function isSignedIn(): boolean {
  try {
    resolveCredentials({ requireKey: true });
    return true;
  } catch {
    return false;
  }
}

export function saveProfile(opts: {
  profile?: string;
  apiKey: string;
  baseURL?: string;
}): void {
  const key = opts.apiKey.trim();
  if (!isValidApiKeyShape(key)) {
    throw new Error(
      'API key must look like "ct_live_…" or "ct_test_…". Create one in Assess → API Keys.',
    );
  }
  let base = (opts.baseURL || DefaultBaseURL).trim().replace(/\/$/, "");
  if (!/^https?:\/\//i.test(base)) {
    throw new Error(`Invalid base URL "${base}". Include https://`);
  }
  const f = loadConfig();
  const name = (opts.profile || f.default_profile || "default").trim() || "default";
  f.profiles ??= {};
  f.profiles[name] = {
    api_key: key,
    base_url: base,
  };
  f.default_profile = name;
  saveConfig(f);
}

export function clearActiveProfile(): void {
  const f = loadConfig();
  const name = f.default_profile || "default";
  if (f.profiles?.[name]) {
    delete f.profiles[name];
  }
  saveConfig(f);
}
