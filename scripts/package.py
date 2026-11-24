#!/usr/bin/env python3
"""Package one built binary together with config/ into an archive.

Usage: python3 scripts/package.py <linux|macos|windows> <version>
"""
import os
import sys
import tarfile
import zipfile

BINARIES = {
    "linux": ("output/webshell-detect-linux", "webshell-detect", "tar.gz"),
    "macos": ("output/webshell-detect-macos", "webshell-detect", "tar.gz"),
    "windows": ("output/webshell-detect-windows.exe", "webshell-detect.exe", "zip"),
}
FILES = ["config/conf.yaml", "config/mock.php", "README.md", "LICENSE"]


def main() -> int:
    if len(sys.argv) != 3 or sys.argv[1] not in BINARIES:
        print(__doc__, file=sys.stderr)
        return 2
    os_name, version = sys.argv[1], sys.argv[2]
    src, arcname, ext = BINARIES[os_name]
    if not os.path.exists(src):
        print(f"missing built binary: {src}", file=sys.stderr)
        return 1

    out = f"output/webshell-detect-{os_name}-default-x86_64-{version}.{ext}"
    if ext == "zip":
        with zipfile.ZipFile(out, "w", zipfile.ZIP_DEFLATED) as zf:
            zf.write(src, arcname)
            for path in FILES:
                zf.write(path, os.path.basename(path))
    else:
        with tarfile.open(out, "w:gz") as tf:
            tf.add(src, arcname=arcname)
            for path in FILES:
                tf.add(path, arcname=os.path.basename(path))
    print(f"wrote {out}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
