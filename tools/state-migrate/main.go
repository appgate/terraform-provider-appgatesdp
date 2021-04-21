package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	// matches block headers, ex:
	// resource "appgate_default_site" "lb-be2" {
	matchBlockHeader = regexp.MustCompile(`(provider|resource|data)(\s+")(appgate)(.*)`)

	// matches resource interpolation strings, ex:
	// default_site = "${appgate_default_site.lb1.id}"
	matchResourceInterpolation = regexp.MustCompile(`(.*?)(\${\s*)(appgate)(_.*?)`)

	// matches datasource interpolation strings, ex:
	// image = "${lookup(resource.appgate_condition.remedy_methods, "id")}"
	matchDatasourceInterpolation = regexp.MustCompile(`(.*?data)(\.)(appgate)(_.*?)`)
	// matches "appgate_ prefixes in statefile
	matchAppgate = regexp.MustCompile(`(.*)(")(appgate)(_.*)`)
)

// replace specific string patterns in template files
func replaceTemplateTokens(str string) string {
	str = matchBlockHeader.ReplaceAllString(str, `$1 "appgatesdp$4`)
	str = matchResourceInterpolation.ReplaceAllString(str, `$1${appgatesdp$4`)
	return matchDatasourceInterpolation.ReplaceAllString(str, `$1.appgatesdp$4`)
}

// replace appgate in statefile
func replaceStatefileTokens(str string) string {
	str = matchDatasourceInterpolation.ReplaceAllString(str, `$1.appgatesdp$4`)
	return matchAppgate.ReplaceAllString(str, `$1"appgatesdp$4`)
}

// FileAction Individual file io strategies for different operations
type FileAction func(string, string) error

// ProcessDirectory traverse a directory, executing the supplied FileAction on each file
func ProcessDirectory(targetDir string, backupDir string, fileActionFn FileAction, targetExtns ...string) (err error) {
	_, err = os.Stat(targetDir)

	if err != nil {
		return fmt.Errorf("Error reading directory\n %s", err)
	}

	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return fmt.Errorf("Error reading directory contents \n %s", err)
	}

	for _, res := range files {
		targetRes := path.Join(targetDir, res.Name())
		backupRes := path.Join(backupDir, res.Name())

		if res.IsDir() {
			if err := ProcessDirectory(targetRes, backupRes, fileActionFn, targetExtns...); err != nil {
				return err
			}
		} else {
			if len(targetExtns) == 0 {
				if err := fileActionFn(targetRes, backupRes); err != nil {
					return err
				}
			} else {
				if contains(targetExtns, filepath.Ext(res.Name())) {
					if err := fileActionFn(targetRes, backupRes); err != nil {
						return err
					}
				} else {
					fmt.Println("Skipping: ", targetDir)
				}
			}
		}
	}

	return
}

// CopyFile from targetFile path to backupFile path
func CopyFile(targetFile string, backupFile string) (err error) {

	// make sure directory structure exists
	bkDir := path.Dir(backupFile)
	_, err = os.Stat(bkDir)
	if err != nil {
		if os.IsNotExist(err) {
			oDir := path.Dir(targetFile)
			fi, err := os.Stat(oDir)

			if err != nil {
				return fmt.Errorf("Error reading original directory %s", err)
			}

			err = os.MkdirAll(bkDir, fi.Mode())

			if err != nil {
				return fmt.Errorf("Error creating directory for file %s", err)
			}
		} else {
			return fmt.Errorf("Unexpected error reading original directory %s", err)
		}
	}

	src, err := os.Open(targetFile)
	if err != nil {
		return fmt.Errorf("Error reading original file\n %s", err)
	}

	defer src.Close()

	dst, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("Error creating backup file\n %s", err)
	}

	defer dst.Close()

	fmt.Printf("Copying %s --> %s", targetFile, backupFile)

	size, err := io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("Error writing file\n %s", err)
	}

	fmt.Printf(", %d bytes\n", size)
	return
}

// MigratePlanFile Read file from backup location, apply transforms and overwrite original file
func MigratePlanFile(targetFile string, backupFile string) (err error) {
	src, err := os.Open(backupFile)
	if err != nil {
		return fmt.Errorf("Error reading file\n %s", err)
	}

	defer src.Close()

	dst, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("Error creating write location\n %s", err)
	}

	defer dst.Close()

	wrtr := bufio.NewWriter(dst)

	var replaceStrategy func(string) string
	if filepath.Ext(backupFile) == ".tf" {
		replaceStrategy = replaceTemplateTokens
	} else {
		replaceStrategy = replaceStatefileTokens
	}

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		str := scanner.Text()
		str = replaceStrategy(str)
		fmt.Fprintln(wrtr, str)
	}
	wrtr.Flush()

	return
}

// find a string in a slice of strings
func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

// CreateBackup copies target directory and append .backup
func CreateBackup(targetDir string, backupDir string) (err error) {
	fmt.Println("Creating backup...", targetDir, "-->", backupDir)

	fi, err := os.Stat(targetDir)
	if err != nil {
		return fmt.Errorf("Error reading directory\n %s", err)
	}

	if !fi.IsDir() {
		return fmt.Errorf("File targeted for migration")
	}

	_, err = os.Stat(backupDir)
	if err == nil {
		return fmt.Errorf("Attempting to overwrite backups")
	}

	fmt.Println("Copying", targetDir, "-->", backupDir)

	if err = ProcessDirectory(targetDir, backupDir, CopyFile); err != nil {
		return err
	}

	bfi, err := os.Stat(backupDir)
	if err != nil {
		return err
	}
	if fi.Size() != bfi.Size() {
		return fmt.Errorf("Backup corrupt")
	}

	fmt.Println("Complete")
	return
}

// RestoreBackup Overwrite target directory with contents of .backup directory
func RestoreBackup(backupDir string, targetDir string) (err error) {
	fmt.Println("Restoring from backup...")

	fi, err := os.Stat(backupDir)
	if err != nil {
		return fmt.Errorf("Error reading backup\n %s", err)
	}

	err = os.RemoveAll(targetDir)

	if err != nil {
		return fmt.Errorf("Error removing original directory\n %s", err)
	}

	os.MkdirAll(targetDir, fi.Mode())

	err = ProcessDirectory(backupDir, targetDir, CopyFile)
	if err != nil {
		return fmt.Errorf("Error restoring from backup directory\n %s", err)
	}

	fmt.Println("Complete")
	return
}

// DeleteBackup removes  .backup directory
func DeleteBackup(backupDir string) (err error) {
	fmt.Println("Purging backup...")

	err = os.RemoveAll(backupDir)
	if err != nil {
		return fmt.Errorf("Error removing backup directory\n %s", err)
	}

	fmt.Println("Complete")
	return
}

// Migrate Traverse all .tf files and apply transforms
func Migrate(targetDir string, backupDir string) (err error) {
	fmt.Println("Migrating plan directory...")
	err = CreateBackup(targetDir, backupDir)

	if err != nil {
		return fmt.Errorf("Error backing up directory before migration\n %s", err)
	}

	err = ProcessDirectory(targetDir, backupDir, MigratePlanFile, ".tf", ".tfstate")
	if err != nil {
		return fmt.Errorf("Error removing backup directory\n %s", err)
	}

	fmt.Println("Complete")
	return
}

func main() {
	if os.Args[1] == "backup" {
		backup := flag.NewFlagSet("backup", flag.PanicOnError)
		backup.Usage = func() {
			backup.PrintDefaults()
			os.Exit(0)
		}
		dir := backup.String("dir", "", "Required, specify the plan directory to operate on")
		purge := backup.Bool("purge", false, "Optional, whether to purge the backup directory")
		restore := backup.Bool("restore", false, "Optional, whether to restore from the backup directory")

		err := backup.Parse(os.Args[2:])

		if *dir == "" {
			fmt.Println("Missing required directory flag\nCommand flags:")
			backup.PrintDefaults()
			os.Exit(1)
		}

		if err != nil {
			panic(err)
		}

		targetDir := path.Clean(*dir)
		backupDir := targetDir + ".backup"

		fmt.Println(targetDir)

		if *purge {
			err := DeleteBackup(backupDir)

			if err != nil {
				panic(err)
			}

			return
		}

		if *restore {
			err := RestoreBackup(backupDir, targetDir)

			if err != nil {
				panic(err)
			}

			return
		}

		err = CreateBackup(targetDir, backupDir)

		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	if os.Args[1] == "migrate" {
		migrate := flag.NewFlagSet("migrate", flag.PanicOnError)
		migrate.Usage = func() {
			migrate.PrintDefaults()
			os.Exit(0)
		}
		dir := migrate.String("dir", "", "Required, specify the plan directory to operate on")
		err := migrate.Parse(os.Args[2:])

		if *dir == "" {
			fmt.Println("Missing required directory flag\nCommand flags:")
			migrate.PrintDefaults()
			os.Exit(1)
		}

		if err != nil {
			panic(err)
		}

		targetDir := path.Clean(*dir)
		backupDir := targetDir + ".backup"

		if err = Migrate(targetDir, backupDir); err != nil {
			panic(err)
		}

		if err = updateTerraformBlock(targetDir, backupDir); err != nil {
			panic(err)
		}
		printMessage()

		os.Exit(0)
	}

	fmt.Println("Unknown command")
	os.Exit(1)
}

// Traverse all .tf files and update `terraform.required_providers.appgate` to use the new name
func updateTerraformBlock(targetDir string, backupDir string) (err error) {
	fmt.Println("Scanning plans for terraform.required_providers.appgate...")

	err = ProcessDirectory(targetDir, backupDir, addAppGateSDPToRequiredProviders, ".tf")
	if err != nil {
		return fmt.Errorf("Error scanning terraform.required_providers.appgate\n %s", err)
	}

	fmt.Println("Complete")
	return
}

// Scan TF files for provider blocks and inject correct provider name value if not specified
func addAppGateSDPToRequiredProviders(targetFile string, backupFile string) error {
	fmt.Printf("Scanning %s\n", targetFile)

	fileInfo, err := os.Stat(targetFile)
	if err != nil {
		return fmt.Errorf("Error os stat provider block\n %s", err)
	}

	const maxSize = 1024 * 1024
	if fileInfo.Size() > maxSize {
		return fmt.Errorf("File too large to process")
	}

	fileBytes, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("Error read file provider block\n %s", err)
	}
	str, err := scanAndUpdateTerraform(string(fileBytes))
	if err != nil {
		return fmt.Errorf("Error updating provider block\n %s", err)
	}

	ioutil.WriteFile(targetFile, []byte(str), fileInfo.Mode())

	return err
}

// find all provider blocks in a string and update appgate block if applicable
func scanAndUpdateTerraform(content string) (string, error) {
	for start, i := 0, -1; ; {
		i, _ = findTokenAfter("terraform {", content, start)

		if i == -1 {
			return content, nil
		}

		start += i

		blockStart, blockEnd := indexOpenCloseTokens('{', '}', content[start:])
		if blockStart == -1 {
			return content, fmt.Errorf("terraform block detected, block start not found")
		}

		if blockEnd == -1 {
			return content, fmt.Errorf("terraform block detected, block end not found")
		}

		end := start + blockEnd + 1

		res, err := insertNewAppgateInProviderBlock(content[start:end])
		if err != nil {
			return content, fmt.Errorf("Problem parsing terraform block\n %s", err)
		}

		content = content[:start] + res + content[end:]

		start = end
	}
}

// rewrite matching required_providers.appgate blocks with new name and namespace
func insertNewAppgateInProviderBlock(content string) (string, error) {
	start, end, isAppgate, hasRequired := terraformBlockHasRequiredProviders(content)
	if start == -1 {
		return content, fmt.Errorf("terraform block start not detected")
	}

	if end == -1 {
		return content, fmt.Errorf("terraform block end not detected")
	}

	if isAppgate && hasRequired {
		start, end, isRight, _ := updateTerraformAppgateName(content)
		if start == -1 {
			return content, fmt.Errorf("required_providers block start not detected")
		}

		if end == -1 {
			return content, fmt.Errorf("required_providers block end not detected")
		}
		if isRight {
			content = strings.Replace(content, "appgate/appgate-sdp", "appgate/appgatesdp", 1)
			content = strings.Replace(content, "appgate =", "appgatesdp =", 1)
		}
	}

	return content, nil
}

var (
	matchTerraform         = regexp.MustCompile(`(.*terraform\s\{)`)
	matchRequiredProviders = regexp.MustCompile(`(.*required_providers\s\{)`)
	// matches `appgate =`
	matchOldAppgate = regexp.MustCompile(`\s*appgate\s*=`)
)

func updateTerraformAppgateName(content string) (start, end int, isAppgate, hasAppgate bool) {
	idx, _ := findToken("appgate", content)

	if idx == -1 {
		return -1, -1, false, false
	}
	subStr := content[idx:]                             // ignore everything before required_providers
	start, end = indexOpenCloseTokens('{', '}', subStr) // limit search to logical required_providers block
	if start == -1 || end == -1 {
		return start, end, false, false
	}
	isTarget := matchOldAppgate.MatchString(subStr[:end]) // make sure it's the right provider
	if !isTarget {
		return start, end, false, false
	}

	blkContents := subStr[start:end]                                                     // get just the logical block
	return idx + start, idx + end, true, matchRequiredProviders.MatchString(blkContents) // check for required_providers field
}

// find first provider block in a string, determine if it's the right provider, find start and end brace indices
func terraformBlockHasRequiredProviders(content string) (start, end int, isAppgate, hasRequired bool) {
	idx, _ := findToken("terraform {", content)

	if idx == -1 {
		return -1, -1, false, false
	}

	subStr := content[idx:]                             // ignore everything before required_providers
	start, end = indexOpenCloseTokens('{', '}', subStr) // limit search to logical required_providers block
	if start == -1 || end == -1 {
		return start, end, false, false
	}

	isTarget := matchTerraform.MatchString(subStr[:end]) // make sure it's the right provider
	if !isTarget {
		return start, end, false, false
	}

	blkContents := subStr[start:end]                                                     // get just the logical block
	return idx + start, idx + end, true, matchRequiredProviders.MatchString(blkContents) // check for required_providers field
}

func findToken(token string, content string) (start int, end int) {
	idx := strings.Index(content, token)
	return idx, idx + len(token)
}

// return the text extent of a token match in a string after a specified index
func findTokenAfter(token string, content string, begin int) (start int, end int) {
	newStr := content[begin:]
	idx := strings.Index(newStr, token)

	if idx == -1 {
		return -1, -1
	}

	return idx, idx + len(token)
}

// parse logical terraform blocks to find open and closing braces
func indexOpenCloseTokens(open rune, close rune, content string) (start int, end int) {
	ct := 0
	start = -1
	for idx := 0; idx < len(content); {
		rn, rnWidth := utf8.DecodeRuneInString(content[idx:])

		// keep track of opening brackets to account for nesting
		if rn == open {
			ct++
			if start < 0 { // start index still -1, record the first opening bracket
				start = idx
			}
		}

		// closing brackets decrement nest level
		if rn == close {
			ct--
			if ct == 0 { // bracket count back to 0, record the final closing bracket
				return start, idx
			}
		}

		idx += rnWidth
		nextRn, nextRnWidth := utf8.DecodeRuneInString(content[idx:])

		// match " and advance idx to closing "
		if rn == '"' {
			for idx < len(content)-1 {
				rn1, w1 := utf8.DecodeRuneInString(content[idx:])
				rn2, w2 := utf8.DecodeRuneInString(content[idx+w1:])

				if rn1 == '\\' && rn2 == '"' {
					idx += w1 + w2
					continue
				}

				idx += w1
				if rn1 == '"' {
					break
				}
			}
			continue
		}

		// match '#' and advance idx to line end
		if rn == '#' {
			for idx < len(content) {
				rn1, w1 := utf8.DecodeRuneInString(content[idx:])
				idx += w1

				if rn1 == '\n' {
					break
				}
			}
			continue
		}

		// match '//' and advance idx to line end
		if rn == '/' && nextRn == '/' {
			idx += nextRnWidth
			for idx < len(content) {
				rn1, w1 := utf8.DecodeRuneInString(content[idx:])
				if rn1 == '\n' {
					break
				}
				idx += w1
			}
			continue
		}

		// match '/*' and advance idx to closing '*/'
		if rn == '/' && nextRn == '*' {
			idx += nextRnWidth
			for idx < len(content)-1 {
				rn1, w1 := utf8.DecodeRuneInString(content[idx:])
				rn2, w2 := utf8.DecodeRuneInString(content[idx+w1:])
				idx += w1
				if rn1 == '*' && rn2 == '/' {
					idx += w2
					break
				}
			}
			continue
		}

		// match '${' and advance idx to closing '}'
		if rn == '$' && nextRn == '{' {
			idx += rnWidth + nextRnWidth
			for idx < len(content)-1 {
				rn1, w1 := utf8.DecodeRuneInString(content[idx:])
				idx += w1
				if rn1 == '}' {
					break
				}
			}
			continue
		}
	}

	return start, -1
}

func printMessage() {
	fmt.Println(`Migration Successful.`)
}
