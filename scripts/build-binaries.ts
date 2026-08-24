#!/usr/bin/env bun
/**
 * Cross-compile standalone binaries for GitHub Releases (install.sh).
 * Asset names: praxicraft-assess-{os}-{arch}[+.exe] + SHA256SUMS.txt
 */
import { mkdirSync, chmodSync, createWriteStream, readdirSync } from "node:fs";
import { rename, rm, writeFile } from "node:fs/promises";
import { createHash } from "node:crypto";
import { readFileSync } from "node:fs";
import solidPlugin from "@opentui/solid/bun-plugin";

const targets = [
  { bun: "bun-darwin-x64", os: "darwin", arch: "x64" },
  { bun: "bun-darwin-arm64", os: "darwin", arch: "arm64" },
  { bun: "bun-linux-x64", os: "linux", arch: "x64" },
  { bun: "bun-linux-arm64", os: "linux", arch: "arm64" },
  { bun: "bun-windows-x64", os: "windows", arch: "x64" },
] as const;

mkdirSync("dist", { recursive: true });

const sums: string[] = [];

for (const t of targets) {
  const tempDir = `./dist/temp-${t.os}-${t.arch}`;
  const build = await Bun.build({
    entrypoints: ["./src/index.ts"],
    outdir: tempDir,
    target: "bun",
    plugins: [solidPlugin],
    compile: {
      target: t.bun as any,
    },
  });

  const artifact = build.outputs[0];
  if (!build.success || !artifact?.path) {
    console.error(`Failed ${t.bun}`, build.logs);
    process.exit(1);
  }

  const assetName =
    t.os === "windows"
      ? `praxicraft-assess-${t.os}-${t.arch}.exe`
      : `praxicraft-assess-${t.os}-${t.arch}`;
  const outPath = `./dist/${assetName}`;
  await rename(artifact.path, outPath);
  if (t.os !== "windows") chmodSync(outPath, 0o755);

  const hash = createHash("sha256").update(readFileSync(outPath)).digest("hex");
  sums.push(`${hash}  ${assetName}`);

  await rm(tempDir, { recursive: true, force: true });
  console.log(`Built ${assetName}`);
}

await writeFile("./dist/SHA256SUMS.txt", sums.join("\n") + "\n");
console.log("Wrote dist/SHA256SUMS.txt");
console.log("Done — binaries in ./dist");
