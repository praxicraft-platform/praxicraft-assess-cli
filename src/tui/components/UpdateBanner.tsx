import { Show } from "solid-js";
import { version as currentVersion } from "../../../package.json";
import { colors, glyphs } from "../theme";

type Props = {
  latestVersion: string | null;
};

export const UpdateBanner = (props: Props) => (
  <Show when={props.latestVersion}>
    {(latest) => (
      <box
        flexShrink={0}
        flexDirection="row"
        marginLeft={2}
        marginRight={2}
        marginTop={1}
        marginBottom={1}
        paddingLeft={1}
        paddingRight={1}
        borderStyle="rounded"
        borderColor={colors.brandLime}
        backgroundColor={colors.brandBlack}
      >
        <text fg={colors.brandLime}>{`${glyphs.arrow} `}</text>
        <text fg={colors.textPrimary}>
          {`Update available: v${latest()} (you have v${currentVersion}). Run `}
        </text>
        <text fg={colors.brand}>npm i -g @praxicraft/assess-cli</text>
        <text fg={colors.textDim}>{` ${glyphs.separator} or `}</text>
        <text fg={colors.brand}>curl -fsSL https://praxicraft.com/install.sh | sh</text>
      </box>
    )}
  </Show>
);
