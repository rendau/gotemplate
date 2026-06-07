import shutil
import subprocess


def run(cmd):
    """Run a command, return True on success, swallow any failure."""
    try:
        subprocess.run(cmd, check=True)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False


print("\nProject '{{ cookiecutter.project_slug }}' generated.")
print("Module: {{ cookiecutter.__go_module }}\n")

# Initialize a fresh git repository with an initial commit (best-effort).
if shutil.which("git"):
    if run(["git", "init", "-q"]):
        run(["git", "add", "-A"])
        run(["git", "commit", "-q", "-m", "Initial commit from gotemplate"])

# Tidy go modules (best-effort: needs `go` and network access).
if shutil.which("go"):
    print("Running 'go mod tidy' (best-effort)...")
    run(["go", "mod", "tidy"])

print("\nNext steps:")
print("  1. make refresh-vendor-proto   # refresh shared protos (needs SSH to the protos repo)")
print("  2. add your .proto files under api/proto/{{ cookiecutter.__proto_pkg }}/")
print("  3. make generate-proto")
print("  4. make build\n")
