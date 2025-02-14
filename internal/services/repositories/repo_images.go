package repositories

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/turplespace/portos/internal/models"
)

// ReadImages reads the images from the JSON file
func ReadImages() (models.ImagesConfig, error) {
	var images models.ImagesConfig
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := fmt.Sprintf("%s_conf/images.json", ex)
	file, err := os.Open(path)
	if err != nil {
		return images, err
	}
	defer file.Close()

	// Decode the JSON data into a map

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {

		return images, err
	}

	// Unmarshal the JSON data into the Images struct

	err = json.Unmarshal(byteValue, &images)
	if err != nil {

		return images, err
	}
	return images, nil
}

// AppendImages appends the image to the images JSON file
func AppendImages(image models.Image) error {
	fmt.Println("Reading images")
	images, err := ReadImages()
	if err != nil {
		fmt.Println("Error reading images", err)
		return err
	}

	// Check if the image already exists
	for _, img := range images.CustomImages {
		if img.Image == image.Image {
			fmt.Println("Image already exists")
			return nil
		}
	}

	fmt.Println("Appending image")
	images.CustomImages = append(images.CustomImages, image)
	err = WriteImages(images)
	if err != nil {
		fmt.Println("Error writing images", err)
		return err
	}
	return nil
}

// WriteImages writes the images to the JSON file
func WriteImages(images models.ImagesConfig) error {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := fmt.Sprintf("%s_conf/images.json", ex)
	file, err := json.MarshalIndent(images, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, file, 0644)
	if err != nil {
		return err
	}
	return nil

}
