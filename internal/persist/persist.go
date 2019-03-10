package persist

import (
	"bufio"
	"github.com/vmihailenco/msgpack"
	"io/ioutil"
	"os"
	"sync"
)

// https://medium.com/@matryer/golang-advent-calendar-day-eleven-persisting-go-objects-to-disk-7caf1ee3d11d

var lock sync.Mutex

func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	r, err := Marshal(v)
	if err != nil {
		return err
	}
	_, err = writer.Write(r)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return err
}

func Load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return Unmarshal(fileContent, v)
}

func Marshal(v interface{}) ([]byte, error) {
	b, err := msgpack.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Unmarshal(b []byte, v interface{}) error {
	return msgpack.Unmarshal(b, v)
}
