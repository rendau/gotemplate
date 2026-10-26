import os
import shutil
import subprocess
import sys
from pathlib import Path

SKILLS_REPO = "https://github.com/rendau/go-service-skill.git"


def run(cmd):
    """Run a command, return True on success, swallow any failure."""
    sys.stdout.flush()  # keep our messages and the command's output in order
    try:
        subprocess.run(cmd, check=True)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False


def ask(question, default):
    """Ask in the terminal; without one (CI, --no-input) take the default."""
    if not sys.stdin.isatty():
        return default
    return input("%s [%s]: " % (question, default)).strip() or default


def is_skills_clone(path):
    return (path / "install.sh").is_file() and (path / "skills/crud/SKILL.md").is_file() and (path / ".git").exists()


def installed_skills_repo():
    """Clone of go-service-skill the installed skills are linked to, if any."""
    crud = Path(os.environ.get("CLAUDE_SKILLS_DIR", "~/.claude/skills")).expanduser() / "crud"
    if crud.is_symlink() and crud.exists():
        repo = crud.resolve().parent.parent  # <repo>/skills/crud
        if is_skills_clone(repo):
            return repo
    return None


def install_skills():
    """Update the installed skills or clone go-service-skill and link them (best-effort)."""
    if not shutil.which("git") or not shutil.which("bash"):
        return False
    print("\nClaude Code skills (go-service-skill):")

    repo = installed_skills_repo()
    if repo:
        print("  already installed from %s" % repo)
        if ask("Update them (git pull + install.sh)?", "y").lower() not in ("y", "yes"):
            return True
        if not run(["git", "-C", str(repo), "pull", "-q", "--ff-only"]):
            print("  could not update %s, linking the current version" % repo)
        return run(["bash", str(repo / "install.sh")])

    base = Path.cwd().parent  # relative paths — from where the service was generated
    default = os.environ.get("GO_SERVICE_SKILL_DIR") or str(base / "go-service-skill")
    repo = (base / Path(ask("Where to clone go-service-skill", default)).expanduser()).resolve()
    if is_skills_clone(repo):
        print("  %s is already a clone, updating it" % repo)
        run(["git", "-C", str(repo), "pull", "-q", "--ff-only"])
    elif repo.exists() and any(repo.iterdir()):
        print("  %s exists and is not empty, choose another folder" % repo)
        return False
    else:
        print("  cloning into %s" % repo)
        if not run(["git", "clone", "-q", SKILLS_REPO, str(repo)]):
            return False
    return run(["bash", str(repo / "install.sh")])


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

# Claude Code skills for this stack (best-effort, asked at generation: install_skills).
skills_wanted = "{{ cookiecutter.install_skills }}".strip().lower() in ("true", "yes", "y", "1")
if skills_wanted and not install_skills():
    print("\nSkills were not installed. To install them later:")
    print("  git clone %s <folder> && <folder>/install.sh" % SKILLS_REPO)

print("\nNext steps:")
print("  1. make refresh-vendor-proto   # refresh shared protos (needs SSH to the protos repo)")
print("  2. add your .proto files under api/proto/{{ cookiecutter.__proto_pkg }}/")
print("  3. make generate-proto")
print("  4. make build\n")
