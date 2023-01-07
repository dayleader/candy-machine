package appconfig

// MetadataConfiguration contains metadata configuration.
type MetadataConfiguration struct {

	// Represents the unique item name prefix
	//
	NamePrefix string `json:"namePrefix"`

	// A human readable description of the collection. Markdown is supported
	//
	Description string `json:"description"`

	// The base image URL
	//
	BaseUri string `json:"baseUri"`

	// This is the URL that will appear below the asset's image
	//
	ExternalURL string `json:"externalUrl"`
}
