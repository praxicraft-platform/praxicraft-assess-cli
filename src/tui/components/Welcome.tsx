import { Show } from "solid-js";
import { useTerminalDimensions } from "@opentui/solid";
import { colors } from "../theme";
import { TipLine } from "./TipLine";

const WIDE_MIN = 72;

export const Welcome = () => {
  const dims = useTerminalDimensions();
  const wide = () => dims().width >= WIDE_MIN;

  return (
    <box
      flexGrow={1}
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
      paddingTop={2}
      paddingBottom={2}
    >
      <ascii_font
        text="PRAXI"
        font="block"
        color={[colors.brand, colors.brandLime]}
        selectable={false}
      />
      <Show when={wide()}>
        <ascii_font
          text="CRAFT"
          font="block"
          color={[colors.brandLime, colors.brand]}
          selectable={false}
        />
      </Show>
      <box flexShrink={0} paddingTop={2}>
        <TipLine />
      </box>
    </box>
  );
};
