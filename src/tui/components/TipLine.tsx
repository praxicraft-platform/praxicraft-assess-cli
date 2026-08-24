import { colors, glyphs } from "../theme";
import { isSignedIn } from "../../utils/config";

interface Tip {
  prefix: string;
  command: string;
  suffix: string;
}

const TIPS_OUT: Tip[] = [
  { prefix: "Type ", command: "/help", suffix: " to see every command" },
  { prefix: "Use ", command: "/login", suffix: " to start using the API" },
  { prefix: "Hit ", command: "/", suffix: " to open the command palette" },
  { prefix: "Press ", command: "Tab", suffix: " to autocomplete commands" },
];

const TIPS_IN: Tip[] = [
  { prefix: "Type ", command: "/help", suffix: " to see every command" },
  { prefix: "Run ", command: "/assessments list", suffix: " to view assessments" },
  { prefix: "Try ", command: "/org billing", suffix: " to check entitlements" },
  { prefix: "Press ", command: "Tab", suffix: " to autocomplete commands" },
  { prefix: "Hit ", command: "/", suffix: " to open the command palette" },
];

const pool = isSignedIn() ? TIPS_IN : TIPS_OUT;
const tip = pool[Math.floor(Math.random() * pool.length)]!;

export const TipLine = () => (
  <box flexDirection="row" justifyContent="center" alignItems="center" flexShrink={0}>
    <text fg={colors.accentAmber}>{glyphs.dot}</text>
    <text fg={colors.accentAmber}> Tip</text>
    <text fg={colors.textPrimary}>{` ${tip.prefix}`}</text>
    <box paddingLeft={1} paddingRight={1} backgroundColor={colors.brandBlack} flexShrink={0}>
      <text fg={colors.accentLime}>{tip.command}</text>
    </box>
    <text fg={colors.textPrimary}>{tip.suffix}</text>
  </box>
);
