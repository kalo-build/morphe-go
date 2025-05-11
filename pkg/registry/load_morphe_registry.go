package registry

import (
	"github.com/kalo-build/morphe-go/pkg/registry/cfg"
)

func LoadMorpheRegistry(hooks LoadMorpheRegistryHooks, config cfg.MorpheLoadRegistryConfig) (*Registry, error) {
	config, loadStartErr := triggerLoadRegistryStart(hooks, config)
	if loadStartErr != nil {
		return nil, triggerLoadRegistryFailure(hooks, config, nil, loadStartErr)
	}

	r := NewRegistry()

	loadErr := loadConfiguredRegistry(config, r)
	if loadErr != nil {
		return nil, triggerLoadRegistryFailure(hooks, config, r, loadErr)
	}

	// Validate the registry to ensure consistency
	validateErr := r.ValidateRegistry()
	if validateErr != nil {
		return nil, triggerLoadRegistryFailure(hooks, config, r, validateErr)
	}

	r, loadSuccessErr := triggerLoadRegistrySuccess(hooks, r)
	if loadSuccessErr != nil {
		return nil, triggerLoadRegistryFailure(hooks, config, r, loadSuccessErr)
	}

	return r, nil
}

func loadConfiguredRegistry(config cfg.MorpheLoadRegistryConfig, r *Registry) error {
	enumsErr := r.LoadEnumsFromDirectory(config.RegistryEnumsDirPath)
	if enumsErr != nil {
		return enumsErr
	}

	modelsErr := r.LoadModelsFromDirectory(config.RegistryModelsDirPath)
	if modelsErr != nil {
		return modelsErr
	}

	structuresErr := r.LoadStructuresFromDirectory(config.RegistryStructuresDirPath)
	if structuresErr != nil {
		return structuresErr
	}

	entitiesErr := r.LoadEntitiesFromDirectory(config.RegistryEntitiesDirPath)
	if entitiesErr != nil {
		return entitiesErr
	}

	return nil
}

func triggerLoadRegistryStart(hooks LoadMorpheRegistryHooks, config cfg.MorpheLoadRegistryConfig) (cfg.MorpheLoadRegistryConfig, error) {
	if hooks.OnRegistryLoadStart == nil {
		return config, nil
	}

	updatedConfig, startErr := hooks.OnRegistryLoadStart(config)
	if startErr != nil {
		return cfg.MorpheLoadRegistryConfig{}, startErr
	}

	return updatedConfig, nil
}

func triggerLoadRegistrySuccess(hooks LoadMorpheRegistryHooks, r *Registry) (*Registry, error) {
	if hooks.OnRegistryLoadSuccess == nil {
		return r, nil
	}
	if r == nil {
		return nil, ErrRegistryNotInitialized
	}
	registry, successErr := hooks.OnRegistryLoadSuccess(*r.DeepClone())
	if successErr != nil {
		return nil, successErr
	}
	r = &registry
	return r, nil
}

func triggerLoadRegistryFailure(hooks LoadMorpheRegistryHooks, config cfg.MorpheLoadRegistryConfig, r *Registry, failureErr error) error {
	if hooks.OnRegistryLoadFailure == nil {
		return failureErr
	}

	if r == nil {
		return hooks.OnRegistryLoadFailure(config, Registry{}, failureErr)
	}

	return hooks.OnRegistryLoadFailure(config, *r.DeepClone(), failureErr)
}
