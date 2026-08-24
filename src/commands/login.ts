import type { CommandContext } from "../tui/CommandContext";
import { isValidApiKeyShape, saveProfile, DefaultBaseURL } from "../utils/config";

export async function handleLogin(ctx: CommandContext): Promise<boolean> {
  try {
    const apiKey = (await ctx.promptInput("API key (ct_live_… / ct_test_…)", true)).trim();
    if (!apiKey) {
      ctx.addBlock({ type: "error", message: "API key required." });
      return false;
    }
    if (!isValidApiKeyShape(apiKey)) {
      ctx.addBlock({
        type: "error",
        message:
          'API key must look like "ct_live_…" or "ct_test_…". Create one in Assess → API Keys.',
      });
      return false;
    }
    const baseRaw = (await ctx.promptInput(`API base URL [${DefaultBaseURL}]`)).trim();
    const baseURL = baseRaw || DefaultBaseURL;
    const profileRaw = (await ctx.promptInput("Profile name [default]")).trim();
    const profile = profileRaw || "default";
    saveProfile({ profile, apiKey, baseURL });
    ctx.addBlock({
      type: "success",
      message: `Signed in as profile "${profile}" (${apiKey.startsWith("ct_live_") ? "live" : "test"}).`,
    });
    return true;
  } catch (e: any) {
    ctx.addBlock({ type: "error", message: e?.message || String(e) });
    return false;
  }
}
