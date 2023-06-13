package scripts

import (
	"fmt"
	"os"

	"github.com/siravan/fits"
)

// processFITSFile reads a FITS file, performs some preprocessing
// (currently placeholder), and returns an error if any occurred
func processFITSFile(filename string) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the FITS file
	units, err := fits.Open(file)
	if err != nil {
		return err
	}

	// Loop over each unit in the FITS file
	for i, unit := range units {
		// Access and print the BITPIX header key
		bitpix, ok := unit.Keys["BITPIX"].(int)
		if !ok {
			return fmt.Errorf("could not access BITPIX key in unit %d", i)
		}
		fmt.Printf("BITPIX in unit %d: %d\n", i, bitpix)

		// Access the image data points
		// This is a placeholder and would need to be adjusted based on the actual structure of the data
		if len(unit.Naxis) >= 2 {
			for x := 0; x < unit.Naxis[0]; x++ {
				for y := 0; y < unit.Naxis[1]; y++ {
					// Access the pixel value at location (x, y)
					// Type assertion would need to be adjusted based on the actual data type of the pixel values
					pixelValue, ok := unit.At(x, y).(float64)
					if !ok {
						return fmt.Errorf("could not access pixel value at (%d, %d) in unit %d", x, y, i)
					}

					// Perform some cleaning on the pixel value
					// This is a placeholder and would need to be replaced with actual preprocessing steps
					cleanedPixelValue := pixelValue // Placeholder

					fmt.Printf("Cleaned pixel value at (%d, %d) in unit %d: %f\n", x, y, i, cleanedPixelValue)
				}
			}
		}
	}

	return nil
}

func datasetPreprocessing() {

	// Read a single FITS file
	err := processFITSFile("tess2023096110322-s0064-0000000001668887-0257-s_lc.fits")
	if err != nil {
		fmt.Println(err)
	}
}
