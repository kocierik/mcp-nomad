const childProcess = require("child_process");

// Lookup table for all platforms and binary distribution packages
const BINARY_DISTRIBUTION_PACKAGES = {
  darwin_x64: "mcp-nomad-darwin-x64",
  darwin_arm64: "mcp-nomad-darwin-arm64",
  linux_x64: "mcp-nomad-linux-x64",
  linux_arm64: "mcp-nomad-linux-arm64",
  win32_x64: "mcp-nomad-win32-x64",
  win32_arm64: "mcp-nomad-win32-arm64",
};

function getBinaryPath(platform = process.platform, arch = process.arch) {
  const platformSpecificPackageName = BINARY_DISTRIBUTION_PACKAGES[`${platform}_${arch}`];
  if (!platformSpecificPackageName) {
    throw new Error(`Unsupported platform/architecture: ${platform}/${arch}`);
  }
  const binaryName = platform === "win32" ? "mcp-nomad.exe" : "mcp-nomad";
  const binaryRequest = `@kocierik/${platformSpecificPackageName}/bin/${binaryName}`;
  try {
    return require.resolve(binaryRequest);
  } catch (cause) {
    throw new Error(
      `Unable to resolve ${binaryRequest}. Reinstall @kocierik/mcp-nomad with optional dependencies enabled.`,
      { cause },
    );
  }
}

module.exports.getBinaryPath = getBinaryPath;
module.exports.runBinary = function (...args) {
  childProcess.execFileSync(getBinaryPath(), args, {
    stdio: "inherit",
  });
};