const assert = require("node:assert/strict");
const fs = require("node:fs");
const { createRequire } = require("node:module");
const path = require("node:path");
const test = require("node:test");
const vm = require("node:vm");

const platformPairs = [
  ["darwin", "x64"],
  ["darwin", "arm64"],
  ["linux", "x64"],
  ["linux", "arm64"],
  ["win32", "x64"],
  ["win32", "arm64"],
];
const launcherPath = path.join(__dirname, "index.js");
const cliPath = path.join(__dirname, "bin", "cli");
const launcherSource = fs.readFileSync(launcherPath, "utf8");
const cliSource = fs.readFileSync(cliPath, "utf8").replace(/^#![^\n]*\n/, "");
const args = ["-transport", "stdio", "-namespace", "team with spaces", ""];

function binaryRequest(platform, arch) {
  const binaryName = platform === "win32" ? "mcp-nomad.exe" : "mcp-nomad";
  return `@kocierik/mcp-nomad-${platform}-${arch}/bin/${binaryName}`;
}

function loadLauncher(platform, arch, { missing = false, normalResolver = false, spawnError } = {}) {
  const resolutionCalls = [];
  const executionCalls = [];
  const apiCalls = [];
  const fixturePath = path.join(__dirname, "fixture binaries", platform, arch,
    platform === "win32" ? "mcp-nomad.exe" : "mcp-nomad");
  const resolutionError = Object.assign(new Error("Cannot find module"), {
    code: "MODULE_NOT_FOUND",
  });
  const childProcess = {
    execFileSync(binaryPath, binaryArgs, options) {
      executionCalls.push({ binaryPath, args: Array.from(binaryArgs), stdio: options.stdio });
      if (spawnError) throw spawnError;
    },
  };
  const nativeRequire = createRequire(launcherPath);
  const module = { exports: {} };
  const injectedRequire = (request) => {
    if (request === "child_process") return childProcess;
    return nativeRequire(request);
  };
  injectedRequire.resolve = (request) => {
    resolutionCalls.push(request);
    if (normalResolver) return nativeRequire.resolve(request);
    if (missing || request !== binaryRequest(platform, arch)) throw resolutionError;
    return fixturePath;
  };
  const processStub = { platform, arch, argv: ["node", cliPath, ...args] };
  vm.runInNewContext(launcherSource, {
    require: injectedRequire, module, __dirname, process: processStub,
  }, { filename: launcherPath });
  const cliRequire = (request) => {
    if (request === "../index.js") {
      return {
        runBinary(...binaryArgs) {
          apiCalls.push(binaryArgs);
          return module.exports.runBinary(...binaryArgs);
        },
      };
    }
    return injectedRequire(request);
  };
  cliRequire.resolve = injectedRequire.resolve;
  return {
    exports: module.exports, resolutionCalls, executionCalls, apiCalls,
    fixturePath, resolutionError,
    runCli() {
      vm.runInNewContext(cliSource, {
        require: cliRequire, __dirname: path.dirname(cliPath), process: processStub,
      }, { filename: cliPath });
    },
  };
}

for (const [platform, arch] of platformPairs) {
  test(`API resolves the exact ${platform}/${arch} package binary`, () => {
    const launcher = loadLauncher(platform, arch);
    assert.equal(launcher.exports.runBinary(...args), undefined);
    assert.deepEqual(launcher.resolutionCalls, [binaryRequest(platform, arch)]);
    assert.deepEqual(launcher.executionCalls, [{
      binaryPath: launcher.fixturePath, args, stdio: "inherit",
    }]);
    assert.equal(launcher.exports.getBinaryPath(platform, arch), launcher.fixturePath);
  });

  test(`CLI shares API resolution and arguments on ${platform}/${arch}`, () => {
    const launcher = loadLauncher(platform, arch);
    launcher.runCli();
    assert.deepEqual(launcher.apiCalls, [args]);
    assert.deepEqual(launcher.resolutionCalls, [binaryRequest(platform, arch)]);
    assert.deepEqual(launcher.executionCalls, [{
      binaryPath: launcher.fixturePath, args, stdio: "inherit",
    }]);
  });

  test(`missing ${platform}/${arch} package fails without a fallback or spawn`, () => {
    const launcher = loadLauncher(platform, arch, { missing: true });
    for (const invoke of [() => launcher.exports.runBinary(...args), () => launcher.runCli()]) {
      assert.throws(invoke, (error) => {
        assert.match(error.message, /resolve|missing|install/i);
        assert.ok(error.message.includes(binaryRequest(platform, arch)));
        assert.equal(error.cause, launcher.resolutionError);
        return true;
      });
    }
    assert.deepEqual(launcher.executionCalls, []);
    assert.deepEqual(launcher.resolutionCalls, Array(2).fill(binaryRequest(platform, arch)));
  });
}

for (const [platform, arch] of [["freebsd", "x64"], ["linux", "ia32"]]) {
  test(`unsupported ${platform}/${arch} fails before resolution or spawning`, () => {
    const launcher = loadLauncher(platform, arch);
    for (const invoke of [() => launcher.exports.runBinary(...args), () => launcher.runCli()]) {
      assert.throws(invoke, (error) => {
        assert.match(error.message, /unsupported/i);
        assert.ok(error.message.includes(platform));
        assert.ok(error.message.includes(arch));
        return true;
      });
    }
    assert.deepEqual(launcher.resolutionCalls, []);
    assert.deepEqual(launcher.executionCalls, []);
  });
}

test("getBinaryPath uses explicit arguments rather than the host platform", () => {
  const launcher = loadLauncher("win32", "arm64");
  assert.throws(() => launcher.exports.getBinaryPath("freebsd", "arm64"), /unsupported/i);
  assert.deepEqual(launcher.resolutionCalls, []);
});

test("execution errors propagate unchanged through API and CLI", () => {
  const spawnError = new Error("binary exited unsuccessfully");
  const launcher = loadLauncher("linux", "x64", { spawnError });
  assert.throws(() => launcher.exports.runBinary(...args), (error) => error === spawnError);
  assert.throws(() => launcher.runCli(), (error) => error === spawnError);
});

test("the six supported platforms match the declared optional dependencies", () => {
  const manifest = JSON.parse(fs.readFileSync(path.join(__dirname, "package.json"), "utf8"));
  assert.deepEqual(Object.keys(manifest.optionalDependencies).sort(),
    platformPairs.map(([platform, arch]) => `@kocierik/mcp-nomad-${platform}-${arch}`).sort());
});

test("normal Node resolution returns an installed binary or a descriptive missing-package error", () => {
  const launcher = loadLauncher(process.platform, process.arch, { normalResolver: true });
  const request = binaryRequest(process.platform, process.arch);
  let resolvedPath;
  try {
    resolvedPath = createRequire(launcherPath).resolve(request);
  } catch (error) {
    assert.equal(error.code, "MODULE_NOT_FOUND");
  }
  if (resolvedPath) {
    assert.equal(launcher.exports.getBinaryPath(process.platform, process.arch), resolvedPath);
  } else {
    assert.throws(() => launcher.exports.getBinaryPath(process.platform, process.arch), (error) => {
      assert.ok(error.message.includes(request));
      assert.equal(error.cause.code, "MODULE_NOT_FOUND");
      return true;
    });
  }
  assert.deepEqual(launcher.resolutionCalls, [request]);
  assert.deepEqual(launcher.executionCalls, []);
});