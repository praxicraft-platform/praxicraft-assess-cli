/**
 * Slash-command palette overlay. Centered modal with sliding-window results,
 * fuzzy ranking, and keyboard handlers (arrow nav, Tab/Enter complete, Esc).
 */

import { For, Show, createMemo } from "solid-js";
import { useTerminalDimensions } from "@opentui/solid";
import { COMMANDS } from "../../lib/commands";
import { colors, glyphs } from "../theme";

const VISIBLE = 9;
const PALETTE_WIDTH = 70;

interface PaletteProps {
  query: string;
  selectedIndex: number;
  visible: boolean;
}

interface RankedCommand {
  command: string;
  description: string;
  score: number;
}

const score = (query: string, command: string): number => {
  if (!query) return 0;
  const q = query.toLowerCase();
  const c = command.toLowerCase();
  if (c.startsWith(q)) return 100 + (q.length / c.length) * 20;
  // Ignore lone leading "/" so "/zzz" does not match every command
  const needle = q.startsWith("/") ? q.slice(1) : q;
  if (!needle) return 1;
  let qi = 0;
  let hits = 0;
  for (let ci = 0; ci < c.length && qi < needle.length; ci++) {
    if (c[ci] === needle[qi]) {
      qi++;
      hits++;
    }
  }
  if (qi === needle.length) return 50 + (hits / c.length) * 20;
  return 0;
};

export const rankCommands = (query: string): RankedCommand[] => {
  if (!query || query === "/") {
    return COMMANDS.map((c) => ({ ...c, score: 1 }));
  }
  return COMMANDS.map((c) => ({ ...c, score: score(query, c.command) }))
    .filter((r) => r.score > 0)
    .sort((a, b) => b.score - a.score || a.command.length - b.command.length);
};

const computeWindow = (selected: number, total: number): { start: number; end: number } => {
  if (total <= VISIBLE) return { start: 0, end: total };
  const half = Math.floor(VISIBLE / 2);
  const start = Math.max(0, Math.min(selected - half, total - VISIBLE));
  return { start, end: start + VISIBLE };
};

export const Palette = (props: PaletteProps) => {
  const dims = useTerminalDimensions();
  const ranked = createMemo(() => rankCommands(props.query));
  const window = createMemo(() => computeWindow(props.selectedIndex, ranked().length));
  const visibleRows = createMemo(() => ranked().slice(window().start, window().end));
  const above = createMemo(() => window().start);
  const below = createMemo(() => Math.max(0, ranked().length - window().end));
  const left = createMemo(() => Math.max(0, Math.floor((dims().width - PALETTE_WIDTH) / 2)));
  const top = 4;

  return (
    <Show when={props.visible}>
      <box
        position="absolute"
        top={top}
        left={left()}
        width={PALETTE_WIDTH}
        borderStyle="single"
        borderColor={colors.accentSky}
        backgroundColor={colors.surfaceDeep}
        paddingLeft={1}
        paddingRight={1}
        flexDirection="column"
      >
        <box flexDirection="row">
          <text fg={colors.accentSky}>{`${glyphs.bullet} command palette `}</text>
          <text fg={colors.textDim}>
            {ranked().length === 0
              ? "no matches"
              : `${window().start + 1}\u2013${window().end} / ${ranked().length}`}
          </text>
        </box>
        <text fg={colors.textDim}>{"\u2500".repeat(PALETTE_WIDTH - 4)}</text>
        <Show when={above() > 0}>
          <text fg={colors.textDim}>{`  \u2191 ${above()} above`}</text>
        </Show>
        <For each={visibleRows()}>
          {(row, i) => {
            const isSelected = () => window().start + i() === props.selectedIndex;
            return (
              <box flexDirection="row">
                <text
                  fg={isSelected() ? colors.accentLime : colors.textPrimary}
                  attributes={isSelected() ? 1 : 0}
                >
                  {`${isSelected() ? `${glyphs.prompt} ` : "  "}${row.command}`}
                </text>
                <text fg={colors.textDim}>{`  ${glyphs.separator} ${row.description}`}</text>
              </box>
            );
          }}
        </For>
        <Show when={below() > 0}>
          <text fg={colors.textDim}>{`  \u2193 ${below()} below`}</text>
        </Show>
      </box>
    </Show>
  );
};
