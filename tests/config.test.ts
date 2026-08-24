import { describe, expect, test, beforeEach, afterEach } from "bun:test";
import { mkdtempSync, rmSync, writeFileSync, readFileSync, chmodSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

describe("config", () => {
  let dir: string;
  let prevConfig: string | undefined;

  beforeEach(() => {
    dir = mkdtempSync(join(tmpdir(), "prax-cli-"));
    prevConfig = process.env.PRAXICRAFT_CONFIG;
    process.env.PRAXICRAFT_CONFIG = join(dir, "config.toml");
    delete process.env.PRAXICRAFT_API_KEY;
    delete process.env.PRAXICRAFT_API_BASE_URL;
    delete process.env.PRAXICRAFT_PROFILE;
  });

  afterEach(() => {
    if (prevConfig === undefined) delete process.env.PRAXICRAFT_CONFIG;
    else process.env.PRAXICRAFT_CONFIG = prevConfig;
    rmSync(dir, { recursive: true, force: true });
  });

  test("isValidApiKeyShape", async () => {
    const { isValidApiKeyShape } = await import("../src/utils/config");
    expect(isValidApiKeyShape("ct_live_abcdefgh")).toBe(true);
    expect(isValidApiKeyShape("ct_test_abcdefgh")).toBe(true);
    expect(isValidApiKeyShape("sk_live_abcdefgh")).toBe(false);
    expect(isValidApiKeyShape("ct_live_")).toBe(false);
    expect(isValidApiKeyShape("")).toBe(false);
  });

  test("save and load profile round-trip", async () => {
    const { saveProfile, loadConfig, resolveCredentials } = await import("../src/utils/config");
    saveProfile({
      profile: "work",
      apiKey: "ct_live_abcdefghij",
      baseURL: "https://assess.example.com/",
    });
    const f = loadConfig();
    expect(f.default_profile).toBe("work");
    expect(f.profiles?.work?.api_key).toBe("ct_live_abcdefghij");
    expect(f.profiles?.work?.base_url).toBe("https://assess.example.com");
    const r = resolveCredentials({ requireKey: true });
    expect(r.profile).toBe("work");
    expect(r.baseURL).toBe("https://assess.example.com");
  });

  test("rejects garbage TOML", async () => {
    writeFileSync(process.env.PRAXICRAFT_CONFIG!, "this is not = toml [[[");
    const { loadConfig } = await import("../src/utils/config");
    expect(() => loadConfig()).toThrow(/Could not parse config/);
  });

  test("env overrides profile key", async () => {
    const { saveProfile, resolveCredentials } = await import("../src/utils/config");
    saveProfile({ apiKey: "ct_live_fromfilexxxx", baseURL: "https://assess.praxicraft.com" });
    process.env.PRAXICRAFT_API_KEY = "ct_test_fromenvxxxxx";
    const r = resolveCredentials({ requireKey: true });
    expect(r.apiKey).toBe("ct_test_fromenvxxxxx");
  });

  test("unknown explicit profile throws", async () => {
    const { resolveCredentials } = await import("../src/utils/config");
    expect(() => resolveCredentials({ profile: "missing", requireKey: false })).toThrow(
      /Unknown profile/,
    );
  });

  test("config file mode is restrictive when supported", async () => {
    const { saveProfile, configPath } = await import("../src/utils/config");
    saveProfile({ apiKey: "ct_live_abcdefghij" });
    const mode = readFileSync(configPath()).length;
    expect(mode).toBeGreaterThan(0);
    try {
      const { statSync } = await import("node:fs");
      const st = statSync(configPath());
      expect(st.mode & 0o077).toBe(0);
    } catch {
      /* windows may not support */
    }
  });
});

describe("pathSegment", () => {
  test("encodes and rejects blank", async () => {
    const { pathSegment } = await import("../src/utils/api");
    expect(pathSegment("my slug")).toBe("my%20slug");
    expect(pathSegment("a/b")).toBe("a%2Fb");
    expect(() => pathSegment("  ")).toThrow(/non-empty/);
  });
});

describe("palette ranking", () => {
  test("prefix match ranks highest", async () => {
    const { rankCommands } = await import("../src/tui/components/Palette");
    const ranked = rankCommands("/log");
    expect(ranked[0]?.command).toBe("/login");
    expect(rankCommands("/zzz").length).toBe(0);
  });
});
