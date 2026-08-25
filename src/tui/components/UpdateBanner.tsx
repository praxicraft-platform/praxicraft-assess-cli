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
        marginBottom={1}
        paddingLeft={1}
        paddingRight={1}
        borderStyle="rounded"
        borderColor={colors.warning}
        backgroundColor={colors.brandBlack}
      >
        <text fg={colors.warning}>{`${glyphs.arrow} `}</text>
        <text fg={colors.textPrimary}>
          {`Update available: v${latest()} (you have v${currentVersion}). Run `}
        </text>
        <text fg={colors.accentSky}>npm i -g @praxicraft/assess-cli</text>
      </box>
    )}
  </Show>
);
