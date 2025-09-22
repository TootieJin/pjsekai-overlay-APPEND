package pjsekaioverlay

import (
	_ "embed"
	"io"
	"os"
	"path/filepath"
	"strings"

	wapi "github.com/iamacarpet/go-win64api"
	so "github.com/iamacarpet/go-win64api/shared"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

//go:embed sekai.obj
var sekaiObj []byte

//go:embed sekai-en.obj
var sekaiObjEn []byte

//go:embed sekai-v1.obj
var sekaiObjv1 []byte

//go:embed sekai-v1-en.obj
var sekaiObjEnv1 []byte

func TryInstallObject() bool {
	processes, _ := wapi.ProcessList()
	var aviutlProcess *so.Process
	for _, process := range processes {
		if process.Executable == "aviutl.exe" {
			aviutlProcess = &process
			break
		}
	}
	if aviutlProcess == nil {
		return false
	}
	aviutlPath := filepath.Dir(aviutlProcess.Fullpath)
	var exeditRoot string
	if _, err := os.Stat(filepath.Join(aviutlPath, "exedit.auf")); err == nil {
		exeditRoot = filepath.Join(aviutlPath)
	} else if _, err := os.Stat(filepath.Join(aviutlPath, "Plugins", "exedit.auf")); err == nil {
		exeditRoot = filepath.Join(aviutlPath, "Plugins")
	} else {
		return false
	}

	os.MkdirAll(filepath.Join(exeditRoot, "script"), 0755)

	var sekaiObjPath = filepath.Join(exeditRoot, "script", "@pjsekai-overlay.obj")
	if _, err := os.Stat(sekaiObjPath); err == nil {
		var sekaiObjFile, _ = os.OpenFile(sekaiObjPath, os.O_RDONLY, 0755)
		defer sekaiObjFile.Close()
		var sekaiObjDecoder = japanese.ShiftJIS.NewDecoder()
		var existingSekaiObj, _ = io.ReadAll(transform.NewReader(sekaiObjFile, sekaiObjDecoder))
		if strings.Contains(string(existingSekaiObj), "--version: "+Version) && Version != "0.0.0" {
			return false
		}
	}
	var sekaiObjPathEn = filepath.Join(exeditRoot, "script", "@pjsekai-overlay-en.obj")
	if _, err := os.Stat(sekaiObjPathEn); err == nil {
		var sekaiObjFileEn, _ = os.OpenFile(sekaiObjPathEn, os.O_RDONLY, 0755)
		defer sekaiObjFileEn.Close()
		var sekaiObjDecoderEn = japanese.ShiftJIS.NewDecoder()
		var existingSekaiObjEn, _ = io.ReadAll(transform.NewReader(sekaiObjFileEn, sekaiObjDecoderEn))
		if strings.Contains(string(existingSekaiObjEn), "--version: "+Version) && Version != "0.0.0" {
			return false
		}
	}
	var sekaiObjPathv1 = filepath.Join(exeditRoot, "script", "@pjsekai-overlay-v1.obj")
	if _, err := os.Stat(sekaiObjPathv1); err == nil {
		var sekaiObjFilev1, _ = os.OpenFile(sekaiObjPathv1, os.O_RDONLY, 0755)
		defer sekaiObjFilev1.Close()
		var sekaiObjDecoderv1 = japanese.ShiftJIS.NewDecoder()
		var existingSekaiObjv1, _ = io.ReadAll(transform.NewReader(sekaiObjFilev1, sekaiObjDecoderv1))
		if strings.Contains(string(existingSekaiObjv1), "--version: "+Version) && Version != "0.0.0" {
			return false
		}
	}
	var sekaiObjPathEnv1 = filepath.Join(exeditRoot, "script", "@pjsekai-overlay-en-v1.obj")
	if _, err := os.Stat(sekaiObjPathEnv1); err == nil {
		var sekaiObjFileEnv1, _ = os.OpenFile(sekaiObjPathEnv1, os.O_RDONLY, 0755)
		defer sekaiObjFileEnv1.Close()
		var sekaiObjDecoderEnv1 = japanese.ShiftJIS.NewDecoder()
		var existingSekaiObjEnv1, _ = io.ReadAll(transform.NewReader(sekaiObjFileEnv1, sekaiObjDecoderEnv1))
		if strings.Contains(string(existingSekaiObjEnv1), "--version: "+Version) && Version != "0.0.0" {
			return false
		}
	}

	err := os.MkdirAll(filepath.Join(exeditRoot, "script"), 0755)
	if err != nil {
		return false
	}
	sekaiObjFile, err := os.Create(sekaiObjPath)
	if err != nil {
		return false
	}
	defer sekaiObjFile.Close()

	sekaiObjFileEn, err := os.Create(sekaiObjPathEn)
	if err != nil {
		return false
	}
	defer sekaiObjFileEn.Close()

	sekaiObjFilev1, err := os.Create(sekaiObjPathv1)
	if err != nil {
		return false
	}
	defer sekaiObjFilev1.Close()

	sekaiObjFileEnv1, err := os.Create(sekaiObjPathEnv1)
	if err != nil {
		return false
	}
	defer sekaiObjFileEnv1.Close()

	var sekaiObjWriter = transform.NewWriter(sekaiObjFile, japanese.ShiftJIS.NewEncoder())
	var sekaiObjWriterEn = transform.NewWriter(sekaiObjFileEn, japanese.ShiftJIS.NewEncoder())
	var sekaiObjWriterv1 = transform.NewWriter(sekaiObjFilev1, japanese.ShiftJIS.NewEncoder())
	var sekaiObjWriterEnv1 = transform.NewWriter(sekaiObjFileEnv1, japanese.ShiftJIS.NewEncoder())

	strings.NewReader(strings.NewReplacer(
		"\r\n", "\r\n",
		"\r", "\r\n",
		"\n", "\r\n",
		"{version}", Version,
	).Replace(string(sekaiObj))).WriteTo(sekaiObjWriter)

	strings.NewReader(strings.NewReplacer(
		"\r\n", "\r\n",
		"\r", "\r\n",
		"\n", "\r\n",
		"{version}", Version,
	).Replace(string(sekaiObjEn))).WriteTo(sekaiObjWriterEn)

	strings.NewReader(strings.NewReplacer(
		"\r\n", "\r\n",
		"\r", "\r\n",
		"\n", "\r\n",
		"{version}", Version,
	).Replace(string(sekaiObjv1))).WriteTo(sekaiObjWriterv1)

	strings.NewReader(strings.NewReplacer(
		"\r\n", "\r\n",
		"\r", "\r\n",
		"\n", "\r\n",
		"{version}", Version,
	).Replace(string(sekaiObjEnv1))).WriteTo(sekaiObjWriterEnv1)
	return true
}

func TryInstallScript() bool {
	processes, _ := wapi.ProcessList()
	var aviutlProcess *so.Process
	for _, process := range processes {
		if process.Executable == "aviutl.exe" {
			aviutlProcess = &process
			break
		}
	}
	if aviutlProcess == nil {
		return false
	}
	aviutlPath := filepath.Dir(aviutlProcess.Fullpath)
	var exeditRoot string
	if _, err := os.Stat(filepath.Join(aviutlPath, "exedit.auf")); err == nil {
		exeditRoot = filepath.Join(aviutlPath)
	} else if _, err := os.Stat(filepath.Join(aviutlPath, "Plugins", "exedit.auf")); err == nil {
		exeditRoot = filepath.Join(aviutlPath, "Plugins")
	} else {
		return false
	}

	scriptDest := filepath.Join(exeditRoot, "script")
	err := os.MkdirAll(scriptDest, 0755)
	if err != nil {
		return false
	}

	depScriptDir := filepath.Join("dependencies", "aviutl script")

	var copyDir func(src, dest string) error
	copyDir = func(src, dest string) error {
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(dest, 0755); err != nil {
			return err
		}
		for _, entry := range entries {
			srcPath := filepath.Join(src, entry.Name())
			destPath := filepath.Join(dest, entry.Name())
			if entry.IsDir() {
				if err := copyDir(srcPath, destPath); err != nil {
					continue
				}
			} else {
				srcFile, err := os.Open(srcPath)
				if err != nil {
					continue
				}
				defer srcFile.Close()

				destFile, err := os.Create(destPath)
				if err != nil {
					srcFile.Close()
					continue
				}
				defer destFile.Close()

				_, err = io.Copy(destFile, srcFile)
				if err != nil {
					continue
				}
			}
		}
		return nil
	}

	if err := copyDir(depScriptDir, scriptDest); err != nil {
		return false
	}
	return true
}
