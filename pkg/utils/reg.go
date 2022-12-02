package utils

import "regexp"

func DeleteExtraSpace(s string) string {
	regstr := "\\s{2,}"
	reg, _ := regexp.Compile(regstr)
	tmpstr := make([]byte, len(s))
	copy(tmpstr, s)
	spc_index := reg.FindStringIndex(string(tmpstr))
	for len(spc_index) > 0 {
		tmpstr = append(tmpstr[:spc_index[0]+1], tmpstr[spc_index[1]:]...)
		spc_index = reg.FindStringIndex(string(tmpstr))
	}
	return string(tmpstr)
}
