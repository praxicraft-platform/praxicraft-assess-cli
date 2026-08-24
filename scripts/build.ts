#!/usr/bin/env bun
import { mkdirSync } from "node:fs";
import { join } from "node:path";
import solidPlugin from "@opentui/solid/bun-plugin";

mkdirSync("dist", { recursive: true });
const result = await Bun.build({
  entrypoints: [join(import.meta.dir, "../src/index.ts")],
  outdir: "dist",
  target: "bun",
  sourcemap: "none",
  minify: false,
  plugins: [solidPlugin],
});

if (!result.success) {
  console.error(result.logs);
  process.exit(1);
}

const banner = "#!/usr/bin/env bun\n";
const out = join(import.meta.dir, "../dist/index.js");
const body = await Bun.file(out).text();
if (!body.startsWith("#!")) {
  await Bun.write(out, banner + body);
}
console.log("Built dist/index.js");
