package artmachine

import (
	"candy-machine/internal/common"
	"log"
	"sync"
)

type shared struct {
	mutex               *sync.RWMutex
	metadata            common.MetadataList
	dnaMap              map[string]int
	collectionSize      int
	uniqueDnaTollerance int
}

func newShared(collectionSize, uniqueDnaTollerance int) *shared {
	return &shared{
		mutex:               &sync.RWMutex{},
		metadata:            make(common.MetadataList, collectionSize),
		dnaMap:              make(map[string]int),
		collectionSize:      collectionSize,
		uniqueDnaTollerance: uniqueDnaTollerance,
	}
}

func (s *shared) addMetadata(index int, md *common.Metadata) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.metadata[index] = md
	if _, ok := s.dnaMap[md.Dna]; ok {
		s.dnaMap[md.Dna] = s.dnaMap[md.Dna] + 1
		log.Printf("dna is not unique, edition number: %d", index+1)
		if s.dnaMap[md.Dna] >= s.uniqueDnaTollerance {
			log.Fatalf("you need more layers to grow your collection to: %d", s.collectionSize)
		}
	} else {
		s.dnaMap[md.Dna] = 1
	}
}
