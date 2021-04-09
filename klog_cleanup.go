package klog

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// cleanup clean up files by age and then by number.
func (sb *syncBuffer) cleanup(tag string, t time.Time) error {
	cleanupFilesByAge(logging.logDir, logNamePrefix(tag), logging.logFileMaxAge)
	cleanupFilesByNumber(logging.logDir, logNamePrefix(tag), int(logging.logFileMaxNumber))
	return nil
}

func logNamePrefix(tag string) string {
	return fmt.Sprintf("%s.%s.%s.log.%s",
		program,
		host,
		getUserName(),
		tag,
	)
}

func cleanupFilesByAge(dir, prefix string, maxAge uint64) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	ds, err := d.Readdir(0)
	if err != nil {
		return err
	}
	for _, f := range ds {
		if !strings.HasPrefix(f.Name(), prefix) {
			continue
		}
		t := time.Now().Sub(f.ModTime())
		expect := time.Duration(maxAge) * time.Minute
		fmt.Printf("name: %#v, t: %v(expect: %v), old: %v\n", f.Name(), t, expect, t >= expect)
	}
	return nil
}

func cleanupFilesByNumber(dir, prefix string, maxNumber int) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	list, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	names := []string{}
	for _, name := range list {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		fmt.Printf("exist name: %v\n", name)
		names = append(names, name)
	}
	n := len(names)
	if n <= maxNumber {
		// no need clean
		return nil
	}

	todelete := names[:n-maxNumber-1]
	for _, name := range todelete {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
