package expr

import (
	"fmt"
	"os"

	"golang.org/x/text/encoding"
)

type TraverseFile struct {
	OpticTypeExpr
	Perm    os.FileMode
	BufSize int
}

func (e TraverseFile) Short() string {
	return fmt.Sprintf("TraverseFile(%v,%v)", e.Perm, e.BufSize)
}

func (e TraverseFile) String() string {
	return fmt.Sprintf("TraverseFile(%v,%v)", e.Perm, e.BufSize)
}

type Overwrite struct {
	OpticTypeExpr
}

func (e Overwrite) Short() string {
	return "Overwrite()"
}

func (e Overwrite) String() string {
	return "Overwrite()"
}

type ExtentData struct {
	OpticTypeExpr
}

func (e ExtentData) Short() string {
	return "ExtentData()"
}

func (e ExtentData) String() string {
	return "ExtentData()"
}

type Bytes struct {
	OpticTypeExpr
}

func (e Bytes) Short() string {
	return "Bytes()"
}

func (e Bytes) String() string {
	return "Bytes()"
}

type Reader struct {
	OpticTypeExpr
}

func (e Reader) Short() string {
	return "Reader()"
}

func (e Reader) String() string {
	return "Reader()"
}

type SplitFileExpr struct {
	OpticTypeExpr
	SplitOn byte
}

func (e SplitFileExpr) Short() string {
	return fmt.Sprintf("Split(%v)", e.SplitOn)
}

func (e SplitFileExpr) String() string {
	return fmt.Sprintf("Split(%v)", e.SplitOn)
}

type DecodeFileExpr struct {
	OpticTypeExpr
	Encoding encoding.Encoding
}

func (e DecodeFileExpr) Short() string {
	return fmt.Sprintf("DecodeFile(%v)", e.Encoding)
}

func (e DecodeFileExpr) String() string {
	return fmt.Sprintf("DecodeFile(%v)", e.Encoding)
}

type StatExpr struct {
	OpticTypeExpr
}

func (e StatExpr) Short() string {
	return "Stat"
}

func (e StatExpr) String() string {
	return "Stat"
}

type TraverseDirExpr struct {
	OpticTypeExpr
}

func (e TraverseDirExpr) Short() string {
	return "TraverseDirExpr"
}

func (e TraverseDirExpr) String() string {
	return "TraverseDirExpr"
}
