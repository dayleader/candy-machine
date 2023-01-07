package artmachine

import (
	"candy-machine/internal/appconfig"
	"fmt"
	"os"
)

// ArtMachine is responsible for creating artwork images based on the provided layers
type ArtMachine struct {
	config  *appconfig.Config
	shared  *shared
	workers []*worker
}

// New creates a new instance of ArtMachine
func New(config *appconfig.Config) (*ArtMachine, error) {
	var (
		workers = make([]*worker, len(config.LayerConfigurations))
		shared  = newShared(config.GetCollectionSize(), config.UniqueDnaTollerance)
	)
	startEditionFrom := 0
	for i, layerConfiguration := range config.LayerConfigurations {
		worker, err := newWorker(config, layerConfiguration, shared, startEditionFrom)
		if err != nil {
			return nil, err
		}
		workers[i] = worker
		startEditionFrom = layerConfiguration.Size
	}
	return &ArtMachine{
		config:  config,
		shared:  shared,
		workers: workers,
	}, nil
}

// Run - start creating
func (machine *ArtMachine) Run() error {
	if err := machine.buildSetup(); err != nil {
		return err
	}
	for _, worker := range machine.workers {
		if err := worker.run(); err != nil {
			return err
		}
	}
	return machine.shared.metadata.WriteToFile(
		fmt.Sprintf("%s/json/_metadata.json", machine.config.BuildPath),
	)
}

func (machine *ArtMachine) buildSetup() error {
	if err := os.RemoveAll(machine.config.BuildPath); err != nil {
		return err
	}
	if err := os.Mkdir(machine.config.BuildPath, 0700); err != nil {
		return err
	}
	if err := os.Mkdir(machine.config.BuildPath+"/images", 0700); err != nil {
		return err
	}
	if err := os.Mkdir(machine.config.BuildPath+"/json", 0700); err != nil {
		return err
	}
	return nil
}
