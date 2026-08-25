import type { CommandContext } from "./CommandContext";
import { ApiError, createClient, pathSegment } from "../utils/api";
import { extractOrgName, formatOutput } from "../utils/output";

function showApiResult(ctx: CommandContext, data: unknown, opts?: { whoami?: boolean }) {
  if (opts?.whoami) {
    ctx.addBlock({ type: "text", text: extractOrgName(data) });
    return;
  }
  ctx.addBlock({ type: "text", text: formatOutput(data, "table") });
}

const API_ROOTS = new Set([
  "/whoami",
  "/org",
  "/assessments",
  "/invites",
  "/results",
  "/cases",
  "/pipelines",
  "/webhooks",
  "/interviews",
  "/integrations",
]);

function formatApiError(e: unknown): string {
  if (e instanceof ApiError) {
    if (e.status === 401) {
      return "Sign in first. Run /login to get started.";
    }
    if (e.status === 403 && e.code === "INSUFFICIENT_SCOPE") {
      return e.message || "This API key is missing a required scope.";
    }
    if (e.status === 403 && e.code === "PLAN_REQUIRED") {
      return e.message || "Your plan does not include this feature.";
    }
    if (e.status === 429) {
      const wait = e.retryAfter ? ` Retry after ${e.retryAfter}s.` : "";
      return `Rate limited.${wait}`;
    }
    return e.message || `HTTP ${e.status}`;
  }
  const msg = e instanceof Error ? e.message : String(e);
  if (/Sign in first|Not signed in/i.test(msg)) {
    return "Sign in first. Run /login to get started.";
  }
  return msg;
}

export async function handleCommand(
  input: string,
  ctx: CommandContext,
  exit: () => void,
  refreshAuth: () => void,
) {
  const trimmedInput = input.trim();
  if (!trimmedInput) return;

  if (!trimmedInput.startsWith("/")) {
    const lower = trimmedInput.toLowerCase();
    if (lower === "exit" || lower === "quit") {
      exit();
      return;
    }
    if (lower === "clear") {
      ctx.clear();
      return;
    }
    const { handleAI } = await import("../commands/ai");
    await handleAI(trimmedInput, ctx);
    return;
  }

  const parts = trimmedInput.split(" ").filter(Boolean);
  const cmd = parts[0]!;
  const subCmd = parts[1];
  const extraArgs = parts.slice(2);

  if (cmd === "/help") {
    ctx.addBlock({ type: "help" });
    return;
  }

  if (cmd === "/clear") {
    ctx.clear();
    return;
  }

  if (cmd === "/exit") {
    exit();
    return;
  }

  if (cmd === "/login" || cmd === "/configure") {
    const { handleLogin } = await import("../commands/login");
    if (await handleLogin(ctx)) refreshAuth();
    return;
  }

  if (cmd === "/logout") {
    const { handleLogout } = await import("../commands/logout");
    await handleLogout(ctx);
    refreshAuth();
    return;
  }

  if (cmd === "/ai") {
    const { handleAI } = await import("../commands/ai");
    await handleAI(parts.slice(1).join(" "), ctx);
    return;
  }

  if (!API_ROOTS.has(cmd)) {
    ctx.addBlock({
      type: "error",
      message: `Unknown command: ${trimmedInput}. Try /help.`,
    });
    return;
  }

  try {
    const client = createClient();

    if (cmd === "/whoami") {
      const data = await client.get("/org/");
      showApiResult(ctx, data, { whoami: true });
      return;
    }

    if (cmd === "/org" && (!subCmd || subCmd === "get")) {
      const data = await client.get("/org/");
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/org" && subCmd === "billing") {
      const data = await client.get("/org/billing/");
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/org" && subCmd === "stats") {
      const data = await client.get("/org/stats/");
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/org") {
      ctx.addBlock({
        type: "error",
        message: "Usage: /org get|billing|stats",
      });
      return;
    }

    if (cmd === "/assessments" && subCmd === "list") {
      const data = await client.get("/assessments/");
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/assessments" && subCmd === "get") {
      const raw = extraArgs[0] || (await ctx.promptInput("Assessment slug"));
      const slug = pathSegment(raw);
      const data = await client.get(`/assessments/${slug}/`);
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/assessments") {
      ctx.addBlock({
        type: "error",
        message: "Usage: /assessments list|get <slug>",
      });
      return;
    }

    if (cmd === "/invites" && subCmd === "list") {
      const data = await client.get("/invites/");
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/invites" && subCmd === "create") {
      const slug = pathSegment(extraArgs[0] || (await ctx.promptInput("Assessment slug")));
      const email = (extraArgs[1] || (await ctx.promptInput("Candidate email"))).trim();
      const name = (extraArgs[2] || (await ctx.promptInput("Candidate name"))).trim();
      if (!email.includes("@")) {
        ctx.addBlock({ type: "error", message: "A valid candidate email is required." });
        return;
      }
      const data = await client.post(`/assessments/${slug}/invites/`, {
        email,
        name: name || email,
      });
      ctx.addBlock({ type: "success", message: "Invite created" });
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/invites") {
      ctx.addBlock({
        type: "error",
        message: "Usage: /invites list|create <slug> <email> <name>",
      });
      return;
    }

    if (cmd === "/results" && subCmd === "list") {
      const slug = pathSegment(extraArgs[0] || (await ctx.promptInput("Assessment slug")));
      const data = await client.get(`/assessments/${slug}/results/`);
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/results") {
      ctx.addBlock({ type: "error", message: "Usage: /results list <slug>" });
      return;
    }

    if (cmd === "/cases" && subCmd === "list") {
      const data = await client.get("/cases/");
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/pipelines" && subCmd === "list") {
      const data = await client.get("/pipelines/");
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/webhooks" && subCmd === "list") {
      const data = await client.get("/webhooks/");
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/interviews" && subCmd === "list") {
      const data = await client.get("/interviews/");
      showApiResult(ctx, data);
      return;
    }

    if (cmd === "/integrations" && subCmd === "list") {
      const data = await client.get("/integrations/");
      showApiResult(ctx, data);
      return;
    }

    ctx.addBlock({
      type: "error",
      message: `Unknown command: ${trimmedInput}. Try /help.`,
    });
  } catch (e: any) {
    ctx.addBlock({
      type: "error",
      message: formatApiError(e),
    });
  }
}
