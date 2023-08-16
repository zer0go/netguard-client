package handler

import (
	"archive/zip"
	"bytes"
	"github.com/rs/zerolog/log"
	"github.com/zer0go/netguard-client/internal/config"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const (
	arch386       = "x86"
	dllFilePath32 = "C:\\Windows\\System32\\wireguard.dll"
)

func install() error {
	arch := getSimplifiedArch()
	dllFilePath := dllFilePath32
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
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}

	for _, zipFile := range zipReader.File {
		if !strings.Contains(zipFile.Name, arch+"/wireguard.dll") {
			continue
		}

		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			return err
		}

		err = os.WriteFile(dllFilePath, unzippedFileBytes, os.ModePerm)
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
