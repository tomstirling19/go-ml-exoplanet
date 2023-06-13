package exoplanet

func main(downlaodData bool,) {
	err := scripts.data()
	if err != nil {
		log.Errorf("Error downloading TESS dataset: %v", err)
	}
}
