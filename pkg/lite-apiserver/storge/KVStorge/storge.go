package KVStorge

import "github.com/syndtr/goleveldb/leveldb"

func (dw *DBWriter) write(key string, value []byte) error {
	db, err := leveldb.OpenFile(dw.StorgeFileName, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Put([]byte(key), value, nil)
}

// return "",ErrNotFound if key is not exist
func (dw *DBWriter) read(key string) ([]byte, error) {
	db, err := leveldb.OpenFile(dw.StorgeFileName, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	data, err := db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Delete will not returns error if key doesn't exist.
func (dw *DBWriter) delete(key string) error {
	db, err := leveldb.OpenFile(dw.StorgeFileName, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Delete([]byte(key), nil)
}

func (dw *DBWriter) has(key string) (bool, error) {
	db, err := leveldb.OpenFile(dw.StorgeFileName, nil)
	if err != nil {
		return false, err
	}
	defer db.Close()

	return db.Has([]byte(key), nil)
}
