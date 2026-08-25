import { version as currentVersion } from "../../package.json";

const NPM_LATEST_URL = "https://registry.npmjs.org/@praxicraft/assess-cli/latest";
const GITHUB_LATEST_URL =
  "https://api.github.com/repos/praxicraft-platform/praxicraft-assess-cli/releases/latest";
const CHECK_TIMEOUT_MS = 5000;

export type UpdateCheckResult = {
  current: string;
  latest: string;
};

/** Compare semver strings (major.minor.patch). Returns 1 if a > b, -1 if a < b, 0 if equal. */
export function compareSemver(a: string, b: string): number {
  const parse = (v: string) =>
    v
      .trim()
      .replace(/^v/i, "")
      .split(/[.-]/)
      .slice(0, 3)
      .map((part) => Number.parseInt(part.replace(/[^\d]/g, ""), 10) || 0);

  const [aMaj = 0, aMin = 0, aPatch = 0] = parse(a);
  const [bMaj = 0, bMin = 0, bPatch = 0] = parse(b);

  if (aMaj !== bMaj) return aMaj > bMaj ? 1 : -1;
  if (aMin !== bMin) return aMin > bMin ? 1 : -1;
  if (aPatch !== bPatch) return aPatch > bPatch ? 1 : -1;
  return 0;
}

function pickNewestVersion(candidates: string[]): string | null {
  const trimmed = candidates.map((v) => v.trim()).filter(Boolean);
  if (trimmed.length === 0) return null;
  return trimmed.reduce((best, next) => (compareSemver(next, best) > 0 ? next : best));
}

async function fetchNpmLatest(fetchImpl: typeof fetch): Promise<string | null> {
  try {
    const res = await fetchImpl(NPM_LATEST_URL, {
      signal: AbortSignal.timeout(CHECK_TIMEOUT_MS),
      headers: { Accept: "application/json" },
    });
    if (!res.ok) return null;
    const body = (await res.json()) as { version?: string };
    return body.version?.trim() || null;
  } catch {
    return null;
  }
}

async function fetchGithubReleaseLatest(fetchImpl: typeof fetch): Promise<string | null> {
  try {
    const res = await fetchImpl(GITHUB_LATEST_URL, {
      signal: AbortSignal.timeout(CHECK_TIMEOUT_MS),
      headers: {
        Accept: "application/vnd.github+json",
        "User-Agent": "praxicraft-assess-cli",
      },
    });
    if (!res.ok) return null;
    const body = (await res.json()) as { tag_name?: string };
    return body.tag_name?.trim().replace(/^v/i, "") || null;
  } catch {
    return null;
  }
}

export async function checkForUpdate(
  fetchImpl: typeof fetch = fetch,
): Promise<UpdateCheckResult | null> {
  const [npmLatest, githubLatest] = await Promise.all([
    fetchNpmLatest(fetchImpl),
    fetchGithubReleaseLatest(fetchImpl),
  ]);
  const latest = pickNewestVersion([npmLatest, githubLatest].filter((v): v is string => !!v));
  if (!latest || compareSemver(latest, currentVersion) <= 0) return null;
  return { current: currentVersion, latest };
}

export function formatUpdateNotice(result: UpdateCheckResult): string {
  return `Update available: v${result.latest} (you have v${result.current}). Run: npm i -g @praxicraft/assess-cli`;
}

/** Retry once after startup in case the first registry fetch is slow or flaky. */
export const UPDATE_RETRY_DELAY_MS = 30_000;

/** Poll interval while the TUI stays open (stops once an update is found). */
export const UPDATE_POLL_INTERVAL_MS = 5 * 60 * 1000;

/**
 * Check npm/GitHub for a newer CLI in the background.
 * Returns a stop function; polls until an update is found or the TUI exits.
 */
export function startBackgroundUpdateCheck(onUpdate: (latest: string) => void): () => void {
  let stopped = false;
  let retryTimer: ReturnType<typeof setTimeout> | undefined;
  let pollTimer: ReturnType<typeof setInterval> | undefined;

  const stop = () => {
    if (stopped) return;
    stopped = true;
    if (retryTimer !== undefined) clearTimeout(retryTimer);
    if (pollTimer !== undefined) clearInterval(pollTimer);
  };

  const run = async () => {
    if (stopped) return;
    const result = await checkForUpdate();
    if (stopped || !result) return;
    onUpdate(result.latest);
    stop();
  };

  void run();
  retryTimer = setTimeout(() => void run(), UPDATE_RETRY_DELAY_MS);
  pollTimer = setInterval(() => void run(), UPDATE_POLL_INTERVAL_MS);

  return stop;
}
