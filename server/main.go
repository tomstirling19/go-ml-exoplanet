package exoplanet

func main() {
	err := scripts.data()
	if err != nil {
		log.Errorf("Error downloading TESS dataset: %v", err)
	}
}
