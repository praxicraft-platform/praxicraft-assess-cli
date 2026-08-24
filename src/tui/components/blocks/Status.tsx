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
  <text fg={colors.textMuted}>
    Type / for commands. Free text asks the Assess assistant. /login to authenticate.
  </text>
);
