package cedar

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"io"
	"os"
)

// Save saves the cedar to an io.Writer,
// where dataType is either "json" or "gob".
func (da *Cedar) Save(out io.Writer, dataType string) error {
	switch dataType {
	case "gob", "GOB":
		dataEecoder := gob.NewEncoder(out)
		return dataEecoder.Encode(da)
	case "json", "JSON":
		dataEecoder := json.NewEncoder(out)
		return dataEecoder.Encode(da)
	}
	return ErrInvalidDataType
}

// SaveToFile saves the cedar to a file,
// where dataType is either "json" or "gob".
func (da *Cedar) SaveToFile(fileName string, dataType string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	out := bufio.NewWriter(file)
	defer out.Flush()
	return da.Save(out, dataType)
}

// Load loads the cedar from an io.Writer,
// where dataType is either "json" or "gob".
func (da *Cedar) Load(in io.Reader, dataType string) error {
	switch dataType {
	case "gob", "GOB":
		dataDecoder := gob.NewDecoder(in)
		return dataDecoder.Decode(da)
	case "json", "JSON":
		dataDecoder := json.NewDecoder(in)
		return dataDecoder.Decode(da)
	}
	return ErrInvalidDataType
}

// LoadFromFile loads the cedar from a file,
// where dataType is either "json" or "gob".
func (da *Cedar) LoadFromFile(fileName string, dataType string) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
	defer file.Close()
	if err != nil {
		return err
	}
	in := bufio.NewReader(file)
	return da.Load(in, dataType)
}
