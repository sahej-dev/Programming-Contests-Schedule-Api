package db

import (
	"fmt"
	"sort"
	"time"

	"snow.sahej.io/utils"
)

const dB_PATH = "./data"           // no trailing slash
const dB_BACKUP_PATH = "./backups" // no trailing slash
const dB_FILENAME = "database"     // no leading slash and file extension
const dB_BACKUPS_TO_KEEP = 20

func GetDbPath() string {
	return fmt.Sprintf("%s/%s.sqlite", dB_PATH, dB_FILENAME)
}

// returns a backup ID which is required to restore backups
func Backup() (string, error) {
	err := utils.EnsureDirExists(dB_BACKUP_PATH)
	if err != nil {
		return "", err
	}

	backupId := time.Now().String()
	backupDbPath := getBackupDbPath(backupId)
	originalDbPath := GetDbPath()

	if !utils.DoesFileExists(originalDbPath) {
		return "", DbDoesNotExist{}
	}

	err = utils.CopyFile(originalDbPath, backupDbPath)

	return backupId, err
}

func Restore(backupId string) error {
	// backup current db, in case the restore fails
	originalBackupId, err := Backup()
	if _, ok := err.(DbDoesNotExist); err != nil && !ok {
		return err
	}

	didOriginalDbExist := err != nil

	err = utils.CopyFile(getBackupDbPath(backupId), GetDbPath())
	if err != nil && didOriginalDbExist {
		Restore(originalBackupId)
	}

	return err

}

func DeleteOldBackupsIfAny() error {
	files, err := utils.ListDirFiles(dB_BACKUP_PATH)
	if err != nil {
		return err
	}

	sort.Strings(files)

	if len(files) <= dB_BACKUPS_TO_KEEP {
		return nil
	}

	filesToDelete := files[:len(files)-dB_BACKUPS_TO_KEEP]

	err = utils.DeleteFiles(filesToDelete)

	return err
}

func getBackupDbPath(backupId string) string {
	return fmt.Sprintf("%s/%s_%s.sqlite", dB_BACKUP_PATH, backupId, dB_FILENAME)
}
