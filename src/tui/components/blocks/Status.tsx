import { For } from "solid-js";
import { COMMANDS } from "../../../lib/commands";
import { colors, glyphs } from "../../theme";

export const Success = (props: { message: string }) => (
  <box flexDirection="row">
    <text fg={colors.success}>{`${glyphs.check} `}</text>
    <text fg={colors.textPrimary}>{props.message}</text>
  </box>
);

export const Error = (props: { message: string }) => (
  <box
    borderStyle="rounded"
    borderColor={colors.error}
    paddingLeft={1}
    paddingRight={1}
    flexDirection="row"
  >
    <text fg={colors.error}>{`${glyphs.cross} `}</text>
    <text fg={colors.textPrimary}>{props.message}</text>
  </box>
);

export const Info = (props: { message: string }) => (
  <box flexDirection="row">
    <text fg={colors.info}>{`${glyphs.bullet} `}</text>
    <text fg={colors.textPrimary}>{props.message}</text>
  </box>
);

export const TextBlock = (props: { text: string }) => (
  <text fg={colors.textPrimary}>{props.text}</text>
);

export const Markdown = (props: { text: string }) => (
  <text fg={colors.textPrimary}>{props.text}</text>
);

export const Help = () => (
  <box flexDirection="column" gap={1} paddingTop={1} paddingBottom={1}>
    <box flexDirection="row" alignItems="center" gap={1}>
      <text fg={colors.brandLime}>{glyphs.arrow}</text>
      <text fg={colors.textPrimary}>Commands</text>
      <text fg={colors.textDim}>· free text asks Assess AI · /login to sign in</text>
    </box>
    <For each={COMMANDS}>
      {(c) => (
        <box flexDirection="row" gap={2}>
          <text fg={colors.brandLime}>{c.command.padEnd(22)}</text>
          <text fg={colors.textMuted}>{c.description}</text>
        </box>
      )}
    </For>
  </box>
);
