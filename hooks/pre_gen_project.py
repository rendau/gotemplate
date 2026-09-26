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

# Манифест сервиса для pulse (internal/app/manifest.go): сведения о сервисе — от человека.
title = {{ cookiecutter.service_title | jsonify }}
description = {{ cookiecutter.service_description | jsonify }}
team = {{ cookiecutter.owner_team | jsonify }}

errors = []
if not team.strip():
    errors.append("owner_team is required: the team that owns the service (same name in all its services)")
if not description.strip() or description.startswith("Что делает сервис"):
    errors.append("service_description is required: what the service does, in words people search by")
# лимиты стандарта: длиннее — pulsekit не даст сервису стартовать
if len(title.strip()) > 100:
    errors.append("service_title is %d characters, the pulse standard allows 100" % len(title.strip()))
if len(description.strip()) > 500:
    errors.append("service_description is %d characters, the pulse standard allows 500" % len(description.strip()))
for name, value in (("service_title", title), ("service_description", description), ("owner_team", team)):
    if '"' in value or "\\" in value or "\n" in value:
        errors.append("%s must not contain quotes, backslashes or line breaks" % name)
if errors:
    sys.stderr.write("\nERROR:\n  " + "\n  ".join(errors) + "\n\n")
    sys.exit(1)
