/**
 * Solid context exposing shared TUI signals to deeply-nested components.
 */

import { createContext, useContext, type Accessor } from "solid-js";

export type AuthInfo = { mode: "test_mode" | "live_mode"; key: string } | null;

export interface TuiContextValue {
  authInfo: Accessor<AuthInfo>;
  input: Accessor<string>;
  setInput: (value: string) => void;
  paletteVisible: Accessor<boolean>;
  paletteIndex: Accessor<number>;
  isProcessing: Accessor<boolean>;
  promptActive: Accessor<boolean>;
  setPromptActive: (active: boolean) => void;
}

const TuiContext = createContext<TuiContextValue>();

export const TuiContextProvider = TuiContext.Provider;

export const useTui = (): TuiContextValue => {
  const ctx = useContext(TuiContext);
  if (!ctx) throw new Error("useTui must be called inside <TuiContextProvider>");
  return ctx;
};
