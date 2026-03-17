// Package repository contains analyzer fixtures for lower-layer packages.
package repository

import "layerguard/service" // want "must not import higher-level package"

func LoadRepositoryName() string {
	return service.LoadServiceName()
}
