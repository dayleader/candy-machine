package appconfig

// LayerConfiguration contains a specific layer configuration.
type LayerConfiguration struct {

	// The size of output variations.
	//
	Size int `json:"size"`

	// The order of the layers.
	//
	LayersOrder []*LayerOrder `json:"layersOrder"`
}

// LayerOrder contains a layer information.
type LayerOrder struct {

	// Represents the name of the folder (in /layers/) that the images reside in.
	//
	Name string `json:"name"`

	// Represents the layer options.
	//
	Options *LayerOptions `json:"options"`
}

// LayerOptions contains a layer options.
type LayerOptions struct {

	// Represents the layer display name.
	//
	DisplayName string `json:"displayName"`

	// Represents the rarity options.
	//
	RarityOptions []*RarityOptions `json:"rarity"`
}

// RarityOptions contains a rarity options.
type RarityOptions struct {

	// Rarity weight.
	// The higher the number, the rarer the layers occur.
	// The basic weight is 1 (the weight cannot be less than this value).
	//
	// For example:
	//		weight = 2, means that layers will be used 2 times less often than the base ones.
	// 		weight = 3, means that layers will be used 3 times less often than the base ones.
	//		weight = 4, means that the layers will be used 4 times less than the base ones.
	//		etc ...
	//
	Weight float32 `json:"weight"`

	// Array of filenames (order doesn't matter).
	//
	FileNames []string `json:"filenames"`
}
