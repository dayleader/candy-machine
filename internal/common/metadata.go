package common

import (
	"encoding/json"
	"os"
)

// Attribute represents the metadata attribute
type Attribute struct {

	// Here trait_type is the name of the trait
	//
	TraitType string `json:"trait_type"`

	// Value is the value of the trait
	//
	Value string `json:"value"`
}

// Metadata is structured according to the official ERC721 metadata standard
type Metadata struct {

	// Represents the unique name of the item
	// Example: Your Collection #1
	//
	Name string `json:"name"`

	// A human readable description of the item. Markdown is supported
	//
	Description string `json:"description"`

	// This is the URL to the image of the item
	// Can be just about any type of image and can be IPFS URLs or paths
	//
	Image string `json:"image"`

	// This is the URL that will appear below the asset's image
	//
	ExternalURL string `json:"external_url"`

	// Unique dna of the item
	//
	Dna string `json:"dna"`

	// These are the attributes for the item
	//
	Attributes []*Attribute `json:"attributes"`
}

// WriteToFile writes metadata to file
func (md *Metadata) WriteToFile(filePath string) error {

	// create metadata file
	metadataFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer metadataFile.Close()

	// marshal to json
	metadataBytes, err := json.Marshal(md)
	if err != nil {
		return err
	}

	// write to file
	_, err = metadataFile.Write(metadataBytes)
	if err != nil {
		return err
	}

	return nil
}

// MetadataList represents a list of metadata
type MetadataList []*Metadata

// WriteToFile writes a list of metadata to file
func (list *MetadataList) WriteToFile(filePath string) error {

	// create metadata file
	metadataFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer metadataFile.Close()

	// marshal to json
	metadataBytes, err := json.Marshal(list)
	if err != nil {
		return err
	}

	// write to file
	_, err = metadataFile.Write(metadataBytes)
	if err != nil {
		return err
	}

	return nil
}
