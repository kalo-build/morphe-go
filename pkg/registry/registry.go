package registry

import (
	"fmt"
	"log"
	"os"
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

// ValidateRegistry checks if the registry state is valid
// Specifically, it ensures that if entities exist, models must also exist
func (r *Registry) ValidateRegistry() error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(r.entities) > 0 && len(r.models) == 0 {
		return fmt.Errorf("invalid registry state: entities exist but no models are defined")
	}

	return nil
}

// HasEnums returns true if the registry has enums defined
func (r *Registry) HasEnums() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.enums) > 0
}

// HasModels returns true if the registry has models defined
func (r *Registry) HasModels() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.models) > 0
}

// HasStructures returns true if the registry has structures defined
func (r *Registry) HasStructures() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.structures) > 0
}

// HasEntities returns true if the registry has entities defined
func (r *Registry) HasEntities() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.entities) > 0
}

// SetEnum is a thread-safe way to write an enum to the registry
func (r *Registry) SetEnum(name string, enum yaml.Enum) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.enums == nil {
		r.enums = make(map[string]yaml.Enum)
	}

	r.enums[name] = enum
}

// GetEnum returns a thread-safe copy of a registry enum
func (r *Registry) GetEnum(name string) (yaml.Enum, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if r.enums == nil {
		return yaml.Enum{}, fmt.Errorf("no enums in registry, enum with name '%s' not found", name)
	}

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

	if r.enums == nil {
		return make(map[string]yaml.Enum)
	}

	enumsClone := clone.DeepCloneMap(r.enums)
	return enumsClone
}

// SetModel is a thread-safe way to write a model to the registry
func (r *Registry) SetModel(name string, model yaml.Model) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.models == nil {
		r.models = make(map[string]yaml.Model)
	}

	r.models[name] = model
}

// GetModel returns a thread-safe copy of a registry model
func (r *Registry) GetModel(name string) (yaml.Model, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if r.models == nil {
		return yaml.Model{}, fmt.Errorf("no models in registry, model with name '%s' not found", name)
	}

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

	if r.models == nil {
		return make(map[string]yaml.Model)
	}

	modelsClone := clone.DeepCloneMap(r.models)
	return modelsClone
}

// SetEntity is a thread-safe way to write an entity to the registry
func (r *Registry) SetEntity(name string, entity yaml.Entity) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.entities == nil {
		r.entities = make(map[string]yaml.Entity)
	}

	r.entities[name] = entity
}

// GetEntity returns a thread-safe copy of a registry entity
func (r *Registry) GetEntity(name string) (yaml.Entity, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if r.entities == nil {
		return yaml.Entity{}, fmt.Errorf("no entities in registry, entity with name '%s' not found", name)
	}

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

	if r.entities == nil {
		return make(map[string]yaml.Entity)
	}

	entitiesClone := clone.DeepCloneMap(r.entities)
	return entitiesClone
}

// SetStructure is a thread-safe way to write a structure to the registry
func (r *Registry) SetStructure(name string, structure yaml.Structure) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.structures == nil {
		r.structures = make(map[string]yaml.Structure)
	}

	r.structures[name] = structure
}

// GetStructure returns a thread-safe copy of a registry structure
func (r *Registry) GetStructure(name string) (yaml.Structure, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if r.structures == nil {
		return yaml.Structure{}, fmt.Errorf("no structures in registry, structure with name '%s' not found", name)
	}

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

	if r.structures == nil {
		return make(map[string]yaml.Structure)
	}

	structuresClone := clone.DeepCloneMap(r.structures)
	return structuresClone
}

func (r *Registry) DeepClone() *Registry {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	registryCopy := &Registry{
		enums:      make(map[string]yaml.Enum),
		models:     make(map[string]yaml.Model),
		structures: make(map[string]yaml.Structure),
		entities:   make(map[string]yaml.Entity),
	}

	if r.enums != nil {
		registryCopy.enums = clone.DeepCloneMap(r.enums)
	}

	if r.models != nil {
		registryCopy.models = clone.DeepCloneMap(r.models)
	}

	if r.structures != nil {
		registryCopy.structures = clone.DeepCloneMap(r.structures)
	}

	if r.entities != nil {
		registryCopy.entities = clone.DeepCloneMap(r.entities)
	}

	return registryCopy
}

func (r *Registry) LoadEnumsFromDirectory(dirPath string) error {
	// Check if directory exists
	if exists, err := directoryExists(dirPath); err != nil {
		return err
	} else if !exists {
		log.Printf("Warning: Enums directory does not exist: %s. Skipping enum loading.", dirPath)
		return nil
	}

	allEnums, unmarshalErr := yamlfile.UnmarshalAllYAMLFiles[yaml.Enum](dirPath, EnumFileSuffix)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	// Normalize whitespace in string fields
	yaml.NormalizeAllEnums(allEnums)

	if len(allEnums) == 0 {
		log.Printf("Warning: No enum files found in directory: %s. Skipping enum loading.", dirPath)
		return nil
	}

	loadErr := r.loadEnumDefinitions(allEnums)
	return loadErr
}

func (r *Registry) LoadModelsFromDirectory(dirPath string) error {
	// Check if directory exists
	if exists, err := directoryExists(dirPath); err != nil {
		return err
	} else if !exists {
		log.Printf("Warning: Models directory does not exist: %s. Skipping model loading.", dirPath)
		return nil
	}

	allModels, unmarshalErr := yamlfile.UnmarshalAllYAMLFiles[yaml.Model](dirPath, ModelFileSuffix)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	// Normalize whitespace in string fields
	yaml.NormalizeAllModels(allModels)

	if len(allModels) == 0 {
		log.Printf("Warning: No model files found in directory: %s. Skipping model loading.", dirPath)
		return nil
	}

	loadErr := r.loadModelDefinitions(allModels)
	return loadErr
}

func (r *Registry) LoadEntitiesFromDirectory(dirPath string) error {
	// Check if directory exists
	if exists, err := directoryExists(dirPath); err != nil {
		return err
	} else if !exists {
		log.Printf("Warning: Entities directory does not exist: %s. Skipping entity loading.", dirPath)
		return nil
	}

	allEntities, unmarshalErr := yamlfile.UnmarshalAllYAMLFiles[yaml.Entity](dirPath, EntityFileSuffix)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	// Normalize whitespace in string fields
	yaml.NormalizeAllEntities(allEntities)

	if len(allEntities) == 0 {
		log.Printf("Warning: No entity files found in directory: %s. Skipping entity loading.", dirPath)
		return nil
	}

	// If entities are present but models are not, we should still enforce the dependency
	if len(r.models) == 0 {
		return fmt.Errorf("attempted to load entities but no models are defined in registry")
	}

	loadErr := r.loadEntityDefinitions(allEntities)
	return loadErr
}

func (r *Registry) LoadStructuresFromDirectory(dirPath string) error {
	// Check if directory exists
	if exists, err := directoryExists(dirPath); err != nil {
		return err
	} else if !exists {
		log.Printf("Warning: Structures directory does not exist: %s. Skipping structure loading.", dirPath)
		return nil
	}

	allStructures, unmarshalErr := yamlfile.UnmarshalAllYAMLFiles[yaml.Structure](dirPath, StructureFileSuffix)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	// Normalize whitespace in string fields
	yaml.NormalizeAllStructures(allStructures)

	if len(allStructures) == 0 {
		log.Printf("Warning: No structure files found in directory: %s. Skipping structure loading.", dirPath)
		return nil
	}

	loadErr := r.loadStructureDefinitions(allStructures)
	return loadErr
}

// Helper function to check if a directory exists
func directoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
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

	// Check if models are defined when entities are loaded
	if len(r.entities) > 0 && len(r.models) == 0 {
		return fmt.Errorf("attempted to load entities but no models are defined in registry")
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
