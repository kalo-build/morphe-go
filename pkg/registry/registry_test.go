package registry_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/morphe-go/internal/testutils"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
)

type RegistryTestSuite struct {
	suite.Suite

	TestDirPath string

	EnumsDirPath      string
	ModelsDirPath     string
	StructuresDirPath string
	EntitiesDirPath   string
}

func TestRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(RegistryTestSuite))
}

func (suite *RegistryTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()

	suite.EnumsDirPath = filepath.Join(suite.TestDirPath, "registry", "verbose", "enums")
	suite.ModelsDirPath = filepath.Join(suite.TestDirPath, "registry", "verbose", "models")
	suite.StructuresDirPath = filepath.Join(suite.TestDirPath, "registry", "verbose", "structures")
	suite.EntitiesDirPath = filepath.Join(suite.TestDirPath, "registry", "verbose", "entities")
}

func (suite *RegistryTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

// TestEmptyRegistryHasNoTypes verifies that an empty registry reports having no types
func (suite *RegistryTestSuite) TestEmptyRegistryHasNoTypes() {
	r := registry.NewRegistry()

	suite.False(r.HasEnums())
	suite.False(r.HasModels())
	suite.False(r.HasStructures())
	suite.False(r.HasEntities())
}

// TestRegistryWithEnumsReportsHavingEnums verifies that a registry with enums reports having enums
func (suite *RegistryTestSuite) TestRegistryWithEnumsReportsHavingEnums() {
	r := registry.NewRegistry()

	r.SetEnum("TestEnum", yaml.Enum{Name: "TestEnum"})

	suite.True(r.HasEnums())
	suite.False(r.HasModels())
	suite.False(r.HasStructures())
	suite.False(r.HasEntities())
}

// TestRegistryWithModelsReportsHavingModels verifies that a registry with models reports having models
func (suite *RegistryTestSuite) TestRegistryWithModelsReportsHavingModels() {
	r := registry.NewRegistry()

	r.SetModel("TestModel", yaml.Model{Name: "TestModel"})

	suite.False(r.HasEnums())
	suite.True(r.HasModels())
	suite.False(r.HasStructures())
	suite.False(r.HasEntities())
}

// TestRegistryWithStructuresReportsHavingStructures verifies that a registry with structures reports having structures
func (suite *RegistryTestSuite) TestRegistryWithStructuresReportsHavingStructures() {
	r := registry.NewRegistry()

	r.SetStructure("TestStructure", yaml.Structure{Name: "TestStructure"})

	suite.False(r.HasEnums())
	suite.False(r.HasModels())
	suite.True(r.HasStructures())
	suite.False(r.HasEntities())
}

// TestRegistryWithEntitiesReportsHavingEntities verifies that a registry with entities reports having entities
func (suite *RegistryTestSuite) TestRegistryWithEntitiesReportsHavingEntities() {
	r := registry.NewRegistry()

	// First add a model since entities require models
	r.SetModel("TestModel", yaml.Model{Name: "TestModel"})
	r.SetEntity("TestEntity", yaml.Entity{Name: "TestEntity"})

	suite.False(r.HasEnums())
	suite.True(r.HasModels())
	suite.False(r.HasStructures())
	suite.True(r.HasEntities())
}

// TestEmptyRegistryValidates verifies that an empty registry passes validation
func (suite *RegistryTestSuite) TestEmptyRegistryValidates() {
	r := registry.NewRegistry()

	suite.Nil(r.ValidateRegistry())
}

// TestRegistryWithOnlyModelsValidates verifies that a registry with only models passes validation
func (suite *RegistryTestSuite) TestRegistryWithOnlyModelsValidates() {
	r := registry.NewRegistry()

	r.SetModel("TestModel", yaml.Model{Name: "TestModel"})

	suite.Nil(r.ValidateRegistry())
}

// TestRegistryWithOnlyEnumsValidates verifies that a registry with only enums passes validation
func (suite *RegistryTestSuite) TestRegistryWithOnlyEnumsValidates() {
	r := registry.NewRegistry()

	r.SetEnum("TestEnum", yaml.Enum{Name: "TestEnum"})

	suite.Nil(r.ValidateRegistry())
}

// TestRegistryWithOnlyStructuresValidates verifies that a registry with only structures passes validation
func (suite *RegistryTestSuite) TestRegistryWithOnlyStructuresValidates() {
	r := registry.NewRegistry()

	r.SetStructure("TestStructure", yaml.Structure{Name: "TestStructure"})

	suite.Nil(r.ValidateRegistry())
}

// TestRegistryWithEntitiesButNoModelsFailsValidation verifies that a registry with entities but no models fails validation
func (suite *RegistryTestSuite) TestRegistryWithEntitiesButNoModelsFailsValidation() {
	r := registry.NewRegistry()

	r.SetEntity("TestEntity", yaml.Entity{Name: "TestEntity"})

	validationErr := r.ValidateRegistry()
	suite.NotNil(validationErr)
	suite.Contains(validationErr.Error(), "invalid registry state: entities exist but no models are defined")
}

// TestRegistryWithEntitiesAndModelsValidates verifies that a registry with both entities and models passes validation
func (suite *RegistryTestSuite) TestRegistryWithEntitiesAndModelsValidates() {
	r := registry.NewRegistry()

	r.SetModel("TestModel", yaml.Model{Name: "TestModel"})
	r.SetEntity("TestEntity", yaml.Entity{Name: "TestEntity"})

	suite.Nil(r.ValidateRegistry())
}

// TestGetAllMethodsReturnEmptyMapsForNilMaps verifies that GetAll methods return empty maps for nil maps
func (suite *RegistryTestSuite) TestGetAllMethodsReturnEmptyMapsForNilMaps() {
	r := registry.NewRegistry()

	suite.Empty(r.GetAllEnums())
	suite.Empty(r.GetAllModels())
	suite.Empty(r.GetAllStructures())
	suite.Empty(r.GetAllEntities())
}

// TestGetMethodsReturnErrorForNonExistentItems verifies that Get methods return error for non-existent items
func (suite *RegistryTestSuite) TestGetMethodsReturnErrorForNonExistentItems() {
	r := registry.NewRegistry()

	_, enumErr := r.GetEnum("NonExistentEnum")
	suite.NotNil(enumErr)

	_, modelErr := r.GetModel("NonExistentModel")
	suite.NotNil(modelErr)

	_, structureErr := r.GetStructure("NonExistentStructure")
	suite.NotNil(structureErr)

	_, entityErr := r.GetEntity("NonExistentEntity")
	suite.NotNil(entityErr)
}

// TestSetMethodsCreateMapsIfNil verifies that Set methods create maps if they are nil
func (suite *RegistryTestSuite) TestSetMethodsCreateMapsIfNil() {
	r := registry.NewRegistry()

	r.SetEnum("TestEnum", yaml.Enum{Name: "TestEnum"})
	r.SetModel("TestModel", yaml.Model{Name: "TestModel"})
	r.SetStructure("TestStructure", yaml.Structure{Name: "TestStructure"})
	r.SetEntity("TestEntity", yaml.Entity{Name: "TestEntity"})

	// Verify the items were added
	enum, enumErr := r.GetEnum("TestEnum")
	suite.Nil(enumErr)
	suite.Equal("TestEnum", enum.Name)

	model, modelErr := r.GetModel("TestModel")
	suite.Nil(modelErr)
	suite.Equal("TestModel", model.Name)

	structure, structureErr := r.GetStructure("TestStructure")
	suite.Nil(structureErr)
	suite.Equal("TestStructure", structure.Name)

	entity, entityErr := r.GetEntity("TestEntity")
	suite.Nil(entityErr)
	suite.Equal("TestEntity", entity.Name)
}

// TestLoadingEntitiesWithoutModelsReturnsError verifies that loading entities without models returns an error
func (suite *RegistryTestSuite) TestLoadingEntitiesWithoutModelsReturnsError() {
	r := registry.NewRegistry()

	entitiesErr := r.LoadEntitiesFromDirectory(suite.EntitiesDirPath)

	suite.NotNil(entitiesErr)
	suite.Contains(entitiesErr.Error(), "attempted to load entities but no models are defined in registry")
}

// TestLoadingModelsFirstAllowsLoadingEntities verifies that loading models first allows entities to be loaded
func (suite *RegistryTestSuite) TestLoadingModelsFirstAllowsLoadingEntities() {
	r := registry.NewRegistry()

	modelsErr := r.LoadModelsFromDirectory(suite.ModelsDirPath)
	suite.Nil(modelsErr)

	entitiesErr := r.LoadEntitiesFromDirectory(suite.EntitiesDirPath)
	suite.Nil(entitiesErr)

	suite.True(r.HasModels())
	suite.True(r.HasEntities())
}

// TestDeepCloneHandlesEmptyRegistry verifies that DeepClone properly handles an empty registry
func (suite *RegistryTestSuite) TestDeepCloneHandlesEmptyRegistry() {
	r := registry.NewRegistry()

	clone := r.DeepClone()

	suite.NotSame(r, clone)
	suite.False(clone.HasEnums())
	suite.False(clone.HasModels())
	suite.False(clone.HasStructures())
	suite.False(clone.HasEntities())
}

// TestDeepCloneHandlesPartiallyPopulatedRegistry verifies that DeepClone properly handles a partially populated registry
func (suite *RegistryTestSuite) TestDeepCloneHandlesPartiallyPopulatedRegistry() {
	r := registry.NewRegistry()

	r.SetEnum("TestEnum", yaml.Enum{Name: "TestEnum"})
	r.SetModel("TestModel", yaml.Model{Name: "TestModel"})

	clone := r.DeepClone()

	suite.NotSame(r, clone)
	suite.True(clone.HasEnums())
	suite.True(clone.HasModels())
	suite.False(clone.HasStructures())
	suite.False(clone.HasEntities())

	// Verify enum was properly cloned
	enum, err := clone.GetEnum("TestEnum")
	suite.Nil(err)
	suite.Equal("TestEnum", enum.Name)

	// Verify model was properly cloned
	model, err := clone.GetModel("TestModel")
	suite.Nil(err)
	suite.Equal("TestModel", model.Name)
}

func (suite *RegistryTestSuite) TestLoadEnumsFromDirectory() {
	registry := registry.NewRegistry()

	enumsErr := registry.LoadEnumsFromDirectory(suite.EnumsDirPath)

	suite.Nil(enumsErr)
	suite.Len(registry.GetAllEnums(), 2)
	suite.Len(registry.GetAllModels(), 0)
	suite.Len(registry.GetAllEntities(), 0)

	enum0, enumErr0 := registry.GetEnum("Country")
	suite.Nil(enumErr0)
	suite.Equal(enum0.Name, "Country")
	suite.Equal(enum0.Type, yaml.EnumTypeString)

	suite.Len(enum0.Entries, 3)

	entry00, entryExists00 := enum0.Entries["US"]
	suite.True(entryExists00)
	suite.Equal(entry00, "United States")

	entry01, entryExists01 := enum0.Entries["DE"]
	suite.True(entryExists01)
	suite.Equal(entry01, "Germany")

	entry02, entryExists02 := enum0.Entries["FR"]
	suite.True(entryExists02)
	suite.Equal(entry02, "France")

	enum1, enumErr1 := registry.GetEnum("Nationality")
	suite.Nil(enumErr1)
	suite.Equal(enum1.Name, "Nationality")
	suite.Equal(enum1.Type, yaml.EnumTypeString)

	suite.Len(enum1.Entries, 3)

	entry10, entryExists10 := enum1.Entries["US"]
	suite.True(entryExists10)
	suite.Equal(entry10, "American")

	entry11, entryExists11 := enum1.Entries["DE"]
	suite.True(entryExists11)
	suite.Equal(entry11, "German")

	entry12, entryExists12 := enum1.Entries["FR"]
	suite.True(entryExists12)
	suite.Equal(entry12, "French")
}

func (suite *RegistryTestSuite) TestLoadEnumsFromDirectory_InvalidDirPath() {
	registry := registry.NewRegistry()

	enumsErr := registry.LoadEnumsFromDirectory("####INVALID/DIR/PATH####")

	suite.Nil(enumsErr)
	suite.Len(registry.GetAllEnums(), 0)
	suite.Len(registry.GetAllModels(), 0)
	suite.Len(registry.GetAllEntities(), 0)
}

func (suite *RegistryTestSuite) TestLoadEnumsFromDirectory_ConflictingName() {
	registry := registry.NewRegistry()

	registry.SetEnum("Country", yaml.Enum{Name: "Country"})

	enumsErr := registry.LoadEnumsFromDirectory(suite.EnumsDirPath)

	suite.NotNil(enumsErr)
	enumsErrMsg := enumsErr.Error()
	suite.Contains(enumsErrMsg, "enum name 'Country' already exists in registry")

	conflictPath := filepath.Join(suite.EnumsDirPath, "country.enum")
	suite.Contains(enumsErrMsg, conflictPath)
}

func (suite *RegistryTestSuite) TestLoadModelsFromDirectory() {
	registry := registry.NewRegistry()

	modelsErr := registry.LoadModelsFromDirectory(suite.ModelsDirPath)

	suite.Nil(modelsErr)
	suite.Len(registry.GetAllModels(), 4)
	suite.Len(registry.GetAllEntities(), 0)

	model0, modelErr0 := registry.GetModel("Company")
	suite.Nil(modelErr0)
	suite.Equal(model0.Name, "Company")

	suite.Len(model0.Fields, 4)

	modelField00, fieldExists00 := model0.Fields["UUID"]
	suite.True(fieldExists00)
	suite.Equal(modelField00.Type, yaml.ModelFieldTypeUUID)
	suite.Len(modelField00.Attributes, 2)
	suite.Equal(modelField00.Attributes[0], "immutable")
	suite.Equal(modelField00.Attributes[1], "mandatory")

	modelField01, fieldExists01 := model0.Fields["ID"]
	suite.True(fieldExists01)
	suite.Equal(modelField01.Type, yaml.ModelFieldTypeAutoIncrement)
	suite.Len(modelField01.Attributes, 1)
	suite.Equal(modelField01.Attributes[0], "mandatory")

	modelField02, fieldExists02 := model0.Fields["FoundedAt"]
	suite.True(fieldExists02)
	suite.Equal(modelField02.Type, yaml.ModelFieldTypeTime)
	suite.Len(modelField02.Attributes, 0)

	modelField03, fieldExists03 := model0.Fields["Name"]
	suite.True(fieldExists03)
	suite.Equal(modelField03.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField03.Attributes, 0)

	suite.Len(model0.Identifiers, 2)

	modelIDs00, idsExist00 := model0.Identifiers["primary"]
	suite.True(idsExist00)
	suite.ElementsMatch(modelIDs00.Fields, []string{"ID"})

	modelIDs01, idsExist01 := model0.Identifiers["entity"]
	suite.True(idsExist01)
	suite.ElementsMatch(modelIDs01.Fields, []string{"UUID"})

	suite.Len(model0.Related, 2)

	modelRelated00, relatedExists00 := model0.Related["Address"]
	suite.True(relatedExists00)
	suite.Equal(modelRelated00.Type, "HasOne")

	modelRelated01, relatedExists01 := model0.Related["Person"]
	suite.True(relatedExists01)
	suite.Equal(modelRelated01.Type, "HasMany")

	model1, modelErr1 := registry.GetModel("Address")
	suite.Nil(modelErr1)
	suite.Equal(model1.Name, "Address")

	suite.Len(model1.Fields, 3)

	modelField10, fieldExists10 := model1.Fields["ID"]
	suite.True(fieldExists10)
	suite.Equal(modelField10.Type, yaml.ModelFieldTypeAutoIncrement)
	suite.Len(modelField10.Attributes, 1)
	suite.Equal(modelField10.Attributes[0], "mandatory")

	modelField11, fieldExists11 := model1.Fields["Street"]
	suite.True(fieldExists11)
	suite.Equal(modelField11.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField11.Attributes, 0)

	modelField12, fieldExists12 := model1.Fields["HouseNumber"]
	suite.True(fieldExists12)
	suite.Equal(modelField12.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField12.Attributes, 0)

	suite.Len(model1.Identifiers, 2)
	modelIDs10, idsExist10 := model1.Identifiers["primary"]
	suite.True(idsExist10)
	suite.ElementsMatch(modelIDs10.Fields, []string{"ID"})

	modelIDs11, idsExist11 := model1.Identifiers["street"]
	suite.True(idsExist11)
	suite.ElementsMatch(modelIDs11.Fields, []string{"Street", "HouseNumber"})

	suite.Len(model1.Related, 1)

	modelRelated10, relatedExists10 := model1.Related["Company"]
	suite.True(relatedExists10)
	suite.Equal(modelRelated10.Type, "ForOne")

	model2, modelErr2 := registry.GetModel("ContactInfo")
	suite.Nil(modelErr2)
	suite.Equal(model2.Name, "ContactInfo")

	suite.Len(model2.Fields, 3)

	modelField20, fieldExists20 := model2.Fields["ID"]
	suite.True(fieldExists20)
	suite.Equal(modelField20.Type, yaml.ModelFieldTypeAutoIncrement)
	suite.Len(modelField20.Attributes, 1)
	suite.Equal(modelField20.Attributes[0], "mandatory")

	modelField21, fieldExists21 := model2.Fields["Email"]
	suite.True(fieldExists21)
	suite.Equal(modelField21.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField21.Attributes, 1)
	suite.Equal(modelField21.Attributes[0], "mandatory")

	modelField22, fieldExists22 := model2.Fields["PhoneNumber"]
	suite.True(fieldExists22)
	suite.Equal(modelField22.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField22.Attributes, 0)

	suite.Len(model2.Identifiers, 1)
	modelID20, idExists20 := model2.Identifiers["primary"]
	suite.True(idExists20)
	suite.ElementsMatch(modelID20.Fields, []string{"ID"})

	suite.Len(model2.Related, 1)

	modelRelated20, relatedExists20 := model2.Related["Person"]
	suite.True(relatedExists20)
	suite.Equal(modelRelated20.Type, "ForOne")

	model3, modelErr3 := registry.GetModel("Person")
	suite.Nil(modelErr3)
	suite.Equal(model3.Name, "Person")

	suite.Len(model3.Fields, 4)

	modelField30, fieldExists30 := model3.Fields["UUID"]
	suite.True(fieldExists30)
	suite.Equal(modelField30.Type, yaml.ModelFieldTypeUUID)
	suite.Len(modelField30.Attributes, 2)
	suite.Equal(modelField30.Attributes[0], "immutable")
	suite.Equal(modelField30.Attributes[1], "mandatory")

	modelField31, fieldExists31 := model3.Fields["ID"]
	suite.True(fieldExists31)
	suite.Equal(modelField31.Type, yaml.ModelFieldTypeAutoIncrement)
	suite.Len(modelField31.Attributes, 1)
	suite.Equal(modelField31.Attributes[0], "mandatory")

	modelField32, fieldExists32 := model3.Fields["FirstName"]
	suite.True(fieldExists32)
	suite.Equal(modelField32.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField32.Attributes, 0)

	modelField33, fieldExists33 := model3.Fields["LastName"]
	suite.True(fieldExists33)
	suite.Equal(modelField33.Type, yaml.ModelFieldTypeString)
	suite.Len(modelField33.Attributes, 0)

	suite.Len(model3.Identifiers, 3)
	modelIDs30, idsExist10 := model3.Identifiers["primary"]
	suite.True(idsExist10)
	suite.ElementsMatch(modelIDs30.Fields, []string{"ID"})

	modelIDs31, idsExist11 := model3.Identifiers["entity"]
	suite.True(idsExist11)
	suite.ElementsMatch(modelIDs31.Fields, []string{"UUID"})

	modelIDs32, idsExist12 := model3.Identifiers["name"]
	suite.True(idsExist12)
	suite.ElementsMatch(modelIDs32.Fields, []string{"FirstName", "LastName"})

	suite.Len(model3.Related, 2)

	modelRelated30, relatedExists30 := model3.Related["Company"]
	suite.True(relatedExists30)
	suite.Equal(modelRelated30.Type, "ForOne")

	modelRelated31, relatedExists31 := model3.Related["ContactInfo"]
	suite.True(relatedExists31)
	suite.Equal(modelRelated31.Type, "HasOne")
}

func (suite *RegistryTestSuite) TestLoadModelsFromDirectory_InvalidDirPath() {
	registry := registry.NewRegistry()

	modelsErr := registry.LoadModelsFromDirectory("####INVALID/DIR/PATH####")

	suite.Nil(modelsErr)
	suite.Len(registry.GetAllModels(), 0)
	suite.Len(registry.GetAllEntities(), 0)
}

func (suite *RegistryTestSuite) TestLoadModelsFromDirectory_ConflictingName() {
	registry := registry.NewRegistry()

	registry.SetModel("Company", yaml.Model{Name: "Company"})

	modelsErr := registry.LoadModelsFromDirectory(suite.ModelsDirPath)

	suite.NotNil(modelsErr)
	modelsErrMsg := modelsErr.Error()
	suite.Contains(modelsErrMsg, "model name 'Company' already exists in registry")

	conflictPath := filepath.Join(suite.ModelsDirPath, "company.mod")
	suite.Contains(modelsErrMsg, conflictPath)
}

func (suite *RegistryTestSuite) TestLoadEntitiesFromDirectory() {
	registry := registry.NewRegistry()

	// Load models first to avoid validation error
	modelsErr := registry.LoadModelsFromDirectory(suite.ModelsDirPath)
	suite.Nil(modelsErr)

	entitiesErr := registry.LoadEntitiesFromDirectory(suite.EntitiesDirPath)

	suite.Nil(entitiesErr)
	suite.Len(registry.GetAllModels(), 4)
	suite.Len(registry.GetAllEntities(), 2)

	entity0, entityErr0 := registry.GetEntity("Company")
	suite.Nil(entityErr0)
	suite.Equal(entity0.Name, "Company")

	suite.Len(entity0.Fields, 6)

	entityField00, fieldExists00 := entity0.Fields["UUID"]
	suite.True(fieldExists00)
	suite.Equal(entityField00.Type, yaml.ModelFieldPath("Company.UUID"))
	suite.Len(entityField00.Attributes, 2)
	suite.Equal(entityField00.Attributes[0], "immutable")
	suite.Equal(entityField00.Attributes[1], "mandatory")

	entityField01, fieldExists01 := entity0.Fields["ID"]
	suite.True(fieldExists01)
	suite.Equal(entityField01.Type, yaml.ModelFieldPath("Company.ID"))
	suite.Len(entityField01.Attributes, 0)

	entityField02, fieldExists02 := entity0.Fields["FoundedAt"]
	suite.True(fieldExists02)
	suite.Equal(entityField02.Type, yaml.ModelFieldPath("Company.FoundedAt"))
	suite.Len(entityField02.Attributes, 0)

	entityField03, fieldExists03 := entity0.Fields["Name"]
	suite.True(fieldExists03)
	suite.Equal(entityField03.Type, yaml.ModelFieldPath("Company.Name"))
	suite.Len(entityField03.Attributes, 0)

	entityField04, fieldExists04 := entity0.Fields["ZipCode"]
	suite.True(fieldExists04)
	suite.Equal(entityField04.Type, yaml.ModelFieldPath("Company.Address.ZipCode"))
	suite.Len(entityField04.Attributes, 0)

	entityField05, fieldExists05 := entity0.Fields["City"]
	suite.True(fieldExists05)
	suite.Equal(entityField05.Type, yaml.ModelFieldPath("Company.Address.City"))
	suite.Len(entityField05.Attributes, 0)

	suite.Len(entity0.Related, 1)

	entityRelated00, relatedExists00 := entity0.Related["Person"]
	suite.True(relatedExists00)
	suite.Equal(entityRelated00.Type, "HasMany")

	entity1, entityErr1 := registry.GetEntity("Person")
	suite.Nil(entityErr1)
	suite.Equal(entity1.Name, "Person")

	suite.Len(entity1.Fields, 5)

	entityField10, fieldExists10 := entity1.Fields["UUID"]
	suite.True(fieldExists10)
	suite.Equal(entityField10.Type, yaml.ModelFieldPath("Person.UUID"))
	suite.Len(entityField10.Attributes, 2)
	suite.Equal(entityField10.Attributes[0], "immutable")
	suite.Equal(entityField10.Attributes[1], "mandatory")

	entityField11, fieldExists11 := entity1.Fields["ID"]
	suite.True(fieldExists11)
	suite.Equal(entityField11.Type, yaml.ModelFieldPath("Person.ID"))
	suite.Len(entityField11.Attributes, 0)

	entityField12, fieldExists12 := entity1.Fields["FirstName"]
	suite.True(fieldExists12)
	suite.Equal(entityField12.Type, yaml.ModelFieldPath("Person.FirstName"))
	suite.Len(entityField12.Attributes, 0)

	entityField13, fieldExists13 := entity1.Fields["LastName"]
	suite.True(fieldExists13)
	suite.Equal(entityField13.Type, yaml.ModelFieldPath("Person.LastName"))
	suite.Len(entityField13.Attributes, 0)

	entityField14, fieldExists14 := entity1.Fields["Email"]
	suite.True(fieldExists14)
	suite.Equal(entityField14.Type, yaml.ModelFieldPath("Person.ContactInfo.Email"))
	suite.Len(entityField14.Attributes, 0)

	suite.Len(entity1.Related, 1)

	entityRelated10, relatedExists10 := entity1.Related["Company"]
	suite.True(relatedExists10)
	suite.Equal(entityRelated10.Type, "ForOne")
}

func (suite *RegistryTestSuite) TestLoadEntitiesFromDirectory_InvalidDirPath() {
	registry := registry.NewRegistry()

	entitiesErr := registry.LoadEntitiesFromDirectory("####INVALID/DIR/PATH####")

	suite.Nil(entitiesErr)
	suite.Len(registry.GetAllModels(), 0)
	suite.Len(registry.GetAllEntities(), 0)
}

func (suite *RegistryTestSuite) TestLoadEntitiesFromDirectory_ConflictingName() {
	registry := registry.NewRegistry()

	// Load models first to avoid validation error
	modelsErr := registry.LoadModelsFromDirectory(suite.ModelsDirPath)
	suite.Nil(modelsErr)

	registry.SetEntity("Company", yaml.Entity{Name: "Company"})

	entitiesErr := registry.LoadEntitiesFromDirectory(suite.EntitiesDirPath)

	suite.NotNil(entitiesErr)
	entitiesErrMsg := entitiesErr.Error()
	suite.Contains(entitiesErrMsg, "entity name 'Company' already exists in registry")

	conflictPath := filepath.Join(suite.EntitiesDirPath, "company.ent")
	suite.Contains(entitiesErrMsg, conflictPath)
}

func (suite *RegistryTestSuite) TestDeepClone() {
	registry := registry.NewRegistry()

	enumsErr := registry.LoadEnumsFromDirectory(suite.EnumsDirPath)
	suite.Nil(enumsErr)

	modelsErr := registry.LoadModelsFromDirectory(suite.ModelsDirPath)
	suite.Nil(modelsErr)

	entitiesErr := registry.LoadEntitiesFromDirectory(suite.EntitiesDirPath)
	suite.Nil(entitiesErr)

	registryClone := registry.DeepClone()

	suite.NotSame(registry, registryClone)
	suite.Equal(registry, registryClone)
}

func (suite *RegistryTestSuite) TestDeepClone_Empty() {
	registry := registry.NewRegistry()

	registryClone := registry.DeepClone()

	suite.NotSame(registry, registryClone)
	suite.Equal(registry, registryClone)
}

func (suite *RegistryTestSuite) TestLoadStructuresFromDirectory() {
	registry := registry.NewRegistry()

	structuresErr := registry.LoadStructuresFromDirectory(suite.StructuresDirPath)

	suite.Nil(structuresErr)
	suite.Len(registry.GetAllStructures(), 1)
	suite.Len(registry.GetAllEnums(), 0)
	suite.Len(registry.GetAllModels(), 0)
	suite.Len(registry.GetAllEntities(), 0)

	structure0, structureErr0 := registry.GetStructure("Address")
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

func (suite *RegistryTestSuite) TestLoadStructuresFromDirectory_InvalidDirPath() {
	registry := registry.NewRegistry()

	structuresErr := registry.LoadStructuresFromDirectory("####INVALID/DIR/PATH####")

	suite.Nil(structuresErr)
	suite.Len(registry.GetAllEnums(), 0)
	suite.Len(registry.GetAllModels(), 0)
	suite.Len(registry.GetAllStructures(), 0)
	suite.Len(registry.GetAllEntities(), 0)
}

func (suite *RegistryTestSuite) TestLoadStructuresFromDirectory_ConflictingName() {
	registry := registry.NewRegistry()

	registry.SetStructure("Address", yaml.Structure{Name: "Address"})

	structuresErr := registry.LoadStructuresFromDirectory(suite.StructuresDirPath)

	suite.NotNil(structuresErr)
	structuresErrMsg := structuresErr.Error()
	suite.Contains(structuresErrMsg, "structure name 'Address' already exists in registry")

	conflictPath := filepath.Join(suite.StructuresDirPath, "address.str")
	suite.Contains(structuresErrMsg, conflictPath)
}
