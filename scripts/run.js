#!/usr/bin/env node
"use strict";

const { spawnSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const BINARY = process.platform === "win32" ? "carbonstop.exe" : "carbonstop";
const binaryPath = path.join(__dirname, "..", "bin", BINARY);

if (!fs.existsSync(binaryPath)) {
  console.error("[carbonstop-cli] binary not found, running install...");
  try {
    require("./install.js");
  } catch (e) {
    console.error("[carbonstop-cli] install failed:", e.message);
  }
  if (!fs.existsSync(binaryPath)) {
    console.error("[carbonstop-cli] could not install binary. Re-run: npm install -g @carbonstopper/cli");
    process.exit(1);
  }
}

const args = process.argv.slice(2);
const result = spawnSync(binaryPath, args, { stdio: "inherit" });

if (result.error) {
  console.error("[carbonstop-cli]", result.error.message);
  process.exit(1);
}

process.exit(result.status != null ? result.status : 1);
