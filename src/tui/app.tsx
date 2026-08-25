import { For, Show, createMemo, createSignal, onMount, untrack } from "solid-js";
import { render, useKeyboard, useRenderer, useSelectionHandler } from "@opentui/solid";
import { ToastViewport } from "./components/Toast";
import { showToast } from "./toast";
import { copyToClipboard } from "./clipboard";
import { Welcome } from "./components/Welcome";
import { InputBar } from "./components/InputBar";
import { HintRow } from "./components/HintRow";
import { Footer } from "./components/Footer";
import { UpdateBanner } from "./components/UpdateBanner";
import { MessageRow } from "./components/MessageRow";
import { Palette, rankCommands } from "./components/Palette";
import { TuiContextProvider, type AuthInfo } from "./context";
import { createMessageStore } from "./CommandContext";
import { handleCommand } from "./router";
import { resolveCredentials } from "../utils/config";
import { checkForUpdate, formatUpdateNotice } from "../utils/version-check";

const App = () => {
  const [authInfo, setAuthInfo] = createSignal<AuthInfo>(null);
  const [input, setInput] = createSignal("");
  const [paletteIndex, setPaletteIndex] = createSignal(0);
  const [paletteDismissed, setPaletteDismissed] = createSignal(false);
  const [history, setHistory] = createSignal<string[]>([]);
  const [historyIndex, setHistoryIndex] = createSignal(-1);
  const [isProcessing, setIsProcessing] = createSignal(false);
  const [updateAvailable, setUpdateAvailable] = createSignal<string | null>(null);
  const [promptActive, setPromptActive] = createSignal(false);
  const [lastEscape, setLastEscape] = createSignal(0);

  const store = createMessageStore();
  const hasMessages = createMemo(() => store.messages().length > 0);
  const paletteVisible = createMemo(() => {
    const v = input();
    if (!v.startsWith("/")) return false;
    if (paletteDismissed() || promptActive()) return false;
    const firstWord = v.split(" ")[0] ?? "";
    if (v.includes(" ") && firstWord.length > 1) {
      const exactMatch = rankCommands(firstWord).some((c) => c.command === firstWord);
      if (exactMatch) return false;
    }
    return true;
  });
  const palettePool = createMemo(() => rankCommands(input()));

  const refreshAuthInfo = () => {
    try {
      const { apiKey } = resolveCredentials({ requireKey: true });
      const mode = apiKey.startsWith("ct_live_") ? "live_mode" : "test_mode";
      const masked = apiKey.slice(0, 10) + "..." + apiKey.slice(-3);
      setAuthInfo({ mode, key: masked });
    } catch {
      setAuthInfo(null);
    }
  };

  onMount(() => {
    refreshAuthInfo();
    void checkForUpdate().then((result) => {
      if (result) {
        setUpdateAvailable(result.latest);
        showToast(`Update available: v${result.latest}`, { variant: "info", durationMs: 10_000 });
      }
    });
  });

  const onInputChange = (next: string) => {
    setInput(next);
    setPaletteIndex(0);
    setHistoryIndex(-1);
    if (paletteDismissed() && next.length > 0) setPaletteDismissed(false);
  };

  const completeWith = (cmd: string) => {
    const next = `${cmd} `;
    setInput(next);
    setPaletteDismissed(true);
  };

  const renderer = useRenderer();

  const exitTui = () => {
    // Defer to let OpenTUI's stdin parser drain in-flight terminal responses
    // before destroy() flips stdin to canonical mode and the kernel echoes them.
    setTimeout(() => {
      renderer.destroy();
      process.exit(0);
    }, 80);
  };

  const onSubmit = async (raw: string) => {
    const trimmed = raw.trim();
    if (!trimmed) return;
    if (paletteVisible()) {
      const pool = untrack(palettePool);
      const pick = pool[paletteIndex()];
      const looksIncomplete =
        !!pick && pick.command !== trimmed && pick.command.startsWith(trimmed);
      if (looksIncomplete) {
        completeWith(pick.command);
        return;
      }
    }
    store.pushUserEcho(trimmed);
    setHistory((prev) => [trimmed, ...prev.filter((h) => h !== trimmed)]);
    setHistoryIndex(-1);
    setInput("");
    setPaletteDismissed(false);
    setIsProcessing(true);
    try {
      await handleCommand(trimmed, store.ctx, exitTui, refreshAuthInfo);
    } catch (e: any) {
      store.ctx.addBlock({
        type: "error",
        message: e?.message ?? String(e),
      });
    } finally {
      setIsProcessing(false);
    }
  };

  useSelectionHandler((selection) => {
    const text = selection.getSelectedText();
    if (!text || text.trim().length === 0) return;
    void copyToClipboard(text);
    showToast("Copied to clipboard", { variant: "success" });
    renderer.clearSelection();
  });

  useKeyboard((key) => {
    if (untrack(promptActive)) return;

    if (key.name === "escape") {
      const now = Date.now();
      const last = untrack(lastEscape);
      if (now - last < 500) {
        exitTui();
      }
      setLastEscape(now);

      if (paletteVisible()) {
        key.preventDefault();
        setPaletteDismissed(true);
      }
      return;
    }

    if (paletteVisible()) {
      const pool = untrack(palettePool);
      if (key.name === "up") {
        key.preventDefault();
        setPaletteIndex((i) => Math.max(0, i - 1));
      } else if (key.name === "down") {
        key.preventDefault();
        setPaletteIndex((i) => Math.min(pool.length - 1, i + 1));
      } else if (key.name === "tab") {
        key.preventDefault();
        const pick = pool[paletteIndex()];
        if (pick) completeWith(pick.command);
      }
      return;
    }

    if (key.name === "up") {
      const hist = untrack(history);
      if (hist.length === 0) return;
      key.preventDefault();
      setHistoryIndex((i) => {
        const next = Math.min(hist.length - 1, i + 1);
        setInput(hist[next] ?? "");
        return next;
      });
    } else if (key.name === "down") {
      const hist = untrack(history);
      if (hist.length === 0) return;
      key.preventDefault();
      setHistoryIndex((i) => {
        const next = i - 1;
        if (next >= 0) {
          setInput(hist[next] ?? "");
          return next;
        }
        setInput("");
        return -1;
      });
    }
  });

  return (
    <TuiContextProvider
      value={{
        authInfo,
        input,
        setInput,
        paletteVisible,
        paletteIndex,
        isProcessing,
        promptActive,
        setPromptActive,
      }}
    >
      <box flexDirection="column" width="100%" height="100%">
        <Show when={hasMessages()} fallback={<Welcome />}>
          <scrollbox
            flexGrow={1}
            stickyScroll={true}
            stickyStart="bottom"
            paddingLeft={2}
            paddingRight={2}
          >
            <For each={store.messages()}>{(m) => <MessageRow message={m} />}</For>
          </scrollbox>
        </Show>
        <box flexShrink={0} flexDirection="column" paddingLeft={2} paddingRight={2}>
          <UpdateBanner latestVersion={updateAvailable()} />
          <InputBar onSubmit={onSubmit} onInput={onInputChange} value={input()} />
          <HintRow />
        </box>
        <Footer />
        <Palette
          query={input()}
          selectedIndex={paletteIndex()}
          visible={paletteVisible()}
        />
        <ToastViewport />
      </box>
    </TuiContextProvider>
  );
};

export async function startTui() {
  // Restore terminal state on any exit, including OpenTUI's own Ctrl+C handler
  process.on("exit", () => {
    process.stdout.write("\x1b[?1000l\x1b[?1002l\x1b[?1003l\x1b[?1006l\x1b[?1015l\x1b[?25h");
  });

  void checkForUpdate().then((result) => {
    if (result) {
      process.stderr.write(`\n${formatUpdateNotice(result)}\n\n`);
    }
  });

  render(() => <App />, { exitOnCtrlC: true, targetFps: 30, useMouse: true });
}
