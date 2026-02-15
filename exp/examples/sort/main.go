package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"slices"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"

	_ "net/http/pprof"
)

var readMode = flag.String("mode", "mmap", "file load mode (mmap|buffer|readat|readall,slicesort)")
var enableProfile = flag.Bool("profile", false, "enable cpu and mem profiling")

const bufSize = 32 * 1024

func main() {

	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatal("expected in-file and out-file paraneters")
	}

	var memProfile *os.File

	if *enableProfile {
		cpuProfile, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		memProfile, err = os.Create("memprofile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(cpuProfile)
	}

	switch *readMode {
	case "mmap":
		opticSort(oio.ReadMMap)
	case "buffer":
		opticSort(oio.ReadBuffer)
	case "readat":
		opticSort(oio.ReadAt)
	case "readall":
		opticSort(oio.ReadAll)
	case "slicesort":
		sliceSort()
	default:
		log.Fatal("unknown mode", *readMode)
	}

	if *enableProfile {
		pprof.StopCPUProfile()
		pprof.WriteHeapProfile(memProfile)
	}

}

func opticSort(mode oio.ReadMode) {
	stringData := Compose(oio.ExtentData(), oio.BytesString())
	fileOptic := Compose(oio.TraverseFile(0644, bufSize, mode), oio.SplitFile('\n'))
	ordrBy := EErr(OrderBy(stringData))
	sort := OrderedCol[int64](ordrBy)

	_, err := Modify(fileOptic, sort, oio.FileNames{
		InFile:  flag.Arg(0),
		OutFile: flag.Arg(1),
	})

	if err != nil {
		log.Fatal(err)
	}
}

func sliceSort() {

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	splits := bytes.Split(data, []byte{'\n'})

	slices.SortFunc(splits, func(a []byte, b []byte) int {
		if string(a) > string(b) {
			return 1
		}

		if string(a) < string(b) {
			return -1
		}

		return 0
	})

	writer, err := os.OpenFile(flag.Arg(1), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	joined := bytes.Join(splits, []byte{'\n'})

	writer.Write(joined)
}
