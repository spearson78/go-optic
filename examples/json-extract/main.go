package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/ojson"
)

func main() {

	//Parse arguments
	set := flag.String("set", "", "JSON value to set")
	modify := flag.Bool("modify", false, "enable modify mode. Write lines of JSON to STDOUT to modify each focused element.")
	flag.Parse()

	//Check arguments
	if len(flag.Args()) != 2 {
		log.Fatal("Usage: json-extract <filename> <query expression>")
	}

	if *set != "" && *modify != false {
		log.Fatal("Expected only one of -set and -modify")
	}

	//modifyOp is the modification, if any, we will apply to the JSON.
	var modifyOp Optic[Void, any, any, any, any, ReturnOne, ReadOnly, UniDir, Err]

	if *set != "" {
		//Set flag was provided

		//Parse the set value as JSON
		setVal, err := Get(ojson.ParseString[any](), *set)
		if err != nil {
			log.Fatalf("Could not parse set value %v", err)
		}

		//Use a constant value for the modifyOp which will set all values to this constant
		modifyOp = EErr(Const[any](setVal))
	}

	if *modify {
		//Modify flag was provided

		//We will read lines of JSON from STDIN
		reader := bufio.NewReader(os.Stdin)

		modifyOp = OpE(func(ctx context.Context, focus any) (any, error) {
			//Display the focus to the user
			fmt.Fprintln(os.Stderr, focus)
			//Read the updated value from the user
			text, err := reader.ReadString('\n')
			if err != nil {
				return "", err
			}

			//Parse the value as JSON
			setVal, err := GetContext(ctx, ojson.ParseString[any](), text)

			//and return it to have the new value set in the JSON
			return setVal, err
		})
	}

	//The query expression is a . separated list of expressions
	// "*" means traverse the focused element.
	// A number means access the nth element of an array
	// Any other value is interpreted as the name of a map key.

	//We will parse the query and build an optic that focuses the element described by the expression
	//We start by Parsing the JSON and convert the index to any and normalise the constraints to ensure we
	//can store the complete expression in the optic variable.

	// Optic[any, []byte, []byte, any, any, ReturnMany, ReadWrite, UniDir, Err]
	optic := normaliseConstraints(
		ReIndexed(
			ojson.Parse[any](),
			UpCast[Void, any](),
			EqDeepT2[any](),
		),
	)

	//These optics are re-used several times so we instantiate them.
	traverseJson := ojson.Traverse()
	parseInt := ParseInt[int](10, 32)

	//Iterate over all the segments of the expression
	for segment := range MustGet(
		SeqOf(
			SplitString(regexp.MustCompile(`\.`)), //Split the query expression on .
		),
		flag.Arg(1), //Arg 1 is the query expression
	) {

		if segment == "*" {
			//Traverse the focus of the current optic.
			optic = normaliseConstraints(
				Compose(
					optic,
					traverseJson,
				),
			)
		} else {

			if num, err := Get(parseInt, segment); err == nil {
				//The value is a number

				optic = normaliseConstraints(
					//Traverse the current optic and focus the element with the numbered index
					Index(
						Compose(
							optic,
							traverseJson,
						),
						any(num),
					),
				)
			} else {
				optic = normaliseConstraints(
					//Traverse the current optic and focus the element with the string index
					Index(
						Compose(
							optic,
							traverseJson,
						),
						any(segment),
					),
				)
			}
		}
	}

	//Optic now parses JSON and focuses the final element of the query expression

	//Read the file
	data, err := os.ReadFile(flag.Arg(0)) //Arg 0 is the file name
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}

	if modifyOp == nil {
		//No modifyOp this is a query

		//Get an iter.Seq of the focused elements
		seq, err := Get(SeqEOf(optic), data)
		if err != nil {
			log.Fatal(err)
		}

		for val := range seq {
			//Extract the value and error
			focus, err := val.Get()
			if err != nil {
				log.Fatal(err)
			}
			//Display the focused element
			fmt.Println(focus)
		}
	} else {
		//we have a modifyOp this is a modification

		//Apply the modifyOp to the json
		//The modify op will display the focus to the user and read the users response.
		newValue, err := Modify(optic, modifyOp, data)
		if err != nil {
			log.Fatal(err)
		}

		//Display the updated JSON to the user.
		fmt.Println(string(newValue))

	}
}

func normaliseConstraints[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, ReturnMany, ReadWrite, UniDir, Err] {
	return RetM(Rw(Ud(EErr(o))))
}
