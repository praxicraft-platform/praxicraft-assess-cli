import type { CommandContext } from "../tui/CommandContext";
import { clearActiveProfile, isSignedIn } from "../utils/config";

export async function handleLogout(ctx: CommandContext): Promise<void> {
  if (!isSignedIn()) {
    ctx.addBlock({ type: "error", message: "No stored profile to sign out from." });
    return;
  }
  const confirmed = await ctx.promptConfirm("Sign out and clear the active profile?");
  if (!confirmed) {
    ctx.addBlock({ type: "error", message: "Logout cancelled." });
    return;
  }
  clearActiveProfile();
  ctx.addBlock({ type: "success", message: "Signed out." });
}
