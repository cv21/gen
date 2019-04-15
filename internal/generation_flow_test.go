package internal

import (
	"path/filepath"
	"testing"
)

func TestGenerationFlow_ValidateResultPath_OutOfPath(t *testing.T) {
	cd, err := filepath.Abs(".")
	if err != nil {
		t.Error(err)
	}

	gf := basicGenerationFlow{
		currentDir: cd,
	}

	xx, err := filepath.Abs("../../../hello_gen.go")
	if err != nil {
		t.Error(err)
	}

	err = gf.ValidateResultPath(xx)
	if err == nil {
		t.Error("wait error, but got nil")
	}
}

func TestGenerationFlow_ValidateResultPath_GenPrefix(t *testing.T) {
	cd, err := filepath.Abs(".")
	if err != nil {
		t.Error(err)
	}

	gf := basicGenerationFlow{
		currentDir: cd,
	}

	xx, err := filepath.Abs("./bla/bla")
	if err != nil {
		t.Error(err)
	}

	err = gf.ValidateResultPath(xx)
	if err == nil {
		t.Error(err)
	}
}

func TestGenerationFlow_ValidateResultPath_Success(t *testing.T) {
	cd, err := filepath.Abs(".")
	if err != nil {
		t.Error(err)
	}

	gf := basicGenerationFlow{
		currentDir: cd,
	}

	xx, err := filepath.Abs("./bla/bla_gen.go")
	if err != nil {
		t.Error(err)
	}
	err = gf.ValidateResultPath(xx)
	if err != nil {
		t.Error(err)
	}
}
