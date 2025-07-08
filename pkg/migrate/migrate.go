package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"taskmgmtsystem/pkg/sqlparser"
)

type Migrate struct {
	path           string
	db             *sql.DB
	migrationFiles []DirEntryWithPrefix
	txn            *sql.Tx
}

func NewMigrate(db *sql.DB, dirPath string) Migrate {
	return Migrate{
		path: dirPath,
		db:   db,
	}
}

func (m *Migrate) RunMigrations() error {
	rawEntries, err := os.ReadDir(m.path)
	if err != nil {
		fmt.Println("ERR1")
		return err
	}

	usableEntries := m.filterSqlFilesWithNumberPrefix(
		m.getFilesFromDirEntries(rawEntries),
	)

	m.sortDirEntryBasedOnPrefix(usableEntries)

	err = m.checkForSamePrefix(usableEntries)
	if err != nil {
		fmt.Println("ERR2")
		return err
	}

	version, err1 := m.getVersion()
	if err1 != nil {
		fmt.Println("ERR3")
		return err1
	}

	// It means we are already on latest db state
	if version == len(usableEntries) {
		fmt.Println("ERR4")
		return nil
	}

	// This means that db is never created
	if version == -1 {
		m.migrationFiles = usableEntries
	} else {
		m.migrationFiles = usableEntries[version:]
	}

	m.txn, err = m.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer m.txn.Rollback()

	err = m.parseFilesAndMigrateDb()
	if err != nil {
		return err
	}

	// clear file
	err = os.Truncate(m.path+"/migrate.log", 0)
	if err != nil {
		return err
	}

	// writing latest db version to file
	latest := []byte(fmt.Sprintf("%d", len(usableEntries)))
	outFile, err2 := os.OpenFile(m.path+"/migrate.log", os.O_RDWR, 0777)
	if err2 != nil {
		return err
	}
	_, err2 = outFile.Write(latest)
	if err2 != nil {
		return err
	}

	err = outFile.Close()
	if err != nil {
		return err
	}

	err = m.txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (m *Migrate) parseFilesAndMigrateDb() error {
	for _, file := range m.migrationFiles {
		filePath := m.path + "/" + file.Dir.Name()
		fmt.Printf("Reading File %s\n", file.Dir.Name())
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		content := string(bytes)
		commands := sqlparser.ParseSqlFile(content)
		for _, command := range commands {
			_, err = m.txn.Exec(command)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
