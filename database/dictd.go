package database

import (
	"os"
	"pault.ag/go/dictd/dictd"
)

// NewDictdDatabase makes a new DictdDatabase using
// the dictionary files and a description.
// Usually file names stem from a configuration file.

func NewDictdDatabase(indexPath string, dbPath string, description string) (*DictdDatabase, error) {

	iFile, err1 := os.Open(indexPath)
	if err1 != nil {
		fmt.Println(err1)
	}
	defer iFile.Close()

	dFile, err2 := os.Open(dbPath)
	if err2 != nil {
		fmt.Println(err2)
	}
	defer dFile.Close()

	databaseBackend := DictdDatabase{
		description: description,
		index:       iFile,
		database:    dFile,
	}

	return &databaseBackend, nil
}

// DictdDatabase container is composed of
// a description of the database
// and file descriptors to the index and database files.

type DictdDatabase struct {
	dictd.Database

	description string
	index       *File
	database    *File
}
