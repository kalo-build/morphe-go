package registry_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/morphe-go/internal/testutils"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/morphe-go/pkg/yaml"
)

type LoadMorpheRegistryTestSuite struct {
	suite.Suite

	TestDirPath string

	EnumsDirPath      string
	ModelsDirPath     string
	StructuresDirPath string
	EntitiesDirPath   string
}

func TestLoadMorpheRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(LoadMorpheRegistryTestSuite))
}

func (suite *LoadMorpheRegistryTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()

	suite.EnumsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "enums")
	suite.ModelsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "models")
	suite.StructuresDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "structures")
	suite.EntitiesDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "entities")
}

func (suite *LoadMorpheRegistryTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *LoadMorpheRegistryTestSuite) TestLoadMorpheRegistry() {
	loadHooks := registry.LoadMorpheRegistryHooks{}
	config := cfg.MorpheLoadRegistryConfig{
		RegistryEnumsDirPath:      suite.EnumsDirPath,
		RegistryModelsDirPath:     suite.ModelsDirPath,
		RegistryStructuresDirPath: suite.StructuresDirPath,
		RegistryEntitiesDirPath:   suite.EntitiesDirPath,
	}

	r, registryErr := registry.LoadMorpheRegistry(loadHooks, config)

	suite.NoError(registryErr)
	suite.NotNil(r)

	enum0, enumErr1 := r.GetEnum("Nationality")
	suite.Nil(enumErr1)
	suite.Equal(enum0.Name, "Nationality")
	suite.Equal(enum0.Type, yaml.EnumTypeString)

	suite.Len(enum0.Entries, 3)

	entry10, entryExists10 := enum0.Entries["US"]
	suite.True(entryExists10)
	suite.Equal(entry10, "American")

	entry11, entryExists11 := enum0.Entries["DE"]
	suite.True(entryExists11)
	suite.Equal(entry11, "German")

	entry12, entryExists12 := enum0.Entries["FR"]
	suite.True(entryExists12)
	suite.Equal(entry12, "French")

	allModels := r.GetAllModels()
	suite.Len(allModels, 2)

	model0, modelExists0 := allModels["Person"]
	suite.True(modelExists0)
	suite.Equal(model0.Name, "Person")

	suite.Len(model0.Fields, 4)

	modelField00, fieldExists00 := model0.Fields["ID"]
	suite.True(fieldExists00)
	suite.Equal(modelField00.Type, yaml.ModelFieldTypeAutoIncrement)
	suite.Len(modelField00.Attributes, 1)
	suite.Equal(modelField00.Attributes[0], "mandatory")

	modelField01, fieldExists01 := model0.Fields["FirstName"]
	suite.True(fieldExists01)
	suite.Equal(modelField01.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField01.Attributes, 0)

	modelField02, fieldExists02 := model0.Fields["LastName"]
	suite.True(fieldExists02)
	suite.Equal(modelField02.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField02.Attributes, 0)

	modelField03, fieldExists03 := model0.Fields["Nationality"]
	suite.True(fieldExists03)
	suite.Equal(modelField03.Type, yaml.ModelFieldType("Nationality"))
	suite.Len(modelField03.Attributes, 0)

	suite.Len(model0.Identifiers, 2)
	modelIDs00, idsExist00 := model0.Identifiers["primary"]
	suite.True(idsExist00)
	suite.ElementsMatch(modelIDs00.Fields, []string{"ID"})

	modelIDs01, idsExist01 := model0.Identifiers["name"]
	suite.True(idsExist01)
	suite.ElementsMatch(modelIDs01.Fields, []string{"FirstName", "LastName"})

	suite.Len(model0.Related, 1)

	modelRelated00, relatedExists00 := model0.Related["ContactInfo"]
	suite.True(relatedExists00)
	suite.Equal(modelRelated00.Type, "HasOne")

	model1, modelExists1 := allModels["ContactInfo"]
	suite.True(modelExists1)
	suite.Equal(model1.Name, "ContactInfo")

	suite.Len(model1.Fields, 2)

	modelField10, fieldExists10 := model1.Fields["ID"]
	suite.True(fieldExists10)
	suite.Equal(modelField10.Type, yaml.ModelFieldTypeAutoIncrement)
	suite.Len(modelField10.Attributes, 1)
	suite.Equal(modelField10.Attributes[0], "mandatory")

	modelField11, fieldExists11 := model1.Fields["Email"]
	suite.True(fieldExists11)
	suite.Equal(modelField11.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField11.Attributes, 0)

	suite.Len(model1.Identifiers, 2)
	modelID10, idExists10 := model1.Identifiers["primary"]
	suite.True(idExists10)
	suite.ElementsMatch(modelID10.Fields, []string{"ID"})

	modelIDs11, idsExist11 := model1.Identifiers["email"]
	suite.True(idsExist11)
	suite.ElementsMatch(modelIDs11.Fields, []string{"Email"})

	suite.Len(model1.Related, 1)

	modelRelated10, relatedExists10 := model1.Related["Person"]
	suite.True(relatedExists10)
	suite.Equal(modelRelated10.Type, "ForOne")

	structure0, structureErr0 := r.GetStructure("Address")
	suite.Nil(structureErr0)
	suite.Equal(structure0.Name, "Address")

	suite.Len(structure0.Fields, 4)

	structureField00, fieldExists00 := structure0.Fields["Street"]
	suite.True(fieldExists00)
	suite.Equal(structureField00.Type, yaml.StructureFieldTypeString)

	structureField01, fieldExists01 := structure0.Fields["HouseNr"]
	suite.True(fieldExists01)
	suite.Equal(structureField01.Type, yaml.StructureFieldTypeString)

	structureField02, fieldExists02 := structure0.Fields["ZipCode"]
	suite.True(fieldExists02)
	suite.Equal(structureField02.Type, yaml.StructureFieldTypeString)

	structureField03, fieldExists03 := structure0.Fields["City"]
	suite.True(fieldExists03)
	suite.Equal(structureField03.Type, yaml.StructureFieldTypeString)

	entity0, entityErr0 := r.GetEntity("Person")
	suite.Nil(entityErr0)
	suite.Equal(entity0.Name, "Person")

	suite.Len(entity0.Fields, 5)

	entityField00, fieldExists00 := entity0.Fields["UUID"]
	suite.True(fieldExists00)
	suite.Equal(entityField00.Type, yaml.ModelFieldPath("Person.UUID"))
	suite.Len(entityField00.Attributes, 2)
	suite.Contains(entityField00.Attributes, "immutable")
	suite.Contains(entityField00.Attributes, "mandatory")

	entityField01, fieldExists01 := entity0.Fields["ID"]
	suite.True(fieldExists01)
	suite.Equal(entityField01.Type, yaml.ModelFieldPath("Person.ID"))
	suite.Len(entityField01.Attributes, 0)

	entityField02, fieldExists02 := entity0.Fields["FirstName"]
	suite.True(fieldExists02)
	suite.Equal(entityField02.Type, yaml.ModelFieldPath("Person.FirstName"))
	suite.Len(entityField02.Attributes, 0)

	entityField03, fieldExists03 := entity0.Fields["LastName"]
	suite.True(fieldExists03)
	suite.Equal(entityField03.Type, yaml.ModelFieldPath("Person.LastName"))
	suite.Len(entityField03.Attributes, 0)

	entityField04, fieldExists04 := entity0.Fields["Email"]
	suite.True(fieldExists04)
	suite.Equal(entityField04.Type, yaml.ModelFieldPath("Person.ContactInfo.Email"))
	suite.Len(entityField04.Attributes, 0)

	suite.Len(entity0.Identifiers, 1)
	entityID00, idExists00 := entity0.Identifiers["primary"]
	suite.True(idExists00)
	suite.ElementsMatch(entityID00.Fields, []string{"UUID"})
}

func (suite *LoadMorpheRegistryTestSuite) TestLoadMorpheRegistry_Failure_InvalidPaths() {
	loadHooks := registry.LoadMorpheRegistryHooks{}
	config := cfg.MorpheLoadRegistryConfig{
		RegistryEnumsDirPath:      "invalid path",
		RegistryModelsDirPath:     "invalid path",
		RegistryStructuresDirPath: "invalid path",
		RegistryEntitiesDirPath:   "invalid path",
	}

	r, registryErr := registry.LoadMorpheRegistry(loadHooks, config)

	suite.ErrorContains(registryErr, "error reading directory 'invalid path")
	suite.Nil(r)
}

func (suite *LoadMorpheRegistryTestSuite) TestLoadMorpheRegistry_StartHook_Successful() {
	loadHooks := registry.LoadMorpheRegistryHooks{
		OnRegistryLoadStart: func(config cfg.MorpheLoadRegistryConfig) (cfg.MorpheLoadRegistryConfig, error) {
			config.RegistryEnumsDirPath = suite.EnumsDirPath
			config.RegistryModelsDirPath = suite.ModelsDirPath
			config.RegistryStructuresDirPath = suite.StructuresDirPath
			config.RegistryEntitiesDirPath = suite.EntitiesDirPath
			return config, nil
		},
	}
	config := cfg.MorpheLoadRegistryConfig{
		RegistryEnumsDirPath:      "invalid path",
		RegistryModelsDirPath:     "invalid/path",
		RegistryStructuresDirPath: "invalid/path",
		RegistryEntitiesDirPath:   "invalid/path",
	}

	r, registryErr := registry.LoadMorpheRegistry(loadHooks, config)

	suite.NoError(registryErr)
	suite.NotNil(r)

	allModels := r.GetAllModels()
	suite.Len(allModels, 2)

	model0, modelExists0 := allModels["Person"]
	suite.True(modelExists0)
	suite.Equal(model0.Name, "Person")

	suite.Len(model0.Fields, 4)

	modelField00, fieldExists00 := model0.Fields["ID"]
	suite.True(fieldExists00)
	suite.Equal(modelField00.Type, yaml.ModelFieldTypeAutoIncrement)
	suite.Len(modelField00.Attributes, 1)
	suite.Equal(modelField00.Attributes[0], "mandatory")

	modelField01, fieldExists01 := model0.Fields["FirstName"]
	suite.True(fieldExists01)
	suite.Equal(modelField01.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField01.Attributes, 0)

	modelField02, fieldExists02 := model0.Fields["LastName"]
	suite.True(fieldExists02)
	suite.Equal(modelField02.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField02.Attributes, 0)

	modelField03, fieldExists03 := model0.Fields["Nationality"]
	suite.True(fieldExists03)
	suite.Equal(modelField03.Type, yaml.ModelFieldType("Nationality"))
	suite.Len(modelField03.Attributes, 0)

	suite.Len(model0.Identifiers, 2)
	modelIDs00, idsExist00 := model0.Identifiers["primary"]
	suite.True(idsExist00)
	suite.ElementsMatch(modelIDs00.Fields, []string{"ID"})

	modelIDs01, idsExist01 := model0.Identifiers["name"]
	suite.True(idsExist01)
	suite.ElementsMatch(modelIDs01.Fields, []string{"FirstName", "LastName"})

	suite.Len(model0.Related, 1)

	modelRelated00, relatedExists00 := model0.Related["ContactInfo"]
	suite.True(relatedExists00)
	suite.Equal(modelRelated00.Type, "HasOne")

	model1, modelExists1 := allModels["ContactInfo"]
	suite.True(modelExists1)
	suite.Equal(model1.Name, "ContactInfo")

	suite.Len(model1.Fields, 2)

	modelField10, fieldExists10 := model1.Fields["ID"]
	suite.True(fieldExists10)
	suite.Equal(modelField10.Type, yaml.ModelFieldTypeAutoIncrement)
	suite.Len(modelField10.Attributes, 1)
	suite.Equal(modelField10.Attributes[0], "mandatory")

	modelField11, fieldExists11 := model1.Fields["Email"]
	suite.True(fieldExists11)
	suite.Equal(modelField11.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField11.Attributes, 0)

	suite.Len(model1.Identifiers, 2)
	modelID10, idExists10 := model1.Identifiers["primary"]
	suite.True(idExists10)
	suite.ElementsMatch(modelID10.Fields, []string{"ID"})

	modelIDs11, idsExist11 := model1.Identifiers["email"]
	suite.True(idsExist11)
	suite.ElementsMatch(modelIDs11.Fields, []string{"Email"})

	suite.Len(model1.Related, 1)

	modelRelated10, relatedExists10 := model1.Related["Person"]
	suite.True(relatedExists10)
	suite.Equal(modelRelated10.Type, "ForOne")

	structure0, structureErr0 := r.GetStructure("Address")
	suite.Nil(structureErr0)
	suite.Equal(structure0.Name, "Address")

	suite.Len(structure0.Fields, 4)

	structureField00, fieldExists00 := structure0.Fields["Street"]
	suite.True(fieldExists00)
	suite.Equal(structureField00.Type, yaml.StructureFieldTypeString)

	structureField01, fieldExists01 := structure0.Fields["HouseNr"]
	suite.True(fieldExists01)
	suite.Equal(structureField01.Type, yaml.StructureFieldTypeString)

	structureField02, fieldExists02 := structure0.Fields["ZipCode"]
	suite.True(fieldExists02)
	suite.Equal(structureField02.Type, yaml.StructureFieldTypeString)

	structureField03, fieldExists03 := structure0.Fields["City"]
	suite.True(fieldExists03)
	suite.Equal(structureField03.Type, yaml.StructureFieldTypeString)

	entity0, entityErr0 := r.GetEntity("Person")
	suite.Nil(entityErr0)
	suite.Equal(entity0.Name, "Person")

	suite.Len(entity0.Fields, 5)

	entityField00, fieldExists00 := entity0.Fields["UUID"]
	suite.True(fieldExists00)
	suite.Equal(entityField00.Type, yaml.ModelFieldPath("Person.UUID"))
	suite.Len(entityField00.Attributes, 2)
	suite.Contains(entityField00.Attributes, "immutable")
	suite.Contains(entityField00.Attributes, "mandatory")

	entityField01, fieldExists01 := entity0.Fields["ID"]
	suite.True(fieldExists01)
	suite.Equal(entityField01.Type, yaml.ModelFieldPath("Person.ID"))
	suite.Len(entityField01.Attributes, 0)

	entityField02, fieldExists02 := entity0.Fields["FirstName"]
	suite.True(fieldExists02)
	suite.Equal(entityField02.Type, yaml.ModelFieldPath("Person.FirstName"))
	suite.Len(entityField02.Attributes, 0)

	entityField03, fieldExists03 := entity0.Fields["LastName"]
	suite.True(fieldExists03)
	suite.Equal(entityField03.Type, yaml.ModelFieldPath("Person.LastName"))
	suite.Len(entityField03.Attributes, 0)

	entityField04, fieldExists04 := entity0.Fields["Email"]
	suite.True(fieldExists04)
	suite.Equal(entityField04.Type, yaml.ModelFieldPath("Person.ContactInfo.Email"))
	suite.Len(entityField04.Attributes, 0)

	suite.Len(entity0.Identifiers, 1)
	entityID00, idExists00 := entity0.Identifiers["primary"]
	suite.True(idExists00)
	suite.ElementsMatch(entityID00.Fields, []string{"UUID"})
}

func (suite *LoadMorpheRegistryTestSuite) TestLoadMorpheRegistry_StartHook_Failure() {
	loadHooks := registry.LoadMorpheRegistryHooks{
		OnRegistryLoadStart: func(config cfg.MorpheLoadRegistryConfig) (cfg.MorpheLoadRegistryConfig, error) {
			return config, fmt.Errorf("compile model start hook error")
		},
	}
	config := cfg.MorpheLoadRegistryConfig{
		RegistryEnumsDirPath:    suite.EnumsDirPath,
		RegistryModelsDirPath:   suite.ModelsDirPath,
		RegistryEntitiesDirPath: suite.EntitiesDirPath,
	}

	r, registryErr := registry.LoadMorpheRegistry(loadHooks, config)

	suite.ErrorContains(registryErr, "compile model start hook error")
	suite.Nil(r)
}

func (suite *LoadMorpheRegistryTestSuite) TestLoadMorpheRegistry_SuccessHook_Successful() {
	loadHooks := registry.LoadMorpheRegistryHooks{
		OnRegistryLoadSuccess: func(r registry.Registry) (registry.Registry, error) {
			personModel, personErr := r.GetModel("Person")
			if personErr != nil {
				return r, personErr
			}
			personModel.Fields["NewField"] = yaml.ModelField{
				Type: yaml.ModelFieldTypeBoolean,
			}
			r.SetModel("Person", personModel)
			return r, nil
		},
	}
	config := cfg.MorpheLoadRegistryConfig{
		RegistryEnumsDirPath:      suite.EnumsDirPath,
		RegistryModelsDirPath:     suite.ModelsDirPath,
		RegistryStructuresDirPath: suite.StructuresDirPath,
		RegistryEntitiesDirPath:   suite.EntitiesDirPath,
	}

	r, registryErr := registry.LoadMorpheRegistry(loadHooks, config)

	suite.NoError(registryErr)
	suite.NotNil(r)

	allModels := r.GetAllModels()
	suite.Len(allModels, 2)

	model0, modelExists0 := allModels["Person"]
	suite.True(modelExists0)
	suite.Equal(model0.Name, "Person")

	suite.Len(model0.Fields, 5)

	modelField00, fieldExists00 := model0.Fields["ID"]
	suite.True(fieldExists00)
	suite.Equal(modelField00.Type, yaml.ModelFieldTypeAutoIncrement)
	suite.Len(modelField00.Attributes, 1)
	suite.Equal(modelField00.Attributes[0], "mandatory")

	modelField01, fieldExists01 := model0.Fields["FirstName"]
	suite.True(fieldExists01)
	suite.Equal(modelField01.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField01.Attributes, 0)

	modelField02, fieldExists02 := model0.Fields["LastName"]
	suite.True(fieldExists02)
	suite.Equal(modelField02.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField02.Attributes, 0)

	modelField03, fieldExists03 := model0.Fields["Nationality"]
	suite.True(fieldExists03)
	suite.Equal(modelField03.Type, yaml.ModelFieldType("Nationality"))
	suite.Len(modelField03.Attributes, 0)

	modelField04, fieldExists04 := model0.Fields["NewField"]
	suite.True(fieldExists04)
	suite.Equal(modelField04.Type, yaml.ModelFieldTypeBoolean)
	suite.Len(modelField04.Attributes, 0)

	suite.Len(model0.Identifiers, 2)
	modelIDs00, idsExist00 := model0.Identifiers["primary"]
	suite.True(idsExist00)
	suite.ElementsMatch(modelIDs00.Fields, []string{"ID"})

	modelIDs01, idsExist01 := model0.Identifiers["name"]
	suite.True(idsExist01)
	suite.ElementsMatch(modelIDs01.Fields, []string{"FirstName", "LastName"})

	suite.Len(model0.Related, 1)

	modelRelated00, relatedExists00 := model0.Related["ContactInfo"]
	suite.True(relatedExists00)
	suite.Equal(modelRelated00.Type, "HasOne")

	model1, modelExists1 := allModels["ContactInfo"]
	suite.True(modelExists1)
	suite.Equal(model1.Name, "ContactInfo")

	suite.Len(model1.Fields, 2)

	modelField10, fieldExists10 := model1.Fields["ID"]
	suite.True(fieldExists10)
	suite.Equal(modelField10.Type, yaml.ModelFieldTypeAutoIncrement)
	suite.Len(modelField10.Attributes, 1)
	suite.Equal(modelField10.Attributes[0], "mandatory")

	modelField11, fieldExists11 := model1.Fields["Email"]
	suite.True(fieldExists11)
	suite.Equal(modelField11.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField11.Attributes, 0)

	suite.Len(model1.Identifiers, 2)
	modelID10, idExists10 := model1.Identifiers["primary"]
	suite.True(idExists10)
	suite.ElementsMatch(modelID10.Fields, []string{"ID"})

	modelIDs11, idsExist11 := model1.Identifiers["email"]
	suite.True(idsExist11)
	suite.ElementsMatch(modelIDs11.Fields, []string{"Email"})

	suite.Len(model1.Related, 1)

	modelRelated10, relatedExists10 := model1.Related["Person"]
	suite.True(relatedExists10)
	suite.Equal(modelRelated10.Type, "ForOne")

	structure0, structureErr0 := r.GetStructure("Address")
	suite.Nil(structureErr0)
	suite.Equal(structure0.Name, "Address")

	suite.Len(structure0.Fields, 4)

	structureField00, fieldExists00 := structure0.Fields["Street"]
	suite.True(fieldExists00)
	suite.Equal(structureField00.Type, yaml.StructureFieldTypeString)

	structureField01, fieldExists01 := structure0.Fields["HouseNr"]
	suite.True(fieldExists01)
	suite.Equal(structureField01.Type, yaml.StructureFieldTypeString)

	structureField02, fieldExists02 := structure0.Fields["ZipCode"]
	suite.True(fieldExists02)
	suite.Equal(structureField02.Type, yaml.StructureFieldTypeString)

	structureField03, fieldExists03 := structure0.Fields["City"]
	suite.True(fieldExists03)
	suite.Equal(structureField03.Type, yaml.StructureFieldTypeString)
}

func (suite *LoadMorpheRegistryTestSuite) TestLoadMorpheRegistry_SuccessHook_Failure() {
	loadHooks := registry.LoadMorpheRegistryHooks{
		OnRegistryLoadSuccess: func(registry registry.Registry) (registry.Registry, error) {
			return registry, fmt.Errorf("compile model success hook error")
		},
	}
	config := cfg.MorpheLoadRegistryConfig{
		RegistryEnumsDirPath:      suite.EnumsDirPath,
		RegistryModelsDirPath:     suite.ModelsDirPath,
		RegistryStructuresDirPath: suite.StructuresDirPath,
		RegistryEntitiesDirPath:   suite.EntitiesDirPath,
	}

	r, registryErr := registry.LoadMorpheRegistry(loadHooks, config)

	suite.ErrorContains(registryErr, "compile model success hook error")
	suite.Nil(r)
}
