/**
 * Praxicraft Assess TUI design tokens.
 * Palette: brand blues + neutral text. Semantic red/green reserved for errors/success only.
 */

export const colors = {
  /** Primary brand blue */
  brand: "#0D41FF",
  /** Lighter brand blue (highlights, links, accents) */
  brandLime: "#5B8AFF",
  /** Darker brand blue (rules, TEST badge) */
  brandForest: "#0A2FCC",
  brandBlack: "#1F2023",

  surfaceDeep: "#0D0D0D",
  border: "#212423",

  textPrimary: "#FFFFFF",
  textMuted: "#737470",
  textDim: "#535452",

  /** Errors / destructive only */
  error: "#EF4444",
  /** Success toasts / checkmarks only */
  success: "#22C55E",

  /** LIVE mode — brand highlight */
  liveMode: "#5B8AFF",
  /** TEST mode — muted brand */
  testMode: "#0A2FCC",
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
