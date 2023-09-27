//go:build single

package conformance

import (
	"fmt"
	"testing"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pkg/client"
	"github.com/xh3b4sd/redigo/pkg/sorted"
)

// Test_Client_Single_Sorted_Create_Order ensures that indices are guaranteed to
// be unique. Below the indices c and d cannot be duplicated. Indices may be
// used to ensure unique usernames.
func Test_Client_Single_Sorted_Create_Order(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "bar", 0.7, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "baz", 0.6, "c", "d")
		if !sorted.IsAlreadyExistsError(err) {
			t.Fatal("expected", "alreadyExistsError", "got", err)
		}
	}

	{
		err := cli.Sorted().Delete().Index("ssk", "bar")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "baz", 0.6, "c", "d")
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
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "bar", 0.7, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "zap", 0.8, "e", "f")
		if !sorted.IsAlreadyExistsError(err) {
			t.Fatal("expected", "alreadyExistsError", "got", err)
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8, "g", "h")
		if !sorted.IsAlreadyExistsError(err) {
			t.Fatal("expected", "alreadyExistsError", "got", err)
		}
	}
}

func Test_Client_Single_Sorted_Create_Value(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Value("ssk", "foo", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Value("ssk", "bar", 0.7)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Verify we can create elements with duplicated scores.
	{
		err := cli.Sorted().Create().Value("ssk", "zap", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Verify values must be unique after all.
	{
		err := cli.Sorted().Create().Value("ssk", "foo", 0.8)
		if !sorted.IsAlreadyExistsError(err) {
			t.Fatal("expected", "alreadyExistsError", "got", err)
		}
	}

	{
		res, err := cli.Sorted().Search().Order("ssk", 0, -1)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 3 {
			t.Fatal("expected", 3, "got", len(res))
		}
		if res[0] != "zap" {
			t.Fatal("expected", "zap", "got", res[0])
		}
		if res[1] != "foo" {
			t.Fatal("expected", "foo", "got", res[1])
		}
		if res[2] != "bar" {
			t.Fatal("expected", "bar", "got", res[2])
		}
	}

	// Deleting an element by value should not delete elements with the same
	// score. When foo is deleted, which has score 0.8 then zap must still exist
	// with the same score as we verify in the next step.
	{
		err := cli.Sorted().Delete().Index("ssk", "foo")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Order("ssk", 0, -1)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 2 {
			t.Fatal("expected", 2, "got", len(res))
		}
		if res[0] != "zap" {
			t.Fatal("expected", "zap", "got", res[0])
		}
		if res[1] != "bar" {
			t.Fatal("expected", "bar", "got", res[1])
		}
	}
}

func Test_Client_Single_Sorted_Delete_Empty(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Delete().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Delete().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Delete().Index("ssk", "foo")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Delete().Index("ssk", "foo")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8, "a", "b")
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
		emp, err := cli.Empty()
		if err != nil {
			t.Fatal(err)
		}

		if emp {
			t.Fatal("storage must not be empty")
		}
	}

	{
		err := cli.Sorted().Delete().Clean("ssk")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		emp, err := cli.Empty()
		if err != nil {
			t.Fatal(err)
		}

		if !emp {
			t.Fatal("storage must be empty")
		}
	}
}

func Test_Client_Single_Sorted_Delete_Limit(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	var hea []float64
	{
		hea = []float64{
			0.1,
			0.2,
			0.3,
			0.4,
		}
	}

	var tai []float64
	{
		tai = []float64{
			0.5,
			0.6,
			0.7,
		}
	}

	{
		for _, s := range append(hea, tai...) {
			exi, err := cli.Sorted().Exists().Score("ssk", s)
			if err != nil {
				t.Fatal(err)
			}
			if exi {
				t.Fatal("element must not exist")
			}
		}
	}

	{
		err = cli.Sorted().Create().Score("ssk", "a", 0.1)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "b", 0.2)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "c", 0.3)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "d", 0.4)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "e", 0.5)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "f", 0.6)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "g", 0.7)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		for _, s := range append(hea, tai...) {
			exi, err := cli.Sorted().Exists().Score("ssk", s)
			if err != nil {
				t.Fatal(err)
			}
			if !exi {
				t.Fatal("element must exist")
			}
		}
	}

	{
		err := cli.Sorted().Delete().Limit("ssk", 3)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		for _, s := range hea {
			exi, err := cli.Sorted().Exists().Score("ssk", s)
			if err != nil {
				t.Fatal(err)
			}
			if exi {
				t.Fatal("element must not exist")
			}
		}
	}

	{
		for _, s := range tai {
			exi, err := cli.Sorted().Exists().Score("ssk", s)
			if err != nil {
				t.Fatal(err)
			}
			if !exi {
				t.Fatal("element must exist")
			}
		}
	}

	{
		err := cli.Sorted().Delete().Clean("ssk")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Test_Client_Single_Sorted_Delete_Score(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
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
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8)
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
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8, "a", "b")
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
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Delete().Index("ssk", "foo")
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
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8, "a", "b")
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
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Index("ssk", "a")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
		}
	}

	{
		exi, err := cli.Sorted().Exists().Index("ssk", "b")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
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
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Index("ssk", "a")
		if err != nil {
			t.Fatal(err)
		}
		if !exi {
			t.Fatal("element must exist")
		}
	}

	{
		exi, err := cli.Sorted().Exists().Index("ssk", "b")
		if err != nil {
			t.Fatal(err)
		}
		if !exi {
			t.Fatal("element must exist")
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
		err := cli.Sorted().Delete().Index("ssk", "foo")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Index("ssk", "a")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
		}
	}

	{
		exi, err := cli.Sorted().Exists().Index("ssk", "b")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
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

func Test_Client_Single_Sorted_Floats(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Value("ssk", "a")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
		}
	}

	var flo float64

	{
		flo, err = cli.Sorted().Floats().Score("ssk", "a", +0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if flo != 0.8 {
			t.Fatal("flo must be 0.8")
		}
	}

	{
		flo, err = cli.Sorted().Floats().Score("ssk", "a", +0.6)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if flo != 1.4 {
			t.Fatal("flo must be 1.4")
		}
	}

	{
		flo, err = cli.Sorted().Floats().Score("ssk", "a", -0.4)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Redis returns 0.9999999999999999 so we round for the test.
	{
		if fmt.Sprintf("%.1f", flo) != "1.0" {
			t.Fatal("flo must be 1.0")
		}
	}

	{
		exi, err := cli.Sorted().Exists().Value("ssk", "a")
		if err != nil {
			t.Fatal(err)
		}
		if !exi {
			t.Fatal("element must exist")
		}
	}

	{
		err := cli.Sorted().Delete().Index("ssk", "a")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Value("ssk", "a")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatal("element must not exist")
		}
	}
}

func Test_Client_Single_Sorted_Metric(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		cou, err := cli.Sorted().Metric().Count("ssk")
		if err != nil {
			t.Fatal(err)
		}
		if cou != 0 {
			t.Fatal("count must not be 0")
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "a", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		cou, err := cli.Sorted().Metric().Count("ssk")
		if err != nil {
			t.Fatal(err)
		}
		if cou != 1 {
			t.Fatal("count must not be 1")
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "b", 0.7)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		cou, err := cli.Sorted().Metric().Count("ssk")
		if err != nil {
			t.Fatal(err)
		}
		if cou != 2 {
			t.Fatal("count must not be 2")
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "c", 0.6)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		cou, err := cli.Sorted().Metric().Count("ssk")
		if err != nil {
			t.Fatal(err)
		}
		if cou != 3 {
			t.Fatal("count must not be 3")
		}
	}
}

func Test_Client_Single_Sorted_Search_Index(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Index("ssk", "a")
		if err != nil {
			t.Fatal(err)
		}
		if res != "" {
			t.Fatal("expected", "empty string", "got", res)
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Index("ssk", "a")
		if err != nil {
			t.Fatal(err)
		}
		if res != "foo" {
			t.Fatal("expected", "foo", "got", res)
		}
	}

	{
		res, err := cli.Sorted().Search().Index("ssk", "b")
		if err != nil {
			t.Fatal(err)
		}
		if res != "foo" {
			t.Fatal("expected", "foo", "got", res)
		}
	}
}

func Test_Client_Single_Sorted_Search_Inter(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Inter("k1", "k2")
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 0 {
			t.Fatal("expected", 0, "got", len(res))
		}
	}

	{
		err = cli.Sorted().Create().Score("k1", "v3", 0.3)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("k1", "v4", 0.4)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("k1", "v5", 0.5)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("k1", "v6", 0.6)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Inter("k1")
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 4 {
			t.Fatal("expected", 4, "got", len(res))
		}
		if res[0] != "v3" {
			t.Fatal("expected", "v3", "got", res[0])
		}
		if res[1] != "v4" {
			t.Fatal("expected", "v4", "got", res[1])
		}
		if res[2] != "v5" {
			t.Fatal("expected", "v5", "got", res[2])
		}
		if res[3] != "v6" {
			t.Fatal("expected", "v6", "got", res[3])
		}
	}

	{
		res, err := cli.Sorted().Search().Inter("k1", "k2")
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 0 {
			t.Fatal("expected", 0, "got", len(res))
		}
	}

	{
		err = cli.Sorted().Create().Score("k2", "v2", 0.2)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("k2", "v4", 0.4)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("k2", "v5", 0.5)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("k2", "v7", 0.7)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Inter("k1", "k2")
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 2 {
			t.Fatal("expected", 2, "got", len(res))
		}
		if res[0] != "v4" {
			t.Fatal("expected", "v4", "got", res[0])
		}
		if res[1] != "v5" {
			t.Fatal("expected", "v5", "got", res[1])
		}
	}
}

func Test_Client_Single_Sorted_Search_Order(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Order("ssk", 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 0 {
			t.Fatal("expected", 0, "got", len(res))
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "foo", 6.0, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Order("ssk", 0, 0)
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
		err := cli.Sorted().Create().Score("ssk", "bar", 7.0, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
	}

	// Ensure to get the latest value, that is, the value of the element with
	// the highest score.
	{
		res, err := cli.Sorted().Search().Order("ssk", 0, 0)
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

	// Ensure to get the penultimate value, that is, the value of the element
	// with the second highest score.
	{
		res, err := cli.Sorted().Search().Order("ssk", 1, 1)
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
		res, err := cli.Sorted().Search().Order("ssk", 0, -1)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 2 {
			t.Fatal("expected", 2, "got", len(res))
		}
		if res[0] != "bar" {
			t.Fatal("expected", "bar", "got", res[0])
		}
		if res[1] != "foo" {
			t.Fatal("expected", "foo", "got", res[1])
		}
	}

	{
		res, err := cli.Sorted().Search().Order("ssk", 0, -1, true)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 2 {
			t.Fatal("expected", 2, "got", len(res))
		}
		if res[0] != "7" {
			t.Fatal("expected", "7", "got", res[0])
		}
		if res[1] != "6" {
			t.Fatal("expected", "6", "got", res[1])
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "baz", 8.0, "e", "f")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Order("ssk", 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 1 {
			t.Fatal("expected", 1, "got", len(res))
		}
		if res[0] != "baz" {
			t.Fatal("expected", "baz", "got", res[0])
		}
	}

	{
		res, err := cli.Sorted().Search().Order("ssk", 0, 0, true)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 1 {
			t.Fatal("expected", 1, "got", len(res))
		}
		if res[0] != "8" {
			t.Fatal("expected", "8", "got", res[0])
		}
	}

	{
		res, err := cli.Sorted().Search().Order("ssk", -1, -1)
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
}

func Test_Client_Single_Sorted_Search_Rando(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		val, err := cli.Sorted().Search().Rando("ssk")
		if err != nil {
			t.Fatal(err)
		}
		if len(val) != 0 {
			t.Fatal("expected", "zero values", "got", len(val))
		}
	}

	{
		err = cli.Sorted().Create().Score("ssk", "foo", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "bar", 0.7)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "baz", 0.6)
		if err != nil {
			t.Fatal(err)
		}
	}

	var lis []string
	{
		for i := 0; i < 100; i++ {
			val, err := cli.Sorted().Search().Rando("ssk")
			if err != nil {
				t.Fatal(err)
			}

			if len(val) != 1 {
				t.Fatal("expected", "one value", "got", len(val))
			}

			lis = append(lis, val[0])
		}
	}

	{
		if !contains(lis, "foo") {
			t.Fatal("expected", "lis to contain 'foo'", "got", "not found")
		}
		if !contains(lis, "bar") {
			t.Fatal("expected", "lis to contain 'bar'", "got", "not found")
		}
		if !contains(lis, "baz") {
			t.Fatal("expected", "lis to contain 'baz'", "got", "not found")
		}
	}
}

func Test_Client_Single_Sorted_Search_Rando_Cou(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		val, err := cli.Sorted().Search().Rando("ssk")
		if err != nil {
			t.Fatal(err)
		}
		if len(val) != 0 {
			t.Fatal("expected", "zero values", "got", len(val))
		}
	}

	{
		err = cli.Sorted().Create().Score("ssk", "a", 0.9)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "b", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "c", 0.7)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "d", 0.6)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "e", 0.5)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "f", 0.4)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "g", 0.3)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "h", 0.2)
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Sorted().Create().Score("ssk", "i", 0.1)
		if err != nil {
			t.Fatal(err)
		}
	}

	var lis []string
	{
		lis, err = cli.Sorted().Search().Rando("ssk", 4)
		if err != nil {
			t.Fatal(err)
		}
	}

	if len(lis) != 4 {
		t.Fatal("expected", "four values", "got", len(lis))
	}

	var cou int
	{
		if contains(lis, "a") {
			cou++
		}
		if contains(lis, "b") {
			cou++
		}
		if contains(lis, "c") {
			cou++
		}
		if contains(lis, "d") {
			cou++
		}
		if contains(lis, "e") {
			cou++
		}
		if contains(lis, "f") {
			cou++
		}
		if contains(lis, "g") {
			cou++
		}
		if contains(lis, "h") {
			cou++
		}
		if contains(lis, "i") {
			cou++
		}
	}

	if cou != 4 {
		t.Fatal("expected", "four values", "got", cou)
	}
}

func Test_Client_Single_Sorted_Search_Value(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Value("ssk", 0.8, 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 0 {
			t.Fatal("expected", 0, "got", len(res))
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Value("ssk", 0.8, 0.8)
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
		err := cli.Sorted().Create().Score("ssk", "bar", 0.7, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Value("ssk", 0.7, 0.7)
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
		res, err := cli.Sorted().Search().Value("ssk", 0.8, 0.7)
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
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Score("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Sorted().Search().Value("ssk", 0.8, 0.8)
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
		err := cli.Sorted().Create().Score("ssk", "baz", 0.6, "a", "b")
		if !sorted.IsAlreadyExistsError(err) {
			t.Fatal("expected", "alreadyExistsError", "got", err)
		}
	}

	{
		res, err := cli.Sorted().Update().Index("ssk", "bar", 0.8, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
		if !res {
			t.Fatal("element must be updated")
		}
	}

	{
		res, err := cli.Sorted().Update().Index("ssk", "bar", 0.8, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
		if res {
			t.Fatal("element must not be updated")
		}
	}

	{
		res, err := cli.Sorted().Search().Value("ssk", 0.8, 0.8)
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
		err := cli.Sorted().Create().Score("ssk", "baz", 0.6, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func contains(lis []string, itm string) bool {
	for _, l := range lis {
		if l == itm {
			return true
		}
	}

	return false
}
