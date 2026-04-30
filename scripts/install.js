"use strict";

const https = require("https");
const fs = require("fs");
const path = require("path");
const os = require("os");
const { execSync } = require("child_process");

const BINARY = process.platform === "win32" ? "carbonstop.exe" : "carbonstop";
const VERSION = "0.3.9"; // Go binary release version on GitHub

function mapPlatform(platform) {
  switch (platform) {
    case "darwin": return "darwin";
    case "linux":  return "linux";
    case "win32":  return "windows";
    default:       throw new Error("Unsupported platform: " + platform);
  }
}

function mapArch(arch) {
  switch (arch) {
    case "x64":   return "x86_64";
    case "arm64": return "arm64";
    default:      throw new Error("Unsupported arch: " + arch);
  }
}

function getAssetName() {
  const p = mapPlatform(process.platform);
  const a = mapArch(os.arch());
  const ext = process.platform === "win32" ? "zip" : "tar.gz";
  return `carbonstop-cli_${VERSION}_${p}_${a}.${ext}`;
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    const get = (u) => {
      https.get(u, (res) => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          return get(res.headers.location);
        }
        if (res.statusCode !== 200) {
          file.close();
          fs.unlinkSync(dest);
          return reject(new Error(`HTTP ${res.statusCode} downloading ${u}`));
        }
        res.pipe(file);
        file.on("finish", () => { file.close(); resolve(); });
      }).on("error", reject);
    };
    get(url);
  });
}

function extract(archivePath, destDir) {
  if (process.platform === "win32") {
    execSync(`powershell -Command "Expand-Archive -Path '${archivePath}' -DestinationPath '${destDir}' -Force"`, { stdio: "pipe" });
  } else {
    execSync(`tar xzf "${archivePath}" -C "${destDir}"`, { stdio: "pipe" });
  }
}

async function install() {
  const binDir = path.join(__dirname, "..", "bin");
  const binaryPath = path.join(binDir, BINARY);

  if (fs.existsSync(binaryPath)) {
    console.log(`[carbonstop-cli] binary already installed`);
    return;
  }

  const assetName = getAssetName();
  const url = `https://github.com/carbonstop-hub/carbonstop-cli/releases/download/v${VERSION}/${assetName}`;

  console.log(`[carbonstop-cli] downloading ${assetName} ...`);

  const tmpDir = path.join(os.tmpdir(), "carbonstop-cli-download");
  fs.mkdirSync(tmpDir, { recursive: true });
  fs.mkdirSync(binDir, { recursive: true });

  const archivePath = path.join(tmpDir, assetName);
  await download(url, archivePath);

  extract(archivePath, binDir);

  fs.unlinkSync(archivePath);

  if (process.platform !== "win32") {
    fs.chmodSync(binaryPath, 0o755);
    if (process.platform === "darwin") {
      try {
        execSync(`xattr -d com.apple.quarantine "${binaryPath}" 2>/dev/null || true`, { stdio: "pipe" });
      } catch (_) {}
    }
  }

  console.log(`[carbonstop-cli] installed successfully`);
}

install().catch((err) => {
  console.error("[carbonstop-cli] install failed:", err.message);
  process.exit(1);
});
