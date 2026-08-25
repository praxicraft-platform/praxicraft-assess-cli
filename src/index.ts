#!/usr/bin/env bun
import { spawnSync } from "node:child_process";
import { homedir } from "node:os";
import { join } from "node:path";
import { fileURLToPath } from "node:url";
import { version } from "../package.json";
import { createClient } from "./utils/api";
import {
  clearActiveProfile,
  resolveCredentials,
  saveProfile,
  DefaultBaseURL,
} from "./utils/config";
import {
  defaultCliFormat,
  extractOrgName,
  formatOutput,
  parseOutputFlag,
} from "./utils/output";

const parsed = parseOutputFlag(process.argv.slice(2));
const args = parsed.args;
const cliOutputFormat = parsed.format ?? defaultCliFormat();

function printApiData(data: unknown, opts?: { whoami?: boolean }) {
  if (opts?.whoami) {
    console.log(extractOrgName(data));
    return;
  }
  console.log(formatOutput(data, cliOutputFormat));
}

function printHelp() {
  console.log(`praxicraft-assess v${version}

Usage:
  praxicraft-assess                  Interactive TUI (Ask AI + command palette)
  praxicraft-assess version
  praxicraft-assess login|configure  Save API key profile
  praxicraft-assess logout
  praxicraft-assess whoami
  praxicraft-assess org get|billing|stats
  praxicraft-assess assessments list|get <slug>
  praxicraft-assess invites list|create <slug> <email> <name>
  praxicraft-assess results list <slug>
  praxicraft-assess cases|pipelines|webhooks|interviews|integrations list

Output: table in terminal (default), JSON when piped. Override with --output json|table

Env: PRAXICRAFT_API_KEY, PRAXICRAFT_API_BASE_URL, PRAXICRAFT_PROFILE
Config: ~/.config/praxicraft/config.toml
Docs: https://docs.praxicraft.com/sdks/cli
`);
}

function isBunRuntime(): boolean {
  return typeof (globalThis as any).Bun !== "undefined";
}

function relaunchTuiWithBun(): never {
  const executable = process.platform === "win32" ? "bun.exe" : "bun";
  const candidates = [
    executable,
    join(homedir(), ".bun", "bin", executable),
    "/usr/local/bin/bun",
    "/opt/homebrew/bin/bun",
  ];
  const self = fileURLToPath(import.meta.url);
  for (const bunPath of candidates) {
    const child = spawnSync(bunPath, [self, ...process.argv.slice(2)], {
      stdio: "inherit",
      env: process.env,
    });
    if (!child.error) {
      if (child.signal) process.kill(process.pid, child.signal);
      process.exit(child.status ?? 1);
    }
    if ((child.error as NodeJS.ErrnoException).code !== "ENOENT") {
      console.error(`Failed to launch Bun: ${child.error.message}`);
      process.exit(1);
    }
  }
  console.error(
    "The interactive Assess TUI requires the Bun runtime.\n" +
      "\n" +
      "Install Bun:    https://bun.com/docs/installation\n" +
      "Subcommands such as `praxicraft-assess whoami` work without the TUI.\n" +
      "Or download a standalone binary:\n" +
      "                https://github.com/praxicraft-platform/praxicraft-assess-cli/releases\n",
  );
  process.exit(1);
}

async function runCli(): Promise<void> {
  const wantsTui = args.length === 0 || args[0] === "interactive";
  if (wantsTui && !isBunRuntime()) {
    relaunchTuiWithBun();
  }

  if (wantsTui) {
    const { startTui } = await import("./tui/app");
    await startTui();
    return;
  }

  const cmd = args[0];

  if (cmd === "version" || cmd === "--version" || cmd === "-V") {
    console.log(`praxicraft-assess version ${version}`);
    return;
  }

  if (cmd === "help" || cmd === "--help" || cmd === "-h") {
    printHelp();
    return;
  }

  if (cmd === "configure" || cmd === "login") {
    const { isValidApiKeyShape } = await import("./utils/config");
    const apiKey = (prompt("API key (ct_live_… / ct_test_…):") || "").trim();
    if (!apiKey) {
      console.error("API key required");
      process.exit(1);
    }
    if (!isValidApiKeyShape(apiKey)) {
      console.error('API key must look like "ct_live_…" or "ct_test_…"');
      process.exit(1);
    }
    const baseURL = (prompt(`API base URL [${DefaultBaseURL}]:`) || DefaultBaseURL).trim();
    const profile = (prompt("Profile name [default]:") || "default").trim() || "default";
    try {
      saveProfile({ profile, apiKey, baseURL });
    } catch (e: any) {
      console.error(e?.message || e);
      process.exit(1);
    }
    console.log(`Saved profile "${profile}"`);
    return;
  }

  if (cmd === "logout") {
    clearActiveProfile();
    console.log("Signed out.");
    return;
  }

  try {
    if (cmd === "whoami") {
      const client = createClient();
      printApiData(await client.get("/org/"), { whoami: true });
      return;
    }

    if (cmd === "org") {
      const client = createClient();
      const sub = args[1] || "get";
      if (!["get", "billing", "stats"].includes(sub)) {
        throw new Error("usage: org get|billing|stats");
      }
      const path =
        sub === "billing" ? "/org/billing/" : sub === "stats" ? "/org/stats/" : "/org/";
      printApiData(await client.get(path));
      return;
    }

    if (cmd === "assessments") {
      const client = createClient();
      const sub = args[1];
      if (sub === "list") {
        printApiData(await client.get("/assessments/"));
        return;
      }
      if (sub === "get") {
        const { pathSegment } = await import("./utils/api");
        const slug = pathSegment(args[2] || "");
        printApiData(await client.get(`/assessments/${slug}/`));
        return;
      }
      throw new Error("usage: assessments list|get <slug>");
    }

    if (cmd === "invites") {
      const client = createClient();
      const sub = args[1];
      if (sub === "list") {
        printApiData(await client.get("/invites/"));
        return;
      }
      if (sub === "create") {
        const { pathSegment } = await import("./utils/api");
        const [, , slugArg, email, name] = args;
        if (!slugArg || !email) throw new Error("usage: invites create <slug> <email> [name]");
        if (!email.includes("@")) throw new Error("a valid candidate email is required");
        const slug = pathSegment(slugArg);
        printApiData(
          await client.post(`/assessments/${slug}/invites/`, {
            email,
            name: name || email,
          }),
        );
        return;
      }
      throw new Error("usage: invites list|create <slug> <email> [name]");
    }

    const listMap: Record<string, string> = {
      cases: "/cases/",
      pipelines: "/pipelines/",
      webhooks: "/webhooks/",
      interviews: "/interviews/",
      integrations: "/integrations/",
    };

    if (cmd === "results") {
      if (args[1] !== "list") throw new Error("usage: results list <slug>");
      const { pathSegment } = await import("./utils/api");
      const slug = pathSegment(args[2] || "");
      const client = createClient();
      printApiData(await client.get(`/assessments/${slug}/results/`));
      return;
    }

    if (cmd in listMap) {
      if (args[1] !== "list") throw new Error(`usage: ${cmd} list`);
      const client = createClient();
      printApiData(await client.get(listMap[cmd]!));
      return;
    }

    resolveCredentials({ requireKey: false });
    console.error(`Unknown command: ${cmd}`);
    printHelp();
    process.exit(1);
  } catch (e: any) {
    console.error(e?.message || e);
    process.exit(1);
  }
}

await runCli();
