package artmachine

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"time"

	"candy-machine/internal/appconfig"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// represents a single layer element
type element struct {

	// Represents the layer order
	// Higher order -> upper element
	//
	Order int `json:"order"`

	// Represents the layer name or display name if specified
	//
	LayerName string `json:"layer_name"`

	// Represents the item name
	//
	Name string `json:"name"`

	// Represents the element filename (with extension)
	//
	Filename string `json:"filename"`

	// Represents the path of the element (in /layers/) that the images reside in
	//
	Path string `json:"path"`

	// Represents the element weight
	// The higher the number, the rarer the element occur
	// The basic weight is 1 (the weight cannot be less than this value)
	//
	// For example:
	//		weight = 2, means that element will be used 2 times less often than the base ones
	// 		weight = 3, means that element will be used 3 times less often than the base ones
	//		weight = 4, means that the element will be used 4 times less than the base ones
	//		etc ...
	//
	Weight float32 `json:"weight"`
}

func sample(cdf []float32) int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Float32()
	bucket := 0
	for r > cdf[bucket] {
		bucket++
	}
	return bucket
}

func checkProbabilityVector(vector []float32) error {
	var (
		sum         float32 = 0
		controllSum float32 = 1
	)
	for _, p := range vector {
		sum += p
	}
	if !(sum >= controllSum-0.1 || sum >= controllSum+0.1) {
		return fmt.Errorf("expected probability vector controll sum %v but got %v", controllSum, sum)
	}
	return nil
}

func sortElementsByOrder(elements []*element) []*element {
	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Order < elements[j].Order
	})
	return elements
}

func getWeight(filename string, layerOrder *appconfig.LayerOrder) float32 {
	if layerOrder.Options != nil && layerOrder.Options.RarityOptions != nil {
		for _, opt := range layerOrder.Options.RarityOptions {
			for _, f := range opt.FileNames {
				if f == filename {
					return opt.Weight
				}
			}
		}
	}
	return 1
}

func getLayerName(layerOrder *appconfig.LayerOrder) string {
	if layerOrder.Options != nil && len(layerOrder.Options.DisplayName) > 0 {
		return layerOrder.Options.DisplayName
	}
	return layerOrder.Name
}

func setupLayers(layersBasePath string, layerOrders []*appconfig.LayerOrder) (map[string][]*element, error) {
	var (
		layers = make(map[string][]*element)
		err    error
	)
	for order, layerOrder := range layerOrders {
		var (
			layerPath = fmt.Sprintf("%s/%s", layersBasePath, layerOrder.Name)
		)
		if _, err := os.Stat(layerPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("layer folder: %s does not exist, path: %s", layerOrder.Name, layerPath)
		}
		elements, err := getElements(layerPath, order, layerOrder)
		if err != nil {
			return nil, err
		}
		layers[layerOrder.Name] = elements
	}
	return layers, err
}

func getElements(layerPath string, order int, layerOrder *appconfig.LayerOrder) ([]*element, error) {
	var (
		elements = make([]*element, 0)
	)
	err := filepath.Walk(layerPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".png" {
			return fmt.Errorf("bad file extension: %s", filepath.Ext(info.Name()))
		}

		elements = append(elements, &element{
			Order:     order,
			LayerName: getLayerName(layerOrder),
			Name:      info.Name()[:len(info.Name())-4],
			Filename:  info.Name(),
			Path:      path,
			Weight:    getWeight(info.Name(), layerOrder),
		})
		return nil
	})
	return elements, err
}

func formatString(str string) string {
	return cases.Title(language.Und, cases.NoLower).String(str)
}

func calculateDna(seletedElements []*element) string {
	var buffer bytes.Buffer
	for _, element := range seletedElements {
		buffer.WriteString(fmt.Sprintf("%s/%s", element.LayerName, element.Name))
	}
	hasher := sha1.New()
	hasher.Write(buffer.Bytes())
	return hex.EncodeToString(hasher.Sum(nil))
}
