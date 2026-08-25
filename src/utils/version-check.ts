import { version as currentVersion } from "../../package.json";

const NPM_LATEST_URL = "https://registry.npmjs.org/@praxicraft/assess-cli/latest";
const CHECK_TIMEOUT_MS = 4000;

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

export async function checkForUpdate(
  fetchImpl: typeof fetch = fetch,
): Promise<UpdateCheckResult | null> {
  try {
    const res = await fetchImpl(NPM_LATEST_URL, {
      signal: AbortSignal.timeout(CHECK_TIMEOUT_MS),
      headers: { Accept: "application/json" },
    });
    if (!res.ok) return null;

    const body = (await res.json()) as { version?: string };
    const latest = body.version?.trim();
    if (!latest || compareSemver(latest, currentVersion) <= 0) return null;

    return { current: currentVersion, latest };
  } catch {
    return null;
  }
}
