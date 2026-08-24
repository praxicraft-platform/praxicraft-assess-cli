export type BlockVariant =
  | { type: "text"; text: string }
  | { type: "markdown"; text: string; streaming?: boolean }
  | { type: "error"; message: string }
  | { type: "success"; message: string }
  | { type: "info"; message: string }
  | { type: "spinner"; label: string }
  | { type: "help" }
  | {
      type: "inline-input";
      label: string;
      secure?: boolean;
      onSubmit: (value: string) => void;
    }
  | {
      type: "inline-select";
      label: string;
      options: { label: string; value: string }[];
      onSubmit: (value: string) => void;
    }
  | {
      type: "confirm";
      message: string;
      onConfirm: () => void;
      onCancel: () => void;
    };

export type BlockType = BlockVariant & { id: string };

export type Message = {
  id: string;
  role: "user" | "system";
  text?: string;
  blocks: BlockType[];
};
