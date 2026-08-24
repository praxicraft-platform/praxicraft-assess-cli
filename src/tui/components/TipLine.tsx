import { Show } from "solid-js";
import { colors, glyphs } from "../theme";
import { isSignedIn } from "../../utils/config";

interface Tip {
  prefix: string;
  command: string;
  suffix: string;
}

const TIPS_OUT: Tip[] = [
  { prefix: "Type ", command: "/help", suffix: " to see every command" },
  { prefix: "Use ", command: "/login", suffix: " to connect your API key" },
  { prefix: "Hit ", command: "/", suffix: " to open the command palette" },
  { prefix: "Press ", command: "Tab", suffix: " to autocomplete commands" },
];

const TIPS_IN: Tip[] = [
  { prefix: "Type ", command: "/help", suffix: " to see every command" },
  { prefix: "Run ", command: "/assessments list", suffix: " to view assessments" },
  { prefix: "Try ", command: "/org billing", suffix: " to check entitlements" },
  { prefix: "Ask in plain English", command: "", suffix: " — or hit / for commands" },
  { prefix: "Press ", command: "Tab", suffix: " to autocomplete commands" },
];

const pool = isSignedIn() ? TIPS_IN : TIPS_OUT;
const tip = pool[Math.floor(Math.random() * pool.length)]!;

export const TipLine = () => (
  <box flexDirection="row" justifyContent="center" alignItems="center" flexShrink={0}>
    <text fg={colors.brandLime}>{glyphs.arrow}</text>
    <text fg={colors.textMuted}>{` ${tip.prefix}`}</text>
    <Show when={tip.command}>
      <box paddingLeft={1} paddingRight={1} backgroundColor={colors.surfaceDeep} flexShrink={0}>
        <text fg={colors.brandLime}>{tip.command}</text>
      </box>
    </Show>
    <text fg={colors.textMuted}>{tip.suffix}</text>
  </box>
);
