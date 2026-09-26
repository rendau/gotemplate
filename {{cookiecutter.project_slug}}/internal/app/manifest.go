package app

import (
	"{{ cookiecutter.__go_module }}/internal/constant"
	"{{ cookiecutter.__go_module }}/internal/infra/pulsekit"
)

// newPulsekit: service manifest for pulse (github.com/mechta-market/pulse, docs/service-manifest.md).
// Dependencies: a.pulsekit.Depend in app.go; endpoints: pulsekit.Handle.
// Сведения ниже — от человека: агент уточняет их у разработчика.
func newPulsekit() *pulsekit.Kit {
	return pulsekit.New(pulsekit.Config{}, pulsekit.Service{
		Name:        constant.ServiceName,
		Title:       "{{ cookiecutter.service_title }}",
		Description: "{{ cookiecutter.service_description }}",
		OwnerTeam:   "{{ cookiecutter.owner_team }}",
		Criticality: "{{ cookiecutter.criticality }}",
		RepoUrl:     "https://{{ cookiecutter.__go_module }}",
	}, pulsekit.Build{Version: constant.Version, Commit: constant.Commit, BuiltAt: constant.BuiltAt})
}
