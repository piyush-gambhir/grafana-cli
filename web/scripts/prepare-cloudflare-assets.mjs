import { cp, mkdir, readFile, rm, writeFile } from "node:fs/promises";
import path from "node:path";

const sitePath = process.argv[2];

if (!sitePath || !/^[a-z0-9-]+$/.test(sitePath)) {
  throw new Error("Expected a URL-safe site path argument.");
}

const source = path.resolve("out");
const destinationRoot = path.resolve(".cloudflare/assets");
const destination = path.join(destinationRoot, sitePath);

await rm(destinationRoot, { recursive: true, force: true });
await mkdir(destination, { recursive: true });
await cp(source, destination, { recursive: true });

// Wrangler reads _redirects only from the assets root, while the site lives
// under /<sitePath>: move the file up and prefix every rule's paths.
const siteRedirects = path.join(destination, "_redirects");
let rules;
try {
  rules = await readFile(siteRedirects, "utf8");
} catch (error) {
  if (error.code !== "ENOENT") throw error;
}
if (rules !== undefined) {
  const prefix = (target) => (target.startsWith("/") ? `/${sitePath}${target}` : target);
  const lines = rules.split("\n").map((line) => {
    const rule = line.trim();
    if (!rule || rule.startsWith("#")) return line;
    const [from, to, ...rest] = rule.split(/\s+/);
    if (!from.startsWith("/") || !to) throw new Error(`Unexpected _redirects rule: ${rule}`);
    return [prefix(from), prefix(to), ...rest].join(" ");
  });
  await writeFile(path.join(destinationRoot, "_redirects"), lines.join("\n"));
  await rm(siteRedirects);
}
