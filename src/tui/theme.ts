/**
 * Design tokens for the Assess TUI. Sole source of truth for brand colors and glyphs.
 */

export const colors = {
  brand: "#0D41FF",
  brandLime: "#5B8AFF",
  brandForest: "#0A2FCC",
  brandBlack: "#1F2023",

  surfaceDeep: "#0D0D0D",
  border: "#212423",

  textPrimary: "#FFFFFF",
  textMuted: "#737470",
  textDim: "#535452",

  success: "#22C55E",
  error: "#EF4444",
  warning: "#F5A623",
  info: "#38BDF8",

  accentSky: "#38BDF8",
  accentAmber: "#F5A623",
  accentMagenta: "#E85BCF",
  accentCyan: "#7FC4D4",
  accentLime: "#C6FE1E",

  testMode: "#F5A623",
  liveMode: "#22C55E",
} as const;

export const glyphs = {
  prompt: "❯",
  bullet: "◆",
  dot: "●",
  check: "✓",
  cross: "✗",
  arrow: "→",
  separator: "·",
} as const;

export const spacing = { xs: 0, sm: 1, md: 2, lg: 3 } as const;
