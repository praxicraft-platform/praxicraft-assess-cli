import { Show } from "solid-js";
import { useTerminalDimensions } from "@opentui/solid";
import { colors, glyphs } from "../theme";
import { TipLine } from "./TipLine";

/** slick ≈ 60 cols for PRAXICRAFT; tiny ≈ 37 for narrow panes */
const SLICK_MIN = 68;

export const Welcome = () => {
  const dims = useTerminalDimensions();
  const useSlick = () => dims().width >= SLICK_MIN;

  return (
    <box
      flexGrow={1}
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
      paddingTop={1}
      paddingBottom={2}
    >
      <box flexDirection="column" alignItems="center" flexShrink={0}>
        <Show
          when={useSlick()}
          fallback={
            <ascii_font
              text="PRAXICRAFT"
              font="tiny"
              color={[colors.brand, colors.brandLime]}
              selectable={false}
            />
          }
        >
          <ascii_font
            text="PRAXICRAFT"
            font="slick"
            color={[colors.brand, colors.brandLime]}
            selectable={false}
          />
        </Show>

        <box
          flexDirection="row"
          alignItems="center"
          justifyContent="center"
          paddingTop={1}
          flexShrink={0}
        >
          <text fg={colors.brand}>{glyphs.bullet}</text>
          <text fg={colors.textMuted}>  Assess CLI  </text>
          <text fg={colors.brand}>{glyphs.bullet}</text>
        </box>

        <box paddingTop={1} flexShrink={0}>
          <text fg={colors.brandForest}>{"━".repeat(useSlick() ? 28 : 18)}</text>
        </box>
      </box>

      <box flexShrink={0} paddingTop={2}>
        <TipLine />
      </box>
    </box>
  );
};
