package appconfig

// Config contains configuration that needs to be applied to the application.
type Config struct {

	// Path to the layers directory.
	//
	LayersPath string `json:"layersPath"`

	// Path to the build directory.
	//
	BuildPath string `json:"buildPath"`

	// The outputted image format.
	//
	Format *Format `json:"format"`

	// Collection starts from index.
	// Example:
	// 	ETH - collection starts from 1
	// 	SOL - collection starts from 0
	//
	CollectionStartsFrom int `json:"collectionStartsFrom"`

	// Number of allowed dna duplicates.
	// Example:
	//	0 - zero tolerance for duplicates
	//  5 - in the created collection, 5 DNA duplicates are allowed
	UniqueDnaTollerance int `json:"uniqueDnaTollerance"`

	// Represents a metadata configuration.
	//
	MetadataConfiguration *MetadataConfiguration `json:"metadataConfiguration"`

	// This is a list of layer configurations.
	// Each configuration can be unique and have a different layer order,
	// use the same layers, or introduce new ones.
	//
	LayerConfigurations []*LayerConfiguration `json:"layerConfigurations"`
}

// GetCollectionSize - returns total collection size.
func (config *Config) GetCollectionSize() int {
	collectionSize := 0
	for _, layerConfiguration := range config.LayerConfigurations {
		collectionSize += layerConfiguration.Size
	}
	return collectionSize
}
