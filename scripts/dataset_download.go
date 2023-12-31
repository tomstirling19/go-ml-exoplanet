package scripts

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

// get script for TESS dataset sector 64 (latest) lightcurves
// a full list can be seen here: https://archive.stsci.edu/tess/bulk_downloads/bulk_downloads_ffi-tp-lc-dv.html
func getDatasetScript(outputDir string) (string, error) {
	url := "https://archive.stsci.edu/missions/tess/download_scripts/sector/tesscurl_sector_64_lc.sh"

	urlParts := strings.Split(url, "/")
	filename := urlParts[len(urlParts)-1]
	scriptFile := filepath.Join(outputDir, filename)

	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("unable to download dataset script file: %v", err)
	}
	defer res.Body.Close()

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return "", fmt.Errorf("unable to create output directory: %v", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %v", err)
	}

	err = ioutil.WriteFile(scriptFile, body, 0644)
	if err != nil {
		return "", fmt.Errorf("unable to write output file: %v", err)
	}

	err = os.Chmod(scriptFile, 0755)
	if err != nil {
		return "", fmt.Errorf("unable to set script file as executable: %v", err)
	}

	fmt.Printf("successfully downloaded TESS .sh file to %s.\n", outputDir)

	return scriptFile, nil
}

// execute TESS sector 64 script to obtain all sector 64 exoplanet data
func getDataset(ctx context.Context, scriptFile string) error {
    absScriptFile, err := filepath.Abs(scriptFile)
    if err != nil {
        return fmt.Errorf("unable to get absolute path of script file: %v", err)
    }

    cmd := exec.CommandContext(ctx, "/bin/sh", "-c", absScriptFile)
    cmd.Dir = filepath.Join("data", "raw", "fits")

    s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
    s.Prefix = "downloading... "
    s.Start()
    output, err := cmd.CombinedOutput()
    s.Stop()

    if err != nil {
        return fmt.Errorf("script execution failed: %v\nOutput: %s", err, output)
    }

    fmt.Println("successfully executed the script.")

    return nil
}

// remove the .sh file that is no longer required
func removeScriptFile(scriptPath string) error {
	err := os.Remove(scriptPath)
	if err != nil {
		return fmt.Errorf("unable to remove script file: %v", err)
	}

	fmt.Println("successfully removed the script file.")

	return nil
}

func datasetDownload() {
	outputDir := filepath.Join("data", "raw", "fits")
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		cancel()
	}()

	scriptFile, err := getDatasetScript(outputDir)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	err = getDataset(ctx, scriptFile)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	err = removeScriptFile(scriptFile)
	if err != nil {
		fmt.Println("error: ", err)
	}

	cancel()
}
