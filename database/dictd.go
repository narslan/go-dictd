package database

import (
	"bufio"
	"errors"
	_ "fmt"
	"github.com/pebbe/dictzip"
	"os"
	"pault.ag/go/dictd/dictd"
	"strings"
)

type IndexPos struct {
	position string
	length   string
}

// NewDictdDatabase makes a new DictdDatabase using
// the dictionary files and a description.
// Usually file names stem from a configuration file.

func NewDictdDatabase(name string, indexPath string, dbPath string) (*DictdDatabase, error) {

	if _, err := os.Stat(dbPath); err != nil {
		return nil, err
	}
	if _, err := os.Stat(indexPath); err != nil {
		return nil, err
	}

	databaseBackend := DictdDatabase{
		Name:  name,
		Index: indexPath,
		File:  dbPath,
	}

	return &databaseBackend, nil
}

// DictdDatabase container is composed of
// a description of the database
// and file descriptors to the index and database files.

type DictdDatabase struct {
	dictd.Database
	Name  string
	Index string
	File  string
}

/* Get all valid Strategies */
func (d *DictdDatabase) Strategies(name string) map[string]string {
	return map[string]string{
		"substring": "Match a substring",
		"suffix":    "Suffix Match",
		"re":        "POSIX regular expression search",
		"soundex":   "The Soundex algorithm",
		"lev":       "The Levenshtein algorithm",
	}
}

//TODO:
/* Handle the information call (SHOW INFO `name`) for this database. */
func (d *DictdDatabase) Info(text string) string {
	return "Not Impelemented"
}

//TODO:
/* Handle the short description of what this database does (for
 * inline `SHOW DB` output) */
func (d *DictdDatabase) Description(s string) string {
	return s
}

/* Handle incoming `DEFINE` calls. */
func (this *DictdDatabase) Define(name string, query string) []*dictd.Definition {
	query = strings.ToLower(query)
	data, err := this.get(query)
	if err != nil {
		/* If we don't have the key, let's bail out. */
		return make([]*dictd.Definition, 0)
	}
	els := make([]*dictd.Definition, 1)
	els[0] = &dictd.Definition{
		DictDatabase:     this,
		DictDatabaseName: name,
		Word:             query,
		Definition:       string(data),
	}
	return els
}

func (d *DictdDatabase) get(query string) (value string, err error) {

	iFile, err1 := os.Open(d.Index)
	if err1 != nil {
		return "", err1
	}
	defer iFile.Close()

	dbFile, err2 := os.Open(d.File)
	if err2 != nil {
		return "", err2
	}
	defer dbFile.Close()

	// inFile is a map of word positons and length
	indexMap := make(map[string]IndexPos, 0)

	scanner := bufio.NewScanner(iFile)

	for scanner.Scan() {
		s := scanner.Text()
		si := strings.Split(s, "\t")

		//		fmt.Printf("%s %d\n", s, len(si))
		if len(si) == 3 {
			indexMap[si[0]] = IndexPos{
				position: si[1],
				length:   si[2],
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	w, ok := indexMap[query]
	if ok == false {
		return "", errors.New("No definition found")
	}

	r, rerr := dictzip.NewReader(dbFile)
	if rerr != nil {
		return "", rerr
	}

	data, gerr := r.GetB64(w.position, w.length)
	if gerr != nil {
		return "", gerr
	}
	//data = []byte("")
	return string(data), err

}
