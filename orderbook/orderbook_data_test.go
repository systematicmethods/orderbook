package orderbook

import (
	"encoding/csv"
	"github.com/gocarina/gocsv"
	"io"
	"orderbook/assert"
	"runtime"
	"strings"
	"testing"
)

type execCSV struct {
	Id        string    `csv:"id"`
	Clientid  string    `csv:"clientid"`
	Clordid   string    `csv:"clordid"`
	Side      Side      `csv:"side"`
	Lastprice float64   `csv:"lastprice"`
	Lastqty   int64     `csv:"lastqty"`
	Status    OrdStatus `csv:"status"`
}

func (ex *Side) MarshalCSV() (string, error) {
	return ex.String(), nil
}

func (ex *Side) UnmarshalCSV(field string) (err error) {
	*ex = SideConv(field)
	return nil
}

func (ex *OrdStatus) MarshalCSV() (string, error) {
	return ex.String(), nil
}

func (ex *OrdStatus) UnmarshalCSV(field string) (err error) {
	*ex = OrdStatusConv(field)
	return nil
}

func loadExecCSV(csvstr string) []*execCSV {
	execs := []*execCSV{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '|'
		return r
	})

	if gocsv.Unmarshal(strings.NewReader(csvstr), &execs) == nil {
		return execs
	}
	return execs
}

func containsExecCSV(t *testing.T, execs []ExecutionReport, ex *execCSV, msg string) {
	var found = 0
	for _, v := range execs {
		if v.ClientID() == ex.Clientid &&
			v.ClOrdID() == ex.Clordid &&
			v.OrdStatus() == ex.Status &&
			v.LastPrice() == ex.Lastprice &&
			v.LastQty() == ex.Lastqty {
			found++
		}
	}
	if found == 0 {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: not found %s id:'%v' %s:%s stat:%v qty:%d price:%f",
			assert.AssertionAt(file), line, msg, ex.Id, ex.Clientid, ex.Clordid, ex.Status, ex.Lastqty, ex.Lastprice)
	}
}

func Test_CSVReader(t *testing.T) {
	var expected = `id|clientid|clordid|side|lastprice|lastqty|status
e0|cli2|id2|sell|1|50|PartiallyFilled
e1|cli1|id21|buy|1|50|PartiallyFilled
e2|cli2|id2|sell|1.01|50|Filled
e3|cli1|id21|buy|1.01|50|PartiallyFilled
e4|cli2|id3|sell|1|1|PartiallyFilled
e5|cli1|id21|buy|1|1|Filled
`
	execs := []*execCSV{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '|'
		return r
	})

	err := gocsv.Unmarshal(strings.NewReader(expected), &execs)
	assert.AssertNilT(t, err, "no errors")
	assert.AssertEqualT(t, 6, len(execs), "len not equal")
	assert.AssertEqualT(t, "e0", execs[0].Id, "Id")
	assert.AssertEqualT(t, "cli2", execs[0].Clientid, "Clientid")
	assert.AssertEqualT(t, "id2", execs[0].Clordid, "Clordid")
	assert.AssertEqualT(t, SideSell, execs[0].Side, "Side")
	assert.AssertEqualT(t, float64(1), execs[0].Lastprice, "Lastprice")
	assert.AssertEqualT(t, int64(50), execs[0].Lastqty, "Lastqty")
	assert.AssertEqualT(t, OrdStatusPartiallyFilled, execs[0].Status, "Status")

	//for _, exec := range execs {
	//	fmt.Printf("clientid %v\n", exec.Clientid)
	//}

}

func Test_CSVReaderExtraCols(t *testing.T) {
	var expected = `id|clientid|clordid|side|lastprice|lastqty|status|other
e0|cli2|id2|sell|1|50|PartiallyFilled|aa
e1|cli1|id21|buy|1|50|PartiallyFilled|bb
`
	execs := []*execCSV{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '|'
		return r
	})

	err := gocsv.Unmarshal(strings.NewReader(expected), &execs)
	assert.AssertNilT(t, err, "no errors")
	assert.AssertEqualT(t, 2, len(execs), "len not equal")
	assert.AssertEqualT(t, "cli2", execs[0].Clientid, "Clientid")
	assert.AssertEqualT(t, "id2", execs[0].Clordid, "Clordid")
	assert.AssertEqualT(t, SideSell, execs[0].Side, "Side")
	assert.AssertEqualT(t, float64(1), execs[0].Lastprice, "Lastprice")
	assert.AssertEqualT(t, int64(50), execs[0].Lastqty, "Lastqty")
	assert.AssertEqualT(t, OrdStatusPartiallyFilled, execs[0].Status, "Status")

}

func Test_CSVReaderEmptyCols(t *testing.T) {
	var expected = `id|clientid|clordid|side|lastprice|lastqty|status|other
e0|cli2|id2|sell||50|PartiallyFilled|
e1|cli1|id21|abc|1||PartiallyFilled|
e1|cli1|id21||1|50|PartiallyFilled|bb
`
	execs := []*execCSV{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '|'
		return r
	})

	err := gocsv.Unmarshal(strings.NewReader(expected), &execs)
	assert.AssertNilT(t, err, "no errors")
	assert.AssertEqualT(t, 3, len(execs), "len not equal")
	assert.AssertEqualT(t, "cli2", execs[0].Clientid, "Clientid")
	assert.AssertEqualT(t, "id2", execs[0].Clordid, "Clordid")
	assert.AssertEqualT(t, SideSell, execs[0].Side, "Side")
	assert.AssertEqualT(t, float64(0), execs[0].Lastprice, "Lastprice")
	assert.AssertEqualT(t, int64(50), execs[0].Lastqty, "Lastqty")
	assert.AssertEqualT(t, OrdStatusPartiallyFilled, execs[0].Status, "Status")

	assert.AssertEqualT(t, int64(0), execs[1].Lastqty, "Lastqty")
	assert.AssertEqualT(t, SideUnknown, execs[1].Side, "Side")
	assert.AssertEqualT(t, SideUnknown, execs[2].Side, "Side")
}

func Test_CSVWriter_ok(t *testing.T) {
	var expected = `id|clientid|clordid|side|lastprice|lastqty|status
e0|cli2|id2|sell|1|50|PartiallyFilled
e1|cli1|id21|buy|1|50|PartiallyFilled
`
	execs := []*execCSV{}
	execs = append(execs, &execCSV{"e0", "cli2", "id2", SideSell, 1, 50, OrdStatusConv("PartiallyFilled")},
		&execCSV{"e1", "cli1", "id21", SideBuy, 1, 50, OrdStatusConv("PartiallyFilled")})

	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = '|'
		return gocsv.NewSafeCSVWriter(writer)
	})

	csvContent, _ := gocsv.MarshalString(&execs)

	assert.AssertEqualT(t, expected, csvContent, "csv")
}

func Test_CSVWriter_error(t *testing.T) {
	var expected = `id|clientid|clordid|side|lastprice|lastqty|status
e0|cli2|id2|sell|1|50|PartiallyFilled
e1|cli1|id21|buy|1|50|PartiallyFilled
e3|cli1|id21|buy|1|50|PartiallyFilled
`
	execs := []*execCSV{}
	execs = append(execs, &execCSV{"e0", "cli2", "id2", SideSell, 1, 50, OrdStatusConv("PartiallyFilled")},
		&execCSV{"e1", "cli1", "id21", SideBuy, 1, 50, OrdStatusConv("PartiallyFilled")})

	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = '|'
		return gocsv.NewSafeCSVWriter(writer)
	})

	csvContent, _ := gocsv.MarshalString(&execs)

	assert.AssertNotEqualT(t, expected, csvContent, "csv")

}
