package get_minimized

import (
    "testing"
	"github.com/deckarep/golang-set"
	"fmt"
	"sort"
    "strings"
)

type testpair struct {
    sigma  map[string]map[byte]string
    Q mapset.Set
    F mapset.Set
	E mapset.Set
	P mapset.Set
	q0 string
	expectedP string
	expectedDfa string
}

func SortString(w string) string {
    s := strings.Split(w, "")
    sort.Strings(s)
    return strings.Join(s, "")
}


var ParseRegTestSet = []testpair{
    {
        map[string]map[byte]string {
		"1": map[byte]string{'0': "2", '1': "3"},
		"2": map[byte]string{'0': "2", '1': "3"},
		"3": map[byte]string{'0': "4", '1': "5"},
		"4": map[byte]string{'0': "6", '1': "7"},
		"5": map[byte]string{'0': "2", '1': "3"},
		"6": map[byte]string{'0': "8", '1': "9"},
		"7": map[byte]string{'0': "6", '1': "7"},
		"8": map[byte]string{'0': "6", '1': "7"},
		"9": map[byte]string{'0': "10", '1': "11"},
		"10": map[byte]string{'0': "10", '1': "11"},
		"11": map[byte]string{'0': "8", '1': "9"}},
	    mapset.NewSet("1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"),
	    mapset.NewSet("4", "7", "8"),
		mapset.NewSet(byte('0'),byte('1')),
		mapset.NewSetFromSlice([]interface{}{
			mapset.NewSetFromSlice([]interface{}{"4", "8", "7"}),
			mapset.NewSetFromSlice([]interface{}{"10", "2", "5", "9", "1:"}), 
			mapset.NewSetFromSlice([]interface{}{"3", "6", "11"}),
		}),
		"1",
		SortString("Set{Set{4, 8, 7}, Set{10, 2, 5, 9, 1}, Set{3, 6, 11}}"),
		"Set{48, 49}\nSet{8;7;4;, 5;9;1:;10;2;, 3;6;11;}\n\nSet{7;4;8;}\nmap[11;3;6;:map[48:7;4;8; 49:1:;10;2;5;9;] 1:;10;2;5;9;:map[48:1:;10;2;5;9; 49:11;3;6;] 7;4;8;:map[48:11;3;6; 49:7;4;8;]]\n",
    },
}

func TestGetClassesRegex(t *testing.T) {
	for num, val := range ParseRegTestSet {

		P := fmt.Sprintf("%v", Minimize(val.Q, val.F, val.E, val.sigma))
        if SortString(P) != SortString(val.expectedP) {
            t.Error(
                "For", num, 
                "expected", val.expectedP,
                "got", P,
            )
        }
    }
}

func TestGetNewDfa(t *testing.T) {
	for num, val := range ParseRegTestSet {
		nQ, nSigma, nq0, nF := NewDfa(val.sigma, val.q0, val.F , val.P, val.E)
		res := DFAToString(nQ, nSigma, nq0, nF, val.E)
        if SortString(res) != SortString(val.expectedDfa) {
            t.Error(
                "For", num, 
                "expected", val.expectedDfa,
                "got", res,
            )
        }
    }
}