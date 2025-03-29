package registry

import (
	"fmt"
	"sync"

	"github.com/kalo-build/clone"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlfile"
)

const EnumFileSuffix = ".enum"
const ModelFileSuffix = ".mod"
const EntityFileSuffix = ".ent"
const StructureFileSuffix = ".str"

type Registry struct {
	mutex sync.RWMutex

	enums      map[string]yaml.Enum      `yaml:"enums"`
	models     map[string]yaml.Model     `yaml:"models"`
	structures map[string]yaml.Structure `yaml:"structures"`
	entities   map[string]yaml.Entity    `yaml:"entities"`
}

// SetEnum is a thread-safe way to write an enum to the registry
func (r *Registry) SetEnum(name string, enum yaml.Enum) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.enums[name] = enum
}

// GetEnum returns a thread-safe copy of a registry enum
func (r *Registry) GetEnum(name string) (yaml.Enum, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	enum, enumFound := r.enums[name]
	if !enumFound {
		return yaml.Enum{}, fmt.Errorf("enum with name '%s' not found registry", name)
	}
	enumClone := enum.DeepClone()
	return enumClone, nil
}

// GetAllEnums returns a thread-safe copy of all registry enums
func (r *Registry) GetAllEnums() map[string]yaml.Enum {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	enumsClone := clone.DeepCloneMap(r.enums)
	return enumsClone
}

// SetModel is a thread-safe way to write a model to the registry
func (r *Registry) SetModel(name string, model yaml.Model) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.models[name] = model
}

// GetModel returns a thread-safe copy of a registry model
func (r *Registry) GetModel(name string) (yaml.Model, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	model, modelFound := r.models[name]
	if !modelFound {
		return yaml.Model{}, fmt.Errorf("model with name '%s' not found registry", name)
	}
	modelClone := model.DeepClone()
	return modelClone, nil
}

// GetAllModels returns a thread-safe copy of all registry models
func (r *Registry) GetAllModels() map[string]yaml.Model {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	modelsClone := clone.DeepCloneMap(r.models)
	return modelsClone
}

// SetEntity is a thread-safe way to write an entity to the registry
func (r *Registry) SetEntity(name string, entity yaml.Entity) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.entities[name] = entity
}

// GetEntity returns a thread-safe copy of a registry entity
func (r *Registry) GetEntity(name string) (yaml.Entity, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	entity, entityFound := r.entities[name]
	if !entityFound {
		return yaml.Entity{}, fmt.Errorf("entity with name '%s' not found registry", name)
	}
	entityClone := entity.DeepClone()
	return entityClone, nil
}

// GetAllEntities returns a thread-safe copy of all registry entities
func (r *Registry) GetAllEntities() map[string]yaml.Entity {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	entitiesClone := clone.DeepCloneMap(r.entities)
	return entitiesClone
}

// SetStructure is a thread-safe way to write a structure to the registry
func (r *Registry) SetStructure(name string, structure yaml.Structure) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.structures[name] = structure
}

// GetStructure returns a thread-safe copy of a registry structure
func (r *Registry) GetStructure(name string) (yaml.Structure, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	structure, structureFound := r.structures[name]
	if !structureFound {
		return yaml.Structure{}, fmt.Errorf("structure with name '%s' not found registry", name)
	}
	structureClone := structure.DeepClone()
	return structureClone, nil
}

// GetAllStructures returns a thread-safe copy of all registry structures
func (r *Registry) GetAllStructures() map[string]yaml.Structure {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	structuresClone := clone.DeepCloneMap(r.structures)
	return structuresClone
}

func (r *Registry) DeepClone() *Registry {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	registryCopy := &Registry{
		enums:      clone.DeepCloneMap(r.enums),
		models:     clone.DeepCloneMap(r.models),
		structures: clone.DeepCloneMap(r.structures),
		entities:   clone.DeepCloneMap(r.entities),
	}

	return registryCopy
}

func (r *Registry) LoadEnumsFromDirectory(dirPath string) error {
	allEnums, unmarshalErr := yamlfile.UnmarshalAllYAMLFiles[yaml.Enum](dirPath, EnumFileSuffix)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	loadErr := r.loadEnumDefinitions(allEnums)
	return loadErr
}

func (r *Registry) LoadModelsFromDirectory(dirPath string) error {
	allModels, unmarshalErr := yamlfile.UnmarshalAllYAMLFiles[yaml.Model](dirPath, ModelFileSuffix)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	loadErr := r.loadModelDefinitions(allModels)
	return loadErr
}

func (r *Registry) LoadEntitiesFromDirectory(dirPath string) error {
	allEntities, unmarshalErr := yamlfile.UnmarshalAllYAMLFiles[yaml.Entity](dirPath, EntityFileSuffix)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	loadErr := r.loadEntityDefinitions(allEntities)
	return loadErr
}

func (r *Registry) LoadStructuresFromDirectory(dirPath string) error {
	allStructures, unmarshalErr := yamlfile.UnmarshalAllYAMLFiles[yaml.Structure](dirPath, StructureFileSuffix)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	loadErr := r.loadStructureDefinitions(allStructures)
	return loadErr
}

func (r *Registry) loadEnumDefinitions(allEnums map[string]yaml.Enum) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.enums == nil {
		r.enums = make(map[string]yaml.Enum)
	}

	for enumPathAbs, enum := range allEnums {
		_, nameConflict := r.enums[enum.Name]
		if nameConflict {
			return fmt.Errorf("enum name '%s' already exists in registry (conflict: %s)", enum.Name, enumPathAbs)
		}

		r.enums[enum.Name] = enum
	}
	return nil
}

func (r *Registry) loadModelDefinitions(allModels map[string]yaml.Model) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.models == nil {
		r.models = make(map[string]yaml.Model)
	}

	for modelPathAbs, model := range allModels {
		_, nameConflict := r.models[model.Name]
		if nameConflict {
			return fmt.Errorf("model name '%s' already exists in registry (conflict: %s)", model.Name, modelPathAbs)
		}

		r.models[model.Name] = model
	}
	return nil
}

func (r *Registry) loadEntityDefinitions(allEntities map[string]yaml.Entity) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.entities == nil {
		r.entities = make(map[string]yaml.Entity)
	}

	for entityPathAbs, entity := range allEntities {
		_, nameConflict := r.entities[entity.Name]
		if nameConflict {
			return fmt.Errorf("entity name '%s' already exists in registry (conflict: %s)", entity.Name, entityPathAbs)
		}

		r.entities[entity.Name] = entity
	}
	return nil
}

func (r *Registry) loadStructureDefinitions(allStructures map[string]yaml.Structure) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.structures == nil {
		r.structures = make(map[string]yaml.Structure)
	}

	for structurePathAbs, structure := range allStructures {
		_, nameConflict := r.structures[structure.Name]
		if nameConflict {
			return fmt.Errorf("structure name '%s' already exists in registry (conflict: %s)", structure.Name, structurePathAbs)
		}

		r.structures[structure.Name] = structure
	}
	return nil
}
