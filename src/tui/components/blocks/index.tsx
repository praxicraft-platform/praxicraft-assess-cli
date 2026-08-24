/**
 * Block dispatcher — switches on block.type and renders the matching component.
 */

import { Match, Switch } from "solid-js";
import type { BlockType } from "../../types";
import { Spinner } from "./Spinner";
import { Success, Error, Info, TextBlock, Markdown, Help } from "./Status";
import { InlineInput } from "./InlineInput";
import { InlineSelect } from "./InlineSelect";
import { Confirm } from "./Confirm";

type Variant<T extends BlockType["type"]> = Extract<BlockType, { type: T }>;

export const renderBlock = (block: BlockType) => (
  <Switch fallback={null}>
    <Match when={block.type === "spinner" && block}>
      {(b: () => Variant<"spinner">) => <Spinner label={b().label} />}
    </Match>
    <Match when={block.type === "success" && block}>
      {(b: () => Variant<"success">) => <Success message={b().message} />}
    </Match>
    <Match when={block.type === "error" && block}>
      {(b: () => Variant<"error">) => <Error message={b().message} />}
    </Match>
    <Match when={block.type === "info" && block}>
      {(b: () => Variant<"info">) => <Info message={b().message} />}
    </Match>
    <Match when={block.type === "text" && block}>
      {(b: () => Variant<"text">) => <TextBlock text={b().text} />}
    </Match>
    <Match when={block.type === "markdown" && block}>
      {(b: () => Variant<"markdown">) => <Markdown text={b().text} />}
    </Match>
    <Match when={block.type === "help" && block}>{() => <Help />}</Match>
    <Match when={block.type === "inline-input" && block}>
      {(b: () => Variant<"inline-input">) => (
        <InlineInput label={b().label} secure={b().secure} onSubmit={b().onSubmit} />
      )}
    </Match>
    <Match when={block.type === "inline-select" && block}>
      {(b: () => Variant<"inline-select">) => (
        <InlineSelect label={b().label} options={b().options} onSubmit={b().onSubmit} />
      )}
    </Match>
    <Match when={block.type === "confirm" && block}>
      {(b: () => Variant<"confirm">) => (
        <Confirm
          message={b().message}
          onConfirm={b().onConfirm}
          onCancel={b().onCancel}
        />
      )}
    </Match>
  </Switch>
);
