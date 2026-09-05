# containers

A minimal container runtime written in Go, built to explore how Linux
containers actually work under the hood. It isolates a process using Linux
namespaces, builds a throwaway root filesystem by copying only the binaries you
ask for (plus their shared-library dependencies), and `chroot`s into it before
executing your command.

This is an educational project, not a production runtime. It is a good way to
see `unshare(2)`, `chroot(2)`, and dynamic-linker resolution work together
without the abstraction of Docker or `runc`.

## How it works

Running a command goes through four steps (`container/cmd/main.go`):

1. **Namespace isolation** (`namespace.go`) — the process calls `unshare` with
   new mount, UTS, IPC, and network namespaces (`CLONE_NEWNS`, `CLONE_NEWUTS`,
   `CLONE_NEWIPC`, `CLONE_NEWNET`).
2. **Sandbox construction** (`setup.go`) — each entry in `req.binaries` is
   copied into `container/sandbox/`. For every binary, its shared libraries are
   resolved with `ldd` and copied into `lib/` and `lib64/` so it can run in
   isolation.
3. **Filesystem isolation** (`chroot.go`) — the process `chroot`s into
   `container/sandbox/` and changes its working directory to `/`.
4. **Execution** — the requested command is run inside the sandbox with stdio
   wired to the host terminal.

## Requirements

- Linux (uses `unshare` and `chroot` syscalls)
- Go 1.23+
- Root privileges (namespace creation and `chroot` require them)

## Usage

### 1. Declare the binaries your command needs

Edit `req.binaries` (one path per line). You can list individual executables or
whole directories:

```
../../iso/toRun
/bin/ls
/bin/pwd
/bin/tree
```

Listing `/bin` and `/usr/bin` has the same effect — every executable in the
directory is copied in.

### 2. Run a command in the container

From `container/cmd/`:

```sh
sudo go run . ./pwd
sudo go run . ./ls
sudo go run . ./tree /some/dir
```

Or build a binary and run that:

```sh
sudo go build -o cmd .
sudo ./cmd ./ls
```

## Repository layout

| Path                 | Purpose                                                        |
| -------------------- | ------------------------------------------------------------- |
| `container/cmd/`     | The runtime: namespace setup, sandbox build, chroot, exec.   |
| `container/sandbox/` | Generated root filesystem (populated at runtime).            |
| `iso/toRun.go`       | A probe program that demonstrates remaining isolation gaps.  |
| `req.binaries`       | List of binaries/directories to include in the sandbox.      |

## Known limitations

- No PID or user namespace, so host processes remain visible via `/proc` and
  UIDs are not remapped (see `iso/toRun.go`).
- No cgroup resource limits.
- The sandbox is not cleaned up between runs.
- Flattens all copied binaries into the sandbox root rather than preserving
  their original paths.

## License

No license has been specified yet.
