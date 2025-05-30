package greenplum

import (
	"github.com/wal-g/wal-g/internal"
	"github.com/wal-g/wal-g/internal/multistorage"
	"github.com/wal-g/wal-g/pkg/storages/storage"
	"github.com/wal-g/wal-g/utility"
)

func newBackupObject(incrementBase, incrementFrom string,
	isFullBackup bool, object storage.Object, storageName string) BackupObject {
	return BackupObject{
		BackupObject:      internal.NewDefaultBackupObject(object),
		isFullBackup:      isFullBackup,
		baseBackupName:    incrementBase,
		incrementFromName: incrementFrom,
		storageName:       storageName,
	}
}

type BackupObject struct {
	internal.BackupObject
	isFullBackup      bool
	baseBackupName    string
	incrementFromName string
	storageName       string
}

var _ internal.BackupObject = BackupObject{}

func (o BackupObject) IsFullBackup() bool {
	return o.isFullBackup
}

func (o BackupObject) GetBaseBackupName() string {
	return o.baseBackupName
}

func (o BackupObject) GetIncrementFromName() string {
	return o.incrementFromName
}

func (o BackupObject) GetStorage() string {
	return o.storageName
}

func makeBackupObjects(folder storage.Folder, objects []storage.Object) ([]internal.BackupObject, error) {
	backupObjects := make([]internal.BackupObject, 0, len(objects))
	for _, object := range objects {
		storageName := multistorage.GetStorage(object)
		incrementBase, incrementFrom, isFullBackup, err := getIncrementInfo(folder, object, storageName)
		if err != nil {
			return nil, err
		}
		gpBackup := newBackupObject(incrementBase, incrementFrom, isFullBackup, object, storageName)

		backupObjects = append(backupObjects, gpBackup)
	}
	return backupObjects, nil
}

func getIncrementInfo(folder storage.Folder, object storage.Object, storageName string) (string, string, bool, error) {
	backup, err := NewBackupInStorage(folder, utility.StripRightmostBackupName(object.GetName()), storageName)
	if err != nil {
		return "", "", true, err
	}
	sentinel, err := backup.GetSentinel()
	if err != nil {
		return "", "", true, err
	}
	if !sentinel.IsIncremental() {
		return "", "", true, nil
	}

	return *sentinel.IncrementFullName, *sentinel.IncrementFrom, false, nil
}
