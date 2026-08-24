import { Show } from "solid-js";
import { version } from "../../../package.json";
import { colors, glyphs } from "../theme";
import { useTui, type AuthInfo } from "../context";

interface InputBarProps {
  onSubmit: (value: string) => void;
  onInput?: (value: string) => void;
  value?: string;
  placeholder?: string;
  disabled?: boolean;
}

export const InputBar = (props: InputBarProps) => {
  const { promptActive, authInfo } = useTui();
  const isFocused = () => !props.disabled && !promptActive();

  const handleInput = (next: string) => {
    props.onInput?.(next);
  };

  const handleSubmit = (final: string) => {
    const trimmed = final.trim();
    if (!trimmed) return;
    props.onSubmit(trimmed);
  };

  return (
    <box flexDirection="row" flexShrink={0} alignItems="stretch">
      <box width={1} flexShrink={0} backgroundColor={colors.brand}>
        <text fg={colors.brand}> </text>
      </box>
      <box
        flexDirection="column"
        flexGrow={1}
        paddingLeft={2}
        paddingRight={2}
        paddingTop={1}
        paddingBottom={1}
        backgroundColor={colors.surfaceDeep}
      >
        <box flexDirection="row" flexShrink={0}>
          <input
            flexGrow={1}
            focused={isFocused()}
            value={props.value ?? ""}
            placeholder={
              props.placeholder ??
              `Ask Assess… ${glyphs.separator} type / for commands`
            }
            onInput={handleInput as any}
            onSubmit={handleSubmit as any}
          />
        </box>
        <box height={1} flexShrink={0} />
        <box flexDirection="row" flexShrink={0}>
          <text fg={colors.brand}>praxicraft</text>
          <text fg={colors.textDim}>/</text>
          <text fg={colors.brandLime}>assess</text>
          <text fg={colors.textDim}>{` ${glyphs.separator} `}</text>
          <Show
            when={authInfo()}
            fallback={
              <>
                <text fg={colors.textMuted}>signed out</text>
                <text fg={colors.textDim}>{` ${glyphs.separator} `}</text>
                <text fg={colors.textMuted}>{`v${version}`}</text>
              </>
            }
          >
            {(info: () => NonNullable<AuthInfo>) => (
              <>
                <text fg={info().mode === "test_mode" ? colors.testMode : colors.liveMode}>
                  {info().mode === "test_mode" ? "TEST" : "LIVE"}
                </text>
                <text fg={colors.textDim}>{` ${glyphs.separator} `}</text>
                <text fg={colors.textMuted}>{info().key}</text>
                <text fg={colors.textDim}>{` ${glyphs.separator} `}</text>
                <text fg={colors.textMuted}>{`v${version}`}</text>
              </>
            )}
          </Show>
        </box>
      </box>
    </box>
  );
};
