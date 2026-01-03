package handler

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zer0go/netguard-client/internal/config"
)

const (
	arch386     = "x86"
	dllFilePath = "C:\\Windows\\System32\\wireguard.dll"
)

func install() error {
	arch := getSimplifiedArch()

	if _, err := os.Stat(dllFilePath); err == nil {
		log.Debug().
			Str("arch", arch).
			Str("dll_file", dllFilePath).
			Msg("driver already exists, skipping download")
		return nil
	}

	log.Debug().Msg("downloading wireguard driver...")

	resp, err := http.Get(config.WireGuardNTUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}

	// Look for the correct architecture path
	archPath := arch + "/wireguard.dll"
	if runtime.GOARCH == "amd64" {
		archPath = "amd64/wireguard.dll"
	}

	for _, zipFile := range zipReader.File {
		if !strings.Contains(zipFile.Name, archPath) {
			continue
		}

		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			return err
		}

		err = os.WriteFile(dllFilePath, unzippedFileBytes, 0644)
		if err != nil {
			return err
		}

		log.Debug().
			Str("arch", arch).
			Str("file_name", zipFile.Name).
			Str("dll_file", dllFilePath).
			Msg("windows driver saved")

		break
	}

	log.Debug().Msg("wireguard driver downloaded.")

	return nil
}

func getSimplifiedArch() string {
	if runtime.GOARCH == "386" {
		return arch386
	}

	return runtime.GOARCH
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer func(f io.ReadCloser) {
		_ = f.Close()
	}(f)

	return io.ReadAll(f)
}
