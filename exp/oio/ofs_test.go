package oio_test

import (
	"io"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/exp/oio"
	"github.com/spearson78/go-optic/otree"
)

var testFsMutex sync.Mutex

func initTestFS() func() {

	testFsMutex.Lock()

	os.Remove("test/renamed.txt")
	f, err := os.Create("test/rename.txt")
	if err != nil {
		testFsMutex.Unlock()
		panic(err)
	}
	f.Write([]byte("Hello World"))
	f.Close()

	f, err = os.Create("test/extents.txt")
	if err != nil {
		testFsMutex.Unlock()
		panic(err)
	}
	f.Write([]byte(`gamma
beta
alpha
delta`))
	f.Close()

	os.Remove("test/test.txt")
	f, err = os.Create("test/test.txt")
	if err != nil {
		testFsMutex.Unlock()
		panic(err)
	}
	f.Write([]byte("test"))
	f.Close()

	return testFsMutex.Unlock
}

func TestDirTree(t *testing.T) {

	defer initTestFS()()

	if r, err := Get(SliceOf(Compose3(Stat(), TraverseFileInfo(), FileInfoName()), 10), "./test/"); err != nil || !reflect.DeepEqual(r, []string{"extents.txt", "rename.txt", "sub", "test.txt"}) {
		t.Fatalf(`Children(), FileInfoName() : %v %v`, r, err)
	}

	if r, err := Get(SliceOf(Compose3(Stat(), TraverseFileInfo(), FileInfoFullPath()), 10), "./test/"); err != nil || !reflect.DeepEqual(r, []string{"test/extents.txt", "test/rename.txt", "test/sub", "test/test.txt"}) {
		t.Fatalf(`Children(), FileInfoFullPath() : %v %v`, r, err)
	}

	if r, err := Get(
		SliceOf(
			Compose(
				Stat(),
				Compose(
					WithIndex(
						otree.TopDown(TraverseFileInfo()),
					),
					ValueIIndex[*otree.PathNode[string], FileInfo](),
				),
			),
			10,
		),
		"./test/",
	); err != nil || !MustGet(EqDeepT2[[]*otree.PathNode[string]](), lo.T2(r, []*otree.PathNode[string]{
		nil,
		otree.Path("extents.txt"),
		otree.Path("rename.txt"),
		otree.Path("sub"),
		otree.Path("sub", "sub.txt"),
		otree.Path("test.txt"),
	})) {
		t.Fatalf(`ToSliceOf(TopDownCol(), "./test/") : %v`, r)
	}

	if r, err := Get(
		SliceOf(
			Compose3(
				Stat(),
				otree.TopDown(TraverseFileInfo()),
				FileInfoFullPath(),
			),
			10,
		),
		"./test/",
	); err != nil || !MustGet(EqDeepT2[[]string](), lo.T2(r, []string{
		"test",
		"test/extents.txt",
		"test/rename.txt",
		"test/sub",
		"test/sub/sub.txt",
		"test/test.txt",
	})) {
		t.Fatalf(`ToSliceOf(TopDownCol(), "./test/") : %v`, r)
	}

	if r, err := Get(
		SliceOf(
			Compose3(
				Stat(),
				WithIndex(
					otree.BottomUp(TraverseFileInfo()),
				),
				ValueIIndex[*otree.PathNode[string], FileInfo](),
			),
			10,
		),
		"./test/",
	); err != nil || !MustGet(EqDeepT2[[]*otree.PathNode[string]](), lo.T2(r, []*otree.PathNode[string]{
		otree.Path("extents.txt"),
		otree.Path("rename.txt"),
		otree.Path("sub", "sub.txt"),
		otree.Path("sub"),
		otree.Path("test.txt"),
		nil,
	})) {
		t.Fatalf(`ToSliceOf(BottomUpCol(), "./test/") : %v`, r)
	}
}

func TestFileBytes(t *testing.T) {

	defer initTestFS()()

	optic := ForEach(
		EErr(Compose3(
			Stat(),
			Filtered(
				TraverseFileInfo(),
				Compose(FileInfoIsDir(), Not()),
			),
			FileInfoFullPath(),
		)),
		EErr(Compose(
			FileBytes(
				0644,
				8192,
				ReadAll,
			),
			IsoCast[[]byte, string](),
		)),
	)

	if fileContent, err := Get(SliceOf(Ordered(optic, OrderBy(Identity[string]())), 10), "test"); err != nil || !reflect.DeepEqual(fileContent, []string{"Hello World", "gamma\nbeta\nalpha\ndelta", "test"}) {
		t.Fatalf(`FileBytes : %v`, fileContent)
	}

	modifiedFiles, err := Set(optic, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	if fileContent, err := Get(SliceOf(optic, 10), "test"); err != nil || !reflect.DeepEqual(fileContent, []string{"test", "test", "test"}) {
		t.Fatalf(`Set : %v`, fileContent)
	}

	equal, err := Get(
		EqCol(
			ValColIE[Void, string, Err](IxMatchVoid(), ValIE(Void{}, "test/extents.txt", nil), ValIE(Void{}, "test/rename.txt", nil), ValIE(Void{}, "test/test.txt", nil)),
			EqT2[string](),
		),
		modifiedFiles,
	)

	if err != nil || !equal {
		t.Fatalf(`Set modified files: %v`, modifiedFiles)
	}
}

func testExtents(t *testing.T, mode ReadMode) {

	defer initTestFS()()

	stringData := Compose(ExtentData(), UpCast[[]byte, string]())

	fileLines := Compose3(TraverseFileOverwrite(0666, 8192, mode), SplitFile('\n'), TraverseColE[int64, Extent, Err]())

	optic := Compose(fileLines, stringData)

	if fileContent, err := Get(SliceOf(optic, 10), "./test/extents.txt"); err != nil || !reflect.DeepEqual(fileContent, []string{
		"gamma",
		"beta",
		"alpha",
		"delta",
	}) {
		t.Fatalf(`Extent : %v`, fileContent)
	}

	sortedOptic := Compose(Ordered(fileLines, OrderBy(stringData)), stringData)

	if sortedContent, err := Get(SliceOf(sortedOptic, 10), "./test/extents.txt"); err != nil || !reflect.DeepEqual(sortedContent, []string{
		"alpha",
		"beta",
		"delta",
		"gamma",
	}) {
		t.Fatalf(`Sorted Extent : %v`, sortedContent)
	}

}

func TestExtents(t *testing.T) {
	testExtents(t, ReadAt)
	testExtents(t, ReadAll)
	testExtents(t, ReadMMap)
}

func TestWriteExtents(t *testing.T) {

	defer initTestFS()()

	stringData := Compose(ExtentData(), UpCast[[]byte, string]())
	fileOptic := Compose(TraverseFileOverwrite(0666, 8192, ReadAll), SplitFile('\n'))
	ordrBy := EErr(OrderBy(stringData))
	sort := OrderedCol[int64](ordrBy)

	res, err := Modify(fileOptic, sort, "./test/extents.txt")
	if err != nil || res != "./test/extents.txt" {
		t.Fatal(res, err)
	}

	f, err := os.Open("./test/extents.txt")
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != `alpha
beta
delta
gamma` {
		t.Fatal(string(data))
	}

}
