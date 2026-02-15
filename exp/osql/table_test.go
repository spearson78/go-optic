package osql_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/exp/osql"
	_ "modernc.org/sqlite"
)

//go:generate ../../makelens osql_test table_test.go table_generated_test.go

type Person struct {
	id   int    `db:"ID" osql:"PK"`
	name string `db:"NAME"`
	age  int    `db:"AGE"`
}

type PersonHobby struct {
	id       int `db:"ID" osql:"PK"`
	personId int `db:"PERSONID"`
	hobbyId  int `db:"HOBBYID"`
}

type Hobby struct {
	id   int    `db:"ID" osql:"PK"`
	name string `db:"NAME"`
}

func OpenTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "")
	if err != nil {
		return nil, err
	}

	return db, InitTestDB(db)
}

func InitTestDB(db *sql.DB) error {

	_, err := db.Exec(`create table PERSON (ID INTEGER PRIMARY KEY,NAME TEXT,AGE INT);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`insert into PERSON (ID,NAME,AGE) VALUES(1,'Max',46);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`insert into PERSON (ID,NAME,AGE) VALUES(2,'Erika',44);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`create table HOBBY (ID INT,NAME TEXT);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`insert into HOBBY (ID,NAME) VALUES(1,'Coding');`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`insert into HOBBY (ID,NAME) VALUES(2,'Snowboarding');`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`create table PERSONHOBBY (ID INT,PERSONID INT,HOBBYID INT);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`insert into PERSONHOBBY (ID,PERSONID,HOBBYID) VALUES(1,1,1);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`create table TREEENTRY (ID INT,PARENTID INT,NAME STRING);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`insert into TREEENTRY (ID,PARENTID,NAME) VALUES(1,NULL,"Root");`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`insert into TREEENTRY (ID,PARENTID,NAME) VALUES(2,1,"Child");`)
	if err != nil {
		return err
	}

	return nil
}

type db struct {
}

func (d db) Person() dbPerson {
	return dbPerson{
		Optic: Table[int, Person]("PERSON"),
	}
}

func (d db) Hobby() dbHobby {
	return dbHobby{
		Optic: Table[int, Hobby]("HOBBY"),
	}
}

type dbPerson struct {
	Optic[Void, *sql.DB, *sql.DB, Collection[int, Person, Err], Collection[int, Person, Err], ReturnMany, ReadWrite, UniDir, Err]
}

type dbHobby struct {
	Optic[Void, *sql.DB, *sql.DB, Collection[int, Hobby, Err], Collection[int, Hobby, Err], ReturnMany, ReadWrite, UniDir, Err]
}

func (d dbPerson) Traverse() *lPerson[int, *sql.DB, *sql.DB, ReturnMany, ReadWrite, UniDir, Err] {
	return OPersonOf(
		RetM(Rw(Ud(EErr(Compose(
			d.Optic,
			TraverseColE[int, Person, Err](),
		))))),
	)
}

func (s *lPerson[I, S, T, RET, RW, DIR, ERR]) PersonHobby() *lPersonHobby[Void, S, T, ReturnMany, RW, UniDir, Err] {
	return OPersonHobbyOf(
		RetM(RwL(Ud(EErr(Compose(
			s.Optic,
			JoinM[Person, PersonHobby]("PERSONHOBBY", "PERSONID"),
		))))),
	)
}

func (s *lPersonHobby[I, S, T, RET, RW, DIR, ERR]) Hobby() *lHobby[Void, S, T, ReturnMany, RW, UniDir, Err] {
	return OHobbyOf(
		RetM(RwL(Ud(EErr(Compose(
			s.Optic,
			Join[PersonHobby, Hobby]("HOBBY", "HOBBYID"),
		))))),
	)
}

var DB db

func TestSelect(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SliceOf(
			DB.Person().Traverse(),
			2,
		),
		db,
	)

	if err != nil || !reflect.DeepEqual(res, []Person{
		{
			id:   1,
			name: "Max",
			age:  46,
		},
		{
			id:   2,
			name: "Erika",
			age:  44,
		},
	}) {
		t.Fatal(res, err)
	}
}

func TestSelectGetFirstI(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ix, res, ok, err := GetFirstI(
		ComposeLeft(
			Filtered(
				DB.Person().Traverse(),
				O.Person().Age().Eq(46),
			),
			O.Person().Name(),
		),
		db,
	)

	if !ok || err != nil || !reflect.DeepEqual([]any{ix, res}, []any{
		1, "Max",
	}) {
		t.Fatal(ix, res, ok, err)
	}
}

func TestSelectSeqIEOf(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SeqIEOf(
			ComposeLeft(
				Filtered(
					DB.Person().Traverse(),
					O.Person().Age().Eq(46),
				),
				O.Person().Name(),
			),
		),
		db,
	)
	if err != nil {
		t.Fatal(res, err)
	}

	ok := false
	for v := range res {
		ix, val, err := v.Get()
		if err != nil || !reflect.DeepEqual([]any{ix, val}, []any{
			1, "Max",
		}) {
			t.Fatal(ix, val, err)
		}
		ok = true
	}

	if !ok {
		t.Fatal(ok)
	}
}

func TestSelectSelfIndex(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SeqIEOf(
			ComposeLeft(
				Filtered(
					SelfIndex(DB.Person().Traverse(), EqT2[Person]()),
					O.Person().Age().Eq(46),
				),
				O.Person().Name(),
			),
		),
		db,
	)
	if err != nil {
		t.Fatal(res, err)
	}

	ok := false
	for v := range res {
		ix, val, err := v.Get()
		if err != nil || !reflect.DeepEqual([]any{ix, val}, []any{
			Person{
				id:   1,
				name: "Max",
				age:  46,
			}, "Max",
		}) {
			t.Fatal(ix, val, err)
		}
		ok = true
	}

	if !ok {
		t.Fatal(ok)
	}
}

func TestSelectReIndexed(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SeqIEOf(
			ComposeLeft(
				Filtered(
					ReIndexed(
						DB.Person().Traverse(),
						Mul(10),
						EqT2[int](),
					),
					O.Person().Age().Eq(46),
				),
				O.Person().Name(),
			),
		),
		db,
	)
	if err != nil {
		t.Fatal(res, err)
	}

	ok := false
	for v := range res {
		ix, val, err := v.Get()
		if err != nil || !reflect.DeepEqual([]any{ix, val}, []any{
			10, "Max",
		}) {
			t.Fatal(ix, val, err)
		}
		ok = true
	}

	if !ok {
		t.Fatal(ok)
	}

}

func TestSelectIndex(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SeqIEOf(
			Index(
				DB.Person().Traverse().Name(),
				2,
			),
		),
		db,
	)
	if err != nil {
		t.Fatal(res, err)
	}

	ok := false
	for v := range res {
		ix, val, err := v.Get()
		if err != nil || !reflect.DeepEqual([]any{ix, val}, []any{
			2, "Erika",
		}) {
			t.Fatal(ix, val, err)
		}
		ok = true
	}

	if !ok {
		t.Fatal(ok)
	}

}

func TestSelectSeqEOf(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SeqEOf(
			ComposeLeft(
				Filtered(
					DB.Person().Traverse(),
					O.Person().Age().Eq(46),
				),
				O.Person().Name(),
			),
		),
		db,
	)
	if err != nil {
		t.Fatal(res, err)
	}

	ok := false
	for v := range res {
		val, err := v.Get()
		if err != nil || !reflect.DeepEqual([]any{val}, []any{
			"Max",
		}) {
			t.Fatal(val, err)
		}
		ok = true
	}

	if !ok {
		t.Fatal(ok)
	}

}

func TestSelectMapOf(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		MapOf(
			ComposeLeft(
				Filtered(
					DB.Person().Traverse(),
					O.Person().Age().Eq(46),
				),
				O.Person().Name(),
			),
			2,
		),
		db,
	)
	if err != nil || !reflect.DeepEqual(res, map[int]string{
		1: "Max",
	}) {
		t.Fatal(res, err)
	}

}

func TestSelectFiltered(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SliceOf(
			Filtered(
				DB.Person().Traverse(),
				O.Person().Age().Eq(44),
			),
			2,
		),
		db,
	)

	if err != nil || !reflect.DeepEqual(res, []Person{
		{
			id:   2,
			name: "Erika",
			age:  44,
		},
	}) {
		t.Fatal(res, err)
	}
}

func TestSelectFilteredReshape(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SliceOf(
			Compose(
				Filtered(
					DB.Person().Traverse(),
					O.Person().Age().Eq(44),
				),
				O.Person().Name(),
			),
			2,
		),
		db,
	)

	if err != nil || !reflect.DeepEqual(res, []string{"Erika"}) {
		t.Fatal(res, err)
	}
}

func TestSelectReshapeField(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SliceOf(
			DB.Person().Traverse().Name(),
			2,
		),
		db,
	)

	if err != nil || !reflect.DeepEqual(res, []string{
		"Max",
		"Erika",
	}) {
		t.Fatal(res, err)
	}

}

func TestSelectReshapeT2(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SliceOf(
			Compose(
				DB.Person().Traverse(),
				T2Of(
					O.Person().Name(),
					O.Person().Age(),
				),
			),
			2,
		),
		db,
	)

	if err != nil || !reflect.DeepEqual(res, []lo.Tuple2[string, int]{
		lo.T2(
			"Max",
			46,
		),
		lo.T2(
			"Erika",
			44,
		),
	}) {
		t.Fatal(res, err)
	}

}

func TestSelectJoin(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SliceOf(
			DB.Person().Traverse().PersonHobby().Hobby(),
			2,
		),
		db,
	)

	if err != nil || !reflect.DeepEqual(res, []Hobby{
		{
			id:   1,
			name: "Coding",
		},
	}) {
		t.Fatal(res, err)
	}
}

func TestSelectJoinReshape(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	res, err := Get(
		SliceOf(
			DB.Person().Traverse().PersonHobby().Hobby().Name(),
			2,
		),
		db,
	)

	if err != nil || !reflect.DeepEqual(res, []string{"Coding"}) {
		t.Fatal(res, err)
	}

}

func TestDelete(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = Modify(
		DB.Person(),
		FilteredCol[int, Person, Err](
			EErr(O.Person().Name().Ne("Max")),
		),
		db,
	)

	if err != nil {
		t.Fatal(err)
	}

	res, err := Get(
		SliceOf(
			DB.Person().Traverse().Name(),
			2,
		),
		db,
	)

	if err != nil || !reflect.DeepEqual(res, []string{"Erika"}) {
		t.Fatal(res, err)
	}

}

func TestUpdate(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = Modify(
		Compose(
			Filtered(
				DB.Person().Traverse(),
				O.Person().Name().Eq("Max"),
			),
			O.Person().Age(),
		),
		Add(1),
		db,
	)

	if err != nil {
		t.Fatal(err)
	}

	res, err := Get(
		SliceOf(
			DB.Person().Traverse().Age(),
			2,
		),
		db,
	)

	if err != nil || !reflect.DeepEqual(res, []int{47, 44}) {
		t.Fatal(res, err)
	}

}

func TestUpdateSet(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = Set(
		Compose(
			Filtered(
				DB.Person().Traverse(),
				O.Person().Name().Eq("Max"),
			),
			O.Person().Age(),
		),
		47,
		db,
	)

	if err != nil {
		t.Fatal(err)
	}

	res, err := Get(
		SliceOf(
			DB.Person().Traverse().Age(),
			2,
		),
		db,
	)

	if err != nil || !reflect.DeepEqual(res, []int{47, 44}) {
		t.Fatal(res, err)
	}

}

func TestInsert(t *testing.T) {

	db, err := OpenTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = Modify(
		DB.Person(),
		AppendCol(ColErr(ValCol(
			Person{
				id:   3,
				name: "Alice",
				age:  45,
			},
		))),
		db,
	)

	if err != nil {
		t.Fatal(err)
	}

	res, err := Get(
		SliceOf(
			Filtered(
				DB.Person().Traverse(),
				O.Person().Age().Eq(45),
			),
			1,
		),
		db,
	)

	if err != nil || !reflect.DeepEqual(res, []Person{
		{
			id:   3,
			name: "Alice",
			age:  45,
		},
	}) {
		t.Fatal(res, err)
	}

}
