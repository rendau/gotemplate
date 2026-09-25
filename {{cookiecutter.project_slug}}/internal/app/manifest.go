package app

import (
	"net/url"

	"{{ cookiecutter.__go_module }}/internal/constant"
	"{{ cookiecutter.__go_module }}/internal/infra/pulsekit"
)

// Манифест сервиса для pulse (github.com/mechta-market/pulse, docs/service-manifest.md): сервис
// сам рассказывает о себе, о своих зависимостях (a.pulsekit.Depend рядом с созданием клиента в
// app.go) и диагностических ручках (pulsekit.Handle). Сведения ниже — от человека: агент,
// внедряющий стандарт, уточняет их у разработчика, а не выдумывает.

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

// hostOf — хост адреса для target зависимости: без схемы, пути и учётных данных.
func hostOf(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return raw
	}
	return u.Host
}
