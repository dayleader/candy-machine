# candy-machine

## Installation 🛠️

If you are cloning the project then run this first, otherwise you can download the source code on the release page and skip this step.

```sh
git clone git@github.com:golang-enthusiast/candy-machine.git
```

Go to the root of your folder and run this command if you have go installed.

```sh
go mod download
```

## Usage ℹ️

Create your different layers as folders in the 'layers' directory, and add all the layer assets in these directories. You can name the assets whatever you want without any restrictions.

Once you have all your layers, go into `config.json` and update the `layerConfigurations` objects `layersOrder` array to be your layer folders name in order of the back layer to the front layer.

_Example:_ If you were creating a portrait design, you might have a background, then a face, hair, eyes, lips, so your `layersOrder` would look something like this:

```json
{
  "layerConfigurations": [
    {
      "size": 5,
      "layersOrder": [
        {
          "name": "background"
        },
        {
          "name": "face"
        },
        {
          "name": "hair"
        },
        {
          "name": "eyes"
        },
        {
          "name": "lips"
        }
      ]
    }
  ]
}
```

The `name` of each layer object represents the name of the folder (in `/layers/`) that the images reside in.

Optionally you can now add multiple different `layerConfigurations` to your collection. Each configuration can be unique and have different layer orders, use the same layers or introduce new ones. This gives the artist flexibility when it comes to fine tuning their collections to their needs.

```json
{
  "layerConfigurations": [
    {
      "size": 5,
      "layersOrder": [
        {
          "name": "background"
        },
        {
          "name": "face"
        },
        {
          "name": "hair"
        },
        {
          "name": "eyes"
        },
        {
          "name": "lips"
        }
      ]
    },
    {
      "size": 2,
      "layersOrder": [
        {
          "name": "background"
        },
        {
          "name": "face"
        },
        {
          "name": "hair"
        },
        {
          "name": "eyes"
        },
        {
          "name": "lips"
        },
        {
          "name": "accessories"
        }
      ]
    }
  ]
}
```

Update your `format` size, ie the outputted image size, and the `size` on each `layerConfigurations` object, which is the amount of variation outputted.

The `collectionStartsFrom` represents a starting index for your collection.

_Example:_ 
- ETH - collection starts from 1
- SOL - collection starts from 0

The `uniqueDnaTollerance` represents the number of allowed dna duplicates. The default value is 0, which means zero tolerance for duplicates.

_Example:_ If set to 5, this means that 5 DNA duplicates are allowed in the generated collection.

To use a different metadata attribute name you can add the `displayName: "Awesome Background"` to the `options` object. All options are optional and can be addes on the same layer if you want to.

To make one asset rarer than another, you need to set a `weight` and then define the `filenames` to which this weight will be applied. The default weight value is 1. Setting the weight value to more than 1 means that the layers will be used n times less than the base ones.

_Example:_

- weight = 2, means that layers will be used 2 times less often than the base ones
- weight = 3, means that layers will be used 3 times less often than the base ones
- weight = 4, means that the layers will be used 4 times less than the base ones

Here is an example on how you can play around with filter fields:

```json
{
  "layerConfigurations": [
    {
      "size": 5,
      "layersOrder": [
        {
          "name": "background",
          "options": {
            "displayName": "Background",
            "rarity": [
              {
                "weight": 2,
                "filenames": [
                  "yellow.png"
                ]
              }
            ]
          }
        },
        {
          "name": "face",
          "options": {
            "displayName": "Face"
          }
        },
        {
          "name": "hair",
          "options": {
            "displayName": "Hair",
            "rarity": [
              {
                "weight": 3,
                "filenames": [
                  "silver.png"
                ]
              }
            ]
          }
        },
        {
          "name": "eyes"
        },
        {
          "name": "lips"
        }
      ]
    }
  ]
}
```

When you are ready, run the following command and your outputted art will be in the `build/images` directory and the json in the `build/json` directory:

```sh
go run main.go --generate-art
```

or if you are using binary

```sh
./candy-machine --generate-art
```

The program will output all the images in the `build/images` directory along with the metadata files in the `build/json` directory. Each collection will have a `_metadata.json` file that consists of all the metadata in the collection inside the `build/json` directory. The `build/json` folder also will contain all the single json files that represent each image file. The single json file of a image will look something like this:

```json
{
   "name":"Your Collection #1",
   "description":"Remember to replace this description",
   "image":"ipfs://NewUriToReplace/1.png",
   "external_url":"",
   "dna":"fc0a52e302db9afb69b193a7f4570631378e4969",
   "attributes":[
      {
         "trait_type":"Background",
         "value":"Steelblue"
      },
      {
         "trait_type":"Face",
         "value":"Chocolate"
      },
      {
         "trait_type":"Hair",
         "value":"Blue"
      },
      {
         "trait_type":"Eyes",
         "value":"Blue"
      },
      {
         "trait_type":"Lips",
         "value":"Coral"
      }
   ]
}
```
