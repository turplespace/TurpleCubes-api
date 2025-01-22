package repositories 

import (
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"	
	"github.com/turplespace/portos/internal/models"
)

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
       
        return images,err
    }

    // Unmarshal the JSON data into the Images struct
   
    err = json.Unmarshal(byteValue, &images)
    if err != nil {
     
        return images , err
    }
	return images ,nil
}
