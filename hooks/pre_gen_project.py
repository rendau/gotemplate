import re
import sys

slug = "{{ cookiecutter.project_slug }}"

if not re.match(r"^[a-z][a-z0-9_]*$", slug):
    sys.stderr.write(
        "\nERROR: project_slug '%s' is invalid.\n"
        "Use lowercase letters, digits and underscores; must start with a letter.\n"
        "(it becomes the Go binary name, the ServiceName constant and the proto package).\n\n"
        % slug
    )
    sys.exit(1)
