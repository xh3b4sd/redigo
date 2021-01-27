// +build single

package client

import (
	"testing"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/sorted"
)

func Test_Client_Single_Sorted_Delete_Score(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{
			Kind: KindSingle,
		}

		cli, err = New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if !exi {
			t.Fatal("element must exist")
		}
	}

	{
		err := cli.Sorted().Delete().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if !exi {
			t.Fatal("element must exist")
		}
	}

	// We just created an element that defined the indices a and b. Now we
	// delete this very element only using its score. With this test we ensure
	// that elements as well as their associated indices get automatically
	// purged when deleting elements only using their score.
	{
		err := cli.Sorted().Delete().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	// It should be possible to create the exact same element again including
	// its indizes after it has been deleted before. This verifies that deleting
	// elements including its indizes works as expected.
	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Delete().Value("ssk", "foo", "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if !exi {
			t.Fatal("element must exist")
		}
	}

	{
		err := cli.Sorted().Delete().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Test_Client_Single_Sorted_Exists(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{
			Kind: KindSingle,
		}

		cli, err = New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
		}
	}

	{
		exi, err := cli.Sorted().Exists().Value("ssk", "foo")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if !exi {
			t.Fatal("element must exist")
		}
	}

	{
		exi, err := cli.Sorted().Exists().Value("ssk", "foo")
		if err != nil {
			t.Fatal(err)
		}
		if !exi {
			t.Fatal("element must exist")
		}
	}

	{
		err := cli.Sorted().Delete().Value("ssk", "foo")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
		}
	}

	{
		exi, err := cli.Sorted().Exists().Value("ssk", "foo")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
		}
	}
}

// Test_Client_Single_Sorted_Create_Index ensures that indices are guaranteed to
// be unique. Below the indices c and d cannot be duplicated. Indices may be
// used to ensure unique usernames.
func Test_Client_Single_Sorted_Create_Index(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{
			Kind: KindSingle,
		}

		cli, err = New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "bar", 0.7, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "baz", 0.6, "c", "d")
		if !sorted.IsAlreadyExistsError(err) {
			t.Fatal("expected", alreadyExistsError, "got", err)
		}
	}

	{
		err := cli.Sorted().Delete().Value("ssk", "bar", "c", "d")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "baz", 0.6, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
	}
}

// Test_Client_Single_Sorted_Create_Score ensures that scores are guaranteed to
// be unique. Below the score 0.8 cannot be duplicated. Scores may be used to
// represent IDs as unix timestamps.
func Test_Client_Single_Sorted_Create_Score(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{
			Kind: KindSingle,
		}

		cli, err = New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "bar", 0.7, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "zap", 0.8, "e", "f")
		if !sorted.IsAlreadyExistsError(err) {
			t.Fatal("expected", alreadyExistsError, "got", err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8, "g", "h")
		if !sorted.IsAlreadyExistsError(err) {
			t.Fatal("expected", alreadyExistsError, "got", err)
		}
	}
}

func Test_Client_Single_Sorted_Search_Index(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{
			Kind: KindSingle,
		}

		cli, err = New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Index("ssk", 0, 1)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 1 {
			t.Fatal("expected", 1, "got", len(res))
		}
		if res[0] != "foo" {
			t.Fatal("expected", "foo", "got", res[0])
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "bar", 0.7, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Index("ssk", 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 1 {
			t.Fatal("expected", 1, "got", len(res))
		}
		if res[0] != "bar" {
			t.Fatal("expected", "bar", "got", res[0])
		}
	}

	{
		res, err := cli.Sorted().Search().Index("ssk", 0, -1)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 2 {
			t.Fatal("expected", 2, "got", len(res))
		}
		if res[0] != "foo" {
			t.Fatal("expected", "foo", "got", res[0])
		}
		if res[1] != "bar" {
			t.Fatal("expected", "bar", "got", res[1])
		}
	}
}

func Test_Client_Single_Sorted_Search_Score(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{
			Kind: KindSingle,
		}

		cli, err = New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Score("ssk", 0.8, 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 1 {
			t.Fatal("expected", 1, "got", len(res))
		}
		if res[0] != "foo" {
			t.Fatal("expected", "foo", "got", res[0])
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "bar", 0.7, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Score("ssk", 0.7, 0.7)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 1 {
			t.Fatal("expected", 1, "got", len(res))
		}
		if res[0] != "bar" {
			t.Fatal("expected", "bar", "got", res[0])
		}
	}

	{
		res, err := cli.Sorted().Search().Score("ssk", 0.8, 0.7)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 2 {
			t.Fatal("expected", 2, "got", len(res))
		}
		if res[0] != "foo" {
			t.Fatal("expected", "foo", "got", res[0])
		}
		if res[1] != "bar" {
			t.Fatal("expected", "bar", "got", res[1])
		}
	}
}

func Test_Client_Single_Sorted_Update(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{
			Kind: KindSingle,
		}

		cli, err = New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Score("ssk", 0.8, 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 1 {
			t.Fatal("expected", 1, "got", len(res))
		}
		if res[0] != "foo" {
			t.Fatal("expected", "foo", "got", res[0])
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "baz", 0.6, "a", "b")
		if !sorted.IsAlreadyExistsError(err) {
			t.Fatal("expected", alreadyExistsError, "got", err)
		}
	}

	{
		res, err := cli.Sorted().Update().Value("ssk", "bar", 0.8, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
		if !res {
			t.Fatal("element must be updated")
		}
	}

	{
		res, err := cli.Sorted().Update().Value("ssk", "bar", 0.8, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
		if res {
			t.Fatal("element must not be updated")
		}
	}

	{
		res, err := cli.Sorted().Search().Score("ssk", 0.8, 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 1 {
			t.Fatal("expected", 1, "got", len(res))
		}
		if res[0] != "bar" {
			t.Fatal("expected", "bar", "got", res[0])
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "baz", 0.6, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}
}
