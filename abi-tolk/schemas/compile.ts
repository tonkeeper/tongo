#!/usr/bin/env -S npx tsx

/*
  Usage:
    $ npm install
    then

    compile all .tolk files in this directory tree:
    $ npx tsx compile.ts --all

    compile single contract:
    $ npx tsx compile.ts elector/elector.tolk
    or:
    $ npx tsx compile.ts elector/elector.tolk elector/elector.json
 */

import fs from 'fs';
import path from 'path';
import { runTolkCompiler } from '@ton/tolk-js';

const SCHEMAS_DIR = path.dirname(new URL(import.meta.url).pathname);

export async function convertTolkFileToABI(absFileName: string): Promise<Record<string, unknown> & { code_boc64: string }> {
  let compileResult = await runTolkCompiler({
    entrypointFileName: absFileName,
    fsReadCallback: p => fs.readFileSync(p, 'utf-8'),
  });
  if (compileResult.status === 'error') {
    throw new Error(`Can not compile with tolk-js: ${compileResult.message}`);
  }

  let json = (compileResult as any).abiJson;
  json.code_boc64 = compileResult.codeBoc64;
  return json as any;
}

function collectTolkFiles(dir: string): string[] {
  let result: string[] = [];
  for (let entry of fs.readdirSync(dir, { withFileTypes: true })) {
    if (entry.name.startsWith('.')) continue;
    let fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      for (let child of fs.readdirSync(fullPath, { withFileTypes: true })) {
        if (child.isFile() && child.name.endsWith('.tolk')) {
          result.push(path.join(fullPath, child.name));
        }
      }
    } else if (entry.isFile() && entry.name.endsWith('.tolk')) {
      result.push(fullPath);
    }
  }
  return result;
}

async function compileFile(tolkPath: string, abiPath?: string) {
  let outPath = abiPath || tolkPath.replace(/\.tolk$/, '.json');
  let result = await convertTolkFileToABI(path.resolve(tolkPath));
  let outDir = path.dirname(outPath);
  if (outDir && !fs.existsSync(outDir)) {
    fs.mkdirSync(outDir, { recursive: true });
  }
  fs.writeFileSync(outPath, JSON.stringify(result, null, 2));
  return outPath;
}

async function compileAll() {
  let tolkFiles = collectTolkFiles(SCHEMAS_DIR);
  if (tolkFiles.length === 0) {
    console.log('No .tolk files found.');
    return;
  }
  console.log(`Found ${tolkFiles.length} .tolk file(s).`);
  for (let f of tolkFiles) {
    let outPath = await compileFile(f);
    console.log(`  ${path.relative(SCHEMAS_DIR, f)} -> ${path.relative(SCHEMAS_DIR, outPath)}`);
  }
}

async function main() {
  let args = process.argv.slice(2);
  if (args.length === 0 || args.includes('--help') || args.includes('-h')) {
    console.error('Usage: compile.ts [--all] | <input.tolk> [output.json]');
    process.exit(args.length === 0 ? 1 : 0);
  }

  if (args.includes('--all')) {
    await compileAll();
    return;
  }

  let tolkPath = args[0];
  let abiPath = args[1];
  let outPath = await compileFile(tolkPath, abiPath);
  console.log(`ABI written to ${outPath}`);
}

main().catch(err => {
  console.error(err);
  process.exit(1);
});
