/**
 * Ask AI — tool loop against the Assess AI proxy + dual MCP (docs + Assess).
 */
import { createRequire } from "node:module";
import { createMCPClient } from "@ai-sdk/mcp";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
import type { CommandContext } from "../tui/CommandContext";
import { ApiError, assistantChat } from "../utils/api";
import { resolveCredentials } from "../utils/config";

const require = createRequire(import.meta.url);

let cachedKnowledge: Awaited<ReturnType<typeof createMCPClient>> | null = null;
let cachedExec: Awaited<ReturnType<typeof createMCPClient>> | null = null;
let cachedTools: Record<string, any> | null = null;
let cachedKey: string | null = null;

const MAX_STEPS = 12;
const MCP_TIMEOUT_MS = 60_000;

function withTimeout<T>(promise: Promise<T>, ms: number, label: string): Promise<T> {
  return new Promise((resolve, reject) => {
    const timer = setTimeout(
      () => reject(new Error(`Timed out after ${Math.round(ms / 1000)}s waiting for ${label}`)),
      ms,
    );
    promise.then(
      (v) => {
        clearTimeout(timer);
        resolve(v);
      },
      (e) => {
        clearTimeout(timer);
        reject(e);
      },
    );
  });
}

function classifyError(error: any, phase: string): string {
  if (error instanceof ApiError) {
    if (error.status === 401 || error.code === "INVALID_API_KEY" || error.code === "EXPIRED_API_KEY") {
      return "Authentication failed. Your API key is invalid or expired. Run /login to refresh.";
    }
    if (error.status === 403) {
      if (error.code === "PLAN_REQUIRED") {
        return "Ask AI requires the Starter plan or above. Upgrade in Assess billing, then try again.";
      }
      if (error.code === "INSUFFICIENT_SCOPE") {
        return "This API key needs the assistant:write scope (or full access). Create a new key in Assess → API Keys.";
      }
      if (error.code === "IP_NOT_ALLOWED") {
        return "This API key is restricted by IP allowlist. Use a permitted network or update the key.";
      }
      return error.message || "Access denied for the assistant.";
    }
    if (error.status === 429) {
      const wait = error.retryAfter ? ` Retry after ${error.retryAfter}s.` : " Wait a moment and try again.";
      return `Assistant rate limit reached.${wait}`;
    }
    if (error.status === 503 || error.code === "AI_DISABLED") {
      return "AI features are temporarily disabled. Try again later.";
    }
    if (error.code === "AI_PROVIDER_ERROR" || error.status === 502) {
      return `Assistant provider error. Try again in a moment. (${phase})`;
    }
    if (error.code === "TIMEOUT" || error.status === 0) {
      return error.message.includes("reach")
        ? `${error.message} (${phase})`
        : `Timeout talking to the assistant. (${phase})`;
    }
    if (error.status >= 500) {
      return `Server error (HTTP ${error.status}). Try again in a moment. (${phase})`;
    }
    if (error.status >= 400) {
      return `Request rejected (HTTP ${error.status}). ${error.message} (${phase})`;
    }
    return `${error.message} (${phase})`;
  }

  const raw = error?.message ?? String(error);
  if (/Sign in first|Not signed in/i.test(raw)) {
    return "Sign in first. Run /login to get started.";
  }
  const code = error?.code ?? error?.cause?.code;
  if (code === "ECONNREFUSED" || code === "ENOTFOUND" || code === "EAI_AGAIN") {
    return `Couldn't reach the assistant. Check your network and try again. (${phase})`;
  }
  if (code === "ECONNRESET" || /socket hang up/i.test(raw)) {
    return `Connection reset. Try again in a moment. (${phase})`;
  }
  if (code === "ETIMEDOUT" || /timed? ?out/i.test(raw)) {
    return `Timeout. ${raw} (${phase})`;
  }
  if (/JSON|Unexpected token|parse/i.test(raw)) {
    return `The assistant returned invalid data. Try again. (${phase})`;
  }
  if (/spawn|ENOENT|npx|bunx/i.test(raw)) {
    return `Couldn't start MCP tools (${phase}). Ensure Node/npm (or Bun) is available, or install @praxicraft/assess-mcp.`;
  }
  if (/MCP|transport/i.test(raw)) {
    return `Communication error: ${raw.slice(0, 200)} (${phase})`;
  }
  if (/abort|AbortError/i.test(raw)) {
    return `Request aborted: ${raw.slice(0, 200)} (${phase})`;
  }
  return `${raw} (${phase})`;
}

function toolsToOpenAI(tools: Record<string, any>): Array<Record<string, unknown>> {
  return Object.entries(tools).map(([name, tool]) => ({
    type: "function",
    function: {
      name,
      description: tool.description || name,
      parameters: tool.parameters || { type: "object", properties: {} },
    },
  }));
}

function resolveAssessMcpSpawn(): { command: string; args: string[] } {
  // Prefer local package when present (faster, pinned); else npx.
  try {
    const pkg = require.resolve("@praxicraft/assess-mcp/package.json");
    const entry = require.resolve("@praxicraft/assess-mcp");
    if (pkg && entry) {
      return { command: process.execPath, args: [entry] };
    }
  } catch {
    /* fall through */
  }
  const runner = process.execPath.includes("bun") ? "bunx" : "npx";
  return { command: runner, args: ["-y", "@praxicraft/assess-mcp"] };
}

async function ensureClients(apiKey: string, baseURL: string) {
  if (!cachedKnowledge) {
    let mcpRemote: string;
    try {
      mcpRemote = require.resolve("mcp-remote/dist/proxy.js");
    } catch {
      mcpRemote = require.resolve("mcp-remote");
    }
    cachedKnowledge = await withTimeout(
      createMCPClient({
        transport: new StdioClientTransport({
          command: process.execPath,
          args: [mcpRemote, "https://docs.praxicraft.com/mcp"],
          env: { ...process.env } as Record<string, string>,
          stderr: "pipe",
        }),
      }),
      MCP_TIMEOUT_MS,
      "docs MCP",
    );
  }

  if (!cachedExec || cachedKey !== apiKey) {
    if (cachedExec) {
      try {
        await cachedExec.close();
      } catch {
        /* ignore */
      }
    }
    const spawn = resolveAssessMcpSpawn();
    cachedExec = await withTimeout(
      createMCPClient({
        transport: new StdioClientTransport({
          command: spawn.command,
          args: spawn.args,
          env: {
            ...(process.env as Record<string, string>),
            PRAXICRAFT_API_KEY: apiKey,
            PRAXICRAFT_API_BASE_URL: baseURL,
          },
          stderr: "pipe",
        }),
      }),
      MCP_TIMEOUT_MS,
      "Assess MCP",
    );
    cachedKey = apiKey;
    cachedTools = null;
  }

  if (!cachedTools) {
    const [kTools, eTools] = await withTimeout(
      Promise.all([cachedKnowledge.tools(), cachedExec.tools()]),
      MCP_TIMEOUT_MS,
      "MCP tool discovery",
    );
    const merged: Record<string, any> = { ...eTools };
    for (const [key, value] of Object.entries(kTools)) {
      merged[merged[key] ? `knowledge_${key}` : key] = value;
    }
    cachedTools = merged;
  }
  return cachedTools;
}

export async function handleAI(query: string, ctx: CommandContext) {
  if (!query.trim()) {
    ctx.addBlock({
      type: "error",
      message: "Question required. Type it directly or use /ai <question>.",
    });
    return;
  }

  const spinnerId = ctx.addBlock({ type: "spinner", label: "Initializing assistant…" });

  try {
    let creds;
    try {
      creds = resolveCredentials({ requireKey: true });
    } catch (e: any) {
      throw Object.assign(new Error(classifyError(e, "Auth")), { _classified: true });
    }

    ctx.updateBlock(spinnerId, { label: "Loading MCP tools…" });
    let tools: Record<string, any>;
    try {
      tools = await ensureClients(creds.apiKey, creds.baseURL);
    } catch (e: any) {
      throw Object.assign(new Error(classifyError(e, "MCP Init")), { _classified: true });
    }
    const openaiTools = toolsToOpenAI(tools);

    const messages: Array<Record<string, unknown>> = [
      { role: "user", content: query },
    ];

    let finalText = "";
    for (let step = 0; step < MAX_STEPS; step++) {
      ctx.updateBlock(spinnerId, {
        label: step === 0 ? "Analyzing your question…" : `Analyzing… (step ${step + 1})`,
      });

      let result;
      try {
        result = await withTimeout(
          assistantChat(creds, {
            messages,
            tools: openaiTools,
            max_tokens: 2048,
          }),
          MCP_TIMEOUT_MS,
          `LLM step ${step + 1}`,
        );
      } catch (e: any) {
        throw Object.assign(
          new Error(classifyError(e, step === 0 ? "LLM Request" : `LLM Step ${step + 1}`)),
          { _classified: true },
        );
      }

      const msg = result.message || {};
      if (msg.content) finalText += String(msg.content);

      const toolCalls = msg.tool_calls as
        | Array<{ id: string; function: { name: string; arguments: string } }>
        | undefined;

      if (result.finish_reason !== "tool_calls" || !toolCalls?.length) {
        break;
      }

      messages.push({
        role: "assistant",
        content: msg.content || "",
        tool_calls: toolCalls,
      });

      for (const tc of toolCalls) {
        const name = tc.function?.name;
        const tool = tools[name];
        let args: Record<string, unknown> = {};
        try {
          args = JSON.parse(tc.function?.arguments || "{}");
        } catch {
          args = {};
        }
        let toolResult = "";
        try {
          if (!tool?.execute) {
            toolResult = `Tool ${name} not executable`;
          } else {
            const out = await tool.execute(args);
            toolResult =
              typeof out === "string" ? out : JSON.stringify(out, null, 2);
          }
        } catch (e: any) {
          toolResult = `Tool error: ${e?.message || e}`;
        }
        messages.push({
          role: "tool",
          tool_call_id: tc.id,
          name,
          content: toolResult.slice(0, 12_000),
        });
      }
    }

    if (finalText) {
      ctx.addBlock({ type: "markdown", text: finalText });
    } else {
      ctx.addBlock({
        type: "error",
        message: "No response from the assistant. Try rephrasing.",
      });
    }
  } catch (e: any) {
    const message = e?._classified ? e.message : classifyError(e, "Unknown");
    ctx.addBlock({ type: "error", message });
  } finally {
    ctx.removeBlock(spinnerId);
  }
}
