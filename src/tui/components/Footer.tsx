import { version } from "../../../package.json";
import { colors, glyphs } from "../theme";

export const Footer = () => (
  <box flexDirection="row" flexShrink={0} paddingLeft={2} paddingRight={2} paddingTop={0}>
    <box flexGrow={1} flexDirection="row">
      <text fg={colors.textDim}>docs.praxicraft.com</text>
      <text fg={colors.textDim}>{` ${glyphs.separator} `}</text>
      <text fg={colors.textDim}>Ask AI needs Starter+</text>
    </box>
    <text fg={colors.textDim}>{`v${version}`}</text>
  </box>
);
