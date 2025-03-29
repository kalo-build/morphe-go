package registry

import (
	"github.com/kalo-build/morphe-go/pkg/registry/cfg"
)

type LoadMorpheRegistryHooks struct {
	OnRegistryLoadStart   OnRegistryLoadStartHook
	OnRegistryLoadSuccess OnRegistryLoadSuccessHook
	OnRegistryLoadFailure OnRegistryLoadFailureHook
}

type OnRegistryLoadStartHook = func(config cfg.MorpheLoadRegistryConfig) (cfg.MorpheLoadRegistryConfig, error)
type OnRegistryLoadSuccessHook = func(registry Registry) (Registry, error)
type OnRegistryLoadFailureHook = func(config cfg.MorpheLoadRegistryConfig, registry Registry, loadFailure error) error
