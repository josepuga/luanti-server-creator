// const MY_PATH = "/home/jose/.minetest/worlds/"
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/josepuga/goini"
)

var version string = "<unknow>"

const CONFIG_FILE = "./config.ini"

type config struct {
	srcDir             string
	destDir            string
	destWorldDir       string
	destModsDir        string
	destGamesDir       string
	serversPath        string
	worldName          string
	worldDir           string
	gameID             string
	worldModsDir       []string
	worldOptions       []string
	worldNameDirectory map[string]string
	worldMtFile        string
}

var cfg config

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&CustomTheme{})
	myWindow := myApp.NewWindow("Luanti Server Creator " + version)

	myWindow.Resize(fyne.NewSize(400, 550))

    // Create config.ini file with default values if not exists
    if ! fileExists(CONFIG_FILE) {
        contentBytes, err := othersFS.ReadFile("embed/others/config.ini")
        if err != nil {
            fmt.Fprintln(os.Stderr, "Error reading embeded file config.ini, using default values")
        } else {
			if err := saveToFile(CONFIG_FILE, contentBytes); err != nil {
				fmt.Fprintln(os.Stderr, "Error creating config.ini, using default values")
			}
        }
    }
    
	//
	// Get Luanti data path. Exit if not exists
	//
	homeDir, _ := os.UserHomeDir()
	cfg.srcDir = filepath.Join(homeDir, ".minetest")
	// Try the config file
	ini := goini.NewIni()
	if err := ini.LoadFromFile("config.ini"); err == nil {
		cfg.srcDir = ini.GetString("", "data_path", cfg.srcDir)
        // Expand ~ if present
        cfg.srcDir = expandPath(cfg.srcDir)
	}
	if !isDir(cfg.srcDir) {
		if !isDir(cfg.srcDir) {
			fmt.Fprintln(os.Stderr, "Error Luanti data directory does not exists", cfg.srcDir)
			return
		}
	}

	// Servers location
	cfg.serversPath = ini.GetString("", "servers_path", "./servers")
    cfg.serversPath = expandPath(cfg.serversPath)

	//
	// Gather a map [World Name]Word Directory. Exit if no worlds or error
	//
	err := readWorldNamesAndDirectories()
	if err != nil || len(cfg.worldNameDirectory) == 0 {
		fmt.Fprintln(os.Stderr, "Error reading or no Worlds found at", cfg.srcDir)
		return
	}

	//
	// Populate Worlds ComboBox
	//
	var worlds []string
	for world := range cfg.worldNameDirectory {
		worlds = append(worlds, world)
	}
	sort.Strings(worlds)
	worldCombo := widget.NewSelect(worlds, nil)

	//
	// Server Tags Image ComboBox
	//

	// Server Tags Images availables. Read from Docker Hub API
	serverTagsCombo := widget.NewSelect([]string{"latest"}, nil)
	if imageTags, err := getDockerTags("josepuga", "luanti-server"); err != nil {
		// Not an importarn error
		fmt.Fprintln(os.Stderr, "Error fetching server tags from docker.io. Only latest available")
	} else {
		serverTagsCombo.Options = imageTags
	}
	serverTagsCombo.Selected = "latest"

	//
	// Multiline, minetest.conf
	//
	confMultiline := widget.NewMultiLineEntry()
	confMultiline.SetPlaceHolder("Write here...")
	confMultiline.Resize(fyne.NewSize(300, 120))
	if confContent, err := othersFS.ReadFile("embed/others/minetest.conf"); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading default minetest.conf")
	} else {
		confMultiline.Text = string(confContent)
		confMultiline.Refresh()
	}

	//
	// Button Generate
	//
	// Must be self referenciated to disable it.
	var createServerButton *widget.Button
	createServerButton = widget.NewButton("Create Server!", func() {

        if worldCombo.Selected == "" {
            return
		}
        start := time.Now()
		createServerButton.Disable()
		errorCount := 0

		cfg.worldName = worldCombo.Selected
		cfg.worldDir = cfg.worldNameDirectory[cfg.worldName]
		cfg.worldMtFile = filepath.Join(cfg.srcDir, "worlds", cfg.worldDir, "world.mt")
		cfg.destDir = filepath.Join(cfg.serversPath, cfg.worldDir)
		cfg.destWorldDir = filepath.Join(cfg.destDir, "data/worlds/world")
		cfg.destModsDir = filepath.Join(cfg.destWorldDir, "worldmods")
		cfg.destGamesDir = filepath.Join(cfg.destDir, "data/games")
		cfg.gameID = getGameID()
		// Fix, if gameID is empty all games will be copied
		if cfg.gameID == "" {
			cfg.gameID = "GAMEID_NOT_FOUND"
		} else {
			// Another Fix, minetest game has 'minetest_game' folder ... ¿?¿?¿
			if cfg.gameID == "minetest" {
				cfg.gameID = "minetest_game"
			}
		}

		// Delete destination Dir if exists
		fmt.Printf("Deleting files from %s...\n", cfg.destDir)
		if isDir(cfg.destDir) {
			err := deleteDir(cfg.destDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error deleting", cfg.destDir, err.Error())
				errorCount++
			}
		}

		//
		// Full directory copy
		//

		// World and game
		paths := map[string]string{
			filepath.Join(cfg.srcDir, "worlds", cfg.worldDir): cfg.destWorldDir,
			filepath.Join(cfg.srcDir, "games", cfg.gameID):    filepath.Join(cfg.destGamesDir, cfg.gameID),
		}

		// Active mods
		if cfg.worldModsDir, err = worldMtGetActiveModDirs(); err != nil {
			fmt.Fprintln(os.Stderr, "Error getting mods location")
			errorCount++
		}
		for _, p := range cfg.worldModsDir {
			srcDir := filepath.Join(cfg.srcDir, p)
			dstDir := filepath.Join(cfg.destModsDir, p)
			partToRemove := filepath.Join("mods")
			// Not in world/worldmods/mods/<my mod> --> world/worldmods/<my mod>
			dstDir = removePartOfPath(dstDir, partToRemove)
			//fmt.Printf("%s -> %s\n", srcDir, dstDir)
			paths[srcDir] = dstDir
		}

		fmt.Println("Copying many directories, this may take a while...")
        var wg sync.WaitGroup
		for src, dst := range paths {
            wg.Add(1)
			//fmt.Printf("Copying %s -> %s\n", src, dst)
            go func(src, dst string) {
                defer wg.Done()
                err := copyDir(src, dst)
                if err != nil {
                    fmt.Fprintln(os.Stderr, "Error copying ", src, "to", dst)
                }
            }(src, dst)
		}

		//
		// Create scripts & files
		//
		scriptFiles := []string{
			"start-server.sh",
			"stop-server.sh",
			"start-server.bat",
			"stop-server.bat",
		}

		fmt.Println("Generating scripts...")
		for _, scriptFile := range scriptFiles {
			//fmt.Println(systemDestFile)
			contentBytes, err := scriptsFS.ReadFile("embed/scripts/" + scriptFile)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading embeded file ", scriptFile)
				errorCount++
				continue
			}

			// Server name (sanitized)
			contentString := strings.ReplaceAll(
				string(contentBytes),
				"%server%",
				"server-"+sanitize(cfg.worldDir),
			)

			// Image tags
			contentString = strings.ReplaceAll(
				contentString,
				"%tag%",
				serverTagsCombo.Selected,
			)

			scriptDestFile := filepath.Join(cfg.destDir, scriptFile)
			if err := saveToFile(scriptDestFile, []byte(contentString)); err != nil {
				fmt.Fprintln(os.Stderr, "Error creating system file", scriptFile)
				errorCount++
				continue
			}
			// Check .sh extension for +x
			// No errors handle
			if filepath.Ext(scriptDestFile) == ".sh" {
				os.Chmod(scriptDestFile, 0744)
			}
		}

		// Create minetest.conf, from multiline widget.
		fmt.Println("Generating minetest.conf")
		if err := saveToFile(filepath.Join(cfg.destDir, "data/minetest.conf"), []byte(confMultiline.Text)); err != nil {
			fmt.Fprintln(os.Stderr, "Error copying minetest.conf")
			errorCount++
		}

		// (re)Create world.mt ripping off mod_references
		fmt.Println("Generating world.mt", cfg.destWorldDir)
		woldMtBytes := []byte(worldMtGetOnlyOptions())
		if err := saveToFile(filepath.Join(cfg.destWorldDir, "world.mt"), woldMtBytes); err != nil {
			fmt.Fprintln(os.Stderr, "Error copying world.mt")
		}
		if errorCount == 0 {
			dialog.ShowInformation("",
				fmt.Sprintf("Server created."),
				myWindow,
			)
		} else {
			dialog.ShowInformation("",
				fmt.Sprintf("Process ended up with %d errors", errorCount),
				myWindow,
			)
		}
        elapse := time.Since(start)
		fmt.Println("Done!. Server created in", elapse)
		createServerButton.Enable()
	})

	// Container with all the elements.
	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Select a World:"),
			worldCombo,
			widget.NewSeparator(),
			widget.NewLabel("Server version:"),
			serverTagsCombo,
			widget.NewSeparator(),
			widget.NewLabel("minetest.conf:"),
		), // Top
		createServerButton, // Bottom
		nil,            // Left
		nil,            //Right
		confMultiline,  // Center
	)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func readWorldNamesAndDirectories() error {
	dirList, err := getDirectories(filepath.Join(cfg.srcDir, "worlds"))
	if err != nil {
		return err
	}
	ini := goini.NewIni()
	cfg.worldNameDirectory = make(map[string]string)
	for _, d := range dirList {
		dir := filepath.Join(cfg.srcDir, "worlds", d)
		// Extract World name (I not sure if "always" is the same as directory)
		if err := ini.LoadFromFile(filepath.Join(dir, "world.mt")); err != nil {
			continue
		}
		worldName := ini.GetString("", "world_name", "")
		if worldName == "" {
			continue
		}
		cfg.worldNameDirectory[worldName] = filepath.Base(dir)
	}
	return nil
}

func getGameID() string {
	ini := goini.NewIni()
	if err := ini.LoadFromFile(cfg.worldMtFile); err != nil {
		return ""
	}
	return ini.GetString("", "gameid", "")
}

func sanitize(s string) string {
	s = strings.ToLower(s)

	// Reemplazar todo lo que no sean letras o números por "_"
	re := regexp.MustCompile(`[^a-z0-9]+`)
	return re.ReplaceAllString(s, "_")
}

func expandPath(path string) string {
    if strings.HasPrefix(path, "~") {
        home, err := os.UserHomeDir()
        if err != nil {
            return path
        }
        return filepath.Join(home, strings.TrimPrefix(path, "~"))
    }
    return path
}

func removePartOfPath(fullPath, partToRemove string) string {
	var separator string
	if runtime.GOOS == "windows" {
		separator = `\`
	} else {
		separator = `/`
	}
	target := separator + partToRemove + separator
	return strings.Replace(fullPath, target, separator, -1)
}
