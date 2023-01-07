package artmachine

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"candy-machine/internal/appconfig"
	"candy-machine/internal/common"
	"candy-machine/internal/parallel"
)

type worker struct {
	config           *appconfig.Config
	layerConfig      *appconfig.LayerConfiguration
	layers           map[string][]*element
	shared           *shared
	startEditionFrom int
}

func newWorker(config *appconfig.Config, layerConfig *appconfig.LayerConfiguration, shared *shared, startEditionFrom int) (*worker, error) {

	// check configuration
	if layerConfig == nil {
		return nil, errors.New("configuration required")
	}

	// setup layers
	layers, err := setupLayers(config.LayersPath, layerConfig.LayersOrder)
	if err != nil {
		return nil, err
	}

	// create a new instance
	return &worker{
		config:           config,
		layerConfig:      layerConfig,
		layers:           layers,
		shared:           shared,
		startEditionFrom: startEditionFrom,
	}, nil
}

func (w *worker) run() error {
	const (
		maxParallelInstances = 10 // increate it if you need more parallel instances
	)
	run := parallel.NewRun(maxParallelInstances)
	for i := 0; i < w.layerConfig.Size; i++ {
		var (
			editionNumber = i + w.startEditionFrom + w.config.CollectionStartsFrom
			editionIndex  = editionNumber - 1
		)
		run.Do(func() error {
			// select randon elements
			selectedElements := make([]*element, 0)
			for _, elements := range w.layers {
				element, err := w.pickRandomLayerElement(elements)
				if err != nil {
					return err
				}
				selectedElements = append(selectedElements, element)
			}
			selectedElements = sortElementsByOrder(selectedElements)

			// save image
			if err := w.saveImage(editionNumber, selectedElements); err != nil {
				return err
			}

			// create metadata
			newMetadata := w.createMetadata(editionNumber, selectedElements)

			// save metadata
			if err := newMetadata.WriteToFile(fmt.Sprintf("%s/json/%d.json", w.config.BuildPath, editionNumber)); err != nil {
				return err
			}

			// append metadata to shared
			w.shared.addMetadata(editionIndex, newMetadata)

			return nil
		})
	}
	if err := run.Wait(); err != nil {
		return err
	}

	return nil
}

func (w *worker) pickRandomLayerElement(elements []*element) (*element, error) {

	// get pdf
	pdf, err := w.getProbabilityVector(elements)
	if err != nil {
		return nil, err
	}

	// get cdf
	len := len(elements)
	cdf := make([]float32, len)
	cdf[0] = pdf[0]
	for i := 1; i < len; i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}

	random := sample(cdf)
	if !(len > random) {
		return nil, fmt.Errorf(
			fmt.Sprintf("random generated trait index out of range, max size: %d, generated index: %d", len, random),
		)
	}

	return elements[random], nil
}

func (w *worker) getProbabilityVector(elements []*element) ([]float32, error) {
	var (
		len                   = len(elements)
		vector                = make([]float32, len)
		chanceOffset  float32 = 1.00
		commonCounter         = 0
		baseChance            = float32(100/len) / 100
	)

	for i, elem := range elements {
		if elem.Weight <= 1 {
			commonCounter++
			continue
		}
		chance := baseChance / elem.Weight
		vector[i] = chance
		chanceOffset -= chance
	}

	for i, p := range vector {
		if p == 0 {
			vector[i] = chanceOffset / float32(commonCounter)
		}
	}

	if err := checkProbabilityVector(vector); err != nil {
		return nil, err
	}
	return vector, nil
}

func (w *worker) saveImage(editionNumber int, seletedElements []*element) error {

	// create image's background
	bgImg := image.NewRGBA(image.Rect(0, 0, w.config.Format.Width, w.config.Format.Height))

	// set the background color (transparent)
	draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	// looping image higher position -> upper element
	for _, element := range seletedElements {

		imgBytes, err := ioutil.ReadFile(element.Path)
		if err != nil {
			return err
		}

		img, err := png.Decode(bytes.NewReader(imgBytes))
		if err != nil {
			return err
		}

		// set the image offset
		offset := image.Pt(0, 0)

		// draw the image
		draw.Draw(bgImg, img.Bounds().Add(offset), img, image.Point{}, draw.Over)
	}

	// encode image to buffer
	buff := new(bytes.Buffer)
	if err := png.Encode(buff, bgImg); err != nil {
		return err
	}

	// create image file
	imgFile, err := os.Create(fmt.Sprintf("%s/images/%d.png", w.config.BuildPath, editionNumber))
	if err != nil {
		return err
	}
	defer imgFile.Close()

	// write to file
	_, err = imgFile.Write(buff.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (w *worker) createMetadata(editionNumber int, seletedElements []*element) *common.Metadata {
	var (
		attributes = make([]*common.Attribute, len(seletedElements))
	)
	for i, element := range seletedElements {
		attributes[i] = &common.Attribute{
			TraitType: formatString(element.LayerName),
			Value:     formatString(element.Name),
		}
	}
	return &common.Metadata{
		Name:        fmt.Sprintf("%s #%d", w.config.MetadataConfiguration.NamePrefix, editionNumber),
		Description: w.config.MetadataConfiguration.Description,
		Image:       fmt.Sprintf("%s/%d.png", w.config.MetadataConfiguration.BaseUri, editionNumber),
		ExternalURL: w.config.MetadataConfiguration.ExternalURL,
		Dna:         calculateDna(seletedElements),
		Attributes:  attributes,
	}
}
