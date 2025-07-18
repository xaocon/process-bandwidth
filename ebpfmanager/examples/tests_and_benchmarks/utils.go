package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
)

// recoverAssets - Recover ebpf asset
func recoverAssets() io.ReaderAt {
	buf, err := Asset("/probe.o")
	if err != nil {
		logrus.Fatal(fmt.Errorf("error:%v , couldn't find asset", err))
	}
	return bytes.NewReader(buf)
}
