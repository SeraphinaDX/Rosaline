# Third-Party Notices

Rosaline uses the following separately licensed Go modules. They remain under
their own licenses and are not relicensed under Rosaline's LGPL license.

- `github.com/fogleman/gg` — MIT License
- `github.com/golang/freetype` — FreeType-style/BSD-style license
- `github.com/gen2brain/avif` — MIT License
- `github.com/tetratelabs/wazero` — Apache License 2.0
- `github.com/ebitengine/purego` — Apache License 2.0
- `golang.org/x/image` — BSD 3-Clause License
- `modernc.org/tk9.0` and its ModernC dependencies — licenses supplied by
  their respective modules

The AVIF module contains CGo-free WebAssembly builds of libavif and AV1 codec
components. Their license texts are supplied in that module's `lib` directory.

Go downloads the exact module versions recorded in `go.mod` and `go.sum`.
Complete license texts are present in those downloaded modules and their source
repositories.
